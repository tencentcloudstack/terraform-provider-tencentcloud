package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseRebootInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseRebootInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_reboot_instance.reboot_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_reboot_instance.reboot_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseRebootInstance = `

resource "tencentcloud_lighthouse_reboot_instance" "reboot_instance" {
  instance_ids = 
}

`
