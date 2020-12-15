---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_service_template"
sidebar_current: "docs-tencentcloud-resource-service_template"
description: |-
  Provides a resource to manage service template.
---

# tencentcloud_service_template

Provides a resource to manage service template.

## Example Usage

```hcl
resource "tencentcloud_service_template" "foo" {
  name     = "service-template-test"
  services = ["tcp:80", "udp:all", "icmp:10-30"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the service template.
* `services` - (Required) Service list. Valid protocols are  `tcp`, `udp`, `icmp`, `gre`. Single port(tcp:80), multi-port(tcp:80,443), port range(tcp:3306-20000), all(tcp:all) format are support. Protocol `icmp` and `gre` cannot specify port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CAM user can be imported using the service template, e.g.

```
$ terraform import tencentcloud_service_template.foo ppm-nwrggd14
```

