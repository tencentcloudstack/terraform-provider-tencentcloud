Provides a resource to create a Makers project.

Example Usage

```hcl
resource "tencentcloud_makers_project" "example" {
  zone_id           = "zone-3hzk0ndvnnyp"
  name              = "tf-test5"
  channel           = "Custom"
  area              = "global"
  code_provider     = "Upload"
  node_js_version   = "22.11.0"
}
```

Import

Makers project can be imported using the composite id `zone_id#project_id`, e.g.

```
terraform import tencentcloud_makers_project.example zone-12345#proj-67890
```
