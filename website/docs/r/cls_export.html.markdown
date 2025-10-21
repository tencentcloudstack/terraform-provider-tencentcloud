---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_export"
sidebar_current: "docs-tencentcloud-resource-cls_export"
description: |-
  Provides a resource to create a cls export
---

# tencentcloud_cls_export

Provides a resource to create a cls export

## Example Usage

```hcl
resource "tencentcloud_cls_export" "export" {
  topic_id  = "7e34a3a7-635e-4da8-9005-88106c1fde69"
  log_count = 2
  query     = "select count(*) as count"
  from      = 1607499107000
  to        = 1607499108000
  order     = "desc"
  format    = "json"
}
```

## Argument Reference

The following arguments are supported:

* `from` - (Required, Int, ForceNew) export start time.
* `log_count` - (Required, Int, ForceNew) export amount of log.
* `query` - (Required, String, ForceNew) export query rules.
* `to` - (Required, Int, ForceNew) export end time.
* `topic_id` - (Required, String, ForceNew) topic id.
* `format` - (Optional, String, ForceNew) log export format.
* `order` - (Optional, String, ForceNew) log export time sorting. desc or asc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cls export can be imported using the id, e.g.

```
terraform import tencentcloud_cls_export.export topic_id#export_id
```

