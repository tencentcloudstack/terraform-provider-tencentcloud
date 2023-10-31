---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_resource"
sidebar_current: "docs-tencentcloud-resource-wedata_resource"
description: |-
  Provides a resource to create a wedata resource
---

# tencentcloud_wedata_resource

Provides a resource to create a wedata resource

## Example Usage

```hcl
resource "tencentcloud_wedata_resource" "example" {
  file_path       = "/datastudio/resource/demo"
  project_id      = "1612982498218618880"
  file_name       = "tf_example"
  cos_bucket_name = "wedata-demo-1314991481"
  cos_region      = "ap-guangzhou"
  files_size      = "8165"
}
```

## Argument Reference

The following arguments are supported:

* `cos_bucket_name` - (Required, String) Cos bucket name.
* `cos_region` - (Required, String) Cos bucket region.
* `file_name` - (Required, String) File name.
* `file_path` - (Required, String) For file path:/datastudio/resource/projectId/folderName; for folder path:/datastudio/resource/folderName.
* `files_size` - (Required, String) File size.
* `project_id` - (Required, String) Project ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `resource_id` - Resource ID.


## Import

wedata resource can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_resource.example 1612982498218618880#/datastudio/resource/demo#75431931-7d27-4034-b3de-3dc3348a220e
```

