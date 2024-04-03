#!/bin/bash

IMAGE=starland-account:latest
docker rmi ${IMAGE}
docker build --label project=starland-account -t ${IMAGE} .
docker push ${IMAGE}
