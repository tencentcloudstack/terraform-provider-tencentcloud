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
  zone_id = "zone-2l1zk57u"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `prefetch_quota` - Prefetch quota list.
  * `batch` - Upper limit of quota for each batch submission.
  * `daily_available` - Remaining daily submission quota.
  * `daily` - Upper limit of daily submission quota.
  * `type` - Cache purge type. Values: `purge_prefix` (prefix purge), `purge_url` (URL purge), `purge_host` (hostname purge), `purge_all` (purge all cache), `purge_cache_tag` (cache tag purge), `prefetch_url` (URL prefetch).
* `purge_quota` - Purge quota list.
  * `batch` - Upper limit of quota for each batch submission.
  * `daily_available` - Remaining daily submission quota.
  * `daily` - Upper limit of daily submission quota.
  * `type` - Cache purge type. Values: `purge_prefix` (prefix purge), `purge_url` (URL purge), `purge_host` (hostname purge), `purge_all` (purge all cache), `purge_cache_tag` (cache tag purge), `prefetch_url` (URL prefetch).


