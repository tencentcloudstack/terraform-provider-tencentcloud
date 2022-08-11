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
  environment_id = "en-853mggjm"
  application_id = "app-3j29aa2p"
  name           = "terraform"
  logset_id      = "b5824781-8d5b-4029-a2f7-d03c37f72bdf"
  topic_id       = "a21a488d-d28f-4ac3-8044-bdf8c91b49f2"
  input_type     = "container_stdout"
  log_type       = "minimalist_log"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String) application ID.
* `environment_id` - (Required, String) environment ID.
* `input_type` - (Required, String) container_stdout or container_file.
* `log_type` - (Required, String) minimalist_log or multiline_log.
* `logset_id` - (Required, String) logset.
* `name` - (Required, String) appConfig name.
* `topic_id` - (Required, String) topic.
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

