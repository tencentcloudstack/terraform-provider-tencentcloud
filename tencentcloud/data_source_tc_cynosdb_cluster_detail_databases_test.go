package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbClusterDetailDatabasesDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_detail_databases.cluster_detail_databases"),
				),
			},
		},
	})
}

const testAccCynosdbClusterDetailDatabasesDataSource = `
data "tencentcloud_cynosdb_cluster_detail_databases" "cluster_detail_databases" {
  cluster_id = "cynosdbmysql-bws8h88b"
  db_name    = "users"
}
`
