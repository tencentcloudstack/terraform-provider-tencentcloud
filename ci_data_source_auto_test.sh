#!/usr/bin/env bash

export TF_LOG=""
export TF_ACC=true

old_path=`pwd`
cd tencentcloud;
go test -i; go test -test.run TestAccTencentCloudAvailabilityZonesDataSource_basic -v
go test -i; go test -test.run TestAccTencentCloudInstanceTypesDataSource_basic -v
go test -i; go test -test.run TestAccTencentCloudImagesDataSource_filter -v
cd $old_path