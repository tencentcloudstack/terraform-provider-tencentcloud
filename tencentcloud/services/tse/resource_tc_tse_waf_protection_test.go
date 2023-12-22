package tse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudTseWafProtectionResource_basic -v
func TestAccTencentCloudTseWafProtectionResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseWafProtection_open,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tse_waf_protection.waf_protection", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_waf_protection.waf_protection", "gateway_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_waf_protection.waf_protection", "global_status"),
				),
			},
			{
				Config: testAccTseWafProtection_close,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tse_waf_protection.waf_protection", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_waf_protection.waf_protection", "gateway_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_waf_protection.waf_protection", "global_status"),
				),
			},
		},
	})
}

const testAccTseWafProtection_open = testAccTseWafDomains + `

resource "tencentcloud_tse_waf_protection" "waf_protection" {
  gateway_id = var.gateway_id
  type       = "Route"
  list       = ["17e8ba6a-e136-454b-9cfa-3e541ffd01dd"]
  operate    = "open"

  depends_on = [tencentcloud_tse_waf_domains.waf_domains]
}

`
const testAccTseWafProtection_close = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_waf_protection" "waf_protection" {
  gateway_id = var.gateway_id
  type       = "Route"
  list       = ["17e8ba6a-e136-454b-9cfa-3e541ffd01dd"]
  operate    = "close"
}

`
