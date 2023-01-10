package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCynosdbClusterParamsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterParamsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_params.cluster_params"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_cluster_params.cluster_params", "items.#", "1"),
				),
			},
		},
	})
}

const testAccCynosdbClusterParamsDataSource = CommonCynosdb + `
data "tencentcloud_cynosdb_cluster_params" "cluster_params" {
	cluster_id = var.cynosdb_cluster_id
	param_name = "innodb_checksum_algorithm"
}
`
