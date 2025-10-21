---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_config_attachment"
sidebar_current: "docs-tencentcloud-resource-cls_config_attachment"
description: |-
  Provides a resource to create a cls config attachment
---

# tencentcloud_cls_config_attachment

Provides a resource to create a cls config attachment

## Example Usage



## Argument Reference

The following arguments are supported:

* `config_id` - (Required, String, ForceNew) Collection configuration id.
* `group_id` - (Required, String, ForceNew) Machine group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls config_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cls_config_attachment.attach config_id#group_id
```

```

