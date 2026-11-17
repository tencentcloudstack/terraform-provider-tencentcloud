Provides a resource to manage a VOD AIGC Quota for a specific sub application.

Example Usage

If quota type is Image

```hcl
# create sub application
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_aigc_quota" "example" {
  sub_app_id         = tencentcloud_vod_sub_application.example.sub_app_id
  quota_type         = "Image"
  quota_limit        = 100
  quota_period       = "Month"
  reset_day          = 1
  model_quota_switch = "ON"
  model_quotas {
    model_name    = "model1"
    model_version = "v1"
    quota_limit   = 100
  }

  model_quotas {
    model_name    = "model2"
    model_version = "v2"
    quota_limit   = 50
  }

  model_quotas {
    model_name    = "model3"
    model_version = "v3"
    quota_limit   = 80
  }
}
```

If quota type is Video

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

If quota type is Text

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

Import

VOD AIGC Quota can be imported using the composite id `sub_app_id#quota_type#api_token`, e.g.

```
# If quota type is Image/Video
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Image
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Video

# If quota type is Text
terraform import tencentcloud_vod_aigc_quota.example 1500066373#Text#<YOUR TOKEN>
```
