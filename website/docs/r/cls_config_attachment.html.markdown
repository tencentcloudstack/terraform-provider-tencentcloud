---
subcategory: "CLS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_config_attachment"
sidebar_current: "docs-tencentcloud-resource-cls_config_attachment"
description: |-
  Provides a resource to create a cls config attachment
---

# tencentcloud_cls_config_attachment

Provides a resource to create a cls config attachment

## Example Usage

```hcl
resource "tencentcloud_cls_config_attachment" "attach" {
  config_id = tencentcloud_cls_config.config.id
  group_id  = "27752a9b-9918-440a-8ee7-9c84a14a47ed"
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required, String, ForceNew) Collection configuration id.
* `group_id` - (Required, String, ForceNew) Machine group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



