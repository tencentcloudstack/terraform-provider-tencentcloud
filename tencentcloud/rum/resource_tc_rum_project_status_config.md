Provides a resource to create a rum project_status_config

Example Usage

```hcl
resource "tencentcloud_rum_project_status_config" "project_status_config" {
  project_id = 131407
  operate    = "stop"
}
```

Import

rum project_status_config can be imported using the id, e.g.

```
terraform import tencentcloud_rum_project_status_config.project_status_config project_id
```