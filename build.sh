#!/usr/bin/env bash

version=2.21
githash=`git log --pretty=format:'%h' -n 1`

cat << EOF > ./version.go
package main

//DO NOT EDIT
//Build script overwrites this file.
//DO NOT EDIT

func getFactoCordVersion() (string, string) {
    return "${version}", "${githash}"
}
EOF

go build -o ./FactoCord main.go version.go
