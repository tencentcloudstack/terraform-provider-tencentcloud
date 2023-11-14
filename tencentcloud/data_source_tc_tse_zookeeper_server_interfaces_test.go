package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseZookeeperServerInterfacesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseZookeeperServerInterfacesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_zookeeper_server_interfaces.zookeeper_server_interfaces")),
			},
		},
	})
}

const testAccTseZookeeperServerInterfacesDataSource = `

data "tencentcloud_tse_zookeeper_server_interfaces" "zookeeper_server_interfaces" {
  instance_id = ""
  }

`
