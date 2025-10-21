---
subcategory: "Tag"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tag_attachment"
sidebar_current: "docs-tencentcloud-resource-tag_attachment"
description: |-
  Provides a resource to create a tag attachment
---

# tencentcloud_tag_attachment

Provides a resource to create a tag attachment

## Example Usage

```hcl
resource "tencentcloud_tag_attachment" "attachment" {
  tag_key   = "test3"
  tag_value = "Terraform3"
  resource  = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}
```

## Argument Reference

The following arguments are supported:

* `resource` - (Required, String, ForceNew) [Six-segment description of resources](https://cloud.tencent.com/document/product/598/10606).
* `tag_key` - (Required, String, ForceNew) tag key.
* `tag_value` - (Required, String, ForceNew) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tag attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tag_attachment.attachment attachment_id
```

