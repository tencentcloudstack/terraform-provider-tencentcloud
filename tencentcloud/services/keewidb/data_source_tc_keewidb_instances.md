Use this data source to query KeeWiDB instances.

Example Usage

Query all instances

```hcl
data "tencentcloud_keewidb_instances" "example" {}
```

Query instances by filter

```hcl
data "tencentcloud_keewidb_instances" "example" {
  instance_id     = "kee-4nmzc0ul"
  instance_name   = "tf-example"
  uniq_vpc_ids    = ["vpc-mjwornzj"]
  uniq_subnet_ids = ["subnet-1ed4w7to"]
  billing_mode    = "postpaid"
}
```
