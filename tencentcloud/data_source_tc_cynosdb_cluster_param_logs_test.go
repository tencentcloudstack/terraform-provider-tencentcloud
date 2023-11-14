package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterParamLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterParamLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_param_logs.cluster_param_logs")),
			},
		},
	})
}

const testAccCynosdbClusterParamLogsDataSource = `

data "tencentcloud_cynosdb_cluster_param_logs" "cluster_param_logs" {
  cluster_id = "123"
  instance_ids = 
  order_by = "123"
  order_by_type = "DESC"
  }

`
