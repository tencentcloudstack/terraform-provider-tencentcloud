Use this data source to query detailed information of Mongodb instances.

Example Usage

```hcl
data "tencentcloud_mongodb_instances" "mongodb" {
  instance_id  = "cmgo-l6lwdsel"
  cluster_type = "REPLSET"
}
```