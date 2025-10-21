---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_modify_instance_bundle"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_modify_instance_bundle"
description: |-
  Use this data source to query detailed information of lighthouse modify_instance_bundle
---

# tencentcloud_lighthouse_modify_instance_bundle

Use this data source to query detailed information of lighthouse modify_instance_bundle

## Example Usage

```hcl
data "tencentcloud_lighthouse_modify_instance_bundle" "modify_instance_bundle" {
  instance_id = "lhins-xxxxxx"
  filters {
    name   = "bundle-id"
    values = ["bundle_gen_mc_med2_02"]

  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `filters` - (Optional, List) Filter list.
- `bundle-id`: filter by the bundle ID.
- `support-platform-type`: filter by system type, valid values: `LINUX_UNIX`, `WINDOWS`.
- `bundle-type`: filter according to package type, valid values: `GENERAL_BUNDLE`, `STORAGE_BUNDLE`, `ENTERPRISE_BUNDLE`, `EXCLUSIVE_BUNDLE`, `BEFAST_BUNDLE`.
- `bundle-state`: filter according to package status, valid values: `ONLINE`, `OFFLINE`.
NOTE: The upper limit of Filters per request is 10. The upper limit of Filter.Values is 5. Parameter does not support specifying both BundleIds and Filters.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter value of field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `modify_bundle_set` - Change package details.
  * `bundle` - Package information.
    * `bundle_display_label` - Package tag.Valid values:ACTIVITY: promotional packageNORMAL: regular packageCAREFREE: carefree package.
    * `bundle_id` - Package ID.
    * `bundle_sales_state` - Package sale status. Valid values are AVAILABLE, SOLD_OUT.
    * `bundle_type_description` - Package type description information.
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
  * `modify_bundle_state` - Change the status of the package. Value:
- SOLD_OUT: the package is sold out;
- AVAILABLE: support package changes;
- UNAVAILABLE: package changes are not supported for the time being.
  * `modify_price` - Change the price difference to be made up after the instance package.
    * `instance_price` - Instance price.
      * `currency` - A monetary unit of price. Value range CNY: RMB. USD: us dollar.
      * `discount_price` - Discounted price.
      * `discount` - Discount.
      * `original_bundle_price` - Original unit price of the package.
      * `original_price` - Original price.
  * `not_support_modify_message` - Package change reason information is not supported. When the package status is changed to `AVAILABLE`, the information is empty.


