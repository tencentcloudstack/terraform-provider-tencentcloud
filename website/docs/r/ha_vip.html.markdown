---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ha_vip"
sidebar_current: "docs-tencentcloud-resource-ha_vip"
description: |-
  Provides a resource to create a HA VIP.
---

# tencentcloud_ha_vip

Provides a resource to create a HA VIP.

## Example Usage

```hcl
resource "tencentcloud_ha_vip" "foo" {
  name      = "terraform_test"
  vpc_id    = "vpc-gzea3dd7"
  subnet_id = "subnet-4d4m4cd4s"
  vip       = "10.0.4.16"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the HA VIP. The length of character is limited to 1-60.
* `subnet_id` - (Required, String, ForceNew) Subnet ID.
* `vpc_id` - (Required, String, ForceNew) VPC ID.
* `check_associate` - (Optional, Bool, ForceNew) Whether to enable verification of the submachine or network card range during HAVIP drift. Not enabled by default.
* `vip` - (Optional, String, ForceNew) Virtual IP address, it must not be occupied and in this VPC network segment. If not set, it will be assigned after resource created automatically.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `address_ip` - EIP that is associated.
* `create_time` - Create time of the HA VIP.
* `instance_id` - Instance ID that is associated.
* `network_interface_id` - Network interface ID that is associated.
* `state` - State of the HA VIP. Valid value: `AVAILABLE`, `UNBIND`.


## Import

HA VIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_ha_vip.foo havip-kjqwe4ba
```

