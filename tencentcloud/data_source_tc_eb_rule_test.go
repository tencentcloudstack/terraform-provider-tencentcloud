package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudEbRuleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbRuleDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_rule.rule")),
			},
		},
	})
}

const testAccEbRuleDataSource = `

data "tencentcloud_eb_rule" "rule" {
  event_bus_id = ""
  order_by = ""
  order = ""
    tags = {
    "createdBy" = "terraform"
  }
}

`
