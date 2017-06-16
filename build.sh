#!/bin/bash

export PACKAGE=metrics-collector
export DEST_PATH=/opt/$PACKAGE


env GOOS=linux GOARCH=amd64 go build -v $PACKAGE.go
scp $PACKAGE root@dockerhub.us:$DEST_PATH

