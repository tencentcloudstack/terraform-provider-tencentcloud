package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRedisReplicateGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReplicate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate.replicate", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_replicate.replicate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisReplicate = `

resource "tencentcloud_redis_replicate" "replicate" {
  instance_id = "crs-c1nl9rpv"
  group_name = "group_1"
  remark = &lt;nil&gt;
}

`
