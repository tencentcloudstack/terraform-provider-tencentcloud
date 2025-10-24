Use this data source to query detailed information of WeData list process lineage

Example Usage

```hcl
data "tencentcloud_wedata_list_process_lineage" "example" {
  process_id   = "20241107221758402"
  process_type = "SCHEDULE_TASK"
  platform     = "WEDATA"
}
```
