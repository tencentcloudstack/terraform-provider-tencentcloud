package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafOwaspRuleTypesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWafOwaspRuleTypesDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_waf_owasp_rule_types.example"),
			),
		}},
	})
}

const testAccWafOwaspRuleTypesDataSource = `
data "tencentcloud_waf_owasp_rule_types" "example" {
  domain = "example.qcloud.com"
  filters {
    name        = "RuleId"
    values      = ["10000001"]
    exact_match = true
  }
}
`
