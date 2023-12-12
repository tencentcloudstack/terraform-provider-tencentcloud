Provides a resource to create a pts cron_job_restart, restart a scheduled task whose status is `JobAborting`

Example Usage

```hcl
resource "tencentcloud_pts_cron_job_restart" "cron_job_restart" {
  project_id  = "project-abc"
  cron_job_id = "job-dtm93vx0"
}
```