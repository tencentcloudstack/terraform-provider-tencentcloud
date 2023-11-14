package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsSnapshotByTimeoffsetTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSnapshotByTimeoffsetTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_snapshot_by_timeoffset_template.snapshot_by_timeoffset_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_snapshot_by_timeoffset_template.snapshot_by_timeoffset_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsSnapshotByTimeoffsetTemplate = `

resource "tencentcloud_mps_snapshot_by_timeoffset_template" "snapshot_by_timeoffset_template" {
  name = &lt;nil&gt;
  width = 0
  height = 0
  resolution_adaptive = "open"
  format = "jpg"
  comment = &lt;nil&gt;
  fill_type = "black"
}

`
