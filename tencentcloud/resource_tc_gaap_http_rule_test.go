package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudGaapHttpRule_basic(t *testing.T) {
	id := new(string)
	listenerId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id, listenerId),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.200", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_httpUpdate(t *testing.T) {
	id := new(string)
	listenerId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id, listenerId),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.200", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "wrr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/health"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "HEAD"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.100", "100"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.200", "200"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdateDisableHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_httpUpdateRealservers(t *testing.T) {
	id := new(string)
	listenerId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id, listenerId),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.200", "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				Config: testAccGaapHttpRuleHttpUpdateRealservers,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_noHealth(t *testing.T) {
	id := new(string)
	listenerId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id, listenerId),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleNoHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "HEAD"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_domainRealserver(t *testing.T) {
	id := new(string)
	listenerId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id, listenerId),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleDomainRealserver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id, listenerId),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "delay_loop", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_http_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGaapHttpRuleExists(n string, id, listenerId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no http rule id is set")
		}

		attrListenerId := rs.Primary.Attributes["listener_id"]
		if attrListenerId == "" {
			return errors.New("no listener id is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rule, _, err := service.DescribeHttpRule(context.TODO(), attrListenerId, rs.Primary.ID)
		if err != nil {
			return err
		}

		if rule == nil {
			return errors.New("rule not exist")
		}

		*listenerId = rule.listenerId
		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapHttpRuleDestroy(id, listenerId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		if *id == "" {
			return errors.New("http rule id is nil")
		}

		rule, _, err := service.DescribeHttpRule(context.TODO(), *listenerId, *id)
		if err != nil {
			return err
		}

		if rule != nil {
			return errors.New("http rule still exists")
		}

		return nil
	}
}

const testAccGaapHttpRuleBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
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

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "${tencentcloud_gaap_http_domain.foo.domain}"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }

  realservers {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.ip}"
    port = 80
  }
}
`

const testAccGaapHttpRuleUpdate = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
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

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "${tencentcloud_gaap_http_domain.foo.domain}"
  path                      = "/new"
  realserver_type           = "IP"
  scheduler                 = "wrr"
  health_check              = true
  connect_timeout           = 3
  health_check_path         = "/health"
  health_check_method       = "HEAD"
  health_check_status_codes = [100, 200]
  
  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }

  realservers {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.ip}"
    port = 80
  }
}
`

const testAccGaapHttpRuleUpdateDisableHealth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
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

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "${tencentcloud_gaap_http_domain.foo.domain}"
  path                      = "/new"
  realserver_type           = "IP"
  health_check              = false
  connect_timeout           = 3
  health_check_path         = "/health"
  health_check_method       = "HEAD"
  health_check_status_codes = [100, 200]
  
  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }

  realservers {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.ip}"
    port = 80
  }
}
`

const testAccGaapHttpRuleHttpUpdateRealservers = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
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

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "${tencentcloud_gaap_http_domain.foo.domain}"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }
}
`

const testAccGaapHttpRuleNoHealth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
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

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "${tencentcloud_gaap_http_domain.foo.domain}"
  path            = "/"
  realserver_type = "IP"
  health_check    = false

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }

  realservers {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.ip}"
    port = 80
  }
}
`

const testAccGaapHttpRuleDomainRealserver = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
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

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "${tencentcloud_gaap_http_domain.foo.domain}"
  path            = "/"
  realserver_type = "DOMAIN"
  health_check    = false

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.domain}"
    port = 80
  }

  realservers {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.domain}"
    port = 80
  }
}
`
