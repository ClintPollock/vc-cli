# Overview

Testing CLI against a set of User Scenarios with Aruba

[Aruba](https://github.com/cucumber/aruba) provides a set of easy to use, easy to understand, toolset to write tests in a logical language that can be mapped to a user guide. Aruba allows us to directly run feature tests against the `veracode` CLI application in a Behavior-Driven Development (BDD) approach.


# Context

The `veracode` CLI application is a frontend application written in Golang that would call a set of tools that performs various Container Security scans (eg. IaC Scan, Vuln Scan, Secrets Scan). The set of tools that does the heavy lifting of Container Security is built as a Docker Image, and is currently known as [verascanner](https://gitlab.laputa.veracode.io/policy/veracode-cli/-/tree/main/verascanner).


# Pre-requisites to run the test

A proper Aruba test setup *should* include:

1. The `veracode` CLI binary
2. The `verascanner` Docker image
3. An API backend

> Note: The API backend responses can be mocked in Aruba, as long as it returns what an actual API server would return (eg. Auth success/Auth failure).


The following are **required** to be installed on the host to run the setup:

1. Golang installed with GOPATH setup locally
2. Docker installed, and the Docker daemon is running
3. Make installed

> Note: Gox does not currently support Go 1.19 (See https://github.com/mitchellh/gox/issues/165)


# Running the tests locally

Follow the steps below to run your Aruba tests locally:

1. Change Directory to the *root* directory of this project, instead of in this `./features` folder.
1. Install [gox](https://github.com/mitchellh/gox): `go install github.com/mitchellh/gox@latest`. `gox` will be used to build the binary to the expected format.
2. Run `make build-test`. `build-test` is added to run `gox` to generate the binary in the format suitable for testing environment.
3. Run `./scripts/aruba.sh`. This should bring you into a running Docker container with `bash` ready to execute commands.
5. Run `cucumber -p ci` to run the feature tests.

> Note: You may have to ensure that the location of the `veracode` CLI binary is included in the `PATH` environment variable


# Running the tests in CI/CD Pipeline

Minimally, our pre-requisites for running the tests are the `veracode` binary, and the `verascanner` Docker image. So, we should ensure that the CI/CD pipeline is able to:

1. Build the `veracode` CLI binary
2. Build the `verascanner` Docker image
3. Add the location of the `veracode` CLI binary to the `PATH` environment variable


## Building the veracode CLI binary

To build the `veracode` CLI binary, and to run the `cucumber` Aruba tests, use the [manual publish-base CI/CD stage](../.gitlab/ci/publish-images.gitlab-ci.yml) to build and publish the `base` image to our [Container Registry](https://gitlab.laputa.veracode.io/policy/veracode-cli/container_registry). The `base` image comes with the necessary packages installed to build the `veracode` CLI binary, as well as to run the `cucumber` Aruba tests.

When using the `base` image in the CI/CD pipeline, ensure that the `GOPATH` has been set, and is included in `PATH` to allow `gox` to run. Also ensure that the `PATH` environment variable is updated with the CLI binary's location (eg. export PATH=$PATH:$PWD/bin).


## Building the verascanner Docker image

To build the `verascanner` Docker image for CI/CD testing, use the [manual publish-verascanner CI/CD stage](../.gitlab/ci/publish-images.gitlab-ci.yml) to build and publish the `verascanner` image to our [Container Registry](https://gitlab.laputa.veracode.io/policy/veracode-cli/container_registry).

Whenever required, such as when the required image has been deleted, or if there is a newer update to `verascanner`, the image can be rebuilt through the manual action.

> A rebuild of the verascanner image is required whenever there are changes to the verascanner toolchain, such as updating the version of a tool (eg. Updating Grype to the latest version), or updating its behaviour of the tool (eg. Changes in output format)


Finally, the tests are run in the [test stage](../.gitlab/ci/test.gitlab-ci.yml) where it would do the following:

1. Load up the `base` image
2. Pull in the `verascanner` image
3. Build the `veracode` CLI binary
4. Run `cucumber` Aruba tests

> Note: Currently, a local copy of `alpine:latest` was saved and loaded. This is not required and can be removed deemed unnecessary. It has been added to run the tests against a copy of an immutable image.


# Platform Specific Gotchas

This section describes a little gotcha experienced when writing a few sample tests. For context, one of the early tests written was a JSON output equivalence test when performing `veracode sbom image alpine:latest`. The image was kept the same throughout, enforced by loading in our own image tarball, yet there were differences in SBOM output when run on the local machine, and when run in the CI/CD Gitlab Runner. 

While almost everything runs in a Docker container, there are minor differences in the recognition of Architecture. When run on an M1 Mac, the Architecture might be `arm64`/`aarch64`, while it might appear as `x86_64` when run on the CI/CD pipeline.

As such, we now have multiple `<architecture>` folders in the `features/resources/expectedoutputs` folder. The correct `architecture` can be resolved in `features/steps.rb`, before the `cucumber` tests are being run. That way, we can allow the tests to be run normally without extra flags or consideration.

> Currently, the method used to collect the Docker Host's architecture is through `docker info --format "{{.Architecture}}"`.


# Aruba

[Aruba](https://github.com/cucumber/aruba) runs Feature tests that are written in a form of logical language made through [Gherkin](https://cucumber.io/docs/gherkin/reference/).

Tests are formed from a set of special keywords such as:

- Feature
- Rule (as of Gherkin 6)
- Example (or Scenario)
- Given, When, Then, And, But for steps (or *)
- Background
- Scenario Outline (or Scenario Template)
- Examples (or Scenarios)

## Test directory

Aruba's tests sits in the `features` folder, and each feature can be placed in a `<test>.feature` file.

The folder structure would look like this:

```
projectroot/
  features/
    sbom.feature
    scan.feature
    ...
```

## Structure of a test.feature file

A basic structure should contain a few keywords:

```
Feature:
  Scenario:
    Given
    When
    Then
```

> Note: The indentation of 2 spaces between Feature-Scenario, and Scenario-Given/When/Then/etc, is important

Also, a `Feature` can contain multiple `Scenario`:


```
Feature: Feature X

  Scenario: S1
    Given A
    When B
    Then C

  Scenario: S2
    Given D
    When E
    Then F
```


## Walkthrough of an example Aruba test:

Consider the following examples:

```
Feature: Generate SBOM in a CycloneDX JSON structure

  Scenario: Perform SBOM Scan on an image with --format=cyclonedx-json
    Given an activated CLI agent
    When I run `veracode sbom image alpine:latest --format=cyclonedx-json`
    And the scan should succeed
    Then the output should match json in "features/resources/expectedoutputs/alpine-image.json"

  Scenario: Perform SBOM Scan on a tarball with --format=cyclonedx-json
    Given an activated CLI agent
    When I run `veracode sbom archive alpine-latest.tar.gz --format=cyclonedx-json`
    And the scan should succeed
    Then the output should match json in "features/resources/expectedoutputs/alpine-tarball.json"
```

> For an in-depth description of each keyword, view the [Gherkin Reference](https://cucumber.io/docs/gherkin/reference/)

The sections `Feature:` and `Scenario:` can be written in free-form. 

However, in each of the `Given:`/`When:`/`And:`/`Then:` step of a `Scenario:`, such as `an activated CLI agent`, Aruba would require a intermediate step to bridge between the described task, and the actual code representation that would perform the test. 

These definitions are defined in the `features/steps.rb` file.

For example:

```
features/steps.rb

...

Then("the output should match json in {string}") do |file|
  expected_json = JSON.parse(File.read(file))
  actual_json = JSON.parse(last_command_started.output)
  expect(Normalize.normalize(actual_json)).to match(Normalize.normalize(expected_json))
end
```

For reference to more examples on Aruba Features:
- https://github.com/cucumber/aruba/tree/main/features
- https://gitlab.laputa.veracode.io/sca/srcclr-agent/-/blob/master/features/json.feature


