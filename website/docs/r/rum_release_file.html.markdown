---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_release_file"
sidebar_current: "docs-tencentcloud-resource-rum_release_file"
description: |-
  Provides a resource to create a rum release file
---

# tencentcloud_rum_release_file

Provides a resource to create a rum release file

## Example Usage

```hcl
resource "tencentcloud_rum_release_file" "example" {
  project_id      = 123
  version         = "1.0"
  file_key        = "120000-last-1632921299138-index.js.map"
  file_name       = "index.js.map"
  file_hash       = "b148c43fd81d845ba7cc6907928ce430"
  release_file_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `file_hash` - (Required, String, ForceNew) Release file hash.
* `file_key` - (Required, String, ForceNew) Release file unique key.
* `file_name` - (Required, String, ForceNew) Release file name.
* `project_id` - (Required, Int, ForceNew) Project ID.
* `release_file_id` - (Required, Int, ForceNew) Release file id.
* `version` - (Required, String, ForceNew) Release File version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateReleaseFile](https://www.tencentcloud.com/document/api/1131/45345)
* Read: [DescribeReleaseFiles](https://www.tencentcloud.com/document/api/1131/45305)
* Delete: [DeleteReleaseFile](https://www.tencentcloud.com/document/api/1131/45336)

## Import

rum release_file can be imported using the id, e.g.

```
terraform import tencentcloud_rum_release_file.example projectId#releaseFileId
```

