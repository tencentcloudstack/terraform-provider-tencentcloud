Use this data source to query detailed instance information of Ckafka

Example Usage

Query all Ckafka instances

```hcl
data "tencentcloud_ckafka_instances" "example" {}
```

Query Ckafka instances by filters

```hcl
data "tencentcloud_ckafka_instances" "example" {
  instance_ids = [
    "ckafka-7k5nbnem",
    "ckafka-8j4raxv8"
  ]

  status = [0, 1, 2]

  filters {
    name   = "InstanceType"
    values = ["profession"]
  }
}
```
