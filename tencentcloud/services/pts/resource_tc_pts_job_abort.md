Provides a resource to create a pts job_abort

Example Usage

```hcl
resource "tencentcloud_pts_job_abort" "job_abort" {
  job_id       = "job-my644ozi"
  project_id   = "project-45vw7v82"
  scenario_id  = "scenario-22q19f3k"
}
```