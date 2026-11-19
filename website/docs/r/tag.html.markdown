---
subcategory: "Tag"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tag"
sidebar_current: "docs-tencentcloud-resource-tag"
description: |-
  Provides a resource to create a Tag
---

# tencentcloud_tag

Provides a resource to create a Tag

## Example Usage

```hcl
resource "tencentcloud_tag" "example" {
  tag_key   = "tagKey"
  tag_value = "tagValue"
}
```

## Argument Reference

The following arguments are supported:

* `tag_key` - (Required, String, ForceNew) Tag key.
* `tag_value` - (Required, String, ForceNew) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateTag](https://www.tencentcloud.com/document/api/651/33015)
* Read: [GetTags](https://www.tencentcloud.com/document/api/651/46536)
* Delete: [DeleteTag](https://www.tencentcloud.com/document/api/651/33013)

## Import

Tag can be imported using the tagKey#tagValue, e.g.

```
terraform import tencentcloud_tag.example tagKey#tagValue
```

