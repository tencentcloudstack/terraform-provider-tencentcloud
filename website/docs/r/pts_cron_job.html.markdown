---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_cron_job"
sidebar_current: "docs-tencentcloud-resource-pts_cron_job"
description: |-
  Provides a resource to create a pts cron_job
---

# tencentcloud_pts_cron_job

Provides a resource to create a pts cron_job

## Example Usage

```hcl
resource "tencentcloud_pts_cron_job" "cron_job" {
  name            = "iac-cron_job-update"
  project_id      = "project-7qkzxhea"
  scenario_id     = "scenario-c22lqb1w"
  scenario_name   = "pts-js(2022-11-10 21:53:53)"
  frequency_type  = 1
  cron_expression = "4 22 10 11 ? 2022"
  job_owner       = "userName"
  # end_time = ""
  notice_id = "notice-vp6i38jt"
  note      = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `cron_expression` - (Required, String) Cron expression.
* `frequency_type` - (Required, Int) Execution frequency type, `1`: execute only once; `2`: daily granularity; `3`: weekly granularity; `4`: advanced.
* `job_owner` - (Required, String) Job Owner.
* `name` - (Required, String) Cron Job Name.
* `project_id` - (Required, String) Project Id.
* `scenario_id` - (Required, String) Scenario Id.
* `scenario_name` - (Required, String) Scenario Name.
* `end_time` - (Optional, String) End Time; type: Timestamp ISO8601.
* `note` - (Optional, String) Note.
* `notice_id` - (Optional, String) Notice ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `abort_reason` - Reason for suspension.
* `app_id` - App ID.
* `created_at` - Creation time; type: Timestamp ISO8601.
* `status` - Scheduled task status.
* `sub_account_uin` - Sub-user ID.
* `uin` - User ID.
* `updated_at` - Update time; type: Timestamp ISO8601.


## Import

pts cron_job can be imported using the projectId#cronJobId, e.g.
```
$ terraform import tencentcloud_pts_cron_job.cron_job project-7qkzxhea#scenario-c22lqb1w
```

