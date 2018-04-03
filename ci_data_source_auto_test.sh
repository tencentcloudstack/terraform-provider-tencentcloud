#!/usr/bin/env bash

export TF_LOG=""
export TF_ACC=true

old_path=`pwd`
cd tencentcloud;
go test -i; go test -test.run TestAccTencentCloudAvailabilityZonesDataSource_basic -v
go test -i; go test -test.run TestAccTencentCloudInstanceTypesDataSource_basic -v
go test -i; go test -test.run TestAccTencentCloudImagesDataSource_filter -v
go test -i; go test -test.run TestAccTencentCloudEipDataSource -v
go test -v -run TestAccDataSourceTencentCloudRouteTable_basic
go test -v -run TestAccDataSourceTencentCloudSecurityGroup_basic
go test -v -run TestAccDataSourceTencentCloudSubnet_basic
go test -v -run TestAccDataSourceTencentCloudVpc_basic
go test -v -run TestAccTencentCloudDataSourceContainerClusterInstances
go test -v -run TestAccTencentCloudDataSourceContainerClusters
cd $old_path
