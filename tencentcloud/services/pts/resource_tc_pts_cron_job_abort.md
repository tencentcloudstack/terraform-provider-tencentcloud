Provides a resource to create a pts cron_job_abort

Example Usage

```hcl
resource "tencentcloud_pts_cron_job_abort" "cron_job_abort" {
  project_id  = "project-abc"
  cron_job_id = "job-dtm93vx0"
}
```