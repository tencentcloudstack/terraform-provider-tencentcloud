package tse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudTseWafDomainsResource_basic -v
func TestAccTencentCloudTseWafDomainsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseWafDomains,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tse_waf_domains.waf_domains", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_waf_domains.waf_domains", "domain", "tse.exmaple.com"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_waf_domains.waf_domains",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseWafDomains = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_waf_domains" "waf_domains" {
  domain     = "tse.exmaple.com"
  gateway_id = var.gateway_id
}

`
