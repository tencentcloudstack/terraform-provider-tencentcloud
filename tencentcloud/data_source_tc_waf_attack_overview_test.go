package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafAttackOverviewDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafAttackOverviewDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_attack_overview.attack_overview")),
			},
		},
	})
}

const testAccWafAttackOverviewDataSource = `

data "tencentcloud_waf_attack_overview" "attack_overview" {
  from_time = ""
  to_time = ""
  appid = 
  domain = ""
  edition = ""
  instance_i_d = ""
              }

`
