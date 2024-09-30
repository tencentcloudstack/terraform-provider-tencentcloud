package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlApplyParameterTemplateOperationResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPostgresqlApplyParameterTemplateOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_apply_parameter_template_operation.apply_parameter_template_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_apply_parameter_template_operation.apply_parameter_template_operation", "db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_apply_parameter_template_operation.apply_parameter_template_operation", "template_id"),
				),
			},
		},
	})
}

const testPostgresqlApplyParameterTemplateOperation = testAccPostgresqlInstance + `
resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "tf_test_apply"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "test"

  modify_param_entry_set {
	name = "vacuum_freeze_table_age"
	expected_value = "160000000"
  }
}

resource tencentcloud_postgresql_apply_parameter_template_operation "apply_parameter_template_operation" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  template_id = tencentcloud_postgresql_parameter_template.parameter_template.id
}
`
