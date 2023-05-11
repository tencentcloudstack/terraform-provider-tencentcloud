package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseZookeeperServerInterfacesDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_zookeeper_server_interfaces.zookeeper_server_interfaces"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_server_interfaces.zookeeper_server_interfaces", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_server_interfaces.zookeeper_server_interfaces", "content.#"),
				),
			},
		},
	})
}

const testAccTseZookeeperServerInterfacesDataSource = `

data "tencentcloud_tse_zookeeper_server_interfaces" "zookeeper_server_interfaces" {
	instance_id = "ins-7eb7eea7"
}

`
