package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresqlDbInstanceVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDbInstanceVersionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions")),
			},
		},
	})
}

const testAccPostgresqlDbInstanceVersionsDataSource = `

data "tencentcloud_postgresql_db_instance_versions" "db_instance_versions" {
  }

`
