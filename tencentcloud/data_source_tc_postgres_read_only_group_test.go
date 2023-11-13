package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresReadOnlyGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresReadOnlyGroupDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgres_read_only_group.read_only_group")),
			},
		},
	})
}

const testAccPostgresReadOnlyGroupDataSource = `

data "tencentcloud_postgres_read_only_group" "read_only_group" {
  filters {
		name = "db-master-instance-id"
		values = 

  }
  page_size = 10
  page_number = 0
  order_by = "CreateTime"
  order_by_type = "asc"
}

`
