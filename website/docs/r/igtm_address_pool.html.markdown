---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_address_pool"
sidebar_current: "docs-tencentcloud-resource-igtm_address_pool"
description: |-
  Provides a resource to create a IGTM address pool
---

# tencentcloud_igtm_address_pool

Provides a resource to create a IGTM address pool

~> **NOTE:** Resource `tencentcloud_igtm_instance` needs to be created before using this resource.

## Example Usage

```hcl
resource "tencentcloud_igtm_address_pool" "example" {
  pool_name        = "tf-example"
  traffic_strategy = "WEIGHT"
  address_set {
    addr      = "1.1.1.1"
    is_enable = "ENABLED"
    weight    = 90
  }

  address_set {
    addr      = "2.2.2.2"
    is_enable = "DISABLED"
    weight    = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `address_set` - (Required, Set) Address list.
* `pool_name` - (Required, String) Address pool name, duplicates are not allowed.
* `traffic_strategy` - (Required, String) Traffic strategy: WEIGHT for load balancing, ALL for resolving all healthy addresses.
* `monitor_id` - (Optional, Int) Monitor ID.

The `address_set` object supports the following:

* `addr` - (Required, String) Address value: only supports IPv4, IPv6, and domain name formats.
Loopback addresses, reserved addresses, internal addresses, and Tencent reserved network segments are not supported.
* `is_enable` - (Required, String) Whether to enable: DISABLED for disabled, ENABLED for enabled.
* `location` - (Optional, String) Address name.
* `weight` - (Optional, Int) Weight, required when traffic strategy is WEIGHT; range 1-100.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `pool_id` - Address pool ID.

The `address_set` object exports the following:

* `address_id` - Address ID.
* `created_on` - Creation time.
* `status` - OK for normal, DOWN for failure, WARN for risk, UNKNOWN for probing, UNMONITORED for unknown.
* `updated_on` - Modification time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAddressPool](https://cloud.tencent.com/document/api/1551/120464) - 创建地址池
* Read: [DescribeAddressPoolDetail](https://cloud.tencent.com/document/api/1551/120461) - 地址池详情
* Update: [ModifyAddressPool](https://cloud.tencent.com/document/api/1551/120458) - 修改地址池
* Delete: [DeleteAddressPool](https://cloud.tencent.com/document/api/1551/120463) - 删除地址池

## Import

IGTM address pool can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_address_pool.example 1012132
```

