package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseResetInstancesPasswordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseResetInstancesPassword,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_reset_instances_password.reset_instances_password", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_reset_instances_password.reset_instances_password",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseResetInstancesPassword = `

resource "tencentcloud_lighthouse_reset_instances_password" "reset_instances_password" {
  instance_ids = 
  password = "xxxxx"
  user_name = "root"
}

`
