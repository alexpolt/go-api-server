#!/bin/sh -e
go build -C src -o ..
./api-server
