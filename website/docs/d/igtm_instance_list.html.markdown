---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_instance_list"
sidebar_current: "docs-tencentcloud-datasource-igtm_instance_list"
description: |-
  Use this data source to query detailed information of IGTM instance list
---

# tencentcloud_igtm_instance_list

Use this data source to query detailed information of IGTM instance list

## Example Usage

### Query all igtm instance list

```hcl
data "tencentcloud_igtm_instance_list" "example" {}
```

### Query igtm instance list by filters

```hcl
data "tencentcloud_igtm_instance_list" "example" {
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
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name, supported list as follows:
- InstanceId: IGTM instance ID.
- Domain: IGTM instance domain.
- MonitorId: Monitor ID.
- PoolId: Pool ID. This is a required parameter, not passing it will cause interface query failure.
* `value` - (Required, Set) Filter field value.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query, only supports filter field name as domain.
When fuzzy query is enabled, maximum Value length is 1, otherwise maximum Value length is 5. (Reserved field, not currently used).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_set` - Instance list.
  * `access_domain` - Access domain.
  * `access_sub_domain` - Access subdomain.
  * `access_type` - Cname domain access method
CUSTOM: Custom access domain
SYSTEM: System access domain.
  * `address_pool_num` - Bound address pool count.
  * `created_on` - Instance creation time.
  * `domain` - Business domain.
  * `global_ttl` - Global record expiration time.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `is_cname_configured` - Whether cname access: true accessed; false not accessed.
  * `monitor_num` - Bound monitor count.
  * `package_type` - Package type
FREE: Free version
STANDARD: Standard version
ULTIMATE: Ultimate version.
  * `pool_id` - Address pool ID.
  * `pool_name` - Address pool name.
  * `remark` - Remark.
  * `resource_id` - Resource ID.
  * `status` - Instance status, ENABLED: Normal, DISABLED: Disabled.
  * `strategy_num` - Strategy count.
  * `updated_on` - Instance update time.
  * `working_status` - Instance running status
NORMAL: Healthy
FAULTY: At risk
DOWN: Down
UNKNOWN: Unknown.
* `system_access_enabled` - Whether system domain access is supported: true supported; false not supported.


