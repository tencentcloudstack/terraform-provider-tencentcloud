package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaAclRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaAclRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl_rule.acl_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_ckafka_acl_rule.acl_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCkafkaAclRule = `

resource "tencentcloud_ckafka_acl_rule" "acl_rule" {
  instance_id = "ckafka-xxx"
  resource_type = "Topic"
  pattern_type = "PREFIXED"
  rule_name = "RuleName"
  rule_list {
		operation = "All"
		permission_type = "Deny"
		host = "*"
		principal = "User:*"

  }
  pattern = "prefix"
  is_applied = 1
  comment = "CommentOfRule"
}

`
