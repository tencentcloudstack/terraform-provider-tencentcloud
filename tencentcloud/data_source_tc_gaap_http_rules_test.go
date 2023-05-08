package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudGaapHttpRules_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpRulesDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.domain", "www.qq.com"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.interval"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.connect_timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_status_codes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realservers.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.forward_host"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapHttpRules_path(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpRulesPath,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.path", "/"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.interval"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.connect_timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_status_codes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realservers.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.forward_host"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapHttpRules_forwardHost(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpRulesForwardHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.path", "/"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.interval"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.connect_timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_status_codes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realservers.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.forward_host", "www.qqq.com"),
				),
			},
		},
	})
}

func gaapHttpRulesResources(port int) string {
	return fmt.Sprintf(`
resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = %d
  proxy_id = "%s"
}

resource "tencentcloud_gaap_http_domain" "foo" {
	listener_id = tencentcloud_gaap_layer7_listener.foo.id
	domain      = "www.qq.com"
}

resource "tencentcloud_gaap_http_rule" "foo" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo.id
  domain          = tencentcloud_gaap_http_domain.foo.domain
  path            = "/"
  realserver_type = "IP"
  health_check    = true

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  forward_host = "www.qqq.com"
}
`, port, defaultGaapProxyId2, defaultGaapRealserverIpId1, defaultGaapRealserverIp1)
}

var TestAccDataSourceTencentCloudGaapHttpRulesDomain = gaapHttpRulesResources(8090) + `

data tencentcloud_gaap_http_rules "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_rule.foo.domain
}
`

var TestAccDataSourceTencentCloudGaapHttpRulesPath = gaapHttpRulesResources(8091) + `

data tencentcloud_gaap_http_rules "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  path        = tencentcloud_gaap_http_rule.foo.path
}
`

var TestAccDataSourceTencentCloudGaapHttpRulesForwardHost = gaapHttpRulesResources(8092) + `

data tencentcloud_gaap_http_rules "foo" {
  listener_id  = tencentcloud_gaap_layer7_listener.foo.id
  forward_host = tencentcloud_gaap_http_rule.foo.forward_host
}
`
