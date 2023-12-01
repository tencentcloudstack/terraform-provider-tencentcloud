---
subcategory: "Tag"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tag"
sidebar_current: "docs-tencentcloud-resource-tag"
description: |-
  Provides a resource to create a tag
---

# tencentcloud_tag

Provides a resource to create a tag

## Example Usage

```hcl
resource "tencentcloud_tag" "tag" {
  tag_key   = "test"
  tag_value = "Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `tag_key` - (Required, String, ForceNew) tag key.
* `tag_value` - (Required, String, ForceNew) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tag tag can be imported using the id, e.g.

```
terraform import tencentcloud_tag.tag tag_id
```

