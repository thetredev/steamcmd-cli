#!/bin/bash

server_image=$(grep "BASE_IMAGE=" Dockerfile.daemon | head -1 | cut -d '=' -f 2)
daemon_image=$(echo "${server_image}" | sed 's/:server/:daemon/')

docker build -f Dockerfile.server -t ${server_image} .
docker build -f Dockerfile.daemon -t ${daemon_image} .
