Use this data source to query detailed information of WeData workflow max permission

Example Usage

```hcl
data "tencentcloud_wedata_workflow_max_permission" "example" {
  project_id  = "3108707295180644352"
  entity_id   = "53e78f97-f145-11f0-ba36-b8cef6a5af5c"
  entity_type = "folder"
}
```
