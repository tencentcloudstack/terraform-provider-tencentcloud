package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigRuleEvaluationResultsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRuleEvaluationResultsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_rule_evaluation_results.example"),
				),
			},
		},
	})
}

func TestAccTencentCloudConfigRuleEvaluationResultsDataSource_withFilters(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRuleEvaluationResultsDataSourceWithFilters,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_config_rule_evaluation_results.example_with_filters"),
				),
			},
		},
	})
}

// NOTE: Replace cr-example123456 with a real config rule ID in your environment.
const testAccConfigRuleEvaluationResultsDataSource = `
data "tencentcloud_config_rule_evaluation_results" "example" {
  config_rule_id = "cr-example123456"
}
`

const testAccConfigRuleEvaluationResultsDataSourceWithFilters = `
data "tencentcloud_config_rule_evaluation_results" "example_with_filters" {
  config_rule_id  = "cr-example123456"
  compliance_type = ["NON_COMPLIANT"]
}
`
