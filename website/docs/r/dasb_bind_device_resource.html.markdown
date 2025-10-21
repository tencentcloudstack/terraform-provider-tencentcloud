---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_bind_device_resource"
sidebar_current: "docs-tencentcloud-resource-dasb_bind_device_resource"
description: |-
  Provides a resource to create a dasb bind device resource
---

# tencentcloud_dasb_bind_device_resource

Provides a resource to create a dasb bind device resource

## Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_resource" "example" {
  resource_id   = "bh-saas-weyosfym"
  device_id_set = [17, 18]
}
```

### Or custom domain_id parameters

```hcl
resource "tencentcloud_dasb_bind_device_resource" "example" {
  resource_id   = "bh-saas-lx1pxhli"
  domain_id     = "net-31nssj3n"
  device_id_set = [115, 116]
}
```

## Argument Reference

The following arguments are supported:

* `device_id_set` - (Required, Set: [`Int`]) Asset ID collection.
* `resource_id` - (Required, String, ForceNew) Bastion host service ID.
* `domain_id` - (Optional, String, ForceNew) Network Domain ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



