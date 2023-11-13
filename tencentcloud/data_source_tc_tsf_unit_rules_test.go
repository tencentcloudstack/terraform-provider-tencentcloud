package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfUnitRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUnitRulesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_unit_rules.unit_rules")),
			},
		},
	})
}

const testAccTsfUnitRulesDataSource = `

data "tencentcloud_tsf_unit_rules" "unit_rules" {
  gateway_instance_id = "gw-ins-lvdypq5k"
  search_word = "test"
  status = ""
  }

`
