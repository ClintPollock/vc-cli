# Overview

A PoC swiss army knife of OSS software for container scanning.

This includes for now:

1. `verascan` a lightweight wrapper script around docker (deprecated).
2. A go CLI that is currently a framework, again intended to launch docker.
3. A docker container in `verascanner/` that builds a toolbox of commands and open source libraries.
4. An insecure test container.

The prerequisites are:

1) A MacOS, Windows or Linux environment (e.g. Ubuntu on WSL2).
2) [Go](https://go.dev/doc/install).
3) [Docker runtime](https://docs.docker.com/engine/install/).
4) Make (should already be installed because of prerequisite 1).

Installation instructions for Windows:

- Install [chocolatey](https://chocolatey.org/install)
- On powershell with admin privileges, run `choco install make`
- Make sure docker is running on the machine.
- In Git Bash, Run `make all` from the root of the repository to build the project.
- `cd veracode`
- `go build -o veracode`
- Run `chmod +x veracode` to make the file executable.
- Now the Veracode CLI binary is built into the `/veracode` folder and you can run the below commands.

Installation instructions for Mac or Linux:

- Make sure docker is running on the machine.
- Run `make all` from the root of the repository to build the project.
- cd veracode && go build -o veracode
- Now the Veracode CLI binary is built into the `/veracode` folder and you can run the below commands.

# CLI Commands

The Veracode CLI binary is built into the `/veracode` folder. Change into that folder before running the following commands.

Also note that on Linux systems, you should prefix the command with `./`, for example:

```Shell
./veracode --help
```

## Configuring the Veracode CLI tool

This allows the configuration of the CLI using the veracode credentials which is found in ~/.veracode/credentials INI file as described in https://docs.veracode.com/r/c_configure_api_cred_file.
This command also validates that the HMAC credentials work and pulls the User resource from the Veracode Identity API associated with the credentials.

```Shell
./veracode init
```

## Inspect

Print out data about the target.

```Shell
veracode inspect [ target type ] [ target ]
```

## SBOM Generation

```Shell
./veracode sbom [ target type ] [ target ]
```

## Scanning

Run the default scan workflow on a scan target this can be an image, a directory, an archive or a git repo.

```Shell
./veracode scan [ target type ] [ target ]
```

This will also run the default policy evaluation against the results, and return success or failure;
the CLI's exit code will reflect this evaluation.

## Kitchen Sink

This runs everything we have available against the target, and returns eveything
in a large JSON structure.

```Shell
./veracode scan [ target type ] [ target ] --kitchen-sink
```

## Clear

This clears out the cache directory.

```Shell
./veracode clear
```

## Running Tools Directly

This is useful to bypass the workflows in the logic and just get direct output of the underlying tools used by each command.

```Shell
./veracode run [ tool ] [[ args ]]
```

# Schema

The overall schema of the tool is JSON follows this format:

```json
{
  "target" : {
    "type": "...",
    "raw_name": "...",
    "name": "...",
    "id": "1a2b3c4d..."
  },
  "inventory" : {
    "target_id": "sha256:e66264b98777e12192600bf9b4d663655c98a090072e1bab49e233d7531d1294",
    "target_type": "image",
    ...
    "sbom": { .... },
    "files" : [ ... ]
  },
  "findings": {
    "vulnerabilities" : { ... }
    "docker": { ... }
    "secrets": { ... }
    "iac" : { ... }
    "permissions" : [ ... ]
  },
  "findings-by-tool": {
    "TOOLNAME": {
      "vulnerabilities" :  { ... }
      ...
    }
  },
  "policy-results": [
    {
      "filename": "sha256:e66264b98777e12192600bf9b4d663655c98a090072e1bab49e233d7531d1294.results.json",
      "namespace": "main",
      "successes": 3
    }
  ],
  "policy-passed": true

```

Depending on the command, the full schema or a subset is returned. If a different
format is requested that will be returned, but this still reflects the internal
data schema.

## Target

This primary document includes the target name provided, a normalized or sanitized name,
the type (image, repo, etc.) and a unique ID derived from the object. It may also
includes metadata about a container image under the "docker" heading.

## Inventory

This includes data about the contents of the target, including an SBOM in a json
format. It also includes the full file manifest with file metadata as well.

## Findings

This includes findings against the target objects, including known vulnerabilities,
image/conatiner configuration issues, secrets, IaC issues, and file permissions.

These findings are currently based on a single tool.

## Findings By Tool

For research purposes, the `findings-by-tool` section can be invluded when very
verbose modes are chosen from the command line optinons. This is structured
by tool name then finding type.

## Policy

The final section report out the policy evaluatio failures in the "policy-results"
section, and in "policy-passed" the PASS/FAIL status of the findings section based on the active
policy (which for now is hardcoded but could be configurable in the future).

This last value is intended to be used as the final output for the tool in CI/CD
pipelines and related use cases.

# Internals

The go CLI for now is for the most part a wrapper around calling docker with a specific docker image and setup. The called docker image is located under the `verascanner` directory, is built on Alpine Linux and embeds various tools under `/usr/bin`.

It also has the following volume mounts in place to enable scans that are outside the container itself:

1. The docker socket is mounted.
2. The docker lib cache is mounted to `/root/.cache`.
3. The cache directory under `~/.veracode/cache` is mounted under `/data` and the normal working directory.
4. The current working directory is mounted as `/local/context`.

The entrypoint scripts enables you to call various tools explicitly, and if called as `bash` will let you exec arbitrary scripts as if you set the "command" directive (this is of course not how we want to deliver but useful in this early stage).

For example to see the default working directory you can run `veracode run bash pwd` and to see its contents run `veracode run bash ls -l`.

# Building

The provided Makefile has the following primary targets:

```Shell
make clean
make all
```

This builds the go CLI wrapper, driver routine and docker images. Use `make distClean` to clean out docker images as well. Open the `Makefile` for more details on building specific targets.

To build each components individually you can do the following.

## Build Scanner Image

```Shell
cd verascanner
docker build . -t verascanner:latest
cd -
```

## Build a Test Image

```Shell
cd insecure-test
docker build . -f Dockerfile -t vera-insecure:latest
```

# Testing

## Feature Testing

We are using Aruba to simulate and test user scenarios when interacting with the CLI. The requirements are as follows:

1. Install `gox` https://github.com/mitchellh/gox. `gox` will be used to build the binary to the expected format.
2. Run `make build-test`. `build-test` is added to run `gox` and generate the binary in the format suitable for testing environment.
3. Be sure to have Docker running.
4. Then run `./scripts/aruba.sh` which brings you into a running container.
5. Run `cucumber -p ci` to run the feature tests.

> Note: Gox does not currently support Go 1.19 (See https://github.com/mitchellh/gox/issues/165)

# Using the tool to analyze databases

You can use the tool to analyze results from various vulnerability databases.

For example to list out the unique CVE Ids of the vulnerabilityies in a container image,
you can run the scan command and parse the results using `jq`:

```Shell
veracode scan image alpine:3.10 | jq .findings.vulnerabilities.matches[].vulnerability.id | sort | uniq
```

To print out results from grype and trivy use the `research` subcommand

```Shell
veracode research image alpine:3.9 | jq .vulnerability_ids
```

For example this shows that grype finds more vulnerability in alpine that trivy.

To see what files are present in an image, you can run

```Shell
veracode sbom image gcr.io/distroless/static-debian11  | jq -r .files[].location.path
```

# Examples of Running Underlying Tools

## Scan for Secrets Using Trivy

```Shell
veracode run trivy image --security-checks secret veray-insecure:latest
```

## Create a SBOM Plus Manifest

```Shell
veracode run syft packages veray-insecure:latest -o json
```

# DEMO

```Shell
veracode sbom image alpine:latest
veracode scan repo https://github.com/apache/pulsar-helm-chart.git

veracode scan archive tests/test.zip
veracode sbom archive tests/test.zip

veracode sbom image verascanner:latest --prettyprint
veracode sbom image verascanner:latest --prettyprint --format spdx
veracode sbom image verascanner:latest --prettyprint --format cyclonedx

veracode scan image veray-insecure:latest --prettyprint
echo $?

veracode scan image verascanner:latest --prettyprint
echo $?
```

## Development

### Unit Testing

This repository utilizes [gotestsum](https://pkg.go.dev/gotest.tools/gotestsum#section-readme)
 to run unit tests currently available in the `veracode` directory. Unit tests follow the
 pattern `*_test.go` and are placed in the same directory as the go module they test. To run
 `gotestsum` locally:

```Shell
# first cd into veracode and then install the package
$ go install gotest.tools/gotestsum@latest
go: downloading gotest.tools/gotestsum v1.8.2
...
# run gotestsum
$ gotestsum
✓  internal/cache (156ms)
∅  .
∅  cmd
∅  cmd/inspect
∅  cmd/research
∅  cmd/run
∅  cmd/sbom
∅  cmd/scan
∅  internal/api/identity/user
∅  internal/hmac
∅  internal/verascanner

DONE 1 tests in 1.633s
```

### Linting

This repository utilizes the [golangci-lint](https://golangci-lint.run/) to lint the
 veracode-cli tool. To run the tool locally and review issues caught by the linter,
 run:

```Shell
 # run this in the veracode directory
 $ golangci-lint run ./...
```

The `lint` job will run in the `build` stage to notify developers
 of any style issues within `veracode-cli/veracode`. To resolve some of these issues
 automatically, run:

```Shell
 # veracode/... tells go fmt to recursively look through all packages in the veracode
 # directory.
 $ go fmt veracode/...
 cmd\clear.go
 cmd\root.go
 ...
```
