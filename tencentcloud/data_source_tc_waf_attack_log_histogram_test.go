package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafAttackLogHistogramDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafAttackLogHistogramDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_attack_log_histogram.attack_log_histogram")),
			},
		},
	})
}

const testAccWafAttackLogHistogramDataSource = `

data "tencentcloud_waf_attack_log_histogram" "attack_log_histogram" {
  domain = ""
  start_time = ""
  end_time = ""
  query_string = ""
    }

`
