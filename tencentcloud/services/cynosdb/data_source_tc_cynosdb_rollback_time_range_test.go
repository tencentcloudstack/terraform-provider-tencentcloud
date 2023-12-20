package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbRollbackTimeRangeDataSource_basic -v
func TestAccTencentCloudCynosdbRollbackTimeRangeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbRollbackTimeRangeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_rollback_time_range.rollback_time_range"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_rollback_time_range.rollback_time_range", "rollback_time_ranges.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_rollback_time_range.rollback_time_range", "rollback_time_ranges.0.time_range_start"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_rollback_time_range.rollback_time_range", "rollback_time_ranges.0.time_range_end"),
				),
			},
		},
	})
}

const testAccCynosdbRollbackTimeRangeDataSource = `
data "tencentcloud_cynosdb_rollback_time_range" "rollback_time_range" {
  cluster_id = "cynosdbmysql-bws8h88b"
}
`
