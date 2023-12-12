package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbBindDeviceAccountPasswordResource_basic -v
func TestAccTencentCloudNeedFixDasbBindDeviceAccountPasswordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbBindDeviceAccountPassword,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_bind_device_account_password.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_bind_device_account_password.example", "device_account_id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_bind_device_account_password.example", "password", "TerraformPassword"),
				),
			},
		},
	})
}

const testAccDasbBindDeviceAccountPassword = `
resource "tencentcloud_dasb_bind_device_account_password" "example" {
  device_account_id = 16
  password          = "TerraformPassword"
}
`
