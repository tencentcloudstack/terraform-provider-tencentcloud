package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

const testAccTsfUnitRulesDataSourceVar = `
variable "gateway_instance_id" {
	default = "` + defaultTsfGateway + `"
}
`

const testAccTsfUnitRulesDataSource = testAccTsfUnitRulesDataSourceVar + `

data "tencentcloud_tsf_unit_rules" "unit_rules" {
	gateway_instance_id = var.gateway_instance_id
	status = "disabled"
}

`
