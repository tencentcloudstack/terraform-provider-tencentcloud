package apm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudApmPrometheusRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccApmPrometheusRule,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.apm_prometheus_rule", "id")),
		}, {
			ResourceName:      "tencentcloud_apm_prometheus_rule.apm_prometheus_rule",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccApmPrometheusRule = `

resource "tencentcloud_apm_prometheus_rule" "apm_prometheus_rule" {
}
`
