---
subcategory: "Tag"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tag_resource_tag"
sidebar_current: "docs-tencentcloud-resource-tag_resource_tag"
description: |-
  Provides a resource to manage a single tag key/value binding to a single cloud resource (resource six-segment) for TencentCloud Tag.
---

# tencentcloud_tag_resource_tag

Provides a resource to manage a single tag key/value binding to a single cloud resource (resource six-segment) for TencentCloud Tag.

## Example Usage

```hcl
resource "tencentcloud_tag_resource_tag" "example" {
  tag_key   = "env"
  tag_value = "prod"
  resource  = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}
```

## Argument Reference

The following arguments are supported:

* `resource` - (Required, String, ForceNew) [Six-segment description of resources](https://cloud.tencent.com/document/product/598/10606).
* `tag_key` - (Required, String, ForceNew) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tag resource tag can be imported using the composite id (tag_key joined with the resource six-segment by `#`), e.g.

```
terraform import tencentcloud_tag_resource_tag.example env#qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4
```

