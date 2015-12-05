#!/bin/bash

#### Config Vars ####
# update these to reflect your service
ServiceName="consul-client"
BasePath="/home/core/dev/"
port="8080"


# determine service path for volume mounting
CurrentDir=`pwd`
ServicePath="${CurrentDir/$BasePath/}"

# run the build
docker run -it -v `pwd`:"/go/src/$ServicePath" divideandconquer/godep:1.5.1 /bin/bash -c "cd /go/src/$ServicePath; ./build/build.sh" || { echo 'build failed' ; exit 1; }

# build the docker container with the new binary
docker build -t $ServiceName .

# run the container
echo "$@"
docker run -it $ServiceName "$@"