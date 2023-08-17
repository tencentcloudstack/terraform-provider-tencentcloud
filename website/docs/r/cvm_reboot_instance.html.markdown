---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_reboot_instance"
sidebar_current: "docs-tencentcloud-resource-cvm_reboot_instance"
description: |-
  Provides a resource to create a cvm reboot_instance
---

# tencentcloud_cvm_reboot_instance

Provides a resource to create a cvm reboot_instance

## Example Usage

```hcl
resource "tencentcloud_cvm_reboot_instance" "reboot_instance" {
  instance_id = "ins-f9jr4bd2"
  stop_type   = "SOFT_FIRST"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `force_reboot` - (Optional, Bool, ForceNew, **Deprecated**) It has been deprecated from version 1.81.21. Please use `stop_type` instead. This parameter has been disused. We recommend using StopType instead. Note that ForceReboot and StopType parameters cannot be specified at the same time. Whether to forcibly restart an instance after a normal restart fails. Valid values are `TRUE` and `FALSE`. Default value: FALSE.
* `stop_type` - (Optional, String, ForceNew) Shutdown type. Valid values: `SOFT`: soft shutdown; `HARD`: hard shutdown; `SOFT_FIRST`: perform a soft shutdown first, and perform a hard shutdown if the soft shutdown fails. Default value: SOFT.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



