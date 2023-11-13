package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmRebootInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRebootInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_reboot_instance.reboot_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_reboot_instance.reboot_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmRebootInstance = `

resource "tencentcloud_cvm_reboot_instance" "reboot_instance" {
  instance_ids = 
  force_reboot = false
  stop_type = "SOFT"
}

`
