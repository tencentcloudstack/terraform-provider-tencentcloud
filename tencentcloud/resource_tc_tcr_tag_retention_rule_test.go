package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrTagRetentionRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrTagRetentionRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_rule.tag_retention_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_tag_retention_rule.tag_retention_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrTagRetentionRule = `

resource "tencentcloud_tcr_tag_retention_rule" "tag_retention_rule" {
  registry_id = "tcr-12345"
  namespace_id = 1
  retention_rule {
		key = "latestPushedK"
		value = 1

  }
  cron_setting = "manual"
  disabled = false
}

`
