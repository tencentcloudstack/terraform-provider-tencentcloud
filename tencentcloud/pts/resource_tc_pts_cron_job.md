Provides a resource to create a pts cron_job

Example Usage

```hcl
resource "tencentcloud_pts_cron_job" "cron_job" {
  name = "iac-cron_job-update"
  project_id = "project-7qkzxhea"
  scenario_id = "scenario-c22lqb1w"
  scenario_name = "pts-js(2022-11-10 21:53:53)"
  frequency_type = 2
  cron_expression = "* 1 * * *"
  job_owner = "userName"
  # end_time = ""
  notice_id = "notice-vp6i38jt"
  note = "desc"
}

```
Import

pts cron_job can be imported using the projectId#cronJobId, e.g.
```
$ terraform import tencentcloud_pts_cron_job.cron_job project-7qkzxhea#scenario-c22lqb1w
```