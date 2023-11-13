package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisSslResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisSsl,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_ssl.ssl", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_ssl.ssl",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisSsl = `

resource "tencentcloud_redis_ssl" "ssl" {
  instance_id = "crs-c1nl9rpv"
}

`
