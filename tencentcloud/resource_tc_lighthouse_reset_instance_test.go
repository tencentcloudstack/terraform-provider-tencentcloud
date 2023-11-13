package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseResetInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseResetInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_reset_instance.reset_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_reset_instance.reset_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseResetInstance = `

resource "tencentcloud_lighthouse_reset_instance" "reset_instance" {
  instance_id = "lhins-xxx"
  blueprint_id = "lhbp-xxx"
}

`
