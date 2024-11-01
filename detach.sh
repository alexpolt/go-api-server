#!/bin/bash -e
echo "docker run --rm -dt --name api-server-01 -v./src:/app/src -p8080:8080 --network api-network-01 api-server-01 $1"
docker run --rm -dt --name api-server-01 -v./src:/app/src -p8080:8080 --network api-network-01 api-server-01 $1
