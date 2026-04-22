Provides a resource to import TEO site configuration.

Example Usage

```hcl
resource "tencentcloud_teo_import_zone_config_operation" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  content = jsonencode({
    GlobalAccelerate = {
      AccelerateMainland = {
        Switch = "off"
      }
      OverseeAccelerate = {
        Switch = "off"
      }
    }
    RuleEngine = [
      {
        Rules = [
          {
            Actions = [
              {
                NormalAction = {
                  Action = "UpstreamUrlRewrite"
                  Parameters = {
                    Type  = "all"
                    Value = "/redirect-path"
                  }
                }
              }
            ]
          }
        ]
      }
    ]
  })
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Site ID.
* `content` - (Required, ForceNew) The configuration content to import. It must be in JSON format and encoded in UTF-8. You can obtain the configuration content via the ExportZoneConfig API.

Attribute Reference

The following attributes are exported:

* `task_id` - The task ID of the import configuration operation.
* `status` - The import task status. Valid values: success, failure, doing.
* `message` - The status message of the import task. When the configuration import fails, you can view the failure reason through this field.
* `import_time` - The start time of the import task.
* `finish_time` - The end time of the import task.
