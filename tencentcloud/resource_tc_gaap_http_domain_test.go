package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudGaapHttpDomain_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpDomainDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_http_domain.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_https_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpDomainDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_http_domain.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_httpsMutualAuthentication(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpDomainDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttpsMutualAuthentication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
			{
				Config: testAccGaapHttpDomainHttpsMutualAuthenticationUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", regexp.MustCompile("cert-.")),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", regexp.MustCompile("cert-.")),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "true"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id", regexp.MustCompile("cert-.")),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "true"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth_id", regexp.MustCompile("cert-.")),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "true"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth_id", regexp.MustCompile("cert-.")),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_httpsPolyClientCertificateIds(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpDomainDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttpsPolyClientCertificateIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", regexp.MustCompile("cert-.")),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
			{
				Config: testAccGaapHttpDomainHttpsPolyClientCertificateIdsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", regexp.MustCompile("cert-.")),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_httpsCCIdToPolyIds(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpDomainDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttpsCcId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", regexp.MustCompile("cert-.")),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
			{
				Config: testAccGaapHttpDomainHttpsPolyIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_id", regexp.MustCompile("cert-.")),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_httpsRealserverCertificateIdOldToNew(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpDomainDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttpsRsIdOld,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
			{
				Config: testAccGaapHttpDomainHttpsRsIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapHttpDomain_httpsRealserverCertificateIds(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapHttpDomainDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainHttpsRsIds2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
			{
				Config: testAccGaapHttpDomainHttpsRsIdsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapHttpDomainExists("tencentcloud_gaap_http_domain.foo", id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "domain", "www.qq.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "certificate_id", "default"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_auth", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_http_domain.foo", "realserver_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "realserver_certificate_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "basic_auth", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_http_domain.foo", "gaap_auth", "false"),
				),
			},
		},
	})
}

func testAccCheckGaapHttpDomainExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no domain id is set")
		}

		split := strings.Split(rs.Primary.ID, "+")
		listenerId, domain := split[0], split[2]

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		httpDomain, err := service.DescribeDomain(context.TODO(), listenerId, domain)
		if err != nil {
			return err
		}

		if httpDomain == nil {
			return fmt.Errorf("domain not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapHttpDomainDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		if *id == "" {
			return errors.New("domain id is nil")
		}

		split := strings.Split(*id, "+")
		listenerId, domain := split[0], split[2]

		httpDomain, err := service.DescribeDomain(context.TODO(), listenerId, domain)
		if err != nil {
			return err
		}

		if httpDomain != nil {
			return errors.New("domain still exists")
		}

		return nil
	}
}

var testAccGaapHttpDomainBasic = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 7170
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}
`, defaultGaapProxyId)

var testAccGaapHttpDomainHttps = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 7171
  proxy_id         = "%s"
  certificate_id   = tencentcloud_gaap_certificate.foo.id
  forward_protocol = "HTTP"
  auth_type        = 0
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = "www.qq.com"
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", defaultGaapProxyId)

var testAccGaapHttpDomainHttpsMutualAuthentication = fmt.Sprintf(`
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
  port                  = 7172
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
  forward_protocol      = "HTTPS"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF", defaultGaapProxyId)

var testAccGaapHttpDomainHttpsMutualAuthenticationUpdate = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "server" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "client" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "realserver" {
  type    = "REALSERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "basic" {
  type    = "BASIC"
  content = %s
}

resource tencentcloud_gaap_certificate "gaap" {
  type    = "PROXY"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 7172
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
  forward_protocol      = "HTTPS"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id           = tencentcloud_gaap_layer7_listener.foo.id
  domain                = "www.qq.com"
  certificate_id        = tencentcloud_gaap_certificate.server.id
  client_certificate_id = tencentcloud_gaap_certificate.client.id

  realserver_auth               = true
  realserver_certificate_id     = tencentcloud_gaap_certificate.realserver.id
  realserver_certificate_domain = "qq.com"

  basic_auth    = true
  basic_auth_id = tencentcloud_gaap_certificate.basic.id

  gaap_auth    = true
  gaap_auth_id = tencentcloud_gaap_certificate.gaap.id
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"\"test:tx2KGdo3zJg/.\"",
	"<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF",
	defaultGaapProxyId,
)

