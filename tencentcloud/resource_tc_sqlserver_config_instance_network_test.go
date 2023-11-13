package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigInstanceNetworkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceNetwork,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_network.config_instance_network", "id")),
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
  instance_id = "mssql-i1z41iwd"
  new_vpc_id = "vpc-j90ok"
  new_subnet_id = "sub-ja891"
  old_ip_retain_time = 0
  vip = "10.1.200.11"
}

`
