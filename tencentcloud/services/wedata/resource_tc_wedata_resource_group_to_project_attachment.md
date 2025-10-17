Provides a resource to create a WeData resource group to project attachment

Example Usage

```hcl
resource "tencentcloud_wedata_resource_group_to_project_attachment" "example" {
  resource_group_id  = "20250909161820129828"
  project_id         = "2983848457986924544"
}
```

Import

WeData resource group to project attachment can be imported using the resourceGroupId#projectId, e.g.

```
terraform import tencentcloud_wedata_resource_group_to_project_attachment.example 20250909161820129828#2983848457986924544
```
