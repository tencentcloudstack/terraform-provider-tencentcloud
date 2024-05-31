package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmChcHostsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcHostsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_chc_hosts.chc_hosts"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.chc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.serial_number"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.instance_state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.device_type"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.project_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.host_ids.#", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.host_ips.#", "0"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.placement.0.host_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.as_vpc_gateway"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_virtual_private_cloud.0.private_ip_addresses.#", "0"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chbmc_security_group_idsc_host_set.0.bmc_virtual_private_cloud.0.ipv6_address_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.as_vpc_gateway"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.private_ip_addresses.#", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_virtual_private_cloud.0.ipv6_address_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_security_group_ids.#", "1"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.cvm_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.hardware_description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.disk"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.bmc_mac"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.deploy_mac"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_hosts.chc_hosts", "chc_host_set.0.tenant_type"),
				),
			},
		},
	})
}

const testAccCvmChcHostsDataSource = testAccCvmChcConfig + `

data "tencentcloud_cvm_chc_hosts" "chc_hosts" {
  chc_ids = [tencentcloud_cvm_chc_config.chc_config.chc_id]
  filters {
    name = "zone"
    values = [var.availability_zone]
  }
}
`
