package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveSnapshotTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveSnapshotTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_snapshot_template.snapshot_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_snapshot_template.snapshot_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveSnapshotTemplate = `

resource "tencentcloud_live_snapshot_template" "snapshot_template" {
  template_name = ""
  cos_app_id = 
  cos_bucket = ""
  cos_region = ""
  description = ""
  snapshot_interval = 
  width = 
  height = 
  porn_flag = 
  cos_prefix = ""
  cos_file_name = ""
}

`
