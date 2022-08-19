package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTemEnvironment_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemEnvironment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tem_environment.environment", "id"),
				),
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
  environment_name = "demo"
  description      = "demo for test"
  vpc              = "vpc-2hfyray3"
  subnet_ids       = ["subnet-rdkj0agk", "subnet-r1c4pn5m", "subnet-02hcj95c"]
}

`
