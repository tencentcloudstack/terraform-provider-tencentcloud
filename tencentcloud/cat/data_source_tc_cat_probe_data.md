Use this data source to query detailed information of cat probe data

Example Usage

```hcl
data "tencentcloud_cat_probe_data" "probe_data" {
  begin_time = 1667923200000
  end_time = 1667996208428
  task_type = "AnalyzeTaskType_Network"
  sort_field = "ProbeTime"
  ascending = true
  selected_fields = ["terraform"]
  offset = 0
  limit = 20
  task_id = ["task-knare1mk"]
}
```