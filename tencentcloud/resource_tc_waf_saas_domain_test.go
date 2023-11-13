package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafSaasDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafSaasDomain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_saas_domain.saas_domain", "id")),
			},
			{
				ResourceName:      "tencentcloud_waf_saas_domain.saas_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafSaasDomain = `

resource "tencentcloud_waf_saas_domain" "saas_domain" {
  domain = ""
  cert_type = 
  is_cdn = 
  upstream_type = 
  is_websocket = 
  load_balance = ""
  cert = ""
  private_key = ""
  s_s_l_id = ""
  resource_id = ""
  upstream_scheme = ""
  https_upstream_port = ""
  is_gray = 
  gray_areas = 
  upstream_domain = ""
  src_list = 
  is_http2 = 
  https_rewrite = 
  ports {
		port = ""
		protocol = ""
		upstream_port = ""
		upstream_protocol = ""
		nginx_server_id = ""

  }
  edition = ""
  is_keep_alive = ""
  instance_i_d = ""
  anycast = 
  weights = 
  active_check = 
  t_l_s_version = 
  ciphers = 
  cipher_template = 
  proxy_read_timeout = 
  proxy_send_timeout = 
  sni_type = 
  sni_host = ""
  ip_headers = 
  x_f_f_reset = 
}

`
