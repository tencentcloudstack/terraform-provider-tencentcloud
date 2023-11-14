package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbParamTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbParamTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_param_template.param_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_param_template.param_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbParamTemplate = `

resource "tencentcloud_cynosdb_param_template" "param_template" {
  template_name = ""
  engine_version = "5.7"
  template_description = ""
  template_id = 1000
  db_mode = "NORMAL"
  param_list {
		param_name = ""
		current_value = ""
		old_value = ""

  }
}

`
