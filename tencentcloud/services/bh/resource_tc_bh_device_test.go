package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhDeviceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhDevice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_device.example", "id"),
				),
			},
			{
				Config: testAccBhDeviceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_device.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_device.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhDevice = `
resource "tencentcloud_bh_device" "example" {
  device_set {
    os_name = "Linux"
    ip      = "1.1.1.1"
    port    = 22
    name    = "tf-example"
  }
}
`

const testAccBhDeviceUpdate = `
resource "tencentcloud_bh_device" "example" {
  device_set {
    os_name = "Linux"
    ip      = "1.1.1.1"
    port    = 44
    name    = "tf-example"
  }
}
`
