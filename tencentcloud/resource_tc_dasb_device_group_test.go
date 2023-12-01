package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbDeviceGroupResource_basic -v
func TestAccTencentCloudNeedFixDasbDeviceGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbDeviceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device_group.example", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device_group.example", "department_id", "1.2"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_device_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDasbDeviceGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device_group.example", "name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_device_group.example", "department_id", "1.2.3"),
				),
			},
		},
	})
}

const testAccDasbDeviceGroup = `
resource "tencentcloud_dasb_device_group" "example" {
  name          = "tf_example"
  department_id = "1.2"
}
`

const testAccDasbDeviceGroupUpdate = `
resource "tencentcloud_dasb_device_group" "example" {
  name          = "tf_example_update"
  department_id = "1.2.3"
}
`
