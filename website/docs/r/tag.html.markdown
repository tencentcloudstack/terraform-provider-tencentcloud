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

* Create: [CreateTag](https://cloud.tencent.com/document/api/651/35325) - 创建标签
* Read: [GetTags](https://cloud.tencent.com/document/api/651/72275) - 获取标签列表
* Delete: [DeleteTag](https://cloud.tencent.com/document/api/651/35323) - 删除标签

## Import

Tag can be imported using the tagKey#tagValue, e.g.

```
terraform import tencentcloud_tag.example tagKey#tagValue
```

