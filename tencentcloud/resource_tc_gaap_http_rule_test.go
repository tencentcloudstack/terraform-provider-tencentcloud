package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudGaapHttpRule_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes."+strconv.Itoa(schema.HashInt(200)), "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "forward_host", "default"),
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

func TestAccTencentCloudGaapHttpRule_httpUpdate_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes."+strconv.Itoa(schema.HashInt(200)), "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "forward_host", "default"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "wrr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/health"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "HEAD"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes."+strconv.Itoa(schema.HashInt(100)), "100"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes."+strconv.Itoa(schema.HashInt(200)), "200"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdateDisableHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
				),
			},
			{
				Config: testAccGaapHttpRuleUpdateForwardHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "forward_host", "www.qqq.com"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_httpUpdateRealservers(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes."+strconv.Itoa(schema.HashInt(200)), "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
			{
				Config: testAccGaapHttpRuleHttpUpdateRealservers,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_noHealth(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleNoHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "HEAD"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpRule_domainRealserver(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleDomainRealserver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#"),
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

func TestAccTencentCloudGaapHttpRule_noRealserver(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleNoRealserver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes."+strconv.Itoa(schema.HashInt(200)), "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "forward_host", "default"),
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

func TestAccTencentCloudGaapHttpRule_deleteRealserver(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpRuleBasic3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "health_check_status_codes."+strconv.Itoa(schema.HashInt(200)), "200"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "forward_host", "default"),
				),
			},
			{
				Config: testAccGaapHttpRuleNoRealserver2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpRuleExists("tencentcloud_gaap_http_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_rule.foo", "realservers.#", "0"),
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

func testAccCheckGaapHttpRuleExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no http rule id is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rule, err := service.DescribeHttpRule(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if rule == nil {
			return errors.New("rule not exist")
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapHttpRuleDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		if *id == "" {
			return errors.New("http rule id is nil")
		}

		rule, err := service.DescribeHttpRule(context.TODO(), *id)
		if err != nil {
			return err
		}

		if rule != nil {
			return errors.New("http rule still exists")
		}

		return nil
	}
}

var testAccGaapHttpRuleBasic = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7070
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleBasic1 = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7071
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleUpdate = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7071
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/new"
  realserver_type           = "IP"
  scheduler                 = "wrr"
  health_check              = true
  connect_timeout           = 3
  health_check_path         = "/health"
  health_check_method       = "HEAD"
  health_check_status_codes = [100, 200]
  
  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleUpdateDisableHealth = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7071
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/new"
  realserver_type           = "IP"
  health_check              = false
  connect_timeout           = 3
  health_check_path         = "/health"
  health_check_method       = "HEAD"
  health_check_status_codes = [100, 200]
  
  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleUpdateForwardHost = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7071
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/new"
  realserver_type           = "IP"
  health_check              = false
  connect_timeout           = 3
  health_check_path         = "/health"
  health_check_method       = "HEAD"
  health_check_status_codes = [100, 200]
  
  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  forward_host = "www.qqq.com"
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleBasic2 = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7072
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleBasic3 = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7077
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleHttpUpdateRealservers = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7078
  proxy_id = "%s"
}


resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1)

var testAccGaapHttpRuleNoHealth = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7073
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo.id
  domain          = tencentcloud_gaap_http_domain.foo.domain
  path            = "/"
  realserver_type = "IP"
  health_check    = false

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverIpId1, defaultGaapRealserverIp1, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)

var testAccGaapHttpRuleDomainRealserver = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7074
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo.id
  domain          = tencentcloud_gaap_http_domain.foo.domain
  path            = "/"
  realserver_type = "DOMAIN"
  health_check    = false

  realservers {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realservers {
    id   = "%s"
    ip   = "%s"
	port = 80
  }
}
`, defaultGaapProxyId, defaultGaapRealserverDomainId1, defaultGaapRealserverDomain1, defaultGaapRealserverDomainId2, defaultGaapRealserverDomain2)

var testAccGaapHttpRuleNoRealserver = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7075
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
}
`, defaultGaapProxyId)

var testAccGaapHttpRuleNoRealserver2 = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7077
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_http_rule "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]
}
`, defaultGaapProxyId)
