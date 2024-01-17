package cdwch_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseInstanceNodesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseInstanceNodesDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_clickhouse_instance_nodes.instance_nodes")),
			},
		},
	})
}

const testAccClickhouseInstanceNodesDataSource = `
data "tencentcloud_clickhouse_instance_nodes" "instance_nodes" {
  instance_id    = "cdwch-mvfjh373"
  node_role      = "data"
  display_policy = "all"
  force_all      = true
}
`
