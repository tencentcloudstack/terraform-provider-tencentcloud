package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudClbListenerRulesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListenerRulesDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbListenerRuleExists("tencentcloud_clb_listener_rule.rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listener_rules.rules", "rule_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_listener_rules.rules", "rule_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_listener_rules.rules", "rule_list.0.listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_listener_rules.rules", "rule_list.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listener_rules.rules", "rule_list.0.session_expire_time", "30"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listener_rules.rules", "rule_list.0.scheduler", "WRR"),
				),
			},
		},
	})
}

const testAccClbListenerRulesDataSource = `
resource "tencentcloud_clb_instance" "clb" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rules"
}

resource "tencentcloud_clb_listener" "listener" {
  clb_id        = tencentcloud_clb_instance.clb.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "mylistener1234"
}

resource "tencentcloud_clb_listener_rule" "rule" {
  clb_id              = tencentcloud_clb_instance.clb.id
  listener_id         = tencentcloud_clb_listener.listener.listener_id
  domain              = "abcde.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

data "tencentcloud_clb_listener_rules" "rules" {
  clb_id      = tencentcloud_clb_instance.clb.id
  listener_id = tencentcloud_clb_listener.listener.listener_id
  domain      = tencentcloud_clb_listener_rule.rule.domain
  url         = tencentcloud_clb_listener_rule.rule.url
}
`
