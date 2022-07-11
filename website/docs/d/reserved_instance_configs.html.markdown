---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_reserved_instance_configs"
sidebar_current: "docs-tencentcloud-datasource-reserved_instance_configs"
description: |-
  Use this data source to query reserved instances configuration.
---

# tencentcloud_reserved_instance_configs

Use this data source to query reserved instances configuration.

## Example Usage

```hcl
data "tencentcloud_reserved_instance_configs" "config" {
  availability_zone = "na-siliconvalley-1"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) The available zone that the reserved instance locates at.
* `duration` - (Optional, Int) Validity period of the reserved instance. Valid values are `31536000`(1 year) and `94608000`(3 years).
* `instance_type` - (Optional, String) The type of reserved instance.
* `offering_type` - (Optional, String) Filter by Payment Type. Such as All Upfront.
* `product_description` - (Optional, String) Filter by the Platform Description (that is, operating system) for Reserved Instance billing. Shaped like: linux.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `config_list` - An information list of reserved instance configuration. Each element contains the following attributes:
  * `availability_zone` - Availability zone of the purchasable reserved instance.
  * `config_id` - Configuration ID of the purchasable reserved instance.
  * `currency_code` - Settlement currency of the reserved instance, which is a standard currency code as listed in ISO 4217.
  * `duration` - Validity period of the reserved instance.
  * `instance_type` - Instance type of the reserved instance.
  * `offering_type` - OfferingType of the reserved instance.
  * `platform` - Platform of the reserved instance.
  * `price` - Purchase price of the reserved instance.
  * `usage_price` - UsagePrice of the reserved instance.


