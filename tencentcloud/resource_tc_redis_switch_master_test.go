package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisSwitchMasterResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_switch_master.switch_master", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switch_master.switch_master", "instance_id", "crs-2yypjrnv"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switch_master.switch_master", "group_id", "8925"),
				),
			},
			{
				Config: testAccRedisSwitchMasterUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_switch_master.switch_master", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switch_master.switch_master", "instance_id", "crs-2yypjrnv"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switch_master.switch_master", "group_id", "8924"),
				),
			},
		},
	})
}

const testAccRedisSwitchMaster = `

resource "tencentcloud_redis_switch_master" "switch_master" {
  instance_id = "crs-2yypjrnv"
  group_id = 8925
}

`

const testAccRedisSwitchMasterUp = `

resource "tencentcloud_redis_switch_master" "switch_master" {
  instance_id = "crs-2yypjrnv"
  group_id = 8924
}

`
