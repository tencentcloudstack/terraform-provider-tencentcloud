package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMpsSampleSnapshotTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSampleSnapshotTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_sample_snapshot_template.sample_snapshot_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_sample_snapshot_template.sample_snapshot_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsSampleSnapshotTemplate = `

resource "tencentcloud_mps_sample_snapshot_template" "sample_snapshot_template" {
  sample_type = &lt;nil&gt;
  sample_interval = &lt;nil&gt;
  name = &lt;nil&gt;
  width = 0
  height = 0
  resolution_adaptive = "open"
  format = "jpg"
  comment = &lt;nil&gt;
  fill_type = "black"
}

`
