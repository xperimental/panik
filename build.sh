#!/usr/bin/env bash

readonly PACKAGE=github.com/xperimental/panik
readonly VERSION=0.1.0

docker build \
  --build-arg "VERSION=${VERSION}" \
  --build-arg "PACKAGE=${PACKAGE}" \
  -t xperimental/panik:${VERSION} .
