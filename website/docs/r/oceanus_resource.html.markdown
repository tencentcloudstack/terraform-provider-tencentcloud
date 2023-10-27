---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_resource"
sidebar_current: "docs-tencentcloud-resource-oceanus_resource"
description: |-
  Provides a resource to create a oceanus resource
---

# tencentcloud_oceanus_resource

Provides a resource to create a oceanus resource

## Example Usage

```hcl
resource "tencentcloud_oceanus_resource" "example" {
  resource_loc {
    storage_type = 1
    param {
      bucket = "keep-terraform-1257058945"
      path   = "OceanusResource/junit-4.13.2.jar"
      region = "ap-guangzhou"
    }
  }

  resource_type          = 1
  remark                 = "remark."
  name                   = "tf_example"
  resource_config_remark = "config remark."
  folder_id              = "folder-7ctl246z"
  work_space_id          = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `resource_loc` - (Required, List) Resource location.
* `resource_type` - (Required, Int) Resource type, only support JAR now, value is 1.
* `folder_id` - (Optional, String) Folder id.
* `name` - (Optional, String) Resource name.
* `remark` - (Optional, String) Resource description.
* `resource_config_remark` - (Optional, String) Resource version description.
* `work_space_id` - (Optional, String) Workspace serialId.

The `param` object supports the following:

* `bucket` - (Required, String) Resource bucket.
* `path` - (Required, String) Resource path.
* `region` - (Optional, String) Resource region, if not set, use resource region, note: this field may return null, indicating that no valid values can be obtained.

The `resource_loc` object supports the following:

* `param` - (Required, List) Json to describe resource location.
* `storage_type` - (Required, Int) The available storage types for resource location are currently limited to 1:COS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `resource_id` - Resource ID.
* `version` - Resource Version.


