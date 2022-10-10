#!/bin/bash

GOOS=darwin GOARCH=amd64 go build -o build/client_mac -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(date +'%Y/%m/%d')' "
GOOS=linux GOARCH=amd64 go build -o build/client_linux -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(date +'%Y/%m/%d')' "
GOOS=windows GOARCH=amd64 go build -o build/client_windows.exe -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(date +'%Y/%m/%d')' "
