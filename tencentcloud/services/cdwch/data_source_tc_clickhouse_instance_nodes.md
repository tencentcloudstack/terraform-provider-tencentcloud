Use this data source to query detailed information of clickhouse instance_nodes

Example Usage

```hcl
data "tencentcloud_clickhouse_instance_nodes" "instance_nodes" {
  instance_id    = "cdwch-mvfjh373"
  node_role      = "data"
  display_policy = "all"
  force_all      = true
}
```