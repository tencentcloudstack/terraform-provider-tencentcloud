---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_cross_border_region_bandwidth_limits"
sidebar_current: "docs-tencentcloud-datasource-ccn_cross_border_region_bandwidth_limits"
description: |-
  Use this data source to query detailed information of ccn_cross_border_region_bandwidth_limits
---

# tencentcloud_ccn_cross_border_region_bandwidth_limits

Use this data source to query detailed information of ccn_cross_border_region_bandwidth_limits

-> **NOTE:** This resource is dedicated to Unicom.

## Example Usage

```hcl
data "tencentcloud_ccn_cross_border_region_bandwidth_limits" "ccn_region_bandwidth_limits" {
  filters {
    name   = "source-region"
    values = ["ap-guangzhou"]
  }

  filters {
    name   = "destination-region"
    values = ["ap-shanghai"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter condition. Currently, only one value is supported. The supported fields, 1)source-region, the value is like ap-guangzhou; 2)destination-region, the value is like ap-shanghai; 3)ccn-ids,cloud network ID array, the value is like ccn-12345678; 4)user-account-id,user account ID, the value is like 12345678.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) attribute name.
* `values` - (Required, Set) Value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ccn_bandwidth_set` - Info of cross region ccn instance.
  * `ccn_id` - ccn id.
  * `ccn_region_bandwidth_limit` - bandwidth limit of cross region.
    * `bandwidth_limit` - bandwidth list(Mbps).
    * `destination_region` - destination region, such as.
    * `source_region` - source region, such as &#39;ap-shanghai&#39;.
  * `created_time` - create time.
  * `expired_time` - expired time.
  * `instance_charge_type` - `POSTPAID` or `PREPAID`.
  * `is_cross_border` - if cross region.
  * `is_security_lock` - `true` means locked.
  * `market_id` - market id.
  * `region_flow_control_id` - Id of RegionFlowControl.
  * `renew_flag` - renew flag.
  * `update_time` - update time.
  * `user_account_id` - user account id.


