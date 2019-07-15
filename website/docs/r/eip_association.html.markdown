---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_association"
sidebar_current: "docs-tencentcloud-resource-cvm-eip-association"
description: |-
  Provides a TencentCloud EIP association.
---

# tencentcloud_eip_association

Provides an eip resource associated with other resource like CVM or ENI.

~> **NOTE:** Please DO NOT define `allocate_public_ip` in `tencentcloud_instance` resource when using `tencentcloud_eip_association`.

## Example Usage

Basic Usage

```hcl
resource "tencentcloud_eip_association" "foo" {
  eip_id      = "eip-xxxxxx"
  instance_id = "ins-xxxxxx"
}
```

or

```hcl
resource "tencentcloud_eip_association" "bar" {
  eip_id               = "eip-xxxxxx"
  network_interface_id = "eni-xxxxxx"
  private_ip           = "10.0.1.22"
}
```

## Argument Reference

The following arguments are supported:

* `eip_id` - (Required) The eip's id.
* `instance_id` - (Optional) The instance id going to bind with the EIP. This field is conflict with `network_interface_id` and `private_ip` fields.
* `network_interface_id` - (Optional) Indicates the network interface id like `eni-xxxxxx`. This field is conflict with `instance_id`.
* `private_ip` - (Optional) Indicates an IP belongs to the `network_interface_id`. This field is conflict with `instance_id`.


## Attributes Reference

The following attributes are exported:

* `id` - The association id.
* `eip_id` - The id of the EIP.
* `instance_id` - The instance id of the EIP bound with.
* `network_interface_id` - The network interface id.
* `private_ip` - (Optional) The IP belongs to the `network_interface_id`. 
