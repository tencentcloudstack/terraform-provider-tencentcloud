Use this data source to query detailed information of tdmqRocketmq role

Example Usage

```hcl
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
```