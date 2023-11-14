package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsFileResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsFile,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_file.file", "id")),
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
  file_id = &lt;nil&gt;
  project_id = &lt;nil&gt;
  kind = &lt;nil&gt;
  name = &lt;nil&gt;
  size = &lt;nil&gt;
  type = &lt;nil&gt;
  line_count = &lt;nil&gt;
  head_lines = &lt;nil&gt;
  tail_lines = &lt;nil&gt;
  header_in_file = &lt;nil&gt;
  header_columns = &lt;nil&gt;
  file_infos {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
}

`
