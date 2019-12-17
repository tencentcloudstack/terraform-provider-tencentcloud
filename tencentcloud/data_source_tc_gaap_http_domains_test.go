package tencentcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapHttpDomains_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapHttpDomainsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_http_domains.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_http_domains.foo", "domains.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_http_domains.foo", "domains.0.domain", regexp.MustCompile(`www\.qq\.com`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_domains.foo", "domains.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_domains.foo", "domains.0.client_certificate_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_domains.foo", "domains.0.realserver_auth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_domains.foo", "domains.0.realserver_certificate_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_domains.foo", "domains.0.basic_auth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_http_domains.foo", "domains.0.gaap_auth"),
				),
			},
		},
	})
}

var testAccGaapHttpDomainsBasic = fmt.Sprintf(`
resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "%s"
}

resource "tencentcloud_gaap_http_domain" "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

data "tencentcloud_gaap_http_domains" "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_domain.foo.domain
}
`, defaultGaapProxyId)
