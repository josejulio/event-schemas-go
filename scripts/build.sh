#!/usr/bin/env sh

find apps core -type f -name "*.go" -exec go build {} \;
