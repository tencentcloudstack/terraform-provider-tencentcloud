---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_prefetch_origin_limit"
sidebar_current: "docs-tencentcloud-resource-teo_prefetch_origin_limit"
description: |-
  Provides a resource to create a TEO prefetch origin limit config.
---

# tencentcloud_teo_prefetch_origin_limit

Provides a resource to create a TEO prefetch origin limit config.

## Example Usage

### Set prefetch origin bandwidth limit for overseas area

```hcl
resource "tencentcloud_teo_prefetch_origin_limit" "example" {
  zone_id     = "zone-3edjdliiw3he"
  domain_name = "example.com"
  area        = "Overseas"
  bandwidth   = 200
}
```

### Set prefetch origin bandwidth limit for Mainland China area

```hcl
resource "tencentcloud_teo_prefetch_origin_limit" "example" {
  zone_id     = "zone-3edjdliiw3he"
  domain_name = "example.com"
  area        = "MainlandChina"
  bandwidth   = 500
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String, ForceNew) Acceleration area for prefetch origin limit. Valid values: `Overseas`, `MainlandChina`.
* `bandwidth` - (Required, Int) Prefetch origin bandwidth limit. Value range: 100-100000, in Mbps.
* `domain_name` - (Required, String, ForceNew) Accelerated domain name.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribePrefetchOriginLimit](https://cloud.tencent.com/document/api/1552/126842) - 查询预热回源限速限制
* Create: [ModifyPrefetchOriginLimit](https://cloud.tencent.com/document/api/1552/126841) - 配置预热回源限速限制
* Read: [DescribePrefetchOriginLimit](https://cloud.tencent.com/document/api/1552/126842) - 查询预热回源限速限制
* Update: [ModifyPrefetchOriginLimit](https://cloud.tencent.com/document/api/1552/126841) - 配置预热回源限速限制
* Delete: [ModifyPrefetchOriginLimit](https://cloud.tencent.com/document/api/1552/126841) - 配置预热回源限速限制

## Import

TEO prefetch origin limit config can be imported using the composite ID format `zone_id#domain_name#area`, e.g.

```
terraform import tencentcloud_teo_prefetch_origin_limit.example zone-3edjdliiw3he#example.com#Overseas
```

