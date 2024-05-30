Use this data source to query instances type.

Example Usage

```hcl
data "tencentcloud_instance_types" "example" {
  availability_zone = "ap-guangzhou-6"
  cpu_core_count    = 4
  memory_size       = 8
}
```

Complete Example

```hcl
data "tencentcloud_instance_types" "example" {
  cpu_core_count   = 4
  memory_size      = 8
  exclude_sold_out = true

  filter {
    name   = "instance-family"
    values = ["SA2"]
  }

  filter {
    name   = "zone"
    values = ["ap-guangzhou-6"]
  }
}
```