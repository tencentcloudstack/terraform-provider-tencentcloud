Provides a resource to create a project

~> **NOTE:** Project can not be destroyed. If run `terraform destroy`, project will be set invisible.

Example Usage

```hcl
resource "tencentcloud_project" "project" {
  project_name = "terraform-test"
  info         = "for terraform test"
}
```

Import

tag project can be imported using the id, e.g.

```
terraform import tencentcloud_project.project project_id
```