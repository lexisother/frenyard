#!/bin/sh
# Pre-commit hook to automatically do the busywork
go build . && golint -set_exit_status . design/... frenyard/... middle/...
RES=$?
exit $RES
