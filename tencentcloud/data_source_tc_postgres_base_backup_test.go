package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresBaseBackupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresBaseBackupDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgres_base_backup.base_backup")),
			},
		},
	})
}

const testAccPostgresBaseBackupDataSource = `

data "tencentcloud_postgres_base_backup" "base_backup" {
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
