Provides a resource to create a cvm launch_template_version

Example Usage

```hcl
resource "tencentcloud_cvm_launch_template_version" "foo" {
  placement {
		zone = "ap-guangzhou-6"
		project_id = 0

  }
  launch_template_id = "lt-r9ajalbi"
  launch_template_version_description = "version description"
  disable_api_termination = false
  instance_type = "S5.MEDIUM4"
  image_id = "img-9qrfy1xt"
}
```

Import

cvm launch_template_version can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_launch_template_version.launch_template_version ${launch_template_id}#${launch_template_version}
```