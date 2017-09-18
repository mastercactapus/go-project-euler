#!/bin/sh
set -e

go install
go-project-euler -n $1
