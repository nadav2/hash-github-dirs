#!/usr/bin/env bash

docker build . -t go-gin1
docker run -i -t -p 8081:8081 go-gin1