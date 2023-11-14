package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresParameterTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresParameterTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgres_parameter_templates.parameter_templates")),
			},
		},
	})
}

const testAccPostgresParameterTemplatesDataSource = `

data "tencentcloud_postgres_parameter_templates" "parameter_templates" {
  filters {
		name = "DBEngine"
		values = 

  }
  limit = 20
  offset = 0
  order_by = "CreateTime"
  order_by_type = "desc"
  template_id = "0b3eaa95-dfba-5253-8cdd-1258ae34d596"
  template_name = "test template"
  d_b_major_version = "14"
  d_b_engine = "PostgreSQL"
  template_description = "for test"
}

`
