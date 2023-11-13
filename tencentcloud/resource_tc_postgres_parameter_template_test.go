package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresParameterTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresParameterTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_parameter_template.parameter_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_parameter_template.parameter_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresParameterTemplate = `

resource "tencentcloud_postgres_parameter_template" "parameter_template" {
  template_name = "test_param_template"
  d_b_major_version = "13"
  d_b_engine = "postgresql"
  template_description = "test use"
}

`
