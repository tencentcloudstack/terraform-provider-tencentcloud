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
resource "tencentcloud_oceanus_resource_config" "resource_config" {
  resource_id = "resource-xxx"
  resource_loc {
    storage_type = 1
    param {
      bucket = "scs-online-1257058945"
      path   = "251008563/100000006047/flink-cos-fs-hadoop-oceanus-1-20210304112"
      region = "ap-chengdu"
    }

  }
  remark        = "xxx"
  auto_delete   = 1
  work_space_id = "space-xxx"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Resource ID.
* `resource_loc` - (Required, List) Resource location.
* `remark` - (Optional, String) Resource description.
* `work_space_id` - (Optional, String) Workspace SerialId.

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



