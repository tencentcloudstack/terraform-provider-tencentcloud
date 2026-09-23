---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_aggregator"
sidebar_current: "docs-tencentcloud-resource-config_aggregator"
description: |-
  Provides a resource to create a Config aggregator.
---

# tencentcloud_config_aggregator

Provides a resource to create a Config aggregator.

## Example Usage

```hcl
resource "tencentcloud_config_aggregator" "example" {
  name        = "tf-example-aggregator"
  description = "tf example aggregator"
  type        = "CUSTOM"
  owner_uin   = "100012345678"

  aggregator_accounts {
    member_uin  = 100012345679
    member_name = "member-1"
  }

  aggregator_accounts {
    member_uin  = 100012345680
    member_name = "member-2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, String) Aggregator description.
* `name` - (Required, String) Aggregator name.
* `owner_uin` - (Required, String, ForceNew) Creator UIN of the aggregator. Required by Describe/Update/Delete APIs. Create does not return it.
* `type` - (Required, String, ForceNew) Aggregator type. Valid values: `RD` (global aggregator), `CUSTOM` (custom aggregator).
* `aggregator_accounts` - (Optional, List) Member account list of the aggregator, up to 100 entries.

The `aggregator_accounts` object supports the following:

* `member_name` - (Required, String) Member account name.
* `member_uin` - (Required, Int) Member account ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `account_group_id` - Aggregator ID.
* `aggregator_status` - Aggregator creation status.


## Import

Config aggregator can be imported using the composite id, e.g. `account_group_id#owner_uin`.

```
terraform import tencentcloud_config_aggregator.example ca-xxxxxxxx#100012345678
```

