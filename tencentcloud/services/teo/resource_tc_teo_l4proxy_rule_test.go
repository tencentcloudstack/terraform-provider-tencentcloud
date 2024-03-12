package teo_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTeoL4proxy_ruleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL4proxy_rule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_l4proxy_rule.l4proxy_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_l4proxy_rule.l4proxy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoL4proxy_rule = `

resource "tencentcloud_teo_l4proxy_rule" "l4proxy_rule" {
  zone_id = "zone-21xfqlh4qjee"
  proxy_id = "proxy-00dde483-9a2b-11ec-bcb0"
  l4_proxy_rules {
		rule_id = "rule-2qzkbvx3hxl7"
		protocol = "TCP"
		port_range = &lt;nil&gt;
		origin_type = "IP_DOMAIN"
		origin_value = &lt;nil&gt;
		origin_port_range = "80-90"
		client_ip_pass_through_mode = "TOA"
		session_persist = "off"
		session_persist_time = 300
		rule_tag = "rule-service1	"
		status = "offline"

  }
}

`
