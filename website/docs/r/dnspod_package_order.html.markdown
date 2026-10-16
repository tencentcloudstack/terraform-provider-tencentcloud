---
subcategory: "DNSPod"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_package_order"
sidebar_current: "docs-tencentcloud-resource-dnspod_package_order"
description: |-
  Provides a resource to create a DNSPod package order
---

# tencentcloud_dnspod_package_order

Provides a resource to create a DNSPod package order

## Example Usage

```hcl
resource "tencentcloud_dnspod_package_order" "example" {
  domain = "demo.com"
  grade  = "DPG_ULTIMATE"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `grade` - (Required, String, ForceNew) Valid options for the package version are as follows: `DPG_PROFESSIONAL`; `DPG_ENTERPRISE`; `DPG_ULTIMATE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain_id` - Domain ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Read: [DescribeDomain](https://www.tencentcloud.com/search/?wd=DescribeDomain&type=document&lang=en)

## Import

DNSPod package order can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_package_order.example demo.com
```

