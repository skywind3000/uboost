@echo off

set "SCRIPT_HOME=%~dp0"
cd /D "%SCRIPT_HOME%"

set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
echo Building for linux/amd64
go build -o bin\uboost-linux-amd64 .

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
echo Building for windows/amd64
go build -o bin\uboost-windows-amd64.exe .


