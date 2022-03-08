---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_record"
sidebar_current: "docs-tencentcloud-resource-dnspod_record"
description: |-
  Provide a resource to create a DnsPod record.
---

# tencentcloud_dnspod_record

Provide a resource to create a DnsPod record.

## Example Usage

```hcl
resource "tencentcloud_dnspod_record" "demo" {
  domain      = "mikatong.com"
  record_type = "A"
  record_line = "默认"
  value       = "1.2.3.9"
  sub_domain  = "demo"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The Domain.
* `record_line` - (Required) The record line.
* `record_type` - (Required) The record type.
* `value` - (Required) The record value.
* `mx` - (Optional) MX priority, valid when the record type is MX, range 1-20. Note: must set when record type equal MX.
* `status` - (Optional) Records the initial state, with values ranging from ENABLE and DISABLE. The default is ENABLE, and if DISABLE is passed in, resolution will not take effect and the limits of load balancing will not be verified.
* `sub_domain` - (Optional) The host records, default value is `@`.
* `ttl` - (Optional) TTL, the range is 1-604800, and the minimum value of different levels of domain names is different. Default is 600.
* `weight` - (Optional) Weight information. An integer from 0 to 100. Only enterprise VIP domain names are available, 0 means off, does not pass this parameter, means that the weight information is not set. Default is 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `monitor_status` - The D monitoring status of the record.


