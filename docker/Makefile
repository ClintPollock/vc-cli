# Builds and pushes the docker image defined in the Makefile target
# Please ensure that you are logged into docker-rw.laputa.veracode.io
# registry. For this, create a personal access token on Laputa Gitlab,
# copy it to the clipboard and do:
#
# pbpaste | docker login -u <username> --stdin-password docker-rw.laputa.veracode.io
#
# where <username> is your Laputa Gitlab user name.
#
define build-push
docker build . -f $@/Dockerfile -t policy/veracode-cli/$@ && \
docker tag policy/veracode-cli/$@ docker-rw.laputa.veracode.io/policy/veracode-cli/$@ && \
docker push docker-rw.laputa.veracode.io/policy/veracode-cli/$@
endef

.PHONY: base

all: base

base:
	$(build-push)