package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbClusterDataSource_basic -v
func TestAccTencentCloudCynosdbClusterDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster.cluster"),
				),
			},
		},
	})
}

const testAccCynosdbClusterDataSource = `
data "tencentcloud_cynosdb_cluster" "cluster" {
  cluster_id = "cynosdbmysql-bws8h88b"
  database   = "users"
  table      = "tb_user_name"
  table_type = "all"
}
`
