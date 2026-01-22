---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ha_vip_instance_attachment"
sidebar_current: "docs-tencentcloud-resource-ha_vip_instance_attachment"
description: |-
  Provides a resource to create a vpc ha_vip_instance_attachment
---

# tencentcloud_ha_vip_instance_attachment

Provides a resource to create a vpc ha_vip_instance_attachment

## Example Usage

```hcl
resource "tencentcloud_ha_vip_instance_attachment" "ha_vip_instance_attachment" {
  instance_id   = "eni-xxxxxx"
  ha_vip_id     = "havip-xxxxxx"
  instance_type = "ENI"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The unique ID of the slave machine or network card to which HaVip is bound.
* `ha_vip_id` - (Optional, String, ForceNew) Unique ID of the HaVip instance.
* `instance_type` - (Optional, String, ForceNew) The type of HaVip binding. Values:CVM, ENI.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc ha_vip_instance_attachment can be imported using the id(${haVipId}#${instanceType}#${instanceId}), e.g.

```
terraform import tencentcloud_ha_vip_instance_attachment.ha_vip_instance_attachment ha_vip_instance_attachment_id
```

