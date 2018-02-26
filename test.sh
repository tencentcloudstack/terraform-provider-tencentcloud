#!/usr/bin/env bash

export TF_ACC=true

old_path=`pwd`
cd tencentcloud;
#go test -i; go test . -v
#go test -i; go test -test.run TestAccTencentCloudKeyPair_pubcliKey -v
#go test -i; go test -test.run TestAccTencentCloudKeyPair_basic -v
#go test -i; go test -test.run TestAccTencentCloudInstance_network -v
#go test -i; go test -test.run TestAccTencentCloudInstance_basic -v
#go test -i; go test -test.run TestAccTencentCloudAvailabilityZonesDataSource_basic -v
#go test -i; go test -test.run TestAccTencentCloudInstance_changed -v
#go test -i; go test -test.run TestAccTencentCloudInstance_vpc -v
#go test -i; go test -test.run TestAccTencentCloudInstance_sg -v
#go test -i; go test -test.run TestAccTencentCloudInstance_imageIdChanged -v
#go test -i; go test -test.run TestAccTencentCloudInstance_passwordChanged -v
#go test -i; go test -test.run TestAccTencentCloudInstance_keypair -v
#go test -i; go test -test.run TestAccTencentCloudInstanceTypesDataSource_basic -v
#go test -i; go test -test.run TestAccTencentCloudImagesDataSource_filter -v
go test -v -run TestAccTencentCloudVpc github.com/tencentyun/terraform-provider-tencentcloud/tencentcloud
go test -v -run TestAccTencentCloudSubnet github.com/tencentyun/terraform-provider-tencentcloud/tencentcloud
go test -v -run TestAccTencentCloudRouteTable github.com/tencentyun/terraform-provider-tencentcloud/tencentcloud
go test -v -run TestAccTencentCloudSecurityGroup_ github.com/tencentyun/terraform-provider-tencentcloud/tencentcloud
go test -v -run TestAccTencentCloudRouteEntry github.com/tencentyun/terraform-provider-tencentcloud/tencentcloud
go test -v -run TestAccTencentCloudSecurityGroupRule github.com/tencentyun/terraform-provider-tencentcloud/tencentcloud
cd $old_path
