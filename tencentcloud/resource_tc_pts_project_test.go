package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPtsProject_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsProject,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_pts_project.project", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_project.project",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsProject = `

resource "tencentcloud_pts_project" "project" {
  name = "ptsObjectName"
  description = ""
  tags {
			tag_key = ""
			tag_value = ""

  }
            }

`
