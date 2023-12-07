Provides a resource to create a pts file

~> **NOTE:** Modification is not currently supported, please go to the console to modify.

Example Usage

```hcl
resource "tencentcloud_pts_file" "file" {
    file_id        = "file-de2dbaf8"
    header_in_file = false
    kind           = 3
    line_count     = 0
    name           = "iac.txt"
    project_id     = "project-45vw7v82"
    size           = 10799
    type           = "text/plain"
    # header_columns = ""
    # file_infos {
    # name = ""
    # size = ""
    # type = ""
    # updated_at = ""
    # }
}

```
Import

pts file can be imported using the project_id#file_id, e.g.
```
$ terraform import tencentcloud_pts_file.file project-45vw7v82#file-de2dbaf8
```