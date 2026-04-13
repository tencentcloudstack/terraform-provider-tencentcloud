package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigStartConfigRuleEvaluationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccConfigStartConfigRuleEvaluationOperation,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_config_start_config_rule_evaluation_operation.example", "id"),
			),
		}},
	})
}

// NOTE: Replace the rule_id with a real config rule ID in your environment.
const testAccConfigStartConfigRuleEvaluationOperation = `
resource "tencentcloud_config_start_config_rule_evaluation_operation" "example" {
  rule_id = "cr-xhsd76j603v0a8ma0i73"
}
`
