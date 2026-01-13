package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataQualityRuleGroupExecResultsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataQualityRuleGroupExecResultsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_quality_rule_group_exec_results.wedata_quality_rule_group_exec_results"),
				),
			},
		},
	})
}

const testAccWedataQualityRuleGroupExecResultsDataSource = `
data "tencentcloud_wedata_quality_rule_group_exec_results" "wedata_quality_rule_group_exec_results" {
  project_id = "3016337760439783424"
  filters {
    name   = "InstanceStatus"
    values = ["FAILED"]
  }
  order_fields {
    name      = "UpdateTime"
    direction = "ASC"
  }
}
`
