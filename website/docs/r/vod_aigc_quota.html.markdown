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

### If quota type is Image

```hcl
# create sub application
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_aigc_quota" "example" {
  sub_app_id  = tencentcloud_vod_sub_application.example.sub_app_id
  quota_type  = "Image"
  quota_limit = 100
}
```

### If quota type is Video

```hcl
# create sub application
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_aigc_quota" "example" {
  sub_app_id  = tencentcloud_vod_sub_application.example.sub_app_id
  quota_type  = "Video"
  quota_limit = 200
}
```

### If quota type is Text

```hcl
# create sub application
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "this is sub application"
}

# create api token
resource "tencentcloud_vod_aigc_api_token" "example" {
  sub_app_id = tencentcloud_vod_sub_application.example.sub_app_id
}

resource "tencentcloud_vod_aigc_quota" "example" {
  sub_app_id  = tencentcloud_vod_sub_application.example.sub_app_id
  quota_type  = "Text"
  quota_limit = 50
  api_token   = tencentcloud_vod_aigc_api_token.example.api_token
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
# If quota type is Image/Video
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Image
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Video

# If quota type is Text
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Text#<YOUR TOKEN>
```

