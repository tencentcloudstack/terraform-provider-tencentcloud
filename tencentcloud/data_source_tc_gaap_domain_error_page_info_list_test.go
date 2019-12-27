package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudGaapDomainErrorPageInfoList_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoListBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_domain_error_page_info_list.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.domain", "www.qq.com"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.error_codes.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.body", "bad request"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.new_error_codes", "502"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.clear_headers.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapDomainErrorPageInfoList_Ids(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfoListIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_domain_error_page_info_list.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "domain"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.domain", "www.qq.com"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.error_codes.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.body", "bad request"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.new_error_codes", "502"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.clear_headers.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_domain_error_page_info_list.foo", "error_page_info_list.0.set_headers.X-TEST", "test"),
				),
			},
		},
	})
}

const testAccGaapDomainErrorPageInfoListListenerAndDomain = `
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

var testAccGaapDomainErrorPageInfoListBasic = fmt.Sprintf(testAccGaapDomainErrorPageInfoListListenerAndDomain+`
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

data tencentcloud_gaap_domain_error_page_info_list "foo" {
  listener_id = tencentcloud_gaap_domain_error_page_info.foo.listener_id
  domain      = tencentcloud_gaap_domain_error_page_info.foo.domain
}
`, defaultGaapProxyId)

var testAccGaapDomainErrorPageInfoListIds = fmt.Sprintf(testAccGaapDomainErrorPageInfoListListenerAndDomain+`
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

resource tencentcloud_gaap_domain_error_page_info "bar" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [403]
  new_error_code = 502
  body           = "bad request"
}

data tencentcloud_gaap_domain_error_page_info_list "foo" {
  listener_id = tencentcloud_gaap_domain_error_page_info.foo.listener_id
  domain      = tencentcloud_gaap_domain_error_page_info.foo.domain
  ids         = [tencentcloud_gaap_domain_error_page_info.foo.id]
}
`, defaultGaapProxyId)
