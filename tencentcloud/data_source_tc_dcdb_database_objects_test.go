package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbDatabaseObjectsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbDatabaseObjectsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_database_objects.database_objects")),
			},
		},
	})
}

const testAccDcdbDatabaseObjectsDataSource = `

data "tencentcloud_dcdb_database_objects" "database_objects" {
  instance_id = "dcdbt-ow7t8lmc"
  db_name = &lt;nil&gt;
        }

`
