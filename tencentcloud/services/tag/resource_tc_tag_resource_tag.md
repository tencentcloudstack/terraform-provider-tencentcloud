Provides a resource to manage a single tag key/value binding to a single cloud resource (resource six-segment) for TencentCloud Tag.

Example Usage

```hcl
resource "tencentcloud_tag_resource_tag" "example" {
  tag_key   = "env"
  tag_value = "prod"
  resource  = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}
```

Import

tag resource tag can be imported using the composite id (tag_key joined with the resource six-segment by `#`), e.g.

```
terraform import tencentcloud_tag_resource_tag.example env#qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4
```
