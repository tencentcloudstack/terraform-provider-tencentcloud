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
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group.example", "department_id"),
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
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group.example", "department_id"),
				),
			},
		},
	})
}

const testAccDasbDeviceGroup = `
resource "tencentcloud_dasb_device_group" "example" {
  name          = "tf_example"
  department_id = "1"
}
`

const testAccDasbDeviceGroupUpdate = `
resource "tencentcloud_dasb_device_group" "example" {
  name          = "tf_example_update"
  department_id = "1"
}
`
