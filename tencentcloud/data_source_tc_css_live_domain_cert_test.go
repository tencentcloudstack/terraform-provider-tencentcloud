package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCssLiveDomainCertDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssLiveDomainCertDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_live_domain_cert.live_domain_cert")),
			},
		},
	})
}

const testAccCssLiveDomainCertDataSource = `

data "tencentcloud_css_live_domain_cert" "live_domain_cert" {
  domain_name = ""
  }

`
