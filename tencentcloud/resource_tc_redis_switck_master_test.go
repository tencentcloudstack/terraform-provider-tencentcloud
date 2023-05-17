package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisSwitckMasterResource_basic -v
func TestAccTencentCloudRedisSwitckMasterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisSwitckMaster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_switck_master.switck_master", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switck_master.switck_master", "instance_id", "crs-2yypjrnv"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switck_master.switck_master", "group_id", "8925"),
				),
			},
			{
				Config: testAccRedisSwitckMasterUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_switck_master.switck_master", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switck_master.switck_master", "instance_id", "crs-2yypjrnv"),
					resource.TestCheckResourceAttr("tencentcloud_redis_switck_master.switck_master", "group_id", "8924"),
				),
			},
		},
	})
}

const testAccRedisSwitckMaster = `

resource "tencentcloud_redis_switck_master" "switck_master" {
  instance_id = "crs-2yypjrnv"
  group_id = 8925
}

`

const testAccRedisSwitckMasterUp = `

resource "tencentcloud_redis_switck_master" "switck_master" {
  instance_id = "crs-2yypjrnv"
  group_id = 8924
}

`
