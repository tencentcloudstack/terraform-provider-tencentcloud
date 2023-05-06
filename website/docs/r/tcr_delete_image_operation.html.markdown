---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_delete_image_operation"
sidebar_current: "docs-tencentcloud-resource-tcr_delete_image_operation"
description: |-
  Provides a resource to create a tcr delete_image_operation
---

# tencentcloud_tcr_delete_image_operation

Provides a resource to create a tcr delete_image_operation

## Example Usage

```hcl
resource "tencentcloud_tcr_delete_image_operation" "delete_image_operation" {
  registry_id     = "tcr-xxx"
  repository_name = "repo"
  image_version   = "v1"
  namespace_name  = "ns"
}
```

## Argument Reference

The following arguments are supported:

* `image_version` - (Required, String, ForceNew) image version name.
* `namespace_name` - (Required, String, ForceNew) namespace name.
* `registry_id` - (Required, String, ForceNew) instance id.
* `repository_name` - (Required, String, ForceNew) repository name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



