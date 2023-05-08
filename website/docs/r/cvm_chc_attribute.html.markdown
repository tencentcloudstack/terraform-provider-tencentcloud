---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_chc_attribute"
sidebar_current: "docs-tencentcloud-resource-cvm_chc_attribute"
description: |-
  Provides a resource to create a cvm chc_attribute
---

# tencentcloud_cvm_chc_attribute

Provides a resource to create a cvm chc_attribute

## Example Usage

```hcl
resource "tencentcloud_cvm_chc_attribute" "chc_attribute" {
  chc_id        = "chc-xxxxx"
  instance_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `chc_id` - (Required, String, ForceNew) CHC host ID.
* `bmc_security_group_ids` - (Optional, Set: [`String`], ForceNew) BMC network security group list.
* `bmc_user` - (Optional, String, ForceNew) Valid characters: Letters, numbers, hyphens and underscores.
* `device_type` - (Optional, String, ForceNew) Server type.
* `instance_name` - (Optional, String, ForceNew) CHC host name.
* `password` - (Optional, String, ForceNew) The password can contain 8 to 16 characters, including letters, numbers and special symbols (()`~!@#$%^&amp;amp;*-+=_|{}).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



