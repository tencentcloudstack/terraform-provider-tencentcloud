---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_upstream_task_instances"
sidebar_current: "docs-tencentcloud-datasource-wedata_upstream_task_instances"
description: |-
  Use this data source to query detailed information of wedata upstream task instances
---

# tencentcloud_wedata_upstream_task_instances

Use this data source to query detailed information of wedata upstream task instances

## Example Usage

```hcl
data "tencentcloud_wedata_task_instances" "wedata_task_instances" {
  project_id = "1859317240494305280"
}

locals {
  instance_keys = data.tencentcloud_wedata_task_instances.wedata_task_instances.data[0].items[*].instance_key
}

data "tencentcloud_wedata_upstream_task_instances" "wedata_upstream_task_instances" {
  for_each = toset(local.instance_keys)

  project_id   = "1859317240494305280"
  instance_key = each.value
}
```

## Argument Reference

The following arguments are supported:

* `instance_key` - (Required, String) Unique instance identifier.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `time_zone` - (Optional, String) Time zone, default UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Upstream instance list.


