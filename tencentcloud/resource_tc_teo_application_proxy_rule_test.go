package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoApplicationProxyRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxyRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_application_proxy_rule.application_proxy_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy_rule.application_proxy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoApplicationProxyRule = `

resource "tencentcloud_teo_application_proxy_rule" "application_proxy_rule" {
  zone_id = &lt;nil&gt;
  proxy_id = &lt;nil&gt;
    proto = &lt;nil&gt;
  port = &lt;nil&gt;
  origin_type = &lt;nil&gt;
  origin_value = &lt;nil&gt;
  origin_port = &lt;nil&gt;
  status = &lt;nil&gt;
  forward_client_ip = &lt;nil&gt;
  session_persist = &lt;nil&gt;
}

`
