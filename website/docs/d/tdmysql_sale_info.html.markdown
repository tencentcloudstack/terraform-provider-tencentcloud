---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_sale_info"
sidebar_current: "docs-tencentcloud-data_source-tdmysql_sale_info"
description: |-
  Provides a data source to query the TDSQL-C for MySQL sale information of regions and zones.
---

# tencentcloud_tdmysql_sale_info

Provides a data source to query the TDSQL-C for MySQL sale information of regions and zones.

## Example Usage

```hcl
data "tencentcloud_tdmysql_sale_info" "example" {
  instance_type = "serverless"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, String) Instance ID. Passing this parameter returns the sale information of the region and zones where this instance is located.
* `instance_type` - (Optional, String) Specify the sale information of regions that support a certain instance type. Enumerated values: serverless: return all regions that support serverless instances. Default value: none. Currently only serverless is supported; passing other values or leaving it empty returns the sale information of all regions.
* `src_region` - (Optional, String) Disaster recovery primary instance region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `available_disk_types` - Optional disk types of the zone. Note: this field may return null, indicating that no valid value is obtained.
    * `available_disk_types` - Supported disk types. Enumerated values: CLOUD_TCS: local disk; CLOUD_HSSD: enhanced cloud disk. Note: this field may return null, indicating that no valid value is obtained.
    * `enable` - Whether for sale. 1: for sale; 0: not for sale.
    * `is_default_master` - Whether it is the default primary zone. 0: no; 1: yes. Note: this field may return null, indicating that no valid value is obtained.
    * `is_support_serverless` - Whether serverless is supported.
    * `is_support_serverless` - Whether serverless is supported.
    * `support_types` - Supported types.
    * `support_types` - Whether it is an exclusive zone.
    * `zone_name` - Zone string in Chinese.
    * `zone_num` - Number of zones.
    * `zone` - Zone string in English.
    * `zones` - Zone combination.
  * `available_zone_num` - Number of optional multi-availability zones. Note: this field may return null, indicating that no valid value is obtained.
  * `enable` - Whether for sale. 1: for sale; 0: not for sale.
  * `is_support_serverless` - Whether serverless is supported.
  * `is_supported_log_replica` - Whether log replica is allowed. Note: this field may return null, indicating that no valid value is obtained.
  * `region_name` - Region string in Chinese.
  * `region` - Region string in English.
  * `zone_group` - Zone group combinations.
  * `zone_list` - Sale zone list.
* `id` - ID of the data source, a random token.
* `region_list` - Sale region information.


