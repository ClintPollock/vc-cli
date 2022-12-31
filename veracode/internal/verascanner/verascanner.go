package verascanner

import (
	"context"

	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

const (
	runImageName = "verascanner:latest"
)

func ImagePull() error {

	dockerCtx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(dockerCtx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	imagePresent := false
	for _, image := range images {
		for _, tag := range image.RepoTags {
			//fmt.Println(tag + " " + image.ID)
			if !imagePresent {
				imagePresent = (tag == runImageName)
			}
		}
	}

	if !imagePresent {
		reader, err := cli.ImagePull(dockerCtx, runImageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}

		if reader != nil {
			defer reader.Close()
			io.Copy(os.Stdout, reader)
		}
	}

	return err
}

func ContainerRun(ctx context.Context, args []string) (io.ReadCloser, int64, error) {

	err := ImagePull()

	dockerCtx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// user, err := user.Current()
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// cacheDir := user.HomeDir + "/.veracode/cache" // This should become an API

	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(dockerCtx,
		&container.Config{
			Image: runImageName,
			Cmd:   args,
			Tty:   false,
		},
		&container.HostConfig{
			//AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
				{
					Type:   mount.TypeBind,
					Source: "/var/lib/docker",
					Target: "/root/.cache/",
				},
				{
					Type:   mount.TypeVolume,
					Source: "cache",
					Target: "/data",
				},
				{
					Type: mount.TypeVolume,
					// This has been hardcoded instead of the relative path from previous implementation. Absolute path is required in this field for Volume type.
					Source: "veracode-cli",
					Target: "/local-context",
				},
			},
		},
		nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(dockerCtx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	var retval int64
	statusCh, errCh := cli.ContainerWait(dockerCtx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case status := <-statusCh:
		retval = status.StatusCode
		//fmt.Printf("exit status: %d", status.StatusCode)

	}

	out, err := cli.ContainerLogs(dockerCtx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	err = cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		panic(err)
	}

	return out, retval, err

}
