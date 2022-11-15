---
subcategory: "Performance Testing Service(PTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_pts_file"
sidebar_current: "docs-tencentcloud-resource-pts_file"
description: |-
  Provides a resource to create a pts file
---

# tencentcloud_pts_file

Provides a resource to create a pts file

~> **NOTE:** Modification is not currently supported, please go to the console to modify.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `file_id` - (Required, String) File id.
* `kind` - (Required, Int) File kind, parameter file-1, protocol file-2, request file-3.
* `name` - (Required, String) File name.
* `project_id` - (Required, String) Project id.
* `size` - (Required, Int) File size.
* `type` - (Required, String) File type, folder-folder.
* `file_infos` - (Optional, List) Files in a folder.
* `head_lines` - (Optional, Set: [`String`]) The first few lines of data.
* `header_columns` - (Optional, Set: [`String`]) Meter head.
* `header_in_file` - (Optional, Bool) Whether the header is in the file.
* `line_count` - (Optional, Int) Line count.
* `tail_lines` - (Optional, Set: [`String`]) The last few lines of data.

The `file_infos` object supports the following:

* `file_id` - (Optional, String) File id.
* `name` - (Optional, String) File name.
* `size` - (Optional, Int) File size.
* `type` - (Optional, String) File type.
* `updated_at` - (Optional, String) Update time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

pts file can be imported using the project_id#file_id, e.g.
```
$ terraform import tencentcloud_pts_file.file project-45vw7v82#file-de2dbaf8
```

