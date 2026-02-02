Use this data source to query detailed information of VOD sub-applications.

Example Usage

Query all sub-applications

```hcl
data "tencentcloud_vod_sub_applications" "all" {}

output "app_list" {
  value = data.tencentcloud_vod_sub_applications.all.sub_application_info_set
}
```

Filter by application name

```hcl
data "tencentcloud_vod_sub_applications" "by_name" {
  name = "MyVideoApp"
}

output "app_id" {
  value = data.tencentcloud_vod_sub_applications.by_name.sub_application_info_set[0].sub_app_id
}
```

Filter by tags

```hcl
data "tencentcloud_vod_sub_applications" "by_tags" {
  tags = {
    Environment = "Production"
    Team        = "VideoTeam"
  }
}

output "production_apps" {
  value = data.tencentcloud_vod_sub_applications.by_tags.sub_application_info_set[*].sub_app_id_name
}
```

Reference sub-application in other resources

```hcl
data "tencentcloud_vod_sub_applications" "existing" {
  name = "ProductionApp"
}

resource "tencentcloud_vod_super_player_config" "config" {
  sub_app_id = data.tencentcloud_vod_sub_applications.existing.sub_application_info_set[0].sub_app_id
  name       = "player-config"
  drm_switch = false
  # ... other configuration
}
```

Export results to file

```hcl
data "tencentcloud_vod_sub_applications" "export" {
  result_output_file = "/tmp/vod_apps.json"
}
```
