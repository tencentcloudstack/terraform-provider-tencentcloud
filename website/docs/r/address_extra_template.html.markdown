---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_address_extra_template"
sidebar_current: "docs-tencentcloud-resource-address_extra_template"
description: |-
  Provides a resource to manage address extra template.
---

# tencentcloud_address_extra_template

Provides a resource to manage address extra template.

~> **NOTE:** Compare to `tencentcloud_address_template`, It contains remarks.

## Example Usage

```hcl
resource "tencentcloud_address_extra_template" "foo" {
  name = "demo"

  addresses_extra {
    address     = "10.0.0.1"
    description = "create by terraform"
  }

  addresses_extra {
    address     = "10.0.1.0/24"
    description = "delete by terraform"
  }

  addresses_extra {
    address     = "10.0.0.1-10.0.0.100"
    description = "modify by terraform"
  }

  tags = {
    createBy = "terraform"
    deleteBy = "terraform"
  }

}
```

## Argument Reference

The following arguments are supported:

* `addresses_extra` - (Required, Set) The address information can contain remarks and be presented by the IP, CIDR block or IP address range.
* `name` - (Required, String) IP address template name.
* `tags` - (Optional, Map) Tags of the Addresses.

The `addresses_extra` object supports the following:

* `address` - (Required, String) IP address.
* `description` - (Optional, String) Remarks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Address template can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_extra_template.foo ipm-makf7k9e
```

