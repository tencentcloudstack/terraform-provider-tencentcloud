package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfGroupInstancesDataSource_basic -v
func TestAccTencentCloudTsfGroupInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGroupInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_group_instances.group_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.agent_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.instance_available_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.instance_import_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.instance_pkg_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.instance_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.lan_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.reason"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.service_instance_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.service_sidecar_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_group_instances.group_instances", "result.0.content.0.wan_ip"),
				),
			},
		},
	})
}

const testAccTsfGroupInstancesDataSourceVar = `
variable "group_id" {
	default = "` + defaultTsfGroupId + `"
}
`

const testAccTsfGroupInstancesDataSource = testAccTsfGroupInstancesDataSourceVar + `

data "tencentcloud_tsf_group_instances" "group_instances" {
	group_id = var.group_id
	order_by = "ASC"
	order_type = 0
}

`
