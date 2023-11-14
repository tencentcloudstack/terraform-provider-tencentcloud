package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsSnapshotByTimeoffsetTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSnapshotByTimeoffsetTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_snapshot_by_timeoffset_templates.snapshot_by_timeoffset_templates")),
			},
		},
	})
}

const testAccMpsSnapshotByTimeoffsetTemplatesDataSource = `

data "tencentcloud_mps_snapshot_by_timeoffset_templates" "snapshot_by_timeoffset_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  snapshot_by_time_offset_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		resolution_adaptive = &lt;nil&gt;
		format = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		fill_type = &lt;nil&gt;

  }
}

`
