package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTdmqRocketmqGroupDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdmqRocketmqGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_group.group"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdmq_rocketmq_group.group", "groups.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceTdmqRocketmqGroup = `
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq_datasource_group"
	remark = "test recket mq"
  }
  
  resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
	cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
	namespace_name = "test_namespace_datasource"
	ttl = 65000
	retention_time = 65000
	remark = "test namespace"
  }
  
  resource "tencentcloud_tdmq_rocketmq_group" "group" {
	group_name = "test_rocketmq_group"
	namespace = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
	read_enable = true
	broadcast_enable = true
	cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
	remark = "test rocketmq group"
  }
  
  data "tencentcloud_tdmq_rocketmq_group" "group" {
	cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
	namespace_id = tencentcloud_tdmq_rocketmq_namespace.namespace.namespace_name
	filter_group = tencentcloud_tdmq_rocketmq_group.group.group_name
  }
`
