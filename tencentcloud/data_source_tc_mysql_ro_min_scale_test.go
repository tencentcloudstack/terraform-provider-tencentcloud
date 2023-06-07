package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRoMinScaleDataSource_basic -v
func TestAccTencentCloudMysqlRoMinScaleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRoMinScaleDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_ro_min_scale.ro_min_scale"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_ro_min_scale.ro_min_scale", "memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_ro_min_scale.ro_min_scale", "volume"),
				),
			},
		},
	})
}

const testAccMysqlRoMinScaleDataSource = `

data "tencentcloud_mysql_ro_min_scale" "ro_min_scale" {
	master_instance_id = "cdb-fitq5t9h"
}

`
