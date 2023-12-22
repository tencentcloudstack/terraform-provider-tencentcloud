package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafModifyAccessPeriodResource_basic -v
func TestAccTencentCloudNeedFixWafModifyAccessPeriodResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafModifyAccessPeriod,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_modify_access_period.modify_access_period", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_modify_access_period.modify_access_period", "topic_id", "1ae37c76-df99-4e2b-998c-20f39eba6226"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_modify_access_period.modify_access_period", "period"),
				),
			},
		},
	})
}

const testAccWafModifyAccessPeriod = `
resource "tencentcloud_waf_modify_access_period" "example" {
  topic_id = "1ae37c76-df99-4e2b-998c-20f39eba6226"
  period   = 30
}
`
