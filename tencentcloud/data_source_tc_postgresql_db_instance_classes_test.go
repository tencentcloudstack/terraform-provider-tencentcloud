package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresqlDbInstanceClassesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDbInstanceClassesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes")),
			},
		},
	})
}

const testAccPostgresqlDbInstanceClassesDataSource = `

data "tencentcloud_postgresql_db_instance_classes" "db_instance_classes" {
  zone = ""
  db_engine = ""
  db_major_version = ""
  }

`
