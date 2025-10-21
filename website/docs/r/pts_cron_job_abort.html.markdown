---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_cron_job_abort"
sidebar_current: "docs-tencentcloud-resource-pts_cron_job_abort"
description: |-
  Provides a resource to create a pts cron_job_abort
---

# tencentcloud_pts_cron_job_abort

Provides a resource to create a pts cron_job_abort

## Example Usage

```hcl
resource "tencentcloud_pts_cron_job_abort" "cron_job_abort" {
  project_id  = "project-abc"
  cron_job_id = "job-dtm93vx0"
}
```

## Argument Reference

The following arguments are supported:

* `cron_job_id` - (Required, String, ForceNew) Cron job ID.
* `project_id` - (Required, String, ForceNew) Project ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



