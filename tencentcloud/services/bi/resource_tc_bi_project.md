Provides a resource to create a bi project

Example Usage

```hcl
resource "tencentcloud_bi_project" "project" {
  name               = "terraform_test"
  color_code         = "#7BD936"
  logo               = "TF-test"
  mark               = "project mark"
}
```

Import

bi project can be imported using the id, e.g.

```
terraform import tencentcloud_bi_project.project project_id
```