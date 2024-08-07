#!/bin/sh

docker build --platform linux/x86_64 -f deploy.dockerfile -t tcharlezin/drone-ci-proxy:latest .
docker push tcharlezin/drone-ci-proxy:latest