package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaSnapshotTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaSnapshotTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_snapshot_template.media_snapshot_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaSnapshotTemplate = `

resource "tencentcloud_ci_media_snapshot_template" "media_snapshot_template" {
  name = &lt;nil&gt;
  snapshot {
		mode = &lt;nil&gt;
		start = &lt;nil&gt;
		time_interval = &lt;nil&gt;
		count = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		c_i_param = &lt;nil&gt;
		is_check_count = &lt;nil&gt;
		is_check_black = &lt;nil&gt;
		black_level = &lt;nil&gt;
		pixel_black_threshold = &lt;nil&gt;
		snapshot_out_mode = &lt;nil&gt;
		sprite_snapshot_config {
			cell_width = &lt;nil&gt;
			cell_height = &lt;nil&gt;
			padding = &lt;nil&gt;
			margin = &lt;nil&gt;
			color = &lt;nil&gt;
			columns = &lt;nil&gt;
			lines = &lt;nil&gt;
		}

  }
      }

`
