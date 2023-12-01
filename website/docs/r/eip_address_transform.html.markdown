---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_address_transform"
sidebar_current: "docs-tencentcloud-resource-eip_address_transform"
description: |-
  Provides a resource to create a eip address_transform
---

# tencentcloud_eip_address_transform

Provides a resource to create a eip address_transform

## Example Usage

```hcl
resource "tencentcloud_eip_address_transform" "address_transform" {
  instance_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) the instance ID of a normal public network IP to be operated. eg:ins-23mk45jn.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

eip address_transform can be imported using the id, e.g.

```
terraform import tencentcloud_eip_address_transform.address_transform address_transform_id
```

