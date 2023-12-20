package dcdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbInstanceNodeInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbInstanceNodeInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_instance_node_info.instance_node_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instance_node_info.instance_node_info", "nodes_info.#"),
				),
			},
		},
	})
}

const testAccDcdbInstanceNodeInfoDataSource = tcacctest.CommonPresetDcdb + `

data "tencentcloud_dcdb_instance_node_info" "instance_node_info" {
  instance_id = local.dcdb_id
}

`
