package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlLogBackupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlLogBackupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_log_backups.log_backups")),
			},
		},
	})
}

const testAccPostgresqlLogBackupsDataSource = `

data "tencentcloud_postgresql_log_backups" "log_backups" {
  min_finish_time = ""
  max_finish_time = ""
  filters {
		name = ""
		values = 

  }
  order_by = ""
  order_by_type = ""

}

`
