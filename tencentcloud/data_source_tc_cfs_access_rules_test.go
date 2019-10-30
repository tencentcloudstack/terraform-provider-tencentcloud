package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCfsAccessRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCfsAccessRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessRuleExists("tencentcloud_cfs_access_rule.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.0.access_rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.0.auth_client_ip", "10.10.1.0/24"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_rules.access_rules", "access_rule_list.0.priority", "1"),
				),
			},
		},
	})
}

const testAccCfsAccessRulesDataSource = `
resource "tencentcloud_cfs_access_group" "foo" {
  name = "test_cfs_access_rule"
}

resource "tencentcloud_cfs_access_rule" "foo" {
  access_group_id = "${tencentcloud_cfs_access_group.foo.id}"
  auth_client_ip = "10.10.1.0/24"
  priority = 1
}

data "tencentcloud_cfs_access_rules" "access_rules" {
  access_group_id = "${tencentcloud_cfs_access_group.foo.id}"
  access_rule_id = "${tencentcloud_cfs_access_rule.foo.id}"
}
`
