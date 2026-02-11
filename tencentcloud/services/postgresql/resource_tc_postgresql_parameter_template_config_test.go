package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlParameterTemplateConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlParameterTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_parameter_template_config.example", "id"),
				),
			},
			{
				Config: testAccPostgresqlParameterTemplateConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_parameter_template_config.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_postgresql_parameter_template_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlParameterTemplateConfig = `
resource "tencentcloud_postgresql_parameter_template" "example" {
  template_name        = "tf-example"
  db_major_version     = "18"
  db_engine            = "postgresql"
  template_description = "remark."
}

resource "tencentcloud_postgresql_parameter_template_config" "example" {
  template_id = tencentcloud_postgresql_parameter_template.example.id
  modify_param_entry_set {
    name           = "min_parallel_index_scan_size"
    expected_value = "64"
  }

  modify_param_entry_set {
    name           = "enable_async_append"
    expected_value = "on"
  }

  modify_param_entry_set {
    name           = "enable_group_by_reordering"
    expected_value = "on"
  }
}
`

const testAccPostgresqlParameterTemplateConfigUpdate = `
resource "tencentcloud_postgresql_parameter_template" "example" {
  template_name        = "tf-example"
  db_major_version     = "18"
  db_engine            = "postgresql"
  template_description = "remark."
}

resource "tencentcloud_postgresql_parameter_template_config" "example" {
  template_id = tencentcloud_postgresql_parameter_template.example.id
  modify_param_entry_set {
    name           = "min_parallel_index_scan_size"
    expected_value = "64"
  }

  modify_param_entry_set {
    name           = "enable_async_append"
    expected_value = "off"
  }

  modify_param_entry_set {
    name           = "enable_group_by_reordering"
    expected_value = "off"
  }
}
`
