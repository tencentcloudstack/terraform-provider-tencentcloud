
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdcpgInstancesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdcpgInstances,
				Check: resource.ComposeTestCheckFunc(
				  testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_instances.instances"),
				),
			},
		},
	})
}

const testAccDataSourceTdcpgInstances = `

data "tencentcloud_tdcpg_instances" "instances" {
  cluster_id = ""
  instance_id = ""
  instance_name = ""
  status = ""
  instance_type = ""
  }

`
