---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_instance_package_list"
sidebar_current: "docs-tencentcloud-datasource-igtm_instance_package_list"
description: |-
  Use this data source to query detailed information of IGTM instance package list
---

# tencentcloud_igtm_instance_package_list

Use this data source to query detailed information of IGTM instance package list

## Example Usage

### Query all igtm instance package list

```hcl
data "tencentcloud_igtm_instance_package_list" "example" {}
```

### Query igtm instance package list by filter

```hcl
data "tencentcloud_igtm_instance_package_list" "example" {
  filters {
    name  = "InstanceId"
    value = ["gtm-uukztqtoaru"]
    fuzzy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions.
* `is_used` - (Optional, Int) Whether used: 0 not used 1 used.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name, supported list as follows:
- InstanceId: instance ID.
- InstanceName: instance name.
- ResourceId: package ID.
- PackageType: package type. This is a required parameter, not passing it will cause interface query failure.
* `value` - (Required, Set) Filter field value.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query, only supports filter field name as domain.
When fuzzy query is enabled, maximum Value length is 1, otherwise maximum Value length is 5. (Reserved field, not currently used).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_set` - Instance package list.
  * `auto_renew_flag` - Whether auto-renew 0 no 1 yes.
  * `cost_item_list` - Billing item.
    * `cost_name` - Billing item name.
    * `cost_value` - Billing item value.
  * `create_time` - Package creation time.
  * `current_deadline` - Package expiration time.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `is_expire` - Whether expired 0 no 1 yes.
  * `min_check_interval` - Minimum check interval time s.
  * `min_global_ttl` - Minimum TTL s.
  * `package_type` - Package type
FREE: Free version
STANDARD: Standard version
ULTIMATE: Ultimate version.
  * `remark` - Remark.
  * `resource_id` - Instance package resource ID.
  * `schedule_strategy` - Strategy type: LOCATION schedule by geographic location, DELAY schedule by delay.
  * `status` - Instance status
ENABLED: Normal
DISABLED: Disabled.
  * `traffic_strategy` - Traffic strategy type: ALL return all, WEIGHT weight.


