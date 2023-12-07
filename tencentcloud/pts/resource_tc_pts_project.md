Provides a resource to create a pts project

Example Usage

```hcl
resource "tencentcloud_pts_project" "project" {
  name = "ptsObjectName-1"
  description = "desc"
  tags {
    tag_key = "createdBy"
    tag_value = "terraform"
  }
}

```
Import

pts project can be imported using the id, e.g.
```
$ terraform import tencentcloud_pts_project.project project-1ep27k1m
```