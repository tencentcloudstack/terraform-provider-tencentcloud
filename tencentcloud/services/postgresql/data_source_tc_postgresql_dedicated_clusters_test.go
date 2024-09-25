package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudPostgresqlDedicatedClustersDataSource_basic -v
func TestAccTencentCloudPostgresqlDedicatedClustersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDedicatedClustersDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_dedicated_clusters.example"),
				),
			},
		},
	})
}

const testAccPostgresqlDedicatedClustersDataSource = `
data "tencentcloud_postgresql_dedicated_clusters" "example" {
  filters {
    name = "dedicated-cluster-id"
    values = ["cluster-262n63e8"]
  }
}
`
