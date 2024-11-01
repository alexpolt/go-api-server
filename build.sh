#!/bin/bash -e
echo "docker build -t api-server-01 -f api-server.docker ."
docker build -t api-server-01 -f api-server.docker .
