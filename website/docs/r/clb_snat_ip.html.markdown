---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_snat_ip"
sidebar_current: "docs-tencentcloud-resource-clb_snat_ip"
description: |-
  Provide a resource to create a SnatIp of CLB instance.
---

# tencentcloud_clb_snat_ip

Provide a resource to create a SnatIp of CLB instance.

~> **NOTE:** Target CLB instance must enable `snat_pro` before creating snat ips.
~> **NOTE:** Dynamic allocate IP doesn't support for now.

## Example Usage

```hcl
resource "tencentcloud_clb_instance" "snat_test" {
  network_type = "OPEN"
  clb_name     = "tf-clb-snat-test"
}

resource "tencentcloud_clb_snat_ip" "foo" {
  clb_id = tencentcloud_clb_instance.snat_test.id
  ips {
    subnet_id = "subnet-12345678"
    ip        = "172.16.0.1"
  }
  ips {
    subnet_id = "subnet-12345678"
    ip        = "172.16.0.2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, String) CLB instance ID.
* `ips` - (Optional, Set) Snat IP address config.

The `ips` object supports the following:

* `ip` - (Required, String) Snat IP.
* `subnet_id` - (Required, String) Subnet ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ClbSnatIp instance can be imported by clb instance id, e.g.
```
$ terraform import tencentcloud_clb_snat_ip.test clb_id
```