var testAccGaapHttpDomainHttpsPolyClientCertificateIds = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "client1" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol                    = "HTTPS"
  name                        = "ci-test-gaap-l7-listener"
  port                        = 7173
  proxy_id                    = "%s"
  certificate_id              = tencentcloud_gaap_certificate.foo.id
  client_certificate_ids      = [tencentcloud_gaap_certificate.bar.id]
  forward_protocol            = "HTTPS"
  auth_type                   = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id                 = tencentcloud_gaap_layer7_listener.foo.id
  domain                      = "www.qq.com"
  client_certificate_ids = [tencentcloud_gaap_certificate.client1.id]
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId)

var testAccGaapHttpDomainHttpsPolyClientCertificateIdsUpdate = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "client2" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "client3" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol                    = "HTTPS"
  name                        = "ci-test-gaap-l7-listener"
  port                        = 7173
  proxy_id                    = "%s"
  certificate_id              = tencentcloud_gaap_certificate.foo.id
  client_certificate_ids = [tencentcloud_gaap_certificate.bar.id]
  forward_protocol            = "HTTPS"
  auth_type                   = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id                 = tencentcloud_gaap_layer7_listener.foo.id
  domain                      = "www.qq.com"
  client_certificate_ids = [tencentcloud_gaap_certificate.client2.id, tencentcloud_gaap_certificate.client3.id]
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId)

var testAccGaapHttpDomainHttpsCcId = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "client1" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol                    = "HTTPS"
  name                        = "ci-test-gaap-l7-listener"
  port                        = 7174
  proxy_id                    = "%s"
  certificate_id              = tencentcloud_gaap_certificate.foo.id
  client_certificate_ids = [tencentcloud_gaap_certificate.bar.id]
  forward_protocol            = "HTTPS"
  auth_type                   = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id           = tencentcloud_gaap_layer7_listener.foo.id
  domain                = "www.qq.com"
  client_certificate_id = tencentcloud_gaap_certificate.client1.id
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId)

var testAccGaapHttpDomainHttpsPolyIds = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "client1" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol                    = "HTTPS"
  name                        = "ci-test-gaap-l7-listener"
  port                        = 7174
  proxy_id                    = "%s"
  certificate_id              = tencentcloud_gaap_certificate.foo.id
  client_certificate_ids = [tencentcloud_gaap_certificate.bar.id]
  forward_protocol            = "HTTPS"
  auth_type                   = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id                = tencentcloud_gaap_layer7_listener.foo.id
  domain                     = "www.qq.com"
 client_certificate_ids = [tencentcloud_gaap_certificate.client1.id]
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId)

var testAccGaapHttpDomainHttpsRsIdOld = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "realserver1" {
  type    = "REALSERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 7176
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
  forward_protocol      = "HTTPS"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"

  realserver_auth               = true
  realserver_certificate_id     = tencentcloud_gaap_certificate.realserver1.id
  realserver_certificate_domain = "qq.com"
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId,
)

var testAccGaapHttpDomainHttpsRsIds2 = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "realserver1" {
  type    = "REALSERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 7177
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
  forward_protocol      = "HTTPS"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"

  realserver_auth               = true
  realserver_certificate_ids    = [tencentcloud_gaap_certificate.realserver1.id]
  realserver_certificate_domain = "qq.com"
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId,
)

var testAccGaapHttpDomainHttpsRsIds = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "realserver1" {
  type    = "REALSERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 7176
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
  forward_protocol      = "HTTPS"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"

  realserver_auth               = true
  realserver_certificate_ids    = [tencentcloud_gaap_certificate.realserver1.id]
  realserver_certificate_domain = "qq.com"
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId,
)

var testAccGaapHttpDomainHttpsRsIdsUpdate = fmt.Sprintf(`
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

resource tencentcloud_gaap_certificate "realserver1" {
  type    = "REALSERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "realserver2" {
  type    = "REALSERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 7177
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
  forward_protocol      = "HTTPS"
  auth_type             = 1
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"

  realserver_auth               = true
  realserver_certificate_ids    = [tencentcloud_gaap_certificate.realserver1.id, tencentcloud_gaap_certificate.realserver2.id]
  realserver_certificate_domain = "qq.com"
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	defaultGaapProxyId,
)
