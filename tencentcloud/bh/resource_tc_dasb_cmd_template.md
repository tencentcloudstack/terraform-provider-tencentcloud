Provides a resource to create a dasb cmd_template

Example Usage

```hcl
resource "tencentcloud_dasb_cmd_template" "example" {
  name     = "tf_example"
  cmd_list = "rm -rf*"
}
```

Import

dasb cmd_template can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_cmd_template.example 15
```