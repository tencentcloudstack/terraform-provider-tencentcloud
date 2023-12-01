---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_pending_risk_info"
sidebar_current: "docs-tencentcloud-datasource-antiddos_pending_risk_info"
description: |-
  Use this data source to query detailed information of antiddos pending risk info
---

# tencentcloud_antiddos_pending_risk_info

Use this data source to query detailed information of antiddos pending risk info

## Example Usage

```hcl
data "tencentcloud_antiddos_pending_risk_info" "pending_risk_info" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attacking_count` - Number of resources in the attack.
* `blocking_count` - Number of resources in blockage.
* `expired_count` - Number of expired resources.
* `is_paid_usr` - Is it a paid user? True: paid user, false: regular user.
* `total` - Total number of all pending risk events.


