---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_services"
sidebar_current: "docs-tencentcloud-datasource-organization_services"
description: |-
  Use this data source to query detailed information of organization services
---

# tencentcloud_organization_services

Use this data source to query detailed information of organization services

## Example Usage

### Query all organization services

```hcl
data "tencentcloud_organization_services" "services" {}
```

### Query organization services by filter

```hcl
data "tencentcloud_organization_services" "services" {
  search_key = "KeyWord"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Keyword for search by name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Organization service list.
  * `can_assign_count` - Limit for the number of delegated admins. Note: This field may return null, indicating that no valid values can be obtained.
  * `console_url` - Console path of the organization service product. Note: This field may return null, indicating that no valid values can be obtained.
  * `description` - Organization service description. Note: This field may return null, indicating that no valid values can be obtained.
  * `document` - Help documentation. Note: This field may return null, indicating that no valid values can be obtained.
  * `grant_status` - Enabling status of organization service authorization. This field is valid when ServiceGrant is 1. Valid values: Enabled, Disabled. Note: This field may return null, indicating that no valid values can be obtained.
  * `is_assign` - Whether to support delegation. Valid values: 1 (yes), 2 (no). Note: This field may return null, indicating that no valid values can be obtained.
  * `is_set_management_scope` - Whether to support setting the delegated management scope. Valid values: 1 (yes), 2 (no).
Note: This field may return null, indicating that no valid values can be obtained.
  * `is_usage_status` - Whether to access the usage status. Valid values: 1 (yes), 2 (no). Note: This field may return null, indicating that no valid values can be obtained.
  * `member_num` - Number of the current delegated admins. Note: This field may return null, indicating that no valid values can be obtained.
  * `product_name` - Organization service product name. Note: This field may return null, indicating that no valid values can be obtained.
  * `product` - Organization service product identifier. Note: This field may return null, indicating that no valid values can be obtained.
  * `service_grant` - Whether to support organization service authorization. Valid values: 1 (yes), 2 (no). Note: This field may return null, indicating that no valid values can be obtained.
  * `service_id` - Organization service ID. Note: This field may return null, indicating that no valid values can be obtained.


