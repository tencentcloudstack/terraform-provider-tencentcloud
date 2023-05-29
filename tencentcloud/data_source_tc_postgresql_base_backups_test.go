package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresqlBaseBackupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBaseBackupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_base_backups.base_backups")),
			},
		},
	})
}

const testAccPostgresqlBaseBackupsDataSource = `

data "tencentcloud_postgresql_base_backups" "base_backups" {
  min_finish_time = ""
  max_finish_time = ""
  filters {
		name = ""
		values = 

  }
  order_by = ""
  order_by_type = ""
    tags = {
    "createdBy" = "terraform"
  }
}

`
