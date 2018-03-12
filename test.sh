#!/usr/bin/env bash

export TF_ACC=true

old_path=`pwd`
cd tencentcloud;
#go test -i; go test . -v
#go test -i; go test -test.run TestAccTencentCloudEipAssociationWithInstance -v
#go test -i; go test -test.run TestAccTencentCloudEipAssociationWithNetworkInterface -v
#go test -i; go test -test.run TestAccTencentCloudEip_basic -v
#go test -i; go test -test.run TestAccTencentCloudEipDataSource -v
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
#go test -i; go test -test.run TestAccTencentCloudVpc -v
#go test -i; go test -test.run TestAccTencentCloudSubnet -v
#go test -i; go test -test.run TestAccTencentCloudRouteTable -v
#go test -i; go test -test.run TestAccTencentCloudSecurityGroup_ -v
#go test -i; go test -test.run TestAccTencentCloudRouteEntry -v
#go test -i; go test -test.run TestAccTencentCloudSecurityGroupRule -v
#go test -i; go test -test.run TestAccTencentCloudNatGateway_basic -v
#go test -i; go test -test.run TestAccTencentCloudDnat_basic -v
go test -i; go test -test.run TestAccTencentCloudNatsDataSource -v
cd $old_path
