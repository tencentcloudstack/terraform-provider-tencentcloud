---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_resource_config"
sidebar_current: "docs-tencentcloud-resource-oceanus_resource_config"
description: |-
  Provides a resource to create a oceanus resource_config
---

# tencentcloud_oceanus_resource_config

Provides a resource to create a oceanus resource_config

## Example Usage

```hcl
resource "tencentcloud_oceanus_resource" "example" {
  resource_loc {
    storage_type = 1
    param {
      bucket = "keep-terraform-1257058945"
      path   = "OceanusResource/junit-4.13.1.jar"
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

resource "tencentcloud_oceanus_resource_config" "example" {
  resource_id = tencentcloud_oceanus_resource.example.resource_id
  resource_loc {
    storage_type = 1
    param {
      bucket = "keep-terraform-1257058945"
      path   = "OceanusResource/junit-4.13.2.jar"
      region = "ap-guangzhou"
    }
  }

  remark        = "config remark."
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Resource ID.
* `resource_loc` - (Required, List) Resource location.
* `remark` - (Optional, String) Resource description.
* `work_space_id` - (Optional, String) Workspace SerialId.

The `param` object of `resource_loc` supports the following:

* `bucket` - (Required, String) Resource bucket.
* `path` - (Required, String) Resource path.
* `region` - (Optional, String) Resource region, if not set, use resource region, note: this field may return null, indicating that no valid values can be obtained.

The `resource_loc` object supports the following:

* `param` - (Required, List) Json to describe resource location.
* `storage_type` - (Required, Int) The available storage types for resource location are currently limited to 1:COS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `version` - Resource Config Version.


