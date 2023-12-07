Use this data source to query detailed information of dcdb price

Example Usage

```hcl
data "tencentcloud_dcdb_price" "price" {
	instance_count   = 1
	zone             = var.default_az
	period           = 1
	shard_node_count = 2
	shard_memory     = 2
	shard_storage    = 10
	shard_count      = 2
	paymode          = "postpaid"
	amount_unit      = "pent"
}
```