package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudEbEventRulesDataSource_basic -v
func TestAccTencentCloudEbEventRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_event_rules.event_rules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.add_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.event_bus_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.mod_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.rule_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.rule_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_event_rules.event_rules", "rules.0.status"),
				),
			},
		},
	})
}

const testAccEbEventRulesDataSourceVar = `
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus_rule"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_eb_event_rule" "event_rule" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccEbEventRulesDataSource = testAccEbEventRulesDataSourceVar + `
data "tencentcloud_eb_event_rules" "event_rules" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  order_by     = "AddTime"
  order        = "DESC"
  depends_on = [tencentcloud_eb_event_rule.event_rule]
}
`
