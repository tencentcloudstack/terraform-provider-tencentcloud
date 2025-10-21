---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_zone"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_zone"
description: |-
  Use this data source to query detailed information of cynosdb zone
---

# tencentcloud_cynosdb_zone

Use this data source to query detailed information of cynosdb zone

## Example Usage

```hcl
data "tencentcloud_cynosdb_zone" "zone" {
  include_virtual_zones = true
  show_permission       = true
}
```

## Argument Reference

The following arguments are supported:

* `include_virtual_zones` - (Optional, Bool) Is virtual zone included.
* `result_output_file` - (Optional, String) Used to save results.
* `show_permission` - (Optional, Bool) Whether to display all available zones under the region and display the permissions of each available zone of the user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_set` - Information of region.
  * `db_type` - Database type.
  * `modules` - Regional module support.
    * `is_disable` - Is zone on sale, optional values: yes, no.
    * `module_name` - Module name.
  * `region_id` - Region ID.
  * `region_zh` - Region name in Chinese.
  * `region` - Region in English.
  * `zone_set` - List of available zones for sale.
    * `has_permission` - Whether the user have zone permissionsNote: This field may return null, indicating that no valid value can be obtained.
    * `is_support_normal` - Does it support normal clusters, 0:Not supported 1:Support.
    * `is_support_serverless` - Does it support serverless clusters, 0:Not supported 1:Support.
    * `is_whole_rdma_zone` - Is zone Rdma.
    * `physical_zone` - Physical zone.
    * `zone_id` - ZoneId.
    * `zone_zh` - Zone name in Chinesee.
    * `zone` - Zone name in English.


