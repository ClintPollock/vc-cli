#!/bin/bash

set +x

if [ "help" = "$1" ]; then

  echo "Execute as docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/docker:/root/.cache/ verascanner [ ARGS ] "
  echo ""
  echo Command options
  echo ""
  echo "help: print this message"
  echo "clear cache: remove cached results files regarding images"
  echo "scan image: do a snazzy full scan and show off stuff ..."
  echo "scan directory: do a snazzy full scan and show off stuff ..."
  echo "scan archive: do a snazzy full scan and show off stuff ..."
  echo "scan repo: do a snazzy full scan and show off stuff ..."
  echo "trivy: execute trivy and use following arguments as args to trivy"
  echo "grype: execute grype and use following arguments as args to grype"
  echo "syft: execute syft and use following arguments as args to syft"
  echo "dockle: execute dockle and use following arguments as args to it"
  echo "skopeo: execute skopeo and use following arguments as args to it"
  echo "docker: execute docker and use following arguments as args to it"
  echo "opa: execute opa and use following arguments as args to it"
  echo "conftest: execute conftest and use following arguments as args to it"
  echo "bash: execute bash and use following arguments as args to it"

  exit 0
fi


if [ "clear" = "$1" ]; then

  shift;

  if [ "cache" = "$1" ]; then

    rm -f *.json .*.json
    rm -f *.txt  .*.txt
    rm -f *.xml  .*.xml

    exit $?
  fi
fi

# example workflow

