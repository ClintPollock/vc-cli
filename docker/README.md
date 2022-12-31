# Dockerfiles

This is where we store our dockerfiles that we use for our functional and future tests.

The images are pushed to the registry and the image names follows the way the
directories that contains the Dockerfile are named relative to this directory.
For example, the Dockerfile for dotnet regression tests is contained in
`base` then the image is named with the tag
`base` and it is pushed to the registry with this tag.

## Build

To build and push all the images described by the Dockerfiles in this
directory, execute the following command.

```
make
```
