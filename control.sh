#!/bin/bash

IMAGE=starland-account:latest
CONTAINER_NAME=starland-account
HTTP_PORT=8081
IMAGE_DIR=/data/starland/image
start() {
  docker run -d -p ${HTTP_PORT}:8081 -v ${PWD}/conf:/app/conf \
    -v ${PWD}/logfile:/app/logfile \
    -v ${IMAGE_DIR}:/app/image \
    --name ${CONTAINER_NAME} ${IMAGE}
}

stop() {
  docker rm ${CONTAINER_NAME} --force
}

case C"$1" in
C)
  echo "Usage: $0 {start|stop|restart}"
  ;;
Cstart)
  start
  echo "Start Done!"
  ;;
Cstop)
  stop
  echo "Stop Done!"
  ;;
Crestart)
  stop
  start
  echo "Restart Done!"
  ;;
C*)
  echo "Usage: $0 {start|stop|restart}"
  ;;
esac
