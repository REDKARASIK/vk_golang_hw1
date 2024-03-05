#!/bin/bash

set -exuo pipefail

root=$PWD
arr=$(find . | grep '99_hw$' | grep -v 'ci_cd/99_hw$' | grep -v '99_hw/code' |  grep -v 'conf_monitoring/99_hw' | grep -v '04_net2/99_hw' | grep -v 'common/')
for i in $arr; do golangci-lint -c .golangci.yml run $i/...;done

if [ -d "$root/09_conf_monitoring/99_hw/server" ]
then
    cd $root/09_conf_monitoring/99_hw/server
    golangci-lint -c $root/.golangci.yml run ./...
fi

if [ -d "$root/04_net2/99_hw/taskbot" ]
then
    cd $root/04_net2/99_hw/taskbot
    golangci-lint -c $root/.golangci.yml run ./...
fi

if [ -f "$root/06_databases/99_hw/redditclone/go.mod" ]
then
    cd $root/06_databases/99_hw/redditclone
    golangci-lint -c $root/.golangci.yml run ./...
fi

if [ -f "$root/05_web_app/99_hw/redditclone/go.mod" ]
then
    cd $root/05_web_app/99_hw/redditclone
    golangci-lint -c $root/.golangci.yml run ./...
fi