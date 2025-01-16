#! /bin/sh

SCRIPTHOME=$(dirname "$0")
cd $SCRIPTHOME

# Build the linux/amd64 version
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

echo "Building for linux/amd64"
go build -o ./bin/uboost-linux-amd64 .

# Build the windows/amd64 version
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=0

echo "Building for windows/amd64"
go build -o ./bin/uboost-windows-amd64.exe .


