---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_address_pool_list"
sidebar_current: "docs-tencentcloud-datasource-igtm_address_pool_list"
description: |-
  Use this data source to query detailed information of IGTM address pool list
---

# tencentcloud_igtm_address_pool_list

Use this data source to query detailed information of IGTM address pool list

## Example Usage

### Query all address pool list

```hcl
data "tencentcloud_igtm_address_pool_list" "example" {}
```

### Query address pool list by filter

```hcl
data "tencentcloud_igtm_address_pool_list" "example" {
  filters {
    name  = "PoolName"
    value = ["tf-example"]
    fuzzy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Alert filter conditions.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name, supported list as follows:
- PoolName: Address pool name.
- MonitorId: Monitor ID. This is a required parameter, failure to provide will cause interface query failure.
* `value` - (Required, Set) Filter field value.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query, only supports filter field name as domain.
When fuzzy query is enabled, maximum Value length is 1, otherwise maximum Value length is 5. (Reserved field, currently not used).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address_pool_set` - Resource group list.
  * `addr_type` - Address pool address type: IPV4, IPV6, DOMAIN.
  * `address_num` - Address count.
  * `address_set` - Address pool address information.
    * `addr` - Address value: only supports IPv4, IPv6 and domain name formats;
Loopback addresses, reserved addresses, internal network addresses and Tencent reserved network segments are not supported.
    * `address_id` - Address ID.
    * `created_on` - Creation time.
    * `is_enable` - Whether to enable: DISABLED disabled; ENABLED enabled.
    * `location` - Address name.
    * `status` - OK normal, DOWN failure, WARN risk, UNKNOWN detecting, UNMONITORED unknown.
    * `updated_on` - Modification time.
    * `weight` - Weight, required when traffic strategy is WEIGHT; range 1-100.
  * `created_on` - Creation time.
  * `instance_info` - Instance related information.
    * `instance_id` - Instance ID.
    * `instance_name` - Instance name.
  * `monitor_group_num` - Probe point count.
  * `monitor_id` - Monitor ID.
  * `monitor_task_num` - Detection task count.
  * `pool_id` - Address pool ID.
  * `pool_name` - Address pool name.
  * `status` - OK normal, DOWN failure, WARN risk, UNKNOWN unknown.
  * `traffic_strategy` - Traffic strategy: WEIGHT load balancing, ALL resolve all.
  * `updated_on` - Update time.


