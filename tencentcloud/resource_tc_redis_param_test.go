package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisParamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisParam,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_param.param", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_param.param",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisParam = `

resource "tencentcloud_redis_param" "param" {
  instance_id = "crs-c1nl9rpv"
  instance_params {
		key = &lt;nil&gt;
		value = &lt;nil&gt;

  }
}

`
