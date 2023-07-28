---
subcategory: "Tag"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_resource_tag"
sidebar_current: "docs-tencentcloud-resource-resource_tag"
description: |-
  Provides a resource to create a tag resource_tag
---

# tencentcloud_resource_tag

Provides a resource to create a tag resource_tag

## Example Usage

```hcl
resource "tencentcloud_resource_tag" "resource_tag" {
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

tag resource_tag can be imported using the id, e.g.

```
terraform import tencentcloud_resource_tag.resource_tag resource_tag_id
```