if [ "scan" = "$1" ]; then

  shift;

  #
  # Determine source type
  #
  type="image"
  if [ "image" = "${1}" -o "directory" = "${1}" -o "archive" = "${1}" -o "repo" = "${1}" ]; then
     type=${1}
     shift;
  fi

  #
  #
  #
  source_name_raw=${1}
  if [ "." = "$1" ]; then
    source_name_raw=localdir
  fi
  source_name=$(echo ${source_name_raw} | tr "/" "_" | tr "." "_" | tr ":" "_" )
  shift;

  echo Analyzing structure for ${type} "${source_name_raw}" ...

  source_id="${source_name}"
  if [ "image" = "${type}" ]; then

    docker pull ${source_name_raw}
    docker inspect ${source_name_raw} > ${source_name}.docker-inspect.json
    source_id=$( cat ${source_name}.docker-inspect.json | jq -r .[].Id )
    echo "${source_id}" > ${source_name}.txt

    skopeo inspect docker-daemon:${source_name_raw} > ${source_id}.skopeo-inspect.json

    dockle -f json ${source_name_raw} > ${source_id}.dockle.json
  fi

  if [ "archive" = "${type}" ]; then

    mkdir ${source_name}
    if [[ "${source_name_raw}" == *".tar" ]]; then
      echo "Expanding tar archive ${source_name_raw}"
      tar xf ${source_name_raw} ${source_name}
    fi
    if [[ "${source_name_raw}" == *".tar.gz" ]]; then
      echo "Expanding gzipped tar archive ${source_name_raw}"
      tar xzf ${source_name_raw} ${source_name}
    fi
    if [[ "${source_name_raw}" == *".zip" ]]; then
      echo "Expanding zip archive ${source_name_raw}"
      unzip ${source_name_raw} ${source_name}
    fi
  fi

  # Make a working copy-only checkout of the repo
  if [ "repo" = "${type}" ]; then
    rm -rf ./${source_name}
    git clone --separate-git-dir=${source_name}.git ${source_name_raw} ${source_name}
    rm -f ${source_name}/.git
    rm -rf ${source_name}.git
  fi

  # If we haven't already run this ...
  if ! [ -f "${source_id}.syft.json" ]; then
    echo "Inventorying ${type} ${source_name_raw} (${source_id}) ..."
    syft -q power-user "${source_id}" > ${source_id}.syft.json
  else
    echo "Using cached results for ${source_name} (${source_id})"
  fi

  exit 0
  # If we have results convert them
  if [ -s "${source_id}.syft.json" ]; then

    cat  ${source_id}.syft.json | jq .files > ${source_id}.manifest.json

    #syft version
    #syft convert ${image_id}.syft.json -o spdx-json=${image_id}.sbom-spdx.json
    #syft convert ${image_id}.syft.json -o cyclonedx-json=${image_id}.sbom-cyclonedx.json
    grype ${source_id}.syft.json -o sarif > ${source_id}.sarif.json
    grype ${source_id}.syft.json -o cyclonedx > ${source_id}.cyclonedx.xml

    grype ${source_id}.syft.json -o json > ${source_id}.grype.json

  else

    echo Error generating SBOM ... cleaning up
    rm -f "${source_id}.syft.json"

  fi

  # Run trivy on vulnerabilities
  trivy image --security-checks vuln -f json "${source_id}"    > ${source_id}.vulnerabilities.json

  # Run trivy on secrets and IaC
  trivy image --security-checks secret -f json "${source_id}"  > ${source_id}.secrets.json
  trivy image --security-checks config -f json "${source_id}"  > ${source_id}.image.json
  trivy fs    --security-checks config -f json /local-context > ${source_id}.iac.json

  #
  # Run opa tests against permissions in container image
  #
  conftest test \
           ${source_id}.manifest.json \
           --policy /opa/policy/permissions \
           -o json > \
      ${source_id}.permissions.json

  # Run opa tests against the Dockerfile
  for dockerfile in $( cat ${source_id}.iac.json | jq -r .Results[].Target ); do

    conftest test \
             /local-context/${dockerfile} \
             -p /opa/policy/dockerfile \
             -o json \
             --data /opa/opa-dockerfile-data.yaml >> \
        ${source_id}.dockerfile-owasp.json

  done

  #checkov --directory /data -o json > checkov.json


  cat ${source_id}.vulnerabilities.json | jq -r .Results[].Vulnerabilities[].VulnerabilityID  | sort | uniq  > ${source_id}.trivy-uniq-vulns.txt
  cat ${source_id}.grype.json | jq -r .matches[].vulnerability.id | sort | uniq  > ${source_id}.grype-uniq-vulns.txt

  sdiff ${source_id}.grype-uniq-vulns.txt ${source_id}.trivy-uniq-vulns.txt

  # echo "SBOM and Inventory "
  # echo "-------------------------------------"
  # cat  ${image_id}.syft.json
  # echo

  # echo "SCA Findings (SARIF format)"
  # echo "-------------------------------------"
  # cat  ${image_id}.sarif.json
  # echo

  # echo "Secrets Findings"
  # echo "-------------------------------------"
  # cat  ${image_id}.secrets.json
  # echo


  # echo "Dockerfile Findings"
  # echo "-------------------------------------"
  # cat  ${image_id}.dockerfile-owasp.json
  # echo

  # echo "Image Permissions Findings"
  # echo "-------------------------------------"
  # cat  ${image_id}.permissions.json
  # echo


  exit $?
fi


# Hooks into individual tools

if [ "trivy" = "$1" ]; then
  shift;
  trivy "$@"
  exit $?
fi

if [ "grype" = "$1" ]; then
  shift;
  grype "$@"
  exit $?
fi

if [ "syft" = "$1" ]; then
  shift;
  #syft --config /syft/config.yaml "$@"
  syft "$@"
  exit $?
fi

if [ "dockle" = "$1" ]; then
  shift;
  dockle "$@"
  exit $?
fi

if [ "skopeo" = "$1" ]; then
  shift;
  skopeo "$@"
  exit $?
fi

if [ "opa" = "$1" ]; then
  shift;
  opa "$@"
  exit $?
fi

if [ "conftest" = "$1" ]; then
  shift;
  conftest "$@"
  exit $?
fi

if [ "yor" = "$1" ]; then
  shift;
  yor --directory /data "$@"
  exit $?
fi

if [ "checkov" = "$1" ]; then
  shift;
  checkov --workdir /data --directory /data "$@"
  exit $?
fi

if [ "docker" = "$1" ]; then
  shift;
  docker $@
  exit $?
fi

if [ "echo" = "$1" ]; then
  shift;
  echo $@
  exit $?
fi

if [ "bash" = "$1" ]; then
  shift;
  exec $@
  exit $?
fi

echo "ERROR: do not understand command ${@}"
exit 1


# grype
# syft
# opa
# conftest
