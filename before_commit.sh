#!/bin/bash

echo "go fmt"
find . -path "./_vendor" -prune -o -name "*.go" -exec go fmt {} \;

echo "golint"
find config -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find work/backup -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find work/checker -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
find lib/logger -path "./_vendor" -prune -o -name "*.go" -exec golint {} \;
