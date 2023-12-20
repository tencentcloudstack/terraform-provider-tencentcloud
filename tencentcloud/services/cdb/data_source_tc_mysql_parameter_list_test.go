package cdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlParameterListDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
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
				Config: testAccMysqlParameterListDataSourceConfig(tcacctest.CommonPresetMysql),
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
	return `
data "tencentcloud_mysql_parameter_list" "mysql_default" {
	engine_version = "5.7"
}
	`
}

func testAccMysqlParameterListDataSourceConfig(commonTestCase string) string {
	return fmt.Sprintf(`
%s
data "tencentcloud_mysql_parameter_list" "mysql" {
	mysql_id = local.mysql_id
}
	`, commonTestCase)
}
