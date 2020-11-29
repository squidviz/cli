#!/bin/bash
exec >&2
go test
bats ./test.sh
