FROM aquasec/trivy:0.31.3 as trivy
FROM openpolicyagent/opa:0.44.0-dev-static as opa
FROM openpolicyagent/conftest:v0.34.0 as conftest
FROM goodwithtech/dockle:v0.4.6 as dockle
FROM anchore/syft:v0.56.0 as syft
FROM anchore/grype:v0.50.1 as grype
FROM python:3.10.6-alpine3.16

RUN mkdir /trivy /trivy/contrib
COPY --from=trivy /usr/local/bin/trivy /bin
COPY --from=trivy /contrib/*.tpl /trivy/contrib/

COPY --from=syft /syft /bin
COPY --from=grype /grype /bin
COPY --from=dockle /usr/bin/dockle /bin/

RUN pip install tern

RUN apk add --no-cache bash curl jq docker git openssh-client skopeo
RUN apk del curl

RUN mkdir /opa
COPY --from=opa      /opa /bin/
COPY --from=conftest /conftest /bin/

COPY syft.yaml /etc/syft/config.yaml
COPY grype.yaml /etc/grype/config.yaml
COPY policy /opa/policy
COPY opa-dockerfile-data.yaml /opa/

RUN grype db update

RUN mkdir /data /local-context
WORKDIR /data

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY entrypoint.py /entrypoint.py
ENTRYPOINT ["python", "/entrypoint.py"]
