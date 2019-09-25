package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudGaapLayer7Listener_basic(t *testing.T) {
	id := new(string)
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, proxyId, "HTTP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener-new"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_https(t *testing.T) {
	id := new(string)
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, proxyId, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener-new"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsTwoWayAuthentication(t *testing.T) {
	id := new(string)
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, proxyId, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsTwoWayAuthentication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsForwardHttps(t *testing.T) {
	id := new(string)
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, proxyId, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsForwardHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

func testAccCheckGaapLayer7ListenerExists(n string, id, proxyId *string, protocol string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no listener ID is set")
		}

		attrProxyId := rs.Primary.Attributes["proxy_id"]
		if attrProxyId == "" {
			return errors.New("no proxy ID is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		switch protocol {
		case "HTTP":
			listeners, err := service.DescribeHTTPListeners(context.TODO(), &attrProxyId, &rs.Primary.ID, nil, nil)
			if err != nil {
				return err
			}

			for _, l := range listeners {
				if l.ListenerId == nil {
					return errors.New("listener id is nil")
				}
				if rs.Primary.ID == *l.ListenerId {
					*id = *l.ListenerId
					*proxyId = attrProxyId
					break
				}
			}

		case "HTTPS":
			listeners, err := service.DescribeHTTPSListeners(context.TODO(), &attrProxyId, &rs.Primary.ID, nil, nil)
			if err != nil {
				return err
			}

			for _, l := range listeners {
				if l.ListenerId == nil {
					return errors.New("listener id is nil")
				}
				if rs.Primary.ID == *l.ListenerId {
					*id = *l.ListenerId
					*proxyId = attrProxyId
					break
				}
			}
		}

		if id == nil {
			return fmt.Errorf("listener not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckGaapLayer7ListenerDestroy(id, proxyId *string, protocol string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		switch protocol {
		case "HTTP":
			listeners, err := service.DescribeHTTPListeners(context.TODO(), proxyId, id, nil, nil)
			if err != nil {
				return err
			}
			if len(listeners) > 0 {
				return errors.New("listener still exists")
			}

		case "HTTPS":
			listeners, err := service.DescribeHTTPSListeners(context.TODO(), proxyId, id, nil, nil)
			if err != nil {
				return err
			}
			if len(listeners) > 0 {
				return errors.New("listener still exists")
			}
		}

		return nil
	}
}

var testAccGaapLayer7ListenerBasic = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "%s"
}
`, GAAP_PROXY_ID)

var testAccGaapLayer7ListenerHttpUpdateName = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener-new"
  port     = 80
  proxy_id = "%s"
}
`, GAAP_PROXY_ID)

var testAccGaapLayer7ListenerHttps = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "%s"
  certificate_id   = "${tencentcloud_gaap_certificate.foo.id}"
  forward_protocol = "HTTP"
  auth_type        = 0
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", GAAP_PROXY_ID)

var testAccGaapLayer7ListenerHttpsUpdate = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener-new"
  port             = 80
  proxy_id         = "%s"
  certificate_id   = "${tencentcloud_gaap_certificate.bar.id}"
  forward_protocol = "HTTP"
  auth_type        = 0
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", GAAP_PROXY_ID)

var testAccGaapLayer7ListenerHttpsTwoWayAuthentication = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 80
  proxy_id              = "%s"
  certificate_id        = "${tencentcloud_gaap_certificate.foo.id}"
  forward_protocol      = "HTTP"
  auth_type             = 1
  client_certificate_id = "${tencentcloud_gaap_certificate.bar.id}"
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF", GAAP_PROXY_ID)

var testAccGaapLayer7ListenerHttpsForwardHttps = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "%s"
  certificate_id   = "${tencentcloud_gaap_certificate.foo.id}"
  forward_protocol = "HTTPS"
  auth_type        = 0
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", GAAP_PROXY_ID)
