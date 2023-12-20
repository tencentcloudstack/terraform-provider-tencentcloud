package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmRebootInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRebootInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_reboot_instance.reboot_instance", "id")),
			},
		},
	})
}

const testAccCvmRebootInstance = tcacctest.DefaultRebootCvmVariable + `

resource "tencentcloud_cvm_reboot_instance" "reboot_instance" {
  instance_id = var.cvm_id
  stop_type = "SOFT_FIRST"
}
`
