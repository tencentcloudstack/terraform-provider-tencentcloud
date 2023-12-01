package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTdcpgInstancesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgInstances_id, defaultTdcpgClusterId, defaultTdcpgInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_instances.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.id", "list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.id", "list.0.instance_id", defaultTdcpgInstanceId),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.id", "list.0.instance_name", defaultTdcpgInstanceName),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.id", "list.0.cluster_id", defaultTdcpgClusterId),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.endpoint_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.id", "list.0.region", defaultRegion),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.id", "list.0.zone", defaultTdcpgZone),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.db_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.status_desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.pay_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.db_major_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.id", "list.0.db_kernel_version"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgInstances_name, defaultTdcpgClusterId, defaultTdcpgInstanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_instances.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.name", "list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.name", "list.0.instance_id", defaultTdcpgInstanceId),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.name", "list.0.instance_name", defaultTdcpgInstanceName),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.name", "list.0.cluster_id", defaultTdcpgClusterId),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.name", "list.0.endpoint_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.name", "list.0.region", defaultRegion),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.name", "list.0.zone", defaultTdcpgZone),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgInstances_status, defaultTdcpgClusterId, "running"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_instances.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.status", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.status", "list.0.endpoint_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.status", "list.0.status", "running"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgInstances_type, defaultTdcpgClusterId, "RW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_instances.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.type", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_instances.type", "list.0.endpoint_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_instances.type", "list.0.instance_type", "RW"),
				),
			},
		},
	})
}

const testAccDataSourceTdcpgInstances_id = `

data "tencentcloud_tdcpg_instances" "id" {
  cluster_id = "%s"
  instance_id = "%s"
  }

`

const testAccDataSourceTdcpgInstances_name = `

data "tencentcloud_tdcpg_instances" "name" {
  cluster_id = "%s"
  instance_name = "%s"
  }

`
const testAccDataSourceTdcpgInstances_status = `

data "tencentcloud_tdcpg_instances" "status" {
  cluster_id = "%s"
  status = "%s"
  }

`
const testAccDataSourceTdcpgInstances_type = `

data "tencentcloud_tdcpg_instances" "type" {
  cluster_id = "%s"
  instance_type = "%s"
  }

`
