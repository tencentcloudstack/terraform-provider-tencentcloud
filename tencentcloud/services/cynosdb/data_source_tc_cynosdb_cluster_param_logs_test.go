package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbClusterParamLogsDataSource_basic -v
func TestAccTencentCloudCynosdbClusterParamLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterParamLogsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_param_logs.cluster_param_logs"),
				),
			},
		},
	})
}

const testAccCynosdbClusterParamLogsDataSource = `
data "tencentcloud_cynosdb_cluster_param_logs" "cluster_param_logs" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  instance_ids  = ["cynosdbmysql-ins-afqx1hy0"]
  order_by      = "CreateTime"
  order_by_type = "DESC"
}
`
