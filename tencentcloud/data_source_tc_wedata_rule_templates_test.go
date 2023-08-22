package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixWedataRuleTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataRuleTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_rule_templates.rule_templates")),
			},
		},
	})
}

const testAccWedataRuleTemplatesDataSource = `

data "tencentcloud_wedata_rule_templates" "rule_templates" {
  type                = 2
  source_object_type  = 2
  project_id          = "1840731346428280832"
  source_engine_types = [2, 4, 16]
}

`
