---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_aigc_quota"
sidebar_current: "docs-tencentcloud-resource-vod_aigc_quota"
description: |-
  Provides a resource to manage a VOD AIGC Quota for a specific sub application.
---

# tencentcloud_vod_aigc_quota

Provides a resource to manage a VOD AIGC Quota for a specific sub application.

## Example Usage

```hcl
resource "tencentcloud_vod_aigc_quota" "image" {
  sub_app_id  = 251006666
  quota_type  = "Image"
  quota_limit = 100
}

resource "tencentcloud_vod_aigc_quota" "video" {
  sub_app_id  = 251006666
  quota_type  = "Video"
  quota_limit = 3600
}

resource "tencentcloud_vod_aigc_quota" "text" {
  sub_app_id  = 251006666
  quota_type  = "Text"
  quota_limit = 5000
  api_token   = "my-api-token"
}
```

## Argument Reference

The following arguments are supported:

* `quota_limit` - (Required, Int) Quota limit value. Unit: when `quota_type` is `Image`, count by images; when `quota_type` is `Video`, count by seconds; when `quota_type` is `Text`, count by tokens.
* `quota_type` - (Required, String, ForceNew) Quota type. Valid values: `Image` (AIGC image generation), `Video` (AIGC video generation), `Text` (AIGC text generation).
* `sub_app_id` - (Required, Int, ForceNew) The VOD [sub application](https://intl.cloud.tencent.com/document/product/266/14574) ID. Users who activated VOD service after December 25, 2023 must specify the application ID.
* `api_token` - (Optional, String, ForceNew) API token for quota restriction. Only meaningful when `quota_type` is `Text`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `usage` - Current usage amount. Unit: when `quota_type` is `Image`, count by images; when `quota_type` is `Video`, count by seconds; when `quota_type` is `Text`, count by tokens.


## Import

VOD AIGC Quota can be imported using the composite id `sub_app_id#quota_type#api_token`, e.g.

```
$ terraform import tencentcloud_vod_aigc_quota.image 251006666#Image#
$ terraform import tencentcloud_vod_aigc_quota.text 251006666#Text#my-api-token
```

