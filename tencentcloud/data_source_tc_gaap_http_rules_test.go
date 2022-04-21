package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo-6030"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.domain", "www.qq.com"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.interval"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.connect_timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.health_check_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.health_check_method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.health_check_status_codes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.realservers.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6030", "rules.0.forward_host"),
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
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo-6031"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.path", "/"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.interval"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.connect_timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.health_check_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.health_check_method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.health_check_status_codes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.realservers.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6031", "rules.0.forward_host"),
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
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo-6032"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.path", "/"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.interval"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.connect_timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.health_check_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.health_check_method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.health_check_status_codes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.realservers.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo-6032", "rules.0.forward_host", "www.qqq.com"),
				),
			},
		},
	})
}

func gaapL7ListenerTemplate(port int) string {
	return fmt.Sprintf(`resource "tencentcloud_gaap_layer7_listener" "foo" {
		protocol = "HTTP"
		name     = "ci-test-gaap-l7-listener"
		port     = %d
		proxy_id = "%s"
	  }`, port, defaultGaapProxyId)
}

var gaapHttpRulesResources = fmt.Sprintf(`
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
`, defaultGaapRealserverIpId1, defaultGaapRealserverIp1)

var TestAccDataSourceTencentCloudGaapHttpRulesDomain = fmt.Sprintf(`
resource "tencentcloud_gaap_layer7_listener" "foo-listener-6030" {
	protocol = "HTTP"
	name     = "ci-test-gaap-l7-listener"
	port     = 6030
	proxy_id = "%s"
  }

resource "tencentcloud_gaap_http_domain" "foo-listener-6030-domain" {
	listener_id = tencentcloud_gaap_layer7_listener.foo-listener-6030.id
	domain      = "www.qq.com"
}

resource "tencentcloud_gaap_http_rule" "foo-listener-6030-rule" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo-listener-6030.id
  domain          = tencentcloud_gaap_http_domain.foo-listener-6030-domain.domain
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

data tencentcloud_gaap_http_rules "foo-6030" {
  listener_id = tencentcloud_gaap_layer7_listener.foo-listener-6030.id
  domain      = tencentcloud_gaap_http_rule.foo-listener-6030-rule.domain
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1)

var TestAccDataSourceTencentCloudGaapHttpRulesPath = fmt.Sprintf(`
resource "tencentcloud_gaap_layer7_listener" "foo-listener-6031" {
	protocol = "HTTP"
	name     = "ci-test-gaap-l7-listener"
	port     = 6031
	proxy_id = "%s"
  }

resource "tencentcloud_gaap_http_domain" "foo-listener-6031-domain" {
	listener_id = tencentcloud_gaap_layer7_listener.foo-listener-6031.id
	domain      = "www.qq.com"
}

resource "tencentcloud_gaap_http_rule" "foo-listener-6031-rule" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo-listener-6031.id
  domain          = tencentcloud_gaap_http_domain.foo-listener-6031-domain.domain
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

data tencentcloud_gaap_http_rules "foo-6031" {
  listener_id = tencentcloud_gaap_layer7_listener.foo-listener-6031.id
  path        = tencentcloud_gaap_http_rule.foo-listener-6031-rule.path
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1)

var TestAccDataSourceTencentCloudGaapHttpRulesForwardHost = fmt.Sprintf(`
resource "tencentcloud_gaap_layer7_listener" "foo-listener-6032" {
	protocol = "HTTP"
	name     = "ci-test-gaap-l7-listener"
	port     = 6032
	proxy_id = "%s"
  }

resource "tencentcloud_gaap_http_domain" "foo-listener-6032-domain" {
	listener_id = tencentcloud_gaap_layer7_listener.foo-listener-6032.id
	domain      = "www.qq.com"
}

resource "tencentcloud_gaap_http_rule" "foo-listener-6032-rule" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo-listener-6032.id
  domain          = tencentcloud_gaap_http_domain.foo-listener-6032-domain.domain
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

data tencentcloud_gaap_http_rules "foo-6032" {
  listener_id = tencentcloud_gaap_layer7_listener.foo-listener-6032.id
  forward_host = tencentcloud_gaap_http_rule.foo-listener-6032-rule.forward_host
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1)
