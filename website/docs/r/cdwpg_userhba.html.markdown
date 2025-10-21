---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_userhba"
sidebar_current: "docs-tencentcloud-resource-cdwpg_userhba"
description: |-
  Provides a resource to create a cdwpg cdwpg_userhba
---

# tencentcloud_cdwpg_userhba

Provides a resource to create a cdwpg cdwpg_userhba

## Example Usage

```hcl
resource "tencentcloud_cdwpg_userhba" "cdwpg_userhba" {
  instance_id = "cdwpg-zpiemnyd"
  hba_configs {
    type     = "host"
    database = "all"
    user     = "all"
    address  = "0.0.0.0/0"
    method   = "md5"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `hba_configs` - (Optional, List) HBA configuration array.

The `hba_configs` object supports the following:

* `address` - (Required, String) IP address.
* `database` - (Required, String) Database.
* `method` - (Required, String) Method.
* `type` - (Required, String) Type.
* `user` - (Required, String) User.
* `mask` - (Optional, String) Mask.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cdwpg cdwpg_userhba can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_userhba.cdwpg_userhba cdwpg_userhba_id
```

