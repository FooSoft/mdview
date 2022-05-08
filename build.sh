#!/usr/bin/bash

rm -rf ./build
export GOARCH=amd64

export GOOS=windows
go build -o ./build/$GOOS/mdv/ .
pushd build/$GOOS
7za a mdv.zip mdv
popd

export GOOS=linux
go build -o ./build/$GOOS/mdv/ .
pushd build/$GOOS
tar czvf mdv.tar.gz mdv
popd

export GOOS=darwin
go build -o ./build/$GOOS/mdv/ .
pushd build/$GOOS
tar czvf mdv.tar.gz mdv
popd
