package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafClbInstanceResource_basic -v
func TestAccTencentCloudNeedFixWafClbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafClbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_clb_instance.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_clb_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafClbInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_clb_instance.example", "id"),
				),
			},
		},
	})
}

const testAccWafClbInstance = `
resource "tencentcloud_waf_clb_instance" "example" {
  goods_category   = "ultimate_clb"
  instance_name    = "tf-example-clb-waf"
  time_span        = 1
  time_unit        = "m"
  auto_renew_flag  = 1
  elastic_mode     = 1
  is_cn_mainland   = 1
  domain_pkg_count = 3
  qps_pkg_count    = 3
}
`

const testAccWafClbInstanceUpdate = `
resource "tencentcloud_waf_clb_instance" "example" {
  goods_category   = "ultimate_clb"
  instance_name    = "tf-example-clb-waf-update"
  time_span        = 1
  time_unit        = "m"
  auto_renew_flag  = 0
  elastic_mode     = 0
  is_cn_mainland   = 1
  domain_pkg_count = 3
  qps_pkg_count    = 3
}
`
