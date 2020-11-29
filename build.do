#!/bin/bash
exec >&2
goreleaser build --snapshot --rm-dist
