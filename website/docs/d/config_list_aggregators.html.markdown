---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_list_aggregators"
sidebar_current: "docs-tencentcloud-datasource-config_list_aggregators"
description: |-
  Use this data source to query detailed information of Config aggregators.
---

# tencentcloud_config_list_aggregators

Use this data source to query detailed information of Config aggregators.

## Example Usage

### Query all aggregators

```hcl
data "tencentcloud_config_list_aggregators" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Aggregator list.
  * `account_count` - Number of accounts in the aggregator.
  * `account_group_id` - Aggregator ID.
  * `aggregator_status` - Aggregator status.
  * `create_time` - Creation time.
  * `description` - Aggregator description.
  * `member_name` - Member name.
  * `name` - Aggregator name.
  * `owner_uin` - Owner UIN.
  * `type` - Aggregator type.
* `total` - Total number of aggregators.


