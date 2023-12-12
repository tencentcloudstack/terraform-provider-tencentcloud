package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbDeviceResource_basic -v
func TestAccTencentCloudNeedFixDasbDeviceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbDevice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "os_name", "Linux"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device.example", "port"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "department_id", "1.2.3"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_device.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDasbDeviceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "os_name", "Linux"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "ip", "192.168.0.1"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device.example", "port"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device.example", "department_id", "1.2.3.4"),
				),
			},
		},
	})
}

const testAccDasbDevice = `
resource "tencentcloud_dasb_device" "example" {
  os_name       = "Linux"
  ip            = "192.168.0.1"
  port          = 80
  name          = "tf_example"
  department_id = "1.2.3"
}
`

const testAccDasbDeviceUpdate = `
resource "tencentcloud_dasb_device" "example" {
  os_name       = "Linux"
  ip            = "192.168.0.1"
  port          = 8080
  name          = "tf_example"
  department_id = "1.2.3.4"
}
`
