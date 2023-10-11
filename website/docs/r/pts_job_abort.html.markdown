---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_job_abort"
sidebar_current: "docs-tencentcloud-resource-pts_job_abort"
description: |-
  Provides a resource to create a pts job_abort
---

# tencentcloud_pts_job_abort

Provides a resource to create a pts job_abort

## Example Usage

```hcl
resource "tencentcloud_pts_job_abort" "job_abort" {
  job_id      = "job-my644ozi"
  project_id  = "project-45vw7v82"
  scenario_id = "scenario-22q19f3k"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Job ID.
* `project_id` - (Required, String, ForceNew) Project ID.
* `scenario_id` - (Required, String, ForceNew) Scenario ID.
* `abort_reason` - (Optional, Int, ForceNew) The reason for aborting the job.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



