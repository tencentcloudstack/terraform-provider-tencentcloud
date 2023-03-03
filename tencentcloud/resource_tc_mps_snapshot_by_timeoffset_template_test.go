package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Config: testAccMpsSnapshotByTimeoffsetTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_snapshot_by_timeoffset_template.snapshot_by_timeoffset_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_snapshot_by_timeoffset_template.snapshot_by_timeoffset_template", "name", "terraform-for-test"),
				),
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
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-test"
  resolution_adaptive = "open"
  width               = 140
}

`

const testAccMpsSnapshotByTimeoffsetTemplateUpdate = `

resource "tencentcloud_mps_snapshot_by_timeoffset_template" "snapshot_by_timeoffset_template" {
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-for-test"
  resolution_adaptive = "open"
  width               = 140
}

`
