Use this data source to query detailed information of dcdb shards

Example Usage

```hcl
data "tencentcloud_dcdb_shards" "shards" {
  instance_id = "your_instance_id"
  shard_instance_ids = ["shard1_id"]
  }
```