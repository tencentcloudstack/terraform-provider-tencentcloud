Use this data source to query detailed information of tdmqRocketmq cluster

Example Usage

```hcl
data "tencentcloud_tdmq_rocketmq_cluster" "example" {
  name_keyword = tencentcloud_tdmq_rocketmq_cluster.example.cluster_name
}

resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}
```