Provides a resource to create a dasb device_group

Example Usage

```hcl
resource "tencentcloud_dasb_device_group" "example" {
  name          = "tf_example"
  department_id = "1.2"
}
```

Import

dasb device_group can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_group.example 36
```