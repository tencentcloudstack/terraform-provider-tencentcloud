package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafSaasInstanceResource_basic -v
func TestAccTencentCloudNeedFixWafSaasInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafSaasInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_saas_instance.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_saas_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafSaasInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_saas_instance.example", "id"),
				),
			},
		},
	})
}

const testAccWafSaasInstance = `
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category   = "premium_saas"
  instance_name    = "tf-example-saas-waf"
  time_span        = 1
  time_unit        = "m"
  auto_renew_flag  = 1
  elastic_mode     = 1
  is_cn_mainland   = 1
  real_region      = "gz"
  domain_pkg_count = 3
  qps_pkg_count    = 3
}
`

const testAccWafSaasInstanceUpdate = `
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category   = "premium_saas"
  instance_name    = "tf-example-saas-waf-update"
  time_span        = 1
  time_unit        = "m"
  auto_renew_flag  = 0
  elastic_mode     = 0
  is_cn_mainland   = 1
  real_region      = "gz"
  domain_pkg_count = 3
  qps_pkg_count    = 3
}
`
