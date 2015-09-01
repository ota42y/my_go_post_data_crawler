#!/bin/bash

echo "go fmt"
find . -path "./_vendor" -prune -o -name "*.go" -exec go fmt {} \;

echo "golint"
find config -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find work/backup -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find work/checker -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find work/crawler -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find lib/logger -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find lib/post -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find test -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;

echo "test"
gom test

echo "gom test work/backup/mongodb/*.go"
gom test work/backup/mongodb/*.go

echo "gom test work/backup/dropbox/*.go"
gom test work/backup/dropbox/*.go

echo "gom test work/checker/error/*.go"
gom test work/checker/log/error/*.go

echo "gom test work/crawler/twitter/*.go"
gom test work/crawler/twitter/*.go

