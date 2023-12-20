package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssSnapshotRuleAttachmentResource_basic -v
func TestAccTencentCloudCssSnapshotRuleAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssSnapshotRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_snapshot_rule_attachment.snapshot_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_rule_attachment.snapshot_rule", "domain_name", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_snapshot_rule_attachment.snapshot_rule", "template_id"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_rule_attachment.snapshot_rule", "app_name", "qqq"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_rule_attachment.snapshot_rule", "stream_name", "ppp"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_snapshot_rule_attachment.snapshot_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssSnapshotRule = testAccCssSnapshotTemplate + `

resource "tencentcloud_css_snapshot_rule_attachment" "snapshot_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = tencentcloud_css_snapshot_template.snapshot_template.id
  app_name    = "qqq"
  stream_name = "ppp"
}

`
