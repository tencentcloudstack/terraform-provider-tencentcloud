package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlParamTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
					"param_list",
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
  param_list {
    current_value = "1"
    name          = "auto_increment_increment"
  }
  param_list {
    current_value = "1"
    name          = "auto_increment_offset"
  }
  param_list {
    current_value = "ON"
    name          = "automatic_sp_privileges"
  }
  template_type = "HIGH_STABILITY"
  engine_type   = "InnoDB"
}

`
