package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var (
	gaapHttpCheckCompose = resource.ComposeTestCheckFunc(
		testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.protocol", "HTTP"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.id"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.listener_id"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.domain", "www.qq.com"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.path", "/"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_type", "IP"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.scheduler", "rr"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check", "true"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.delay_loop", "3"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.connect_timeout", "3"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_path", "/"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_method", "GET"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_status_codes.0", "200"),

		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.certificate_id"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.client_certificate_id"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.basic_auth"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.basic_auth_config_id"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_auth"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_certificate_id"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.gaap_auth"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.gaap_certificate_id"),
		resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_certificate_domain"),

		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.realservers.#", "1"),
	)

	gaapHttpsCheckCompose = resource.ComposeTestCheckFunc(
		testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_rules.foo"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.protocol", "HTTPS"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.id"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.listener_id"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.domain", "www.qq.com"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.path", "/"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_type", "IP"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.scheduler", "rr"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check", "true"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.delay_loop", "3"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.connect_timeout", "3"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_path", "/"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_method", "GET"),
		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.health_check_status_codes.0", "200"),

		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.certificate_id"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.client_certificate_id"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.basic_auth"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.basic_auth_config_id"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_auth"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_certificate_id"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.gaap_auth"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.gaap_certificate_id"),
		resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_rules.foo", "rules.0.realserver_certificate_domain"),

		resource.TestCheckResourceAttr("data.tencentcloud_gaap_http_rules.foo", "rules.0.realservers.#", "1"),
	)
)

func TestAccDataSourceTencentCloudGaapHttpRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpRulesBasic,
				Check:  gaapHttpCheckCompose,
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapHttpRules_domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpRulesDomain,
				Check:  gaapHttpCheckCompose,
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
				Check:  gaapHttpCheckCompose,
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapHttpRules_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpRulesHttps,
				Check:  gaapHttpsCheckCompose,
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapHttpRules_httpsDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpsRulesDomain,
				Check:  gaapHttpsCheckCompose,
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapHttpRules_httpsPath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapHttpsRulesPath,
				Check:  gaapHttpsCheckCompose,
			},
		},
	})
}

const TestAccDataSourceTencentCloudGaapHttpRulesBasic = `
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
    }
  ]
}

data tencentcloud_gaap_http_rules "foo" {
  protocol    = "HTTP"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
}
`

const multiGaapHttpRulesResources = `
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

resource tencentcloud_gaap_http_rule "bar" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "qq.com"
  path                      = "/new"
  realserver_type           = "IP"
  health_check              = true
  delay_loop                = 3
  connect_timeout           = 3
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
  realservers = [
    {
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
    }
  ]
}
`

const TestAccDataSourceTencentCloudGaapHttpRulesDomain = multiGaapHttpRulesResources + `

data tencentcloud_gaap_http_rules "foo" {
  protocol    = "HTTP"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}
`

const TestAccDataSourceTencentCloudGaapHttpRulesPath = multiGaapHttpRulesResources + `

data tencentcloud_gaap_http_rules "foo" {
  protocol    = "HTTP"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  path        = "/"
}
`

const TestAccDataSourceTencentCloudGaapHttpRulesHttps = `
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
    }
  ]
}
`

const multiGaapHttpsRulesResources = `
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

resource tencentcloud_gaap_http_rule "bar" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "qq.com"
  path                      = "/new"
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
      id     = "${tencentcloud_gaap_realserver.bar.id}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
    }
  ]
}
`

const TestAccDataSourceTencentCloudGaapHttpsRulesDomain = multiGaapHttpsRulesResources + `

data tencentcloud_gaap_http_rules "foo" {
  protocol    = "HTTPS"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}
`

const TestAccDataSourceTencentCloudGaapHttpsRulesPath = multiGaapHttpsRulesResources + `

data tencentcloud_gaap_http_rules "foo" {
  protocol    = "HTTPS"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  path        = "/"
}
`
