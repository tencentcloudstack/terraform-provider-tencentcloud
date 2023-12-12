Provides a resource to create a oceanus work_space

Example Usage

```hcl
resource "tencentcloud_oceanus_work_space" "example" {
  work_space_name = "tf_example"
  description     = "example description."
}
```

Import

oceanus work_space can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_work_space.example space-0dan3yux#tf_example
```