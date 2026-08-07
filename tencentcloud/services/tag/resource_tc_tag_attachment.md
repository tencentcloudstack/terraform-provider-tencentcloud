Provides a resource to create a tag attachment

Example Usage

```hcl

resource "tencentcloud_tag_attachment" "attachment" {
  tag_key = "test3"
  tag_value = "Terraform3"
  resource = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}

```

`tag_value` can be updated in place; changing it calls the `UpdateResourceTagValue` API to modify the bound tag value in a single request instead of destroying and recreating the attachment. `tag_key` and `resource` are immutable and require the resource to be recreated when changed.

```hcl
resource "tencentcloud_tag_attachment" "attachment" {
  tag_key   = "test3"
  tag_value = "Terraform3_update"
  resource  = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}
```

Import

tag attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tag_attachment.attachment attachment_id
```