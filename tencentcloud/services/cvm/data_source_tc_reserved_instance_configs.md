Use this data source to query reserved instances configuration.

Example Usage

```hcl
data "tencentcloud_reserved_instance_configs" "config" {
  availability_zone = "na-siliconvalley-1"
}
```