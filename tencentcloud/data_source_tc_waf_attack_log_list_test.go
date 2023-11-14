package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafAttackLogListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafAttackLogListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_attack_log_list.attack_log_list")),
			},
		},
	})
}

const testAccWafAttackLogListDataSource = `

data "tencentcloud_waf_attack_log_list" "attack_log_list" {
  domain = ""
  start_time = ""
  end_time = ""
    query_string = ""
  sort = ""
      }

`
