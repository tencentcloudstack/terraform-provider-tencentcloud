---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_address_template"
sidebar_current: "docs-tencentcloud-resource-address_template"
description: |-
  Provides a resource to manage address template.
---

# tencentcloud_address_template

Provides a resource to manage address template.

## Example Usage

```hcl
resource "tencentcloud_address_template" "foo" {
  name      = "cam-user-test"
  addresses = ["10.0.0.1", "10.0.1.0/24", "10.0.0.1-10.0.0.100"]
}
```

## Argument Reference

The following arguments are supported:

* `addresses` - (Required) Address list. IP(10.0.0.1), CIDR(10.0.1.0/24), IP range(10.0.0.1-10.0.0.100) format are supported.
* `name` - (Required, ForceNew) Name of the address template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Address template can be imported using the id, e.g.

```
$ terraform import tencentcloud_address_template.foo ipm-makf7k9e"
```

