package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapHttpRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_httpUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "wr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/health"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "HEAD"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "100"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.1", "200"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdateDisableHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_httpUpdateRealservers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				Config: testAccGaapHttpRuleHttpUpdateRealservers,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_noHealth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleNoHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_updateDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdateDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "qq.com"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_httpsTowWayAuthentication(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleHttpsTowWayAuthentication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "client_certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_httpsUpdateCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleHttpsTowWayAuthentication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.0", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "client_certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				Config: testAccGaapHttpRuleHttpsTowWayAuthenticationUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "certificate_id", "s-id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "client_certificate_id", "c-id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_domainRealserver(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleDomainRealserver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_http_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

const testAccGaapHttpRuleBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "www.qq.com"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  delay_loop                = 3
  connect_timeout           = 3
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleUpdate = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "www.qq.com"
  path                      = "/new"
  realserver_type           = "IP"
  scheduler                 = "wr"
  health_check              = true
  delay_loop                = 5
  connect_timeout           = 5
  health_check_path         = "/health"
  health_check_method       = "HEAD"
  health_check_status_codes = [100, 200]
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleUpdateDisableHealth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "www.qq.com"
  path            = "/new"
  realserver_type = "IP"
  health_check    = false
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleHttpUpdateRealservers = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "www.qq.com"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  delay_loop                = 3
  connect_timeout           = 3
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
  ]
}
`

const testAccGaapHttpRuleNoHealth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "www.qq.com"
  path            = "/"
  realserver_type = "IP"
  health_check    = false
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleUpdateDomain = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "qq.com"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  delay_loop                = 3
  connect_timeout           = 3
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleHttps = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  auth_type        = 0
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "www.qq.com"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  delay_loop                = 3
  connect_timeout           = 3
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
  certificate_id            = "default"
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleHttpsTowWayAuthentication = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  auth_type        = 0
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "www.qq.com"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  delay_loop                = 3
  connect_timeout           = 3
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
  certificate_id            = "default"
  client_certificate_id     = "default"
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleHttpsTowWayAuthenticationUpdate = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  auth_type        = 0
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "www.qq.com"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  delay_loop                = 3
  connect_timeout           = 3
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
  certificate_id            = "s-id" // TODO
  client_certificate_id     = "c-id" // TODO
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapHttpRuleDomainRealserver = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_realserver "foo" {
  domain = "www.qq.com"
  name   = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  domain = "qq.com"
  name   = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "www.qq.com"
  path            = "/"
  realserver_type = "DOMAIN"
  health_check    = false
  realservers = [
    {
      id   = "${tencentcloud_gaap_realserver.foo.id}"
      ip   = "${tencentcloud_gaap_realserver.foo.domain}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.domain}"
      port   = 80
      weight = 2
    }
  ]
}
`
