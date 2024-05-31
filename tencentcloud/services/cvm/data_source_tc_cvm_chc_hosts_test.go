package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmChcHostsDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcHostsDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_chc_hosts.chc_hosts"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.instance_name"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.host_ids.#", "0"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.#", "1"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_security_group_ids.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.device_type"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.subnet_id"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.hardware_description"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.vpc_id"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.vpc_id"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_mac"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_ip"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.as_vpc_gateway"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.cpu"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.disk"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_ip"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.created_time"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.serial_number"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.instance_state"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.project_id"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.zone"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.tenant_type"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.ipv6_address_count"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_security_group_ids.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.memory"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.as_vpc_gateway"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.private_ip_addresses.#", "0"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.private_ip_addresses.#", "0"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.chc_id"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.host_ips.#", "0"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.subnet_id"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_mac")),
			},
		},
	})
}

const testAccCvmChcHostsDataSource_BasicCreate = `

data "tencentcloud_security_groups" "security_groups1" {
    name = "keep-tke"
}
data "tencentcloud_security_groups" "security_groups2" {
    name = "keep-reject-all"
}
data "tencentcloud_cvm_chc_hosts" "chc_hosts" {
    chc_ids = [tencentcloud_cvm_chc_config.chc_config.chc_id]
    
    filters {
        values = ["ap-guangzhou-7"]
        name = "zone"
    }
}
resource "tencentcloud_vpc" "vpc" {
    name = "tf-cvm-vpc"
    cidr_block = "172.16.0.0/16"
}
resource "tencentcloud_subnet" "subnet1" {
    name = "tf_cvm_vpc_subnet1"
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-7"
    cidr_block = "172.16.0.0/20"
    is_multicast = true
}
resource "tencentcloud_subnet" "subnet2" {
    name = "tf_cvm_vpc_subnet2"
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-7"
    cidr_block = "172.16.16.0/20"
    is_multicast = true
}
resource "tencentcloud_cvm_chc_config" "chc_config" {
    chc_id = "chc-mn3l1qf5"
    instance_name = "test"
    bmc_user = "admin"
    password = "123"
    
    bmc_virtual_private_cloud {
        vpc_id = tencentcloud_vpc.vpc.id
        subnet_id = tencentcloud_subnet.subnet1.id
        as_vpc_gateway = false
    }
    bmc_security_group_ids = [data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id]
    
    deploy_virtual_private_cloud {
        as_vpc_gateway = false
        vpc_id = tencentcloud_vpc.vpc.id
        subnet_id = tencentcloud_subnet.subnet1.id
    }
    deploy_security_group_ids = [data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id]
}

`
