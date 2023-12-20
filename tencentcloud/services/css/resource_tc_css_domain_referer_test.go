package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssDomainRefererResource_basic -v
func TestAccTencentCloudCssDomainRefererResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssDomainReferer,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_domain_referer.domain_referer", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "allow_empty", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "domain_name", tcacctest.DefaultCSSPlayDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "enable", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "rules", "example.com"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "type", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_domain_referer.domain_referer",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCssDomainRefererUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_domain_referer.domain_referer", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "allow_empty", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "domain_name", tcacctest.DefaultCSSPlayDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "enable", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "rules", "example.com"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain_referer.domain_referer", "type", "0"),
				),
			},
		},
	})
}

const testAccCssDomainRefererVar = `
variable "domain_name" {
  default = "` + tcacctest.DefaultCSSPlayDomainName + `"
}
`

const testAccCssDomainReferer = testAccCssDomainRefererVar + `

resource "tencentcloud_css_domain_referer" "domain_referer" {
  allow_empty = 1
  domain_name = var.domain_name
  enable      = 1
  rules       = "example.com"
  type        = 0
}

`

const testAccCssDomainRefererUp = testAccCssDomainRefererVar + `

resource "tencentcloud_css_domain_referer" "domain_referer" {
  allow_empty = 1
  domain_name = var.domain_name
  enable      = 0
  rules       = "example.com"
  type        = 0
}

`
