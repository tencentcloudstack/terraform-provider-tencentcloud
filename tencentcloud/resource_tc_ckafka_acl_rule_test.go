package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaAclRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaAclRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl_rule.acl_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl_rule.acl_rule", "is_applied", "0"),
				),
			},
			{
				Config: testAccCkafkaAclRule_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl_rule.acl_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl_rule.acl_rule", "is_applied", "1"),
				),
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
	instance_id = "ckafka-vv7wpvae"
	resource_type = "Topic"
	pattern_type = "PRESET"
	rule_name = "RuleName1"
	rule_list {
		  operation = "All"
		  permission_type = "Deny"
		  host = "*"
		  principal = "User:*"
  
	}
	pattern = "prefix"
	is_applied = 0
  }
`

const testAccCkafkaAclRule_update = `
resource "tencentcloud_ckafka_acl_rule" "acl_rule" {
	instance_id = "ckafka-vv7wpvae"
	resource_type = "Topic"
	pattern_type = "PRESET"
	rule_name = "RuleName1"
	rule_list {
		  operation = "All"
		  permission_type = "Deny"
		  host = "*"
		  principal = "User:*"
  
	}
	pattern = "prefix"
	is_applied = 1
  }
`
