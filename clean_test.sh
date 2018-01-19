#!/usr/bin/env bash

export TF_LOG=""
old_path=`pwd`
cd tencentcloud;
go test -i; go test . -v
cd $old_path