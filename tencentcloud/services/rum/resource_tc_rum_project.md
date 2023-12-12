Provides a resource to create a rum project

Example Usage

```hcl
resource "tencentcloud_rum_project" "project" {
  name = "projectName"
  instance_id = "rum-pasZKEI3RLgakj"
  rate = "100"
  enable_url_group = "0"
  type = "web"
  repo = ""
  url = "iac-tf.com"
  desc = "projectDesc-1"
}

```
Import

rum project can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_project.project project_id
```