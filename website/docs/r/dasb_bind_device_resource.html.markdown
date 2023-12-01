---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_bind_device_resource"
sidebar_current: "docs-tencentcloud-resource-dasb_bind_device_resource"
description: |-
  Provides a resource to create a dasb bind_device_resource
---

# tencentcloud_dasb_bind_device_resource

Provides a resource to create a dasb bind_device_resource

## Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_resource" "example" {
  resource_id   = "bh-saas-ocmzo6lgxiv"
  device_id_set = [17, 18]
}
```

## Argument Reference

The following arguments are supported:

* `device_id_set` - (Required, Set: [`Int`], ForceNew) Asset ID collection.
* `resource_id` - (Required, String, ForceNew) Bastion host service ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



