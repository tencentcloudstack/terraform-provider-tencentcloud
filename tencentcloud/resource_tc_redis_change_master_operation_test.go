package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisChangeMasterOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisChangeMasterOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_change_master_operation.change_master_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_change_master_operation.change_master_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisChangeMasterOperation = `

resource "tencentcloud_redis_change_master_operation" "change_master_operation" {
  instance_id = "crs-c1nl9rpv"
  group_id = "crs-rpl-c1nl9rpv"
}

`
