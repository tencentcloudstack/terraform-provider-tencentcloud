---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_job"
sidebar_current: "docs-tencentcloud-resource-pts_job"
description: |-
  Provides a resource to create a pts job
---

# tencentcloud_pts_job

Provides a resource to create a pts job

## Example Usage

```hcl
resource "tencentcloud_pts_job" "job" {
  scenario_id = "scenario-22q19f3k"
  job_owner   = "username"
  project_id  = "project-45vw7v82"
  # debug = ""
  note = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `job_owner` - (Required, String) Job owner.
* `project_id` - (Required, String) Project ID.
* `scenario_id` - (Required, String) Pts scenario id.
* `debug` - (Optional, Bool) Whether to debug.
* `note` - (Optional, String) Note.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `abort_reason` - Cause of interruption.
* `created_at` - Creation time of the job.
* `cron_id` - Scheduled job ID.
* `datasets` - Dataset file for the job.
* `domain_name_config` - Domain name binding configuration.
* `duration` - Job duration.
* `end_time` - End time of the job.
* `error_rate` - Percentage of error rate.
* `load` - Pressure configuration of job.
* `max_requests_per_second` - Maximum requests per second.
* `max_virtual_user_count` - Maximum number of VU for the job.
* `plugins` - Expansion package file information.
* `protocols` - Protocol script information.
* `request_files` - Request file information.
* `request_total` - Total number of requests.
* `requests_per_second` - Average number of requests per second.
* `response_time_average` - Average response time.
* `response_time_max` - Maximum response time.
* `response_time_min` - Minimum response time.
* `response_time_p90` - 90th percentile response time.
* `response_time_p95` - 95th percentile response time.
* `response_time_p99` - 99th percentile response time.
* `start_time` - Start time of the job.
* `status` - The running status of the task; `0`: JobUnknown, `1`: JobCreated, `2`: JobPending, `3`: JobPreparing, `4`: JobSelectClustering, `5`: JobCreateTasking, `6`: JobSyncTasking, `11`: JobRunning, `12`: JobFinished, `13`: JobPrepareException, `14`: JobFinishException, `15`: JobAborting, `16`: JobAborted, `17`: JobAbortException, `18`: JobDeleted, `19`: JobSelectClusterException, `20`: JobCreateTaskException, `21`: JobSyncTaskException.
* `test_scripts` - Test script information.
* `type` - Scene Type.


## Import

pts job can be imported using the projectId#scenarioId#jobId, e.g.
```
$ terraform import tencentcloud_pts_job.job project-45vw7v82#scenario-22q19f3k#job-dtm93vx0
```

