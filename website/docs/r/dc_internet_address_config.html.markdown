---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_internet_address_config"
sidebar_current: "docs-tencentcloud-resource-dc_internet_address_config"
description: |-
  Provides a resource to create a dc internet address config
---

# tencentcloud_dc_internet_address_config

Provides a resource to create a dc internet address config

## Example Usage

```hcl
resource "tencentcloud_dc_internet_address" "example" {
  mask_len   = 30
  addr_type  = 2
  addr_proto = 0
}

resource "tencentcloud_dc_internet_address_config" "example" {
  instance_id = tencentcloud_dc_internet_address.example.id
  enable      = true
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) whether enable internet address.
* `instance_id` - (Required, String) internet public address id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dc internet address config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address_config.example ipv4-5091pc5v
```

