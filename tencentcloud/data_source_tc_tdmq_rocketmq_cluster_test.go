package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdmqRocketmqClusterDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRocketmqCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_cluster.cluster"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdmq_rocketmq_cluster.cluster", "cluster_list.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceRocketmqCluster = `
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq_datasource"
	remark = "test recket mq"
}

data "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
  	name_keyword = "test_rocketmq_datasource"
}
`
