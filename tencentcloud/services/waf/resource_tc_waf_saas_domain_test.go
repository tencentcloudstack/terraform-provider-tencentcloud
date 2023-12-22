package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafSaasDomainResource_basic -v
func TestAccTencentCloudWafSaasDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafSaasDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_saas_domain.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_saas_domain.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafSaasDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_saas_domain.example", "id"),
				),
			},
		},
	})
}

const testAccWafSaasDomain = `
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id     = "waf_2kxtlbky01b3wceb"
  domain          = "tf.qcloudwaf.com"
  is_cdn          = 3
  cert_type       = 2
  ssl_id          = "3a6B5y8v"
  load_balance    = "2"
  https_rewrite   = 1
  is_http2        = 1
  upstream_scheme = "https"
  src_list        = [
    "1.1.1.1",
    "2.2.2.2"
  ]
  weights = [
    50,
    60
  ]

  ports {
    port              = "80"
    protocol          = "http"
    upstream_port     = "80"
    upstream_protocol = "http"
  }

  ports {
    port              = "443"
    protocol          = "https"
    upstream_port     = "443"
    upstream_protocol = "https"
  }

  ip_headers = [
    "headers_1",
    "headers_2",
    "headers_3",
  ]

  is_keep_alive      = "1"
  active_check       = 1
  tls_version        = 3
  cipher_template    = 1
  proxy_read_timeout = 500
  proxy_send_timeout = 500
  sni_type           = 3
  sni_host           = "3.3.3.3"
  xff_reset          = 1
  bot_status         = 1
  api_safe_status    = 1
}
`

const testAccWafSaasDomainUpdate = `
resource "tencentcloud_waf_saas_domain" "example" {
  instance_id     = "waf_2kxtlbky01b3wceb"
  domain          = "tf.qcloudwaf.com"
  is_cdn          = 3
  cert_type       = 2
  ssl_id          = "3a6B5y8v"
  load_balance    = "2"
  https_rewrite   = 0
  is_http2        = 0
  upstream_scheme = "https"
  src_list        = [
    "1.1.1.1",
    "2.2.2.2",
	"3.3.3.3",
  ]
  weights = [
    50,
    60,
    70,
  ]

  ports {
    port              = "81"
    protocol          = "http"
    upstream_port     = "81"
    upstream_protocol = "http"
  }

  ports {
    port              = "443"
    protocol          = "https"
    upstream_port     = "443"
    upstream_protocol = "https"
  }

  ports {
    port              = "4443"
    protocol          = "https"
    upstream_port     = "4443"
    upstream_protocol = "https"
  }

  ip_headers = [
    "headers_1",
    "headers_2",
    "headers_3",
    "headers_4"
  ]

  is_keep_alive      = "0"
  active_check       = 0
  tls_version        = 3
  cipher_template    = 1
  proxy_read_timeout = 200
  proxy_send_timeout = 200
  sni_type           = 3
  sni_host           = "4.4.4.4"
  xff_reset          = 0
  bot_status         = 0
  api_safe_status    = 0
}
`
