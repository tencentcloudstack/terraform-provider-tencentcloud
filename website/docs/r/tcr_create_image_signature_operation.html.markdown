---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_create_image_signature_operation"
sidebar_current: "docs-tencentcloud-resource-tcr_create_image_signature_operation"
description: |-
  Provides a resource to operate a tcr image signature.
---

# tencentcloud_tcr_create_image_signature_operation

Provides a resource to operate a tcr image signature.

## Example Usage

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_example_ns"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_repository" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  name           = "test"
  brief_desc     = "111"
  description    = "111111111111111111111111111111111111"
}

resource "tencentcloud_tcr_create_image_signature_operation" "example" {
  registry_id     = tencentcloud_tcr_instance.example.id
  namespace_name  = tencentcloud_tcr_namespace.example.name
  repository_name = tencentcloud_tcr_repository.example.name
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

