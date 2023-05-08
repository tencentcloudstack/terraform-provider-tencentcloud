---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_create_image_signature_operation"
sidebar_current: "docs-tencentcloud-resource-tcr_create_image_signature_operation"
description: |-
  Provides a resource to create a tcr image_signature_operation
---

# tencentcloud_tcr_create_image_signature_operation

Provides a resource to create a tcr image_signature_operation

## Example Usage

```hcl
resource "tencentcloud_tcr_create_image_signature_operation" "image_signature_operation" {
  registry_id     = "tcr-xxx"
  namespace_name  = "ns"
  repository_name = "repo"
  image_version   = "v1"

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



## Import

tcr image_signature_operation can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_create_image_signature_operation.image_signature_operation image_signature_operation_id
```

