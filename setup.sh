#!/bin/sh
go install
git config filter.encrypted.clean 'go-project-euler -clean'
git config filter.encrypted.smudge 'go-project-euler -smudge'
git config filter.encrypted.required true
