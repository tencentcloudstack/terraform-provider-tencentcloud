Provides a resource to create a lighthouse blueprint

Example Usage

```hcl
resource "tencentcloud_lighthouse_blueprint" "blueprint" {
  blueprint_name = "blueprint_name_test"
  description = "blueprint_description_test"
  instance_id = "lhins-xxxxxx"
}
```

Import

lighthouse blueprint can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_blueprint.blueprint blueprint_id
```