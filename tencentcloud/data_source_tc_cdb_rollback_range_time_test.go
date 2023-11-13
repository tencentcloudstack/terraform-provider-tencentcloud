package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRollbackRangeTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRollbackRangeTimeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_rollback_range_time.rollback_range_time")),
			},
		},
	})
}

const testAccCdbRollbackRangeTimeDataSource = `

data "tencentcloud_cdb_rollback_range_time" "rollback_range_time" {
  instance_ids = 
  is_remote_zone = ""
  backup_region = ""
  }

`
