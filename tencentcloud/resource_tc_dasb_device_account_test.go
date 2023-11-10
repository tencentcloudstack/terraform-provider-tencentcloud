package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbDeviceAccountResource_basic -v
func TestAccTencentCloudNeedFixDasbDeviceAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbDeviceAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_account.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_account.example", "device_id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device_account.example", "account", "root"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_device_account.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDasbDeviceAccount = `
resource "tencentcloud_dasb_device_account" "example" {
  device_id = 100
  account   = "root"
}
`
