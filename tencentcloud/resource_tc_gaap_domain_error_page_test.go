package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudGaapDomainErrorPage_basic(t *testing.T) {
	t.Parallel()
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageExists("tencentcloud_gaap_domain_error_page.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page.foo", "listener_id"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "404"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "503"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "body", "bad request"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPage_singleErrorCode(t *testing.T) {
	t.Parallel()
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageSingleErrorCode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageExists("tencentcloud_gaap_domain_error_page.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page.foo", "listener_id"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "400"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "body", "bad request"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPage_newErrorCode(t *testing.T) {
	t.Parallel()
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageNewErrorCode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageExists("tencentcloud_gaap_domain_error_page.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page.foo", "listener_id"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "402"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "body", "bad request"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "new_error_code", "501"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPage_clearHeaders(t *testing.T) {
	t.Parallel()
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageClearHeaders,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageExists("tencentcloud_gaap_domain_error_page.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page.foo", "listener_id"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "403"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "body", "bad request"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "clear_headers.*", "Content-Length"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "clear_headers.*", "X-TEST"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPage_setHeaders(t *testing.T) {
	t.Parallel()
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageSetHeaders,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageExists("tencentcloud_gaap_domain_error_page.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page.foo", "listener_id"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "405"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "body", "bad request"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPage_full(t *testing.T) {
	t.Parallel()
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageExists("tencentcloud_gaap_domain_error_page.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page.foo", "listener_id"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "406"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "error_codes.*", "504"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "body", "bad request"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "clear_headers.*", "Content-Length"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_gaap_domain_error_page.foo", "clear_headers.*", "X-TEST"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page.foo", "set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

func testAccCheckGaapDomainErrorPageExists(n string, listenerId, domain, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no domain error page info id is set")
		}

		*listenerId = rs.Primary.Attributes["listener_id"]
		*domain = rs.Primary.Attributes["domain"]

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		info, err := service.DescribeDomainErrorPageInfo(context.TODO(), *listenerId, *domain, rs.Primary.ID)
		if err != nil {
			return err
		}

		if info == nil {
			return errors.New("domain error page info not exist")
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapDomainErrorPageDestroy(listenerId, domain, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		if *id == "" {
			return errors.New("domain error page info id is nil")
		}

		info, err := service.DescribeDomainErrorPageInfo(context.TODO(), *listenerId, *domain, *id)
		if err != nil {
			return err
		}

		if info != nil {
			return errors.New("domain error page info still exists")
		}

		return nil
	}
}

const testAccGaapDomainErrorPageListenerAndDomain = `
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = "%s"
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}`

var testAccGaapDomainErrorPageBasic = fmt.Sprintf(testAccGaapDomainErrorPageListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_domain.foo.domain
  error_codes = [404, 503]
  body        = "bad request"
}
`, "81", defaultGaapProxyId)

var testAccGaapDomainErrorPageSingleErrorCode = fmt.Sprintf(testAccGaapDomainErrorPageListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_domain.foo.domain
  error_codes = [400]
  body        = "bad request"
}
`, "82", defaultGaapProxyId)

var testAccGaapDomainErrorPageNewErrorCode = fmt.Sprintf(testAccGaapDomainErrorPageListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [402]
  body           = "bad request"
  new_error_code = 501
}
`, "83", defaultGaapProxyId)

var testAccGaapDomainErrorPageClearHeaders = fmt.Sprintf(testAccGaapDomainErrorPageListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [403]
  body           = "bad request"
  clear_headers  = ["Content-Length", "X-TEST"]
}
`, "84", defaultGaapProxyId)

var testAccGaapDomainErrorPageSetHeaders = fmt.Sprintf(testAccGaapDomainErrorPageListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [405]
  body           = "bad request"

  set_headers = {
    "X-TEST" = "test"
  }
}
`, "85", defaultGaapProxyId)

var testAccGaapDomainErrorPageFull = fmt.Sprintf(testAccGaapDomainErrorPageListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [406, 504]
  new_error_code = 502
  body           = "bad request"
  clear_headers  = ["Content-Length", "X-TEST"]

  set_headers = {
    "X-TEST" = "test"
  }
}
`, "86", defaultGaapProxyId)
