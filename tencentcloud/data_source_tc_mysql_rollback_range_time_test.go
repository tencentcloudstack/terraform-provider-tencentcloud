package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRollbackRangeTimeDataSource_basic -v
func TestAccTencentCloudMysqlRollbackRangeTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRollbackRangeTimeDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_rollback_range_time.rollback_range_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_rollback_range_time.rollback_range_time", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_rollback_range_time.rollback_range_time", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_rollback_range_time.rollback_range_time", "items.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_rollback_range_time.rollback_range_time", "items.0.times.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_rollback_range_time.rollback_range_time", "items.0.times.0.begin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_rollback_range_time.rollback_range_time", "items.0.times.0.end"),
				),
			},
		},
	})
}

const testAccMysqlRollbackRangeTimeDataSourceVar = `
variable "instance_id" {
	default = "` + defaultDbBrainInstanceId + `"
}
`
const testAccMysqlRollbackRangeTimeDataSource = testAccMysqlRollbackRangeTimeDataSourceVar + `

data "tencentcloud_mysql_rollback_range_time" "rollback_range_time" {
	instance_ids = [var.instance_id]
}

`
