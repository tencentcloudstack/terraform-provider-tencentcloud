---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_device"
sidebar_current: "docs-tencentcloud-resource-bh_device"
description: |-
  Provides a resource to create a BH device
---

# tencentcloud_bh_device

Provides a resource to create a BH device

## Example Usage

```hcl
resource "tencentcloud_bh_device" "example" {
  device_set {
    os_name = "Linux"
    ip      = "1.1.1.1"
    port    = 22
    name    = "tf-example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `device_set` - (Required, List) Asset parameter list.
* `account_id` - (Optional, Int, ForceNew) Cloud account ID to which the asset belongs.

The `device_set` object supports the following:

* `ip` - (Required, String, ForceNew) IP address.
* `os_name` - (Required, String, ForceNew) The operating system name can only be one of the following: Host (Linux, Windows), Database (MySQL, SQL Server, MariaDB, PostgreSQL, MongoDBReplicaSet, MongoDBSharded, Redis), or Container (TKE, EKS).
* `port` - (Required, Int) Management port.
* `ap_code` - (Optional, String, ForceNew) Region to which the asset belongs.
* `ap_name` - (Optional, String, ForceNew) Region name.
* `department_id` - (Optional, String) Department ID to which the asset belongs.
* `enable_ssl` - (Optional, Int, ForceNew) Whether to enable SSL, 1: enable, 0: disable, only supports Redis assets.
* `instance_id` - (Optional, String, ForceNew) Asset instance ID.
* `ip_port_set` - (Optional, Set, ForceNew) Asset multi-node: IP and port fields.
* `name` - (Optional, String, ForceNew) Host name, can be empty.
* `public_ip` - (Optional, String, ForceNew) Public IP.
* `ssl_cert_name` - (Optional, String, ForceNew) SSL certificate name, required when EnableSSL is enabled.
* `ssl_cert` - (Optional, String, ForceNew) SSL certificate, required when EnableSSL is enabled.
* `subnet_id` - (Optional, String, ForceNew) Subnet to which the asset belongs.
* `vpc_id` - (Optional, String, ForceNew) VPC to which the asset belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `device_id` - ID of the device.


## Import

BH device can be imported using the id, e.g.

```
terraform import tencentcloud_bh_device.example 1875
```

