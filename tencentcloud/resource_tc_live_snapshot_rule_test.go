package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveSnapshotRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveSnapshotRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_snapshot_rule.snapshot_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_snapshot_rule.snapshot_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveSnapshotRule = `

resource "tencentcloud_live_snapshot_rule" "snapshot_rule" {
  domain_name = ""
  template_id = 
  app_name = ""
  stream_name = ""
}

`
