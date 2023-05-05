package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlParamTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlParamTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_param_template.param_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_param_template.param_template",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"engine_type",
				},
			},
		},
	})
}

const testAccMysqlParamTemplate = `

resource "tencentcloud_mysql_param_template" "param_template" {
  name           = "terraform-test"
  description    = "terraform-test"
  engine_version = "8.0"
  template_type  = "HIGH_STABILITY"
  engine_type    = "InnoDB"
}

`
