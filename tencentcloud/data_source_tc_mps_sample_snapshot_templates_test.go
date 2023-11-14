package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsSampleSnapshotTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSampleSnapshotTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_sample_snapshot_templates.sample_snapshot_templates")),
			},
		},
	})
}

const testAccMpsSampleSnapshotTemplatesDataSource = `

data "tencentcloud_mps_sample_snapshot_templates" "sample_snapshot_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  sample_snapshot_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		resolution_adaptive = &lt;nil&gt;
		format = &lt;nil&gt;
		sample_type = &lt;nil&gt;
		sample_interval = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		fill_type = &lt;nil&gt;

  }
}

`
