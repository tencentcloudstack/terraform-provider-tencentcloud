package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPtsFile_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_pts_file.file", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_file.file",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsFile = `

resource "tencentcloud_pts_file" "file" {
  file_id = ""
  project_id = ""
  kind = ""
  name = ""
  size = ""
  type = ""
  line_count = ""
  head_lines = ""
  tail_lines = ""
  header_in_file = ""
  header_columns = ""
  file_infos {
			name = ""
			size = ""
			type = ""
			updated_at = ""
			file_id = ""

  }
}

`
