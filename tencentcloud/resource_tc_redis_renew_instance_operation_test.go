package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisRenewInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisRenewInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_renew_instance_operation.renew_instance_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_renew_instance_operation.renew_instance_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisRenewInstanceOperation = `

resource "tencentcloud_redis_renew_instance_operation" "renew_instance_operation" {
  instance_id = "crs-c1nl9rpv"
  period = 1
  modify_pay_mode = "prepaid"
}

`
