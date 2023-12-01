package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcWorkGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcWorkGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_work_group.work_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_work_group.work_group", "work_group_description", "dlc workgroup test"),
				),
			},
			{
				Config: testAccDlcWorkGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_work_group.work_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_work_group.work_group", "work_group_description", "dlc workgroup"),
				),
			},
			{
				ResourceName:      "tencentcloud_dlc_work_group.work_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcWorkGroup = `

resource "tencentcloud_dlc_work_group" "work_group" {
  work_group_name        = "tf-demo"
  work_group_description = "dlc workgroup test"
}

`

const testAccDlcWorkGroupUpdate = `

resource "tencentcloud_dlc_work_group" "work_group" {
  work_group_name        = "tf-demo"
  work_group_description = "dlc workgroup"
}

`
