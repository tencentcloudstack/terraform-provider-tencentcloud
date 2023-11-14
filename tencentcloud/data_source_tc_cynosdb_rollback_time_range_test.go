package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbRollbackTimeRangeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbRollbackTimeRangeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_rollback_time_range.rollback_time_range")),
			},
		},
	})
}

const testAccCynosdbRollbackTimeRangeDataSource = `

data "tencentcloud_cynosdb_rollback_time_range" "rollback_time_range" {
  cluster_id = "cynosdbmysql-oib3wx0i"
      }

`
