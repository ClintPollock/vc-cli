!/bin/sh
# Sets up a container for running the agent aruba tests

set -ex

TMP=".tmp"
DEST="/veracode-cli"

if [ ! -d $TMP ] ; then
    mkdir $TMP
fi

# This INIT will copy a suitable binary to the tmp directory which will be used to run commands when aruba runs tests.
# Then load alpine-latest.tar.gz via curl command to the docker sock - this circumvents the need
# to perform a `docker pull` which has more overheads and credentials requirements.  
# Note that alpine-latest.tar.gz is an image archive generated using `docker save` and will be scanned in our tests.
# With these commands, we have the binary available and an image(alpine:latest) loaded for testing.
INIT=$(cat <<CMD
[ ! -d "$TMP/bin" ] && mkdir $TMP/bin 
cp ./bin/veracode-cli_linux_amd64 $TMP/bin/veracode  
curl --silent -XPOST --unix-socket /run/docker.sock --data-binary "@tests/alpine-latest.tar.gz" -H 'Content-Type: application/x-tar' http://localhost/images/load
export PATH='/usr/local/bundle/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/veracode-cli/.tmp/bin'
CMD
)

INIT_SCRIPT=$TMP/install.sh
echo "$INIT" >| $INIT_SCRIPT

IMAGE_NAME="docker-ro.laputa.veracode.io/policy/veracode-cli/base"

# This runs the image mounting docker daemon sock 
docker run --rm -it \
    -w /veracode-cli  \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $(pwd):/veracode-cli \
 $IMAGE_NAME bash --init-file $INIT_SCRIPT 