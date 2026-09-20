---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service"
sidebar_current: "docs-tencentcloud-resource-teo_inference_service"
description: |-
  Provides a resource to create a TEO inference service.
---

# tencentcloud_teo_inference_service

Provides a resource to create a TEO inference service.

## Example Usage

```hcl
resource "tencentcloud_teo_inference_service" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  name        = "tf-example-inference-service"
  listen_port = 8500
  description = "tf example inference service"

  containers {
    image_type = "TCR"

    tcr_repository_config {
      tcr_type    = "Enterprise"
      image       = "example.tencentcloudcr.com/ns/repo:v1"
      registry_id = "tcr-xxxxxxxx"
      region_name = "ap-guangzhou"
    }

    startup_command = "python app.py"

    environment_variables {
      key   = "MODEL_NAME"
      value = "example-model"
    }
  }

  resource_config {
    scaling_mode     = "Auto"
    hardware_spec_id = "hs-xxxxxxxx"
    concurrency      = 1

    auto_scaling_config {
      min_instance_count = 1

      scaling_policies {
        policy_name = "tf-example-policy"
        policy_type = "ScheduledScaling"

        scheduled_scaling_policy {
          scheduled_actions {
            cron_expression    = "0 0 * * *"
            min_instance_count = 2
          }

          effective_range {
            effective_type = "LongTerm"
          }

          time_zone = "Asia/Shanghai"
        }
      }
    }
  }

  affinity_config {
    switch        = "On"
    affinity_mode = "SessionId"

    session_id_affinity_config {
      source      = "Header"
      header_name = "EO-Infer-Session-Id"
    }
  }

  request_paths = ["/v1/predict"]
}
```

## Argument Reference

The following arguments are supported:

* `containers` - (Required, List) Container config list, only support 1 container currently.
* `listen_port` - (Required, Int) Model service listen port, 1-65535.
* `name` - (Required, String, ForceNew) Inference service name, not accepted by ModifyInferenceService, hence ForceNew.
* `resource_config` - (Required, List) Resource config of inference service.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `affinity_config` - (Optional, List) Affinity config, write-only (not read back from DescribeInferenceServices).
* `description` - (Optional, String) Description, up to 60 chars.
* `request_paths` - (Optional, List: [`String`]) Request path list, up to 20.

The `affinity_config` object supports the following:

* `switch` - (Required, String) Affinity switch: On or Off.
* `affinity_mode` - (Optional, String) Affinity mode, e.g. SessionId.
* `session_id_affinity_config` - (Optional, List) Session ID affinity config, required when affinity_mode is SessionId.

The `auto_scaling_config` object of `resource_config` supports the following:

* `min_instance_count` - (Required, Int) Min instance count.
* `scaling_policies` - (Optional, List) Scaling policy list.

The `containers` object supports the following:

* `image_type` - (Required, String) Image type, e.g. TCR.
* `environment_variables` - (Optional, List) Container environment variables.
* `startup_command` - (Optional, String) Container startup command.
* `tcr_repository_config` - (Optional, List) TCR repository config, required when image_type is TCR.

The `effective_range` object of `scheduled_scaling_policy` supports the following:

* `effective_type` - (Required, String) Effective type: LongTerm or Custom.
* `end_date` - (Optional, String) End date, required when effective_type is Custom.
* `start_date` - (Optional, String) Start date, required when effective_type is Custom.

The `environment_variables` object of `containers` supports the following:

* `key` - (Required, String) Variable name.
* `value` - (Optional, String) Variable value.

The `hardware_config` object of `resource_config` supports the following:

* `cpu_num` - (Optional, Float64) CPU core number.
* `disk_size` - (Optional, Int) Disk size in MB.
* `gpu_num` - (Optional, Float64) GPU card number, only effective on create.
* `mem_size` - (Optional, Int) Memory size in MB.

The `manual_instance_config` object of `resource_config` supports the following:

* `fixed_instance_count` - (Required, Int) Fixed instance count.

The `resource_config` object supports the following:

* `scaling_mode` - (Required, String) Scaling mode: Auto or Manual.
* `auto_scaling_config` - (Optional, List) Auto scaling config, required when scaling_mode is Auto.
* `concurrency` - (Optional, Int) Concurrency per instance, default 1.
* `hardware_config` - (Optional, List) Hardware config.
* `hardware_spec_id` - (Optional, String) Hardware spec unique ID, only effective on create.
* `hardware_spec` - (Optional, String) Hardware spec identifier, deprecated, use hardware_spec_id instead.
* `manual_instance_config` - (Optional, List) Manual instance config, required when scaling_mode is Manual.

The `scaling_policies` object of `auto_scaling_config` supports the following:

* `policy_name` - (Required, String) Policy name.
* `policy_type` - (Required, String) Policy type.
* `scheduled_scaling_policy` - (Optional, List) Scheduled scaling policy.

The `scheduled_actions` object of `scheduled_scaling_policy` supports the following:

* `cron_expression` - (Required, String) Cron expression.
* `min_instance_count` - (Required, Int) Min instance count.

The `scheduled_scaling_policy` object of `scaling_policies` supports the following:

* `effective_range` - (Required, List) Effective range.
* `scheduled_actions` - (Required, List) Scheduled action list.
* `time_zone` - (Optional, String) Time zone, e.g. UTC, Asia/Shanghai.

The `session_id_affinity_config` object of `affinity_config` supports the following:

* `header_name` - (Optional, String) Header name for session id.
* `source` - (Optional, String) Source of session id param, e.g. Header.

The `tcr_repository_config` object of `containers` supports the following:

* `image` - (Required, String) Image address.
* `tcr_type` - (Required, String) TCR service type: Personal or Enterprise.
* `region_name` - (Optional, String) Region name.
* `registry_id` - (Optional, String) Registry instance ID, required when tcr_type is Enterprise.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `current_instance_count` - Current running instance count.
* `inference_url` - Inference access URL.
* `scaling_status` - Scaling status.
* `service_id` - Inference service ID.
* `status` - Inference service status.
* `update_time` - Last update time.


## Import

TEO inference service can be imported using the zoneId#serviceId, e.g.

```
terraform import tencentcloud_teo_inference_service.example zone-3fkff38fyw8s#is-xxxxxxxx
```

