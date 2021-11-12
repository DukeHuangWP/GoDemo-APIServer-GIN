#!/usr/bin/env bash
rm -r ./build/configs
rm -rf ./build/website
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix cgo -o ./build/IntergTest main.go
cp -rf ./configs ./build/configs
cp -rf ./website ./build/website

#更正權限,刪除不需要檔案
find ./build -type f -name ".DS_Store" -depth -exec rm -rfv {} \;
chmod 755 *
chmod -R 755 ./build/*




