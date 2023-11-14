package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisSwitchMasterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisSwitchMaster,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_switch_master.switch_master", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_switch_master.switch_master",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisSwitchMaster = `

resource "tencentcloud_redis_switch_master" "switch_master" {
  instance_id = "crs-c1nl9rpv"
  group_id = 10001
}

`
