package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCfsAccessRulesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCfsAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessRuleExists("tencentcloud_cfs_access_rule.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.0.access_rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.0.auth_client_ip", "172.16.16.0/24"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.0.priority", "1"),
				),
			},
		},
	})
}

const testAccCfsAccessRulesDataSource = defaultCfsAccessGroup + `

resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = local.cfs_access_group_id
  auth_client_ip = "172.16.16.0/24"
  priority = 1
}

data "tencentcloud_cfs_access_rules" "access_rules" {
  access_group_id = local.cfs_access_group_id
  access_rule_id = tencentcloud_cfs_access_rule.foo.id
}
`
