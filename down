#!/bin/bash
if ! [ "`which docker-compose`" ]; then
  echo 'Error: docker-compose must be installed' >&2
  exit 1
fi

if [ -z "$@" ]; then
    sudo docker-compose down
else 
    sudo docker-compose kill $@
fi