Use this data source to query detailed information of pts scenario_with_jobs

Example Usage

```hcl
data "tencentcloud_pts_scenario_with_jobs" "scenario_with_jobs" {
  project_ids    = ["project-45vw7v82"]
  scenario_ids   = ["scenario-koakp3h6"]
  scenario_name  = "pts-jmeter"
  ascend         = true
  ignore_script  = true
  ignore_dataset = true
  scenario_type  = "pts-jmeter"
}
```