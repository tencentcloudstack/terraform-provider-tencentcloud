Provides a resource to create a cvm launch_template_default_version

Example Usage

```hcl
resource "tencentcloud_cvm_launch_template_default_version" "launch_template_default_version" {
  launch_template_id = "lt-34vaef8fe"
  default_version = 2
}
```

Import

cvm launch_template_default_version can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_launch_template_default_version.launch_template_default_version launch_template_id
```