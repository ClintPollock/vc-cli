package greenlight

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type Artifact struct {
	Filename     string
	sha256       string
	filenamePath string
	contents     []byte
	size         int64
}

func (a *Artifact) Load(filename string) error {

	a.Filename = filename

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}
	a.sha256 = fmt.Sprintf("%x", h.Sum(nil))

	fileInfo, _ := f.Stat()
	a.filenamePath = fileInfo.Name()
	a.size = fileInfo.Size()
	// contents, err := ioutil.ReadFile(filename)
	// a.size = cap(contents)

	return err

}
