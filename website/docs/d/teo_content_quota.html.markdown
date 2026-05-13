---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_content_quota"
sidebar_current: "docs-tencentcloud-datasource-teo_content_quota"
description: |-
  Use this data source to query detailed information of TEO content quota
---

# tencentcloud_teo_content_quota

Use this data source to query detailed information of TEO content quota

## Example Usage

```hcl
data "tencentcloud_teo_content_quota" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `prefetch_quota` - Cache prefetch quota list.
  * `batch` - Single batch submission quota limit.
  * `daily_available` - Daily remaining available quota.
  * `daily` - Daily submission quota limit.
  * `type` - Quota type. Valid values: `prefetch_url`.
* `purge_quota` - Cache purge quota list.
  * `batch` - Single batch submission quota limit.
  * `daily_available` - Daily remaining available quota.
  * `daily` - Daily submission quota limit.
  * `type` - Quota type. Valid values: `purge_prefix`, `purge_url`, `purge_host`, `purge_all`, `purge_cache_tag`.


