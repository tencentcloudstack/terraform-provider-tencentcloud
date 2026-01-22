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
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_quality_rule_templates.example", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_quality_rule_templates.example", "data.#"),
				),
			},
		},
	})
}

const testAccWedataQualityRuleTemplatesDataSource = `
data "tencentcloud_wedata_quality_rule_templates" "example" {
  project_id = "3016337760439783424"
  
  order_fields {
    name      = "CitationCount"
    direction = "DESC"
  }
}
`
