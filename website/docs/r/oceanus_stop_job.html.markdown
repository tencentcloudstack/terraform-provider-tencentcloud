---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_stop_job"
sidebar_current: "docs-tencentcloud-resource-oceanus_stop_job"
description: |-
  Provides a resource to create a oceanus stop_job
---

# tencentcloud_oceanus_stop_job

Provides a resource to create a oceanus stop_job

## Example Usage

```hcl
resource "tencentcloud_oceanus_stop_job" "example" {
  stop_job_descriptions {
    job_id    = "cql-4xwincyn"
    stop_type = 1
  }
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `stop_job_descriptions` - (Required, List, ForceNew) The description information for batch job stop.
* `work_space_id` - (Optional, String, ForceNew) Workspace SerialId.

The `stop_job_descriptions` object supports the following:

* `job_id` - (Required, String) Job Id.
* `stop_type` - (Required, Int) Stop type,1 stopped 2 paused.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



