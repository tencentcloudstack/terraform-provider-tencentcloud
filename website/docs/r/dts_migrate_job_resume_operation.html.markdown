---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_job_resume_operation"
sidebar_current: "docs-tencentcloud-resource-dts_migrate_job_resume_operation"
description: |-
  Provides a resource to create a dts migrate_job_resume_operation
---

# tencentcloud_dts_migrate_job_resume_operation

Provides a resource to create a dts migrate_job_resume_operation

## Example Usage

```hcl
resource "tencentcloud_dts_migrate_job_resume_operation" "resume" {
  job_id        = "job_id"
  resume_option = "normal"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) job id.
* `resume_option` - (Required, String, ForceNew) resume mode: 1.clearData-Clear target data; 2.overwrite-The task is executed in overwrite mode; 3.normal-No extra action. Note that clearData and overwrite are valid only for redis links, normal is valid only for non-Redis links.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



