---
subcategory: "Identity and Access Management Open APIs (IOA)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ioa_delivery_strategy"
sidebar_current: "docs-tencentcloud-resource-ioa_delivery_strategy"
description: |-
  Provides a resource to create or edit an IOA (T-Sec Internet of Amount Security Management Center) log delivery strategy
---

# tencentcloud_ioa_delivery_strategy

Provides a resource to create or edit an IOA log delivery strategy.

~> Note: `ModifyDeliveryStrategy` is a combined create/update API. On creation the strategy id is
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

## Import

IOA delivery strategy can be imported using the id, e.g.

```bash
terraform import tencentcloud_ioa_delivery_strategy.example 42
```
