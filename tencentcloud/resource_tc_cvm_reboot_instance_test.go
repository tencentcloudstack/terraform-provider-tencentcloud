package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		},
	})
}

const testAccCvmRebootInstance = defaultRebootCvmVariable + `

resource "tencentcloud_cvm_reboot_instance" "reboot_instance" {
  instance_id = var.cvm_id
  force_reboot = false
}
`
