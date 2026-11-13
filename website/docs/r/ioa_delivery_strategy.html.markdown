---
subcategory: "Identity and Access Management Open APIs(IOA)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ioa_delivery_strategy"
sidebar_current: "docs-tencentcloud-resource-ioa_delivery_strategy"
description: |-
  Provides a resource to create or edit an IOA (T-Sec Internet of Amount Security Management Center) log delivery strategy.
---

# tencentcloud_ioa_delivery_strategy

Provides a resource to create or edit an IOA (T-Sec Internet of Amount Security Management Center) log delivery strategy.

Note: `ModifyDeliveryStrategy` is a combined create/update API. On creation the strategy id is
omitted and returned in the response; on update the id from the state is sent. Removing an optional
attribute from the configuration does not clear the corresponding value on the cloud side.
In addition, the describe API does not return the `enabled` field, so drift detection is not
available for it and the state value is kept on refresh.

## Example Usage

```hcl
resource "tencentcloud_ioa_delivery_strategy" "example" {
  name        = "tf-example"
  description = "delivery strategy for IOA logs"
  type        = 0
  config_id   = 100
  enabled     = true
  log_types   = ["040201000000", "040203000000"]
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Optional, Int) Delivery server id (32-bit only).
* `description` - (Optional, String) Delivery strategy description.
* `enabled` - (Optional, Bool) Whether the strategy is enabled. Note: DescribeDeliveryStrategy does not return this field, so drift detection is not available and the state value is kept on refresh.
* `log_types` - (Optional, Set) List of forwarded log types. Valid values: identity security `040201000000`, resource access log `040203000000`, client operation log `040202000000`, terminal security log `040205000000`, terminal control log `040206000000`, terminal compliance log `040207000000`, license authorization log `040212000000`.
* `name` - (Optional, String) Delivery strategy name.
* `type` - (Optional, Int) Delivery method. Valid values: `0` syslog, `2` kafka, `3` http (32-bit only).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `config_id` - Delivery server id.
    * `config_name` - Delivery server name.
    * `description` - Delivery strategy description.
    * `id` - Delivery strategy id.
    * `log_types` - List of forwarded log types.
    * `name` - Delivery strategy name.
    * `type` - Delivery method. Valid values: `0` syslog, `2` kafka, `3` http (32-bit only).
  * `data` - Delivery strategy detail.
* `data` - Delivery strategy detail data returned by DescribeDeliveryStrategy.
* `id` - ID of the delivery strategy, from the ID returned by ModifyDeliveryStrategy.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeDeliveryStrategy](https://www.tencentcloud.com/search/?wd=DescribeDeliveryStrategy&type=document&lang=en)
  - [ModifyDeliveryStrategy](https://www.tencentcloud.com/search/?wd=ModifyDeliveryStrategy&type=document&lang=en)
* Read: [DescribeDeliveryStrategy](https://www.tencentcloud.com/search/?wd=DescribeDeliveryStrategy&type=document&lang=en)
* Update:
  - [DescribeDeliveryStrategy](https://www.tencentcloud.com/search/?wd=DescribeDeliveryStrategy&type=document&lang=en)
  - [ModifyDeliveryStrategy](https://www.tencentcloud.com/search/?wd=ModifyDeliveryStrategy&type=document&lang=en)
* Delete: [DeleteDeliveryStrategy](https://www.tencentcloud.com/search/?wd=DeleteDeliveryStrategy&type=document&lang=en)

## Import

IOA delivery strategy can be imported using the id, e.g.

```
terraform import tencentcloud_ioa_delivery_strategy.example 42
```

