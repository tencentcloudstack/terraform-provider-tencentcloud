package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseRebootInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseRebootInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_reboot_instance.reboot_instance", "id")),
			},
		},
	})
}

const testAccLighthouseRebootInstance = `

resource "tencentcloud_lighthouse_reboot_instance" "reboot_instance" {
  instance_id = "lhins-hwe21u91"
}

`
