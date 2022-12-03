package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdmqRocketmqRoleDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdmqRocketmqRole,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_role.role"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdmq_rocketmq_role.role", "role_sets.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceTdmqRocketmqRole = `
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq_datasource_role"
	remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = "test_rocketmq_role"
  remark = "test rocketmq role"
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
}

data "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = tencentcloud_tdmq_rocketmq_role.role.role_name
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
}
`
