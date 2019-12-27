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

func TestAccTencentCloudGaapDomainErrorPageInfo_basic(t *testing.T) {
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageInfoDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageInfoExists("tencentcloud_gaap_domain_error_page_info.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page_info.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(404)), "404"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(503)), "503"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "body", "bad request"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPageInfo_singleErrorCode(t *testing.T) {
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageInfoDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoSingleErrorCode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageInfoExists("tencentcloud_gaap_domain_error_page_info.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page_info.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(400)), "400"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "body", "bad request"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPageInfo_newErrorCode(t *testing.T) {
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageInfoDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoNewErrorCode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageInfoExists("tencentcloud_gaap_domain_error_page_info.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page_info.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(402)), "402"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "body", "bad request"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "new_error_code", "501"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPageInfo_clearHeaders(t *testing.T) {
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageInfoDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoClearHeaders,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageInfoExists("tencentcloud_gaap_domain_error_page_info.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page_info.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(403)), "403"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "body", "bad request"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "clear_headers."+strconv.Itoa(schema.HashString("Content-Length")), "Content-Length"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "clear_headers."+strconv.Itoa(schema.HashString("X-TEST")), "X-TEST"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPageInfo_setHeaders(t *testing.T) {
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageInfoDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoSetHeaders,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageInfoExists("tencentcloud_gaap_domain_error_page_info.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page_info.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(405)), "405"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "body", "bad request"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapDomainErrorPageInfo_full(t *testing.T) {
	listenerId := new(string)
	domain := new(string)
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapDomainErrorPageInfoDestroy(listenerId, domain, id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapDomainErrorPageInfoExists("tencentcloud_gaap_domain_error_page_info.foo", listenerId, domain, id),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_domain_error_page_info.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(406)), "406"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "error_codes."+strconv.Itoa(schema.HashInt(504)), "504"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "body", "bad request"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "clear_headers."+strconv.Itoa(schema.HashString("Content-Length")), "Content-Length"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "clear_headers."+strconv.Itoa(schema.HashString("X-TEST")), "X-TEST"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_domain_error_page_info.foo", "set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

func testAccCheckGaapDomainErrorPageInfoExists(n string, listenerId, domain, id *string) resource.TestCheckFunc {
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

func testAccCheckGaapDomainErrorPageInfoDestroy(listenerId, domain, id *string) resource.TestCheckFunc {
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

const testAccGaapDomainErrorPageInfoListenerAndDomain = `
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}`

var testAccGaapDomainErrorPageInfoBasic = fmt.Sprintf(testAccGaapDomainErrorPageInfoListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page_info "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_domain.foo.domain
  error_codes = [404, 503]
  body        = "bad request"
}
`, defaultGaapProxyId)

var testAccGaapDomainErrorPageInfoSingleErrorCode = fmt.Sprintf(testAccGaapDomainErrorPageInfoListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page_info "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_domain.foo.domain
  error_codes = [400]
  body        = "bad request"
}
`, defaultGaapProxyId)

var testAccGaapDomainErrorPageInfoNewErrorCode = fmt.Sprintf(testAccGaapDomainErrorPageInfoListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page_info "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [402]
  body           = "bad request"
  new_error_code = 501
}
`, defaultGaapProxyId)

var testAccGaapDomainErrorPageInfoClearHeaders = fmt.Sprintf(testAccGaapDomainErrorPageInfoListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page_info "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [403]
  body           = "bad request"
  clear_headers  = ["Content-Length", "X-TEST"]
}
`, defaultGaapProxyId)

var testAccGaapDomainErrorPageInfoSetHeaders = fmt.Sprintf(testAccGaapDomainErrorPageInfoListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page_info "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [405]
  body           = "bad request"

  set_headers = {
    "X-TEST" = "test"
  }
}
`, defaultGaapProxyId)

var testAccGaapDomainErrorPageInfoFull = fmt.Sprintf(testAccGaapDomainErrorPageInfoListenerAndDomain+`
resource tencentcloud_gaap_domain_error_page_info "foo" {
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
`, defaultGaapProxyId)
