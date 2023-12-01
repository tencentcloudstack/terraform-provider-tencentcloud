package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainMysqlProcessListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainMysqlProcessListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "product", "mysql"),
					// return
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.user"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.host"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.db"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.command"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list", "process_list.0.info"),
				),
			},
		},
	})
}

const testAccDbbrainMysqlProcessListDataSource = CommonPresetMysql + `

data "tencentcloud_dbbrain_mysql_process_list" "mysql_process_list" {
  instance_id = local.mysql_id
  product     = "mysql"
}

`
