Provides a resource to create a dasb user_group

Example Usage

```hcl
resource "tencentcloud_dasb_user_group" "example" {
  name = "tf_example_update"
}
```

Or

```hcl
resource "tencentcloud_dasb_user_group" "example" {
  name          = "tf_example_update"
  department_id = "1.2"
}
```

Import

dasb user_group can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user_group.example 16
```