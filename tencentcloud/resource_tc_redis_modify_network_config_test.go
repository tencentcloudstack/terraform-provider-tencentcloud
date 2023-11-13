package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisModifyNetworkConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisModifyNetworkConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_modify_network_config.modify_network_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_modify_network_config.modify_network_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisModifyNetworkConfig = `

resource "tencentcloud_redis_modify_network_config" "modify_network_config" {
  instance_id = "crs-c1nl9rpv"
  operation = "changeVip"
  vip = "10.1.1.2"
  vpc_id = "vpc-hu6khgap"
  subnet_id = "subnet-6mt7lir6"
  recycle = 7
  v_port = 6379
}

`
