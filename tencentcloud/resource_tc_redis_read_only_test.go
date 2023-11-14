package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisReadOnlyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReadOnly,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_read_only.read_only", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_read_only.read_only",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisReadOnly = `

resource "tencentcloud_redis_read_only" "read_only" {
  instance_id = "crs-c1nl9rpv"
  input_mode = &lt;nil&gt;
}

`
