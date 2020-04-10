---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_association"
sidebar_current: "docs-tencentcloud-resource-eip_association"
description: |-
  Provides an eip resource associated with other resource like CVM, ENI and CLB.
---

# tencentcloud_eip_association

Provides an eip resource associated with other resource like CVM, ENI and CLB.

~> **NOTE:** Please DO NOT define `allocate_public_ip` in `tencentcloud_instance` resource when using `tencentcloud_eip_association`.

## Example Usage

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

* `eip_id` - (Required, ForceNew) The id of eip.
* `instance_id` - (Optional, ForceNew) The CVM or CLB instance id going to bind with the eip. This field is conflict with `network_interface_id` and `private_ip fields`.
* `network_interface_id` - (Optional, ForceNew) Indicates the network interface id like `eni-xxxxxx`. This field is conflict with `instance_id`.
* `private_ip` - (Optional, ForceNew) Indicates an IP belongs to the `network_interface_id`. This field is conflict with `instance_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



