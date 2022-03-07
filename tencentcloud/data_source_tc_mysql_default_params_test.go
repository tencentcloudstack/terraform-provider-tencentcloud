package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceMysqlDefaultParams(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMysqlDefaultParamBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_default_params.mysql_57", "db_version", "5.7"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.current_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.max"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.min"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.need_reboot"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_default_params.mysql_57", "param_list.0.param_type"),
				),
			},
		},
	})
}

const testAccDataSourceMysqlDefaultParamBasic = `
data "tencentcloud_mysql_default_params" "mysql_57" {
	db_version = "5.7"
}
`
