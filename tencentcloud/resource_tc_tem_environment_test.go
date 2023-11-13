package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTemEnvironmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemEnvironment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_environment.environment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_environment.environment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemEnvironment = `

resource "tencentcloud_tem_environment" "environment" {
  environment_name = "xxx"
  description = "xxx"
  vpc = "vpc-xxx"
  subnet_ids = 
  tags {
		tag_key = "key"
		tag_value = "tag value"

  }
}

`
