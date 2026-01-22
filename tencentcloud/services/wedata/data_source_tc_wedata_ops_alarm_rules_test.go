package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataOpsAlarmRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataOpsAlarmRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_ops_alarm_rules.wedata_ops_alarm_rules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_ops_alarm_rules.wedata_ops_alarm_rules", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_alarm_rules.wedata_ops_alarm_rules", "data.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_ops_alarm_rules.wedata_ops_alarm_rules", "data.0.items.#", "20"),
				),
			},
		},
	})
}

const testAccWedataOpsAlarmRulesDataSource = `

data "tencentcloud_wedata_ops_alarm_rules" "wedata_ops_alarm_rules" {
  project_id = "1859317240494305280"
}
`
