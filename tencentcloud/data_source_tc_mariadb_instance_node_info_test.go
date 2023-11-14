package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbInstanceNodeInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbInstanceNodeInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_instance_node_info.instance_node_info")),
			},
		},
	})
}

const testAccMariadbInstanceNodeInfoDataSource = `

data "tencentcloud_mariadb_instance_node_info" "instance_node_info" {
  instance_id = ""
  }

`
