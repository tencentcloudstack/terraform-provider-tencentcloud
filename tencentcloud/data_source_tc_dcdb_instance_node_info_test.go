package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbInstanceNodeInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbInstanceNodeInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_instance_node_info.instance_node_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instance_node_info.instance_node_info", "nodes_info.#"),
				),
			},
		},
	})
}

const testAccDcdbInstanceNodeInfoDataSource = CommonPresetDcdb + `

data "tencentcloud_dcdb_instance_node_info" "instance_node_info" {
  instance_id = local.dcdb_id
}

`
