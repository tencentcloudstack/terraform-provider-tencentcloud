package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterDetailDatabasesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterDetailDatabasesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_detail_databases.cluster_detail_databases")),
			},
		},
	})
}

const testAccCynosdbClusterDetailDatabasesDataSource = `

data "tencentcloud_cynosdb_cluster_detail_databases" "cluster_detail_databases" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  db_name = "test"
  }

`
