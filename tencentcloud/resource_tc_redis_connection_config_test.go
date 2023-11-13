package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisConnectionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisConnectionConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_connection_config.connection_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_connection_config.connection_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisConnectionConfig = `

resource "tencentcloud_redis_connection_config" "connection_config" {
  instance_id = "crs-c1nl9rpv"
  client_limit = &lt;nil&gt;
  bandwidth = &lt;nil&gt;
}

`
