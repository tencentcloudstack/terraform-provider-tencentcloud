package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainDiagDbInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDiagDbInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "is_supported", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "product", "mysql"),
					resource.TestCheckTypeSetElemAttr("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "instance_names.*", "keep_preset_mysql"),
					// return
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.health_score"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.product"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.instance_conf.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.instance_conf.0.daily_inspection"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.instance_conf.0.overview_display"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_db_instances.diag_db_instances", "items.0.instance_conf.0.key_delimiters.#"),
				),
			},
		},
	})
}

const testAccDbbrainDiagDbInstancesDataSource = `

data "tencentcloud_dbbrain_diag_db_instances" "diag_db_instances" {
	is_supported   = true
	product        = "mysql"
	instance_names = ["keep_preset_mysql"]
}

`
