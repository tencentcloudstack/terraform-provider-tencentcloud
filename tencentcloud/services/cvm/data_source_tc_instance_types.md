Use this data source to query instances type.

Example Usage

```hcl
data "tencentcloud_instance_types" "foo" {
  availability_zone = "ap-guangzhou-2"
  cpu_core_count    = 2
  memory_size       = 4
}

data tencentcloud_instance_types "t1c1g" {
  cpu_core_count    = 1
  memory_size       = 1
  exclude_sold_out=true
  filter {
    name   = "instance-charge-type"
    values = ["POSTPAID_BY_HOUR"]
  }
  filter {
    name   = "zone"
    values = ["ap-shanghai-2"]
  }
}
```