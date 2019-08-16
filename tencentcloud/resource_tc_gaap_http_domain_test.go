package tencentcloud

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapHttpDomain_basic(t *testing.T) {
	id := new(string)
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, proxyId, "HTTP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainBasic,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_domain"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_https(t *testing.T) {
	id := new(string)
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, proxyId, "HTTP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttps,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_domain"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_httpsMutualAuthentication(t *testing.T) {
	id := new(string)
	proxyId := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, proxyId, "HTTP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttpsMutualAuthentication,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", "default"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_domain"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth_id"),
				),
			},
			{
				Config: testAccGaapHttpDomainHttpsMutualAuthenticationUpdate,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, proxyId, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "client_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "realserver_certificate_domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "basic_auth_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "gaap_auth_id"),
				),
			},
		},
	})
}

const testAccGaapHttpDomainBasic = `
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

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}
`

var testAccGaapHttpDomainHttps = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id   = "${tencentcloud_gaap_certificate.foo.id}"
  forward_protocol = "HTTP"
  auth_type        = 0
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id    = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain         = "www.qq.com"
  certificate_id = "default"
}

` + testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF")

var testAccGaapHttpDomainHttpsMutualAuthentication = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 80
  proxy_id              = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id        = "${tencentcloud_gaap_certificate.foo.id}"
  client_certificate_id = "${tencentcloud_gaap_certificate.bar.id}"
  forward_protocol      = "HTTP"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id           = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                = "www.qq.com"
  certificate_id        = "default"
  client_certificate_id = "default"
}

` + testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF") +
	strings.Replace(testAccGaapCertificate(1, "<<EOF\n"+testAccGaapCertificateClientCA+"EOF", "", "<<EOF\n"+testAccGaapCertificateClientCAKey+"EOF"), "foo", "bar", 1)

var testAccGaapHttpDomainHttpsMutualAuthenticationUpdate = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 80
  proxy_id              = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id        = "${tencentcloud_gaap_certificate.foo.id}"
  client_certificate_id = "${tencentcloud_gaap_certificate.bar.id}"
  forward_protocol      = "HTTP"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id           = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                = "www.qq.com"
  certificate_id        = "${tencentcloud_gaap_certificate.server.id}"
  client_certificate_id = "${tencentcloud_gaap_certificate.client.id}"

  realserver_auth               = true
  realserver_certificate_id     = "${tencentcloud_gaap_certificate.realserver.id}"
  realserver_certificate_domain = "test"

  basic_auth    = true
  basic_auth_id = "${tencentcloud_gaap_certificate.basic.id}"

  gaap_auth    = true
  gaap_auth_id = "${tencentcloud_gaap_certificate.gaap.id}"
}

` + testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF") +
	strings.Replace(testAccGaapCertificate(1, "<<EOF\n"+testAccGaapCertificateClientCA+"EOF", "", "<<EOF\n"+testAccGaapCertificateClientCAKey+"EOF"), "foo", "bar", 1) +
	strings.Replace(testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF"), "foo", "server", 1) +
	strings.Replace(testAccGaapCertificate(1, "<<EOF\n"+testAccGaapCertificateClientCA+"EOF", "", "<<EOF\n"+testAccGaapCertificateClientCAKey+"EOF"), "foo", "client", 1) +
	strings.Replace(testAccGaapCertificate(3, "<<EOF\n"+testAccGaapCertificateClientCA+"EOF", "", "<<EOF\n"+testAccGaapCertificateClientCAKey+"EOF"), "foo", "realserver", 1) +
	strings.Replace(testAccGaapCertificate(0, "\"test:tx2KGdo3zJg/.\"", "", ""), "foo", "basic", 1) +
	strings.Replace(testAccGaapCertificate(4, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF"), "foo", "gaap", 1)
