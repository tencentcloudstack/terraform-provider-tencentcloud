Provides a resource to manage a VOD AIGC Quota for a specific sub application.

Example Usage

If quota type is Image

```hcl
resource "tencentcloud_vod_aigc_quota" "example" {
  sub_app_id  = 1500066373
  quota_type  = "Image"
  quota_limit = 100
}
```

If quota type is Video

```hcl
resource "tencentcloud_vod_aigc_quota" "example" {
  sub_app_id  = 1500066373
  quota_type  = "Video"
  quota_limit = 200
}
```

If quota type is Text

```hcl
# create api token
resource "tencentcloud_vod_aigc_api_token" "example" {
  sub_app_id = 1500066373
}

# set quota limit
resource "tencentcloud_vod_aigc_quota" "example" {
  sub_app_id  = 1500066373
  quota_type  = "Text"
  quota_limit = 50
  api_token   = tencentcloud_vod_aigc_api_token.example.api_token
}
```

Import

VOD AIGC Quota can be imported using the composite id `sub_app_id#quota_type#api_token`, e.g.

```
# If quota type is Image/Video
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Image
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Video

# If quota type is Text
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Text#<YOUR TOKEN>
```
