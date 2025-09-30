Provides a resource to create a wedata wedata_workflow

Example Usage

```hcl
resource "tencentcloud_wedata_workflow" "wedata_workflow" {
  project_id = 2905622749543821312
  workflow_name = "test"
  parent_folder_path = "/tfmika"
  workflow_type = "cycle"
}
```

Import

wedata wedata_workflow can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_workflow.wedata_workflow wedata_workflow_id
```
