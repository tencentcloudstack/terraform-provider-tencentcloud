Use this data source to query reserved instances.

Example Usage

```hcl
data "tencentcloud_reserved_instances" "instances" {
  availability_zone = "na-siliconvalley-1"
  instance_type     = "S2.MEDIUM8"
}
```