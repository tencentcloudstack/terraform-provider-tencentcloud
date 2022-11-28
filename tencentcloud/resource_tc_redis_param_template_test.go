package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRedisParam_template_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisParam_template,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_param_template.param_template", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_param_template.param_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisParam_template = `

resource "tencentcloud_redis_param_template" "param_template" {
  name = ""
  description = ""
  product_type = ""
  template_id = ""
  param_list {
			key = ""
			value = ""

  }
}

`
