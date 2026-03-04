Use this data source to query detailed instance information of Ckafka instances

Example Usage

Query all Ckafka instances

```hcl
data "tencentcloud_ckafka_instances_v2" "example" {}
```

Query Ckafka instances by filters

```hcl
data "tencentcloud_ckafka_instances_v2" "example" {
  filters {
    name   = "InstanceId"
    values = ["ckafka-7k5nbnem"]
  }
}
```
