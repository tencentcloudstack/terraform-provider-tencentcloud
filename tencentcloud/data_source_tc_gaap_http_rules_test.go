package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapHttpRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapHttpRules_path(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
				),
			},
		},
	})
}

const gaapHttpRulesResources = `
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "www.qq.com"
  path            = "/"
  realserver_type = "IP"
  health_check    = true

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }
}
`

const TestAccDataSourceTencentCloudGaapHttpRulesDomain = gaapHttpRulesResources + `

data tencentcloud_gaap_http_rules "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "${tencentcloud_gaap_http_rule.foo.domain}"
}
`

const TestAccDataSourceTencentCloudGaapHttpRulesPath = gaapHttpRulesResources + `

data tencentcloud_gaap_http_rules "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  path        = "${tencentcloud_gaap_http_rule.foo.path}"
}
`
