#!/usr/bin/bash

rm -rf ./build
export GOARCH=amd64

export GOOS=windows
go build -o ./build/$GOOS/mdview/ .
pushd build/$GOOS
7za a mdview.$GOOS.zip mdview
popd

export GOOS=linux
go build -o ./build/$GOOS/mdview/ .
pushd build/$GOOS
tar czvf mdview.$GOOS.tar.gz mdview
popd

export GOOS=darwin
go build -o ./build/$GOOS/mdview/ .
pushd build/$GOOS
tar czvf mdview.$GOOS.tar.gz mdview
popd
