#!/bin/bash

server_image=$(grep "CLI_IMAGE=" Dockerfile.daemon | head -1 | cut -d '=' -f 2)
go_image=$(echo "${server_image}" | sed 's/:server/:golang/')
daemon_image=$(echo "${server_image}" | sed 's/:server/:daemon/')

docker build -f Dockerfile.golang -t ${go_image} .
docker build -f Dockerfile.server --build-arg GO_IMAGE=${go_image} -t ${server_image} .
docker build -f Dockerfile.daemon --build-arg GO_IMAGE=${go_image} -t ${daemon_image} .
