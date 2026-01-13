package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataQualityRuleTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataQualityRuleTemplatesDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_quality_rule_templates.wedata_quality_rule_templates", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_quality_rule_templates.wedata_quality_rule_templates", "data.#"),
				),
			},
		},
	})
}

const testAccWedataQualityRuleTemplatesDataSource = `
data "tencentcloud_wedata_quality_rule_templates" "wedata_quality_rule_templates" {
  project_id = "1840731346428280832"
  
  order_fields {
    name      = "CitationCount"
    direction = "DESC"
  }
  
  filters {
    name   = "Type"
    values = ["1"]
  }
}
`
