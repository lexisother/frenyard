#!/bin/sh
# Pre-commit hook to automatically do the busywork
git stash push -k
go build . && golint -set_exit_status ./...
RES=$?
git stash pop
exit $RES
