---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_data_transform"
sidebar_current: "docs-tencentcloud-resource-cls_data_transform"
description: |-
  Provides a resource to create a cls data_transform
---

# tencentcloud_cls_data_transform

Provides a resource to create a cls data_transform

## Example Usage

```hcl
resource "tencentcloud_cls_data_transform" "data_transform" {
  func_type    = 1
  src_topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  name         = "task"
  etl_content  = "xxx"
  task_type    = 1
  enable_flag  = 1
  dst_resources {
    topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
    alias    = "test"

  }
}
```

## Argument Reference

The following arguments are supported:

* `etl_content` - (Required, String) data transform content.
* `func_type` - (Required, Int) task type.
* `name` - (Required, String) task name.
* `src_topic_id` - (Required, String) src topic id.
* `task_type` - (Required, Int) task type.
* `dst_resources` - (Optional, List) data transform des resources.
* `enable_flag` - (Optional, Int) task enable flag.

The `dst_resources` object supports the following:

* `alias` - (Required, String) alias.
* `topic_id` - (Required, String) dst topic id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls data_transform can be imported using the id, e.g.

```
terraform import tencentcloud_cls_data_transform.data_transform data_transform_id
```

