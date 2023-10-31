---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_trigger_job_savepoint"
sidebar_current: "docs-tencentcloud-resource-oceanus_trigger_job_savepoint"
description: |-
  Provides a resource to create a oceanus trigger_job_savepoint
---

# tencentcloud_oceanus_trigger_job_savepoint

Provides a resource to create a oceanus trigger_job_savepoint

## Example Usage

```hcl
resource "tencentcloud_oceanus_trigger_job_savepoint" "example" {
  job_id        = "cql-4xwincyn"
  description   = "description."
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Job SerialId.
* `description` - (Optional, String, ForceNew) Savepoint description.
* `work_space_id` - (Optional, String, ForceNew) Workspace SerialId.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



