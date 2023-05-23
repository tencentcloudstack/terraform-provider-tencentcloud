package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigInstanceNetworkResource_basic -v
func TestAccTencentCloudSqlserverConfigInstanceNetworkResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceNetwork,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_network.config_instance_network", "id"),
				),
			},
			{
				Config: testAccSqlserverConfigInstanceNetworkUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_network.config_instance_network", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_instance_network.config_instance_network",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigInstanceNetwork = `
resource "tencentcloud_sqlserver_config_instance_network" "config_instance_network" {
  instance_id = "mssql-qelbzgwf"
  new_vpc_id = "vpc-1yg5ua6l"
  new_subnet_id = "subnet-h7av55g8"
}
`

const testAccSqlserverConfigInstanceNetworkUpdate = `
resource "tencentcloud_sqlserver_config_instance_network" "config_instance_network" {
  instance_id = "mssql-qelbzgwf"
  new_vpc_id = "vpc-4owdpnwr"
  new_subnet_id = "subnet-ahv6swf2"
}
`
