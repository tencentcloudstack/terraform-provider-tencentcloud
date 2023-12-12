Provides a resource to create a pts job

Example Usage

```hcl
resource "tencentcloud_pts_job" "job" {
  scenario_id = "scenario-22q19f3k"
  job_owner = "username"
  project_id = "project-45vw7v82"
  # debug = ""
  note = "desc"
}

```
Import

pts job can be imported using the projectId#scenarioId#jobId, e.g.
```
$ terraform import tencentcloud_pts_job.job project-45vw7v82#scenario-22q19f3k#job-dtm93vx0
```