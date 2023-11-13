package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisDestroyPostpaidInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisDestroyPostpaidInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_destroy_postpaid_instance.destroy_postpaid_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_destroy_postpaid_instance.destroy_postpaid_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisDestroyPostpaidInstance = `

resource "tencentcloud_redis_destroy_postpaid_instance" "destroy_postpaid_instance" {
  instance_id = "crs-c1nl9rpv"
}

`
