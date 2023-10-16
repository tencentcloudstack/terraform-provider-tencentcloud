package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCssLiveDomainCertDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssLiveDomainCertDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_live_domain_cert.cert"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_live_domain_cert.cert", "domain_cert_info.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_live_domain_cert.cert", "domain_cert_info.0.cert_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_css_live_domain_cert.cert", "domain_cert_info.0.cert_name", "keep_ssl_css_domain_test"),
					resource.TestCheckResourceAttr("data.tencentcloud_css_live_domain_cert.cert", "domain_cert_info.0.domain_name", "test122.jingxhu.top"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_live_domain_cert.cert", "domain_cert_info.0.create_time"),
				),
			},
		},
	})
}

const testAccCssLiveDomainCertDataSource = `

data "tencentcloud_css_live_domain_cert" "cert" {
  domain_name = "test122.jingxhu.top"
}

`
