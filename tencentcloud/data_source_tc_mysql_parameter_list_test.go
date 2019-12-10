package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMysqlParameterListDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlParameterListDataSourceDefaultConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql_default", "parameter_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql_default", "parameter_list.0.parameter_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql_default", "parameter_list.0.parameter_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql_default", "parameter_list.0.default_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql_default", "parameter_list.0.need_reboot"),
				),
			},
			{
				Config: testAccMysqlParameterListDataSourceConfig(mysqlInstanceCommonTestCase),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql", "parameter_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql", "parameter_list.0.parameter_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql", "parameter_list.0.parameter_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql", "parameter_list.0.default_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_parameter_list.mysql", "parameter_list.0.need_reboot"),
				),
			},
		},
	})
}

func testAccMysqlParameterListDataSourceDefaultConfig() string {
	return fmt.Sprintf(`
data "tencentcloud_mysql_parameter_list" "mysql_default" {
	engine_version = "5.7"
}
	`)
}

func testAccMysqlParameterListDataSourceConfig(commonTestCase string) string {
	return fmt.Sprintf(`
%s
data "tencentcloud_mysql_parameter_list" "mysql" {
	mysql_id = "${tencentcloud_mysql_instance.default.id}"
}
	`, commonTestCase)
}
