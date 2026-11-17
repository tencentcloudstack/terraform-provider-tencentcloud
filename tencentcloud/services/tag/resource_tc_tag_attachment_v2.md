Provides a resource to manage a single tag key/value binding to a single cloud resource (resource six-segment) for TencentCloud Tag. This is the v2 resource that supports updating `tag_value` in place. If you need the original behavior where changing `tag_value` forces a new binding (destroy and recreate), please use `tencentcloud_tag_attachment`.

~> **NOTE**: Compared with `tencentcloud_tag_attachment` (v1), the v2 resource supports updating `tag_value` **in place** via the `UpdateResourceTagValue` API, instead of destroying and recreating the binding. The ID format is `tag_key#resource` (without `tag_value`), while v1 uses `tag_key#tag_value#resource`. The `tag_key` and `resource` fields remain immutable (`ForceNew`) as in v1.

Example Usage

```hcl
resource "tencentcloud_tag_attachment_v2" "example" {
  tag_key   = "env"
  tag_value = "prod"
  resource  = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}
```

Import

tag attachment v2 can be imported using the composite id (tag_key joined with the resource six-segment by `#`), e.g.

```
terraform import tencentcloud_tag_attachment_v2.example env#qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4
```
