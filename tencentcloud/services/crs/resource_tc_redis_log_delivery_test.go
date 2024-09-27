package crs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudRedisLogDeliveryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccRedisLogDelivery,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_log_delivery.delivery", "id")),
		}, {
			ResourceName:      "tencentcloud_redis_log_delivery.delivery",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccRedisLogDelivery = `

resource "tencentcloud_redis_log_delivery" "delivery" {
  instance_id  = "crs-dmjj8en7"
  log_region   = "ap-guangzhou"
  logset_name  = "test"
  topic_name   = "test"
  period       = 20
  create_index = true
}
`
