package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRedisReplicateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReplicateInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate_instance.replicate_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_replicate_instance.replicate_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisReplicateInstance = `

resource "tencentcloud_redis_replicate_instance" "replicate_instance" {
  instance_id = "crs-c1nl9rpv"
  group_id = "crs-rpl-c1nl9rpv"
  instance_role = "rw"
}

`
