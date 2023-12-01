package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbDeviceGroupMembersResource_basic -v
func TestAccTencentCloudNeedFixDasbDeviceGroupMembersResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbDeviceGroupMembers,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group_members.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_device_group_members.example", "device_group_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_device_group_members.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDasbDeviceGroupMembers = `
resource "tencentcloud_dasb_device_group_members" "example" {
  device_group_id = 3
  member_id_set   = [1, 2, 3]
}
`
