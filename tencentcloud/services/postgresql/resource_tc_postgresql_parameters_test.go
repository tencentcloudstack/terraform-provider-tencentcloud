package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlParametersResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlParameters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_parameters.postgresql_parameters", "id"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameters.postgresql_parameters", "param_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameters.postgresql_parameters", "param_list.0.expected_value", "on"),
				),
			},
			{
				ResourceName:      "tencentcloud_postgresql_parameters.postgresql_parameters",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPostgresqlParametersUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameters.postgresql_parameters", "param_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_postgresql_parameters.postgresql_parameters", "param_list.0.expected_value", "off"),
				),
			},
		},
	})
}

const testAccPostgresqlParameters = testAccPostgresqlInstancePostpaid + `

resource "tencentcloud_postgresql_parameters" "postgresql_parameters" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  param_list {
    expected_value = "on"
    name           = "check_function_bodies"
  }
}
`

const testAccPostgresqlParametersUp = testAccPostgresqlInstancePostpaid + `

resource "tencentcloud_postgresql_parameters" "postgresql_parameters" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  param_list {
    expected_value = "off"
    name           = "check_function_bodies"
  }
}
`
