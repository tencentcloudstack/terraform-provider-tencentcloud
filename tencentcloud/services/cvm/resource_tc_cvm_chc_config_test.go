package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmChcConfigResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcConfigResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "chc_id"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.#", "1"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.0.as_vpc_gateway"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "password"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.0.subnet_id"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.#", "1"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.0.vpc_id"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.0.subnet_id"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.0.as_vpc_gateway"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "instance_name", "test"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "device_type"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "deploy_security_group_ids.#", "1"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "id"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_user", "admin"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.0.vpc_id"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_security_group_ids.#", "1")),
			},
			{
				Config: testAccCvmChcConfigResource_BasicChange1,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "password"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "id"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "instance_name", "test1"), resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_user", "admin1")),
			},
			{
				ResourceName:            "tencentcloud_cvm_chc_config.chc_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bmc_user", "password"},
			},
		},
	})
}

const testAccCvmChcConfigResource_BasicCreate = `

data "tencentcloud_security_groups" "security_groups1" {
    name = "keep-tke"
}
data "tencentcloud_security_groups" "security_groups2" {
    name = "keep-reject-all"
}
resource "tencentcloud_vpc" "vpc" {
    name = "tf-cvm-vpc"
    cidr_block = "172.16.0.0/16"
}
resource "tencentcloud_subnet" "subnet1" {
    is_multicast = false
    name = "tf_cvm_vpc_subnet1"
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-7"
    cidr_block = "172.16.0.0/20"
}
resource "tencentcloud_subnet" "subnet2" {
    name = "tf_cvm_vpc_subnet2"
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-7"
    cidr_block = "172.16.16.0/20"
    is_multicast = false
}
resource "tencentcloud_cvm_chc_config" "chc_config" {
    
    deploy_virtual_private_cloud {
        vpc_id = tencentcloud_vpc.vpc.id
        subnet_id = tencentcloud_subnet.subnet1.id
        as_vpc_gateway = false
    }
    deploy_security_group_ids = [data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id]
    chc_id = "chc-mn3l1qf5"
    instance_name = "test"
    bmc_user = "admin"
    password = "123"
    
    bmc_virtual_private_cloud {
        subnet_id = tencentcloud_subnet.subnet1.id
        as_vpc_gateway = false
        vpc_id = tencentcloud_vpc.vpc.id
    }
    bmc_security_group_ids = [data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id]
}

`
const testAccCvmChcConfigResource_BasicChange1 = `

data "tencentcloud_security_groups" "security_groups1" {
    name = "keep-tke"
}
data "tencentcloud_security_groups" "security_groups2" {
    name = "keep-reject-all"
}
resource "tencentcloud_vpc" "vpc" {
    name = "tf-cvm-vpc"
    cidr_block = "172.16.0.0/16"
}
resource "tencentcloud_subnet" "subnet1" {
    is_multicast = false
    name = "tf_cvm_vpc_subnet1"
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-7"
    cidr_block = "172.16.0.0/20"
}
resource "tencentcloud_subnet" "subnet2" {
    name = "tf_cvm_vpc_subnet2"
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-7"
    cidr_block = "172.16.16.0/20"
    is_multicast = false
}
resource "tencentcloud_cvm_chc_config" "chc_config" {
    
    deploy_virtual_private_cloud {
        vpc_id = tencentcloud_vpc.vpc.id
        subnet_id = tencentcloud_subnet.subnet1.id
        as_vpc_gateway = false
    }
    deploy_security_group_ids = [data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id]
    chc_id = "chc-mn3l1qf5"
    instance_name = "test1"
    bmc_user = "admin1"
    password = "123"
    
    bmc_virtual_private_cloud {
        subnet_id = tencentcloud_subnet.subnet1.id
        as_vpc_gateway = false
        vpc_id = tencentcloud_vpc.vpc.id
    }
    bmc_security_group_ids = [data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id]
}

`
