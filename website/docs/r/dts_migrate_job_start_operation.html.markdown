---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_job_start_operation"
sidebar_current: "docs-tencentcloud-resource-dts_migrate_job_start_operation"
description: |-
  Provides a resource to start a dts migrate_job
---

# tencentcloud_dts_migrate_job_start_operation

Provides a resource to start a dts migrate_job

## Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_start_operation" "start" {
  job_id = tencentcloud_dts_migrate_job.job.id
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Job Id from `tencentcloud_dts_migrate_job`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



