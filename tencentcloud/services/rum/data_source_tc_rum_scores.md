Use this data source to query detailed information of rum scores

Example Usage

```hcl
data "tencentcloud_rum_scores" "scores" {
  end_time   = "2023082215"
  start_time = "2023082214"
  project_id = 1
  is_demo    = 1
}
```