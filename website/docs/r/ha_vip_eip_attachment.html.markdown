---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ha_vip_eip_attachment"
sidebar_current: "docs-tencentcloud-resource-ha_vip_eip_attachment"
description: |-
  Provides a resource to create a HA VIP EIP attachment.
---

# tencentcloud_ha_vip_eip_attachment

Provides a resource to create a HA VIP EIP attachment.

## Example Usage

```hcl
resource "tencentcloud_ha_vip_eip_attachment" "foo" {
  havip_id   = "havip-kjqwe4ba"
  address_ip = "1.1.1.1"
}
```

## Argument Reference

The following arguments are supported:

* `address_ip` - (Required, ForceNew) Public address of the EIP.
* `havip_id` - (Required, ForceNew) Id of the attached HA VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

HA VIP EIP attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_ha_vip_eip_attachment.foo havip-kjqwe4ba#1.1.1.1
```

