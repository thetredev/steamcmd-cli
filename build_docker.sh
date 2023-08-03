#!/bin/bash

cli_image=$(grep "CLI_IMAGE=" Dockerfile.daemon | head -1 | cut -d '=' -f 2)
go_image=$(echo "${cli_image}" | sed 's/:latest/:golang/')
daemon_image=$(echo "${cli_image}" | sed 's/:latest/:daemon/')

docker build -f Dockerfile.golang -t ${go_image} .
docker build -f Dockerfile.cli --build-arg GO_IMAGE=${go_image} -t ${cli_image} .
docker build -f Dockerfile.daemon --build-arg GO_IMAGE=${go_image} -t ${daemon_image} .
