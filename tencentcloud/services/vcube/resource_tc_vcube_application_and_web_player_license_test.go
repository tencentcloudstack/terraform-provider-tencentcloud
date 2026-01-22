package vcube_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudVcubeApplicationAndWebPlayerLicenseResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVcubeApplicationAndWebPlayerLicense,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vcube_application_and_web_player_license.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vcube_application_and_web_player_license.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVcubeApplicationAndWebPlayerLicense = `
resource "tencentcloud_vcube_application_and_web_player_license" "example" {
  app_name = "tf-example"
  domain_list = [
    "www.example1.com",
    "www.example2.com",
    "www.example3.com",
  ]
}
`
