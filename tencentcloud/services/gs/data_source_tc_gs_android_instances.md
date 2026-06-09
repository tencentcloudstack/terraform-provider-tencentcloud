Use this data source to query GS Android instances.

Example Usage

Query all GS Android instances

```hcl
data "tencentcloud_gs_android_instances" "example" {}
```

Query GS Android instances by filter

```hcl
data "tencentcloud_gs_android_instances" "example" {
  android_instance_ids = [
    "cai-1308726196-0352wk8np9s"
  ]
  android_instance_region = "ap-beijing"
  android_instance_zone   = "ap-beijing-1"
  label_selector {
    key      = "key"
    operator = "IN"
    values   = ["value"]
  }

  filters {
    name   = "Name"
    values = ["tf-example"]
  }
}
```
