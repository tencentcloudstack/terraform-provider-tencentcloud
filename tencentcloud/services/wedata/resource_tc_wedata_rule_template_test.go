package wedata_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixWedataRuleTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataRuleTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_rule_template.rule_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_rule_template.rule_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataRuleTemplate = `

resource "tencentcloud_wedata_rule_template" "rule_template" {
  project_id          = "1840731346428280832"
  type                = 2
  name                = "tf-test"
  quality_dim         = 3
  source_object_type  = 2
  description         = "for tf test"
  source_engine_types = [2, 4, 16]
  multi_source_flag   = false
  sql_expression      = base64encode("select * from db")
  where_flag          = false
}

`
