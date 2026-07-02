Provides a resource to manage a VOD AIGC Quota for a specific sub application.

Example Usage

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

Import

VOD AIGC Quota can be imported using the composite id `sub_app_id#quota_type#api_token`, e.g.

```
$ terraform import tencentcloud_vod_aigc_quota.image 251006666#Image#
$ terraform import tencentcloud_vod_aigc_quota.text 251006666#Text#my-api-token
```