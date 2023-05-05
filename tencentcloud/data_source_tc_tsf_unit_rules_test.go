package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfUnitRulesDataSource_basic -v
func TestAccTencentCloudTsfUnitRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUnitRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_unit_rules.unit_rules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.gateway_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.updated_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.dest_namespace_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.dest_namespace_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.priority"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.relationship"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.unit_rule_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.unit_rule_tag_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.unit_rule_tag_list.0.tag_field"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.unit_rule_tag_list.0.tag_operator"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.unit_rule_tag_list.0.tag_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.unit_rule_tag_list.0.tag_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_unit_rules.unit_rules", "result.0.content.0.unit_rule_item_list.0.unit_rule_tag_list.0.unit_rule_item_id"),
				),
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
