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
		Steps: []resource.TestStep{
			{
				Config: testAccApmPrometheusRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "service_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "metric_match_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "metric_name_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "status"),
				),
			},
			{
				Config: testAccApmPrometheusRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "service_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "metric_match_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "metric_name_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_prometheus_rule.example", "status"),
				),
			},
			{
				ResourceName:      "tencentcloud_apm_prometheus_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApmPrometheusRule = `
resource "tencentcloud_apm_prometheus_rule" "example" {
  instance_id       = "apm-lhqHyRBuA"
  name              = "tf-example"
  service_name      = "java-market-service"
  metric_match_type = 0
  metric_name_rule  = "task.duration"
  status            = 1
}
`

const testAccApmPrometheusRuleUpdate = `
resource "tencentcloud_apm_prometheus_rule" "example" {
  instance_id       = "apm-lhqHyRBuA"
  name              = "tf-example-update"
  service_name      = "java-market-service"
  metric_match_type = 0
  metric_name_rule  = "task.duration"
  status            = 2
}
`
