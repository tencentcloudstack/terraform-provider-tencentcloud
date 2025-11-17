package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafOwaspRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWafOwaspRulesDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_waf_owasp_rules.example"),
			),
		}},
	})
}

const testAccWafOwaspRulesDataSource = `
data "tencentcloud_waf_owasp_rules" "example" {
  domain = "example.qcloud.com"
  by     = "RuleId"
  order  = "desc"
  filters {
    name        = "RuleId"
    values      = ["106251141"]
    exact_match = true
  }
}
`
