---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_protocol_template"
sidebar_current: "docs-tencentcloud-resource-protocol_template"
description: |-
  Provides a resource to manage protocol template.
---

# tencentcloud_protocol_template

Provides a resource to manage protocol template.

## Example Usage

```hcl
resource "tencentcloud_protocol_template" "foo" {
  name      = "protocol-template-test"
  protocols = ["tcp:80", "udp:all", "icmp:10-30"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Name of the protocol template.
* `protocols` - (Required, Set: [`String`]) Protocol list. Valid protocols are  `tcp`, `udp`, `icmp`, `gre`. Single port(tcp:80), multi-port(tcp:80,443), port range(tcp:3306-20000), all(tcp:all) format are support. Protocol `icmp` and `gre` cannot specify port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Protocol template can be imported using the id, e.g.

```
$ terraform import tencentcloud_protocol_template.foo ppm-nwrggd14
```

