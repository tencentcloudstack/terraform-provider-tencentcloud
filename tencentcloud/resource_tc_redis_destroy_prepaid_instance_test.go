package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisDestroyPrepaidInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisDestroyPrepaidInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_destroy_prepaid_instance.destroy_prepaid_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_destroy_prepaid_instance.destroy_prepaid_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisDestroyPrepaidInstance = `

resource "tencentcloud_redis_destroy_prepaid_instance" "destroy_prepaid_instance" {
  instance_id = "crs-c1nl9rpv"
}

`
