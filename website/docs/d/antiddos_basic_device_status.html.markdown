---
subcategory: "Anti-DDoS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_basic_device_status"
sidebar_current: "docs-tencentcloud-datasource-antiddos_basic_device_status"
description: |-
  Use this data source to query detailed information of antiddos basic_device_status
---

# tencentcloud_antiddos_basic_device_status

Use this data source to query detailed information of antiddos basic_device_status

## Example Usage

```hcl
data "tencentcloud_antiddos_basic_device_status" "basic_device_status" {
  ip_list = [
    "127.0.0.1"
  ]
  filter_region = 1
}
```

## Argument Reference

The following arguments are supported:

* `filter_region` - (Optional, Int) Region Id.
* `id_list` - (Optional, Set: [`String`]) Named resource transfer ID.
* `ip_list` - (Optional, Set: [`String`]) Ip resource list.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clb_data` - Note: This field may return null, indicating that a valid value cannot be obtained.
  * `key` - Properties name.
  * `value` - Properties value.
* `data` - Return resources and status, status code: 1- Blocking status 2- Normal status 3- Attack status.
  * `key` - Properties name.
  * `value` - Properties value.


