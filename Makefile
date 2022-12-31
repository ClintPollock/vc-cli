ifdef CI_COMMIT_TAG
	VERSION := ${CI_COMMIT_TAG}
else
	VERSION := "0.0.0"
endif

ifdef CI_COMMIT_SHORT_SHA
	GIT_HASH := ${CI_COMMIT_SHORT_SHA}
else
	GIT_HASH := $(shell git rev-parse --short HEAD)
endif

all:  insecure-test/.done verascanner/.done veracode/veracode

veracode/veracode:
	cd veracode && \
	go build -o veracode -ldflags "-X github.com/veracode/veracode-cli/cmd/version.Version=$(VERSION) -X github.com/veracode/veracode-cli/cmd/version.GitHash=$(GIT_HASH)"

insecure-test: insecure-test/.done

verascanner: verascanner/.done

insecure-test/.done:
	cd insecure-test; docker build . -f Dockerfile -t veray-insecure; touch .done ; cd -

verascanner/.done:
	cd verascanner; docker build . -t verascanner; touch .done ; cd -

test:
	./verascan trivy image  --security-checks secret vera-insecure:latest

clean:
	rm -f insecure-test/.done verascanner/.done veracode/veracode

distclean: clean
	docker rmi verascanner veray-insecure

build-push-aruba:
	cd docker && \
	docker build . -f base/Dockerfile -t policy/veracode-cli/base && \
	docker tag policy/veracode-cli/base docker-rw.laputa.veracode.io/policy/veracode-cli/base && \
	docker push docker-rw.laputa.veracode.io/policy/veracode-cli/base

build-test:
	cd veracode && gox -osarch="linux/amd64" -output "../bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -ldflags "-X github.com/veracode/veracode-cli/cmd/version.Version=$(VERSION) -X github.com/veracode/veracode-cli/cmd/version.GitHash=$(GIT_HASH)"
