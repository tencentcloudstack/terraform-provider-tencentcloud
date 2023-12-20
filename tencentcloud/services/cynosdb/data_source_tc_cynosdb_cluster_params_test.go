package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbClusterParamsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterParamsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_params.cluster_params"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_cluster_params.cluster_params", "items.#", "1"),
				),
			},
		},
	})
}

const testAccCynosdbClusterParamsDataSource = tcacctest.CommonCynosdb + `
data "tencentcloud_cynosdb_cluster_params" "cluster_params" {
	cluster_id = var.cynosdb_cluster_id
	param_name = "innodb_checksum_algorithm"
}
`
