package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudChdfsAccessRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsAccessRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_access_rule.access_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_chdfs_access_rule.access_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsAccessRule = `

resource "tencentcloud_chdfs_access_rule" "access_rule" {
  access_rule {
		address = &lt;nil&gt;
		access_mode = &lt;nil&gt;
		priority = &lt;nil&gt;

  }
  access_group_id = &lt;nil&gt;
}

`
