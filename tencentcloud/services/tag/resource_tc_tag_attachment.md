Provides a resource to create a tag attachment

Example Usage

```hcl

resource "tencentcloud_tag_attachment" "attachment" {
  tag_key = "test3"
  tag_value = "Terraform3"
  resource = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}

```

The `tag_value` field can be updated in place; changing it calls the `UpdateResourceTagValue` cloud API to modify the tag value without destroying and recreating the attachment.

Import

tag attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tag_attachment.attachment attachment_id
```