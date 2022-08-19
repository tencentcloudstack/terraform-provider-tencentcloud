---
subcategory: "TEM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_log_config"
sidebar_current: "docs-tencentcloud-resource-tem_log_config"
description: |-
  Provides a resource to create a tem logConfig
---

# tencentcloud_tem_log_config

Provides a resource to create a tem logConfig

## Example Usage

```hcl
resource "tencentcloud_tem_log_config" "logConfig" {
  environment_id = "en-o5edaepv"
  application_id = "app-3j29aa2p"
  workload_id    = resource.tencentcloud_tem_workload.workload.id
  name           = "terraform"
  logset_id      = "b5824781-8d5b-4029-a2f7-d03c37f72bdf"
  topic_id       = "5a85bb6d-8e41-4e04-b7bd-c05e04782f94"
  input_type     = "container_stdout"
  log_type       = "minimalist_log"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String, ForceNew) application ID.
* `environment_id` - (Required, String, ForceNew) environment ID.
* `input_type` - (Required, String) container_stdout or container_file.
* `log_type` - (Required, String) minimalist_log or multiline_log.
* `logset_id` - (Required, String) logset.
* `name` - (Required, String, ForceNew) appConfig name.
* `topic_id` - (Required, String) topic.
* `workload_id` - (Required, String, ForceNew) application ID, which is combined by environment ID and application ID, like `en-o5edaepv#app-3j29aa2p`.
* `beginning_regex` - (Optional, String) regex pattern.
* `file_pattern` - (Optional, String) file name pattern if container_file.
* `log_path` - (Optional, String) directory if container_file.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tem logConfig can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_log_config.logConfig environmentId#applicationId#name
```

