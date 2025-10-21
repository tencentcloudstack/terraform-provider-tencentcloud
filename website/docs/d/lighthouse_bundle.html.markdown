---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_bundle"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_bundle"
description: |-
  Use this data source to query detailed information of lighthouse bundle
---

# tencentcloud_lighthouse_bundle

Use this data source to query detailed information of lighthouse bundle

## Example Usage

```hcl
data "tencentcloud_lighthouse_bundle" "bundle" {
}
```

## Argument Reference

The following arguments are supported:

* `bundle_ids` - (Optional, Set: [`String`]) Bundle ID list.
* `filters` - (Optional, List) Filter list.
- `bundle-id`: filter by the bundle ID.
- `support-platform-type`: filter by system type, valid values: `LINUX_UNIX`, `WINDOWS`.
- `bundle-type`: filter according to package type, valid values: `GENERAL_BUNDLE`, `STORAGE_BUNDLE`, `ENTERPRISE_BUNDLE`, `EXCLUSIVE_BUNDLE`, `BEFAST_BUNDLE`.
- `bundle-state`: filter according to package status, valid values: `ONLINE`, `OFFLINE`.
NOTE: The upper limit of Filters per request is 10. The upper limit of Filter.Values is 5. Parameter does not support specifying both BundleIds and Filters.
* `limit` - (Optional, Int) Number of returned results. Default value is 20. Maximum value is 100.
* `offset` - (Optional, Int) Offset. Default value is 0.
* `result_output_file` - (Optional, String) Used to save results.
* `zones` - (Optional, Set: [`String`]) Zone list, which contains all zones by default.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter value of field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bundle_set` - List of bundle details.
  * `bundle_display_label` - Package tag.Valid values:ACTIVITY: promotional packageNORMAL: regular packageCAREFREE: carefree package.
  * `bundle_id` - Package ID.
  * `bundle_sales_state` - Package sale status. Valid values are AVAILABLE, SOLD_OUT.
  * `bundle_type` - Package type.Valid values:GENERAL_BUNDLE: generalSTORAGE_BUNDLE: Storage.
  * `cpu` - CPU.
  * `internet_charge_type` - Network billing mode.
  * `internet_max_bandwidth_out` - Peak bandwidth in Mbps.
  * `memory` - Memory size in GB.
  * `monthly_traffic` - Monthly network traffic in Gb.
  * `price` - Current package unit price information.
    * `instance_price` - Instance price.
      * `currency` - Currency unit. Valid values: CNY and USD.
      * `discount_price` - Discounted price.
      * `discount` - Discount.
      * `original_bundle_price` - Original package unit price.
      * `original_price` - Original price.
  * `support_linux_unix_platform` - Whether Linux/Unix is supported.
  * `support_windows_platform` - Whether Windows is supported.
  * `system_disk_size` - System disk size.
  * `system_disk_type` - System disk type.


