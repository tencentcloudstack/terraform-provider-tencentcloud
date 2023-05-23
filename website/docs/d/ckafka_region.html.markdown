---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_region"
sidebar_current: "docs-tencentcloud-datasource-ckafka_region"
description: |-
  Use this data source to query detailed information of ckafka region
---

# tencentcloud_ckafka_region

Use this data source to query detailed information of ckafka region

## Example Usage

```hcl
data "tencentcloud_ckafka_region" "region" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Return a list of region enumeration results.
  * `area_name` - area name.
  * `ipv6` - Whether to support ipv6, 0: means not supported, 1: means supported.
  * `multi_zone` - Whether to support cross-availability zones, 0: means not supported, 1: means supported.
  * `region_code_v3` - Region Code(V3 version).
  * `region_code` - Region Code.
  * `region_id` - region ID.
  * `region_name` - geographical name.
  * `support` - NONE: The default value does not support any special models CVM: Supports CVM types.


