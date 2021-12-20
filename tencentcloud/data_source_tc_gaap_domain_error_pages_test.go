package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudGaapDomainErrorPages_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPagesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_domain_error_pages.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.domain", "www.qq.com"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.error_codes.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.body", "bad request"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.new_error_codes", "502"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.clear_headers.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapDomainErrorPages_Ids(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPagesIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_domain_error_pages.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.domain", "www.qq.com"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.error_codes.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.body", "bad request"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.new_error_codes", "502"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.clear_headers.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_pages.foo", "error_page_info_list.0.set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

const testAccGaapDomainErrorPagesListenerAndDomain = `
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

var testAccGaapDomainErrorPagesBasic = fmt.Sprintf(testAccGaapDomainErrorPagesListenerAndDomain+`
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

data tencentcloud_gaap_domain_error_pages "foo" {
  listener_id = tencentcloud_gaap_domain_error_page.foo.listener_id
  domain      = tencentcloud_gaap_domain_error_page.foo.domain
}
`, defaultGaapProxyId)

var testAccGaapDomainErrorPagesIds = fmt.Sprintf(testAccGaapDomainErrorPagesListenerAndDomain+`
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

resource tencentcloud_gaap_domain_error_page "bar" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [403]
  new_error_code = 502
  body           = "bad request"
}

data tencentcloud_gaap_domain_error_pages "foo" {
  listener_id = tencentcloud_gaap_domain_error_page.foo.listener_id
  domain      = tencentcloud_gaap_domain_error_page.foo.domain
  ids         = [tencentcloud_gaap_domain_error_page.foo.id]
}
`, defaultGaapProxyId)
