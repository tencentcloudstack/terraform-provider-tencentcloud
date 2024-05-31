package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmChcConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "chc_id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "instance_name", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "device_type"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_user", "admin"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "password"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.0.subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "bmc_virtual_private_cloud.0.as_vpc_gateway"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.0.subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "deploy_virtual_private_cloud.0.as_vpc_gateway"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "deploy_security_group_ids.#", "1"),
				),
			},
			{
				Config: testAccCvmChcConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "instance_name", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "bmc_user", "admin1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "password"),
				),
			},
			{
				ResourceName:            "tencentcloud_cvm_chc_config.chc_config",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"bmc_user", "password"},
			},
		},
	})
}

const testAccCvmChcConfigBasis = `
variable "availability_zone" {
  default = "ap-guangzhou-7"
}

variable "vpc_cidr" {
  default = "172.16.0.0/16"
}

variable "subnet_cidr1" {
  default = "172.16.0.0/20"
}

variable "subnet_cidr2" {
  default = "172.16.16.0/20"
}

variable "tke_cidr_a" {
  default = [
    "10.31.0.0/23",
    "10.31.2.0/24",
    "10.31.3.0/24",
    "10.31.16.0/24",
    "10.31.32.0/24"
  ]
}

variable "default_img_id" {
  default = "img-2lr9q49h"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "tf-cvm-vpc"
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "subnet1" {
  name              = "tf_cvm_vpc_subnet1"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr1
  is_multicast      = false
}

resource "tencentcloud_subnet" "subnet2" {
  name              = "tf_cvm_vpc_subnet2"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr2
  is_multicast      = false
}

data "tencentcloud_security_groups" "security_groups1" {
  name = "keep-tke"
}

data "tencentcloud_security_groups" "security_groups2" {
  name = "keep-reject-all"
}

locals {
  vpc_id     = tencentcloud_vpc.vpc.id
  subnet_id1 = tencentcloud_subnet.subnet1.id
  subnet_id2 = tencentcloud_subnet.subnet2.id

  sg_id1 = data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id
  sg_id2 = data.tencentcloud_security_groups.security_groups2.security_groups.0.security_group_id
}
`

const testAccCvmChcConfig = testAccCvmChcConfigBasis + `
resource "tencentcloud_cvm_chc_config" "chc_config" {
  chc_id        = "chc-mn3l1qf5"
  instance_name = "test"
  bmc_user      = "admin"
  password      = "123"
  bmc_virtual_private_cloud {
    vpc_id         = local.vpc_id
    subnet_id      = local.subnet_id1
    as_vpc_gateway = false
  }
  bmc_security_group_ids = [local.sg_id1]

  deploy_virtual_private_cloud {
    vpc_id         = local.vpc_id
    subnet_id      = local.subnet_id1
    as_vpc_gateway = false
  }
  deploy_security_group_ids = [local.sg_id1]
}
`

const testAccCvmChcConfig_update = testAccCvmChcConfigBasis + `
resource "tencentcloud_cvm_chc_config" "chc_config" {
  chc_id        = "chc-mn3l1qf5"
  instance_name = "test1"
  bmc_user      = "admin1"
  password      = "123456"
  bmc_virtual_private_cloud {
    vpc_id         = local.vpc_id
    subnet_id      = local.subnet_id1
    as_vpc_gateway = false
  }
  bmc_security_group_ids = [local.sg_id1]

  deploy_virtual_private_cloud {
    vpc_id         = local.vpc_id
    subnet_id      = local.subnet_id1
    as_vpc_gateway = false
  }
  deploy_security_group_ids = [local.sg_id1]
}
`
