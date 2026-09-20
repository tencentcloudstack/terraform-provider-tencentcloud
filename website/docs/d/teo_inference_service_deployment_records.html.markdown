---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_deployment_records"
sidebar_current: "docs-tencentcloud-datasource-teo_inference_service_deployment_records"
description: |-
  Use this data source to query TEO inference service deployment records, such as the operation type, status, duration and configuration snapshot of each deployment.
---

# tencentcloud_teo_inference_service_deployment_records

Use this data source to query TEO inference service deployment records, such as the operation type, status, duration and configuration snapshot of each deployment.

## Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_records" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-xxxxxxxxxxxx"
  sort_by    = "create-time"
  sort_order = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) Inference service ID.
* `zone_id` - (Required, String) Site ID.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sort field, value: create-time. Default value: create-time.
* `sort_order` - (Optional, String) Sort method, value: asc / desc. Default value: desc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_set` - Inference service deployment record list.
  * `active_status` - Whether the deployment configuration is the currently effective configuration. Value: active/inactive.
  * `create_time` - Deployment initiation time, using ISO date format.
  * `duration` - Deployment duration, unit: seconds.
  * `inference_service_config` - Inference service deployment configuration.
    * `affinity_config` - Affinity configuration of the inference service.
      * `affinity_mode` - Inference service affinity mode. Value: SessionId.
      * `session_id_affinity_config` - Session ID affinity configuration.
        * `header_name` - The request header name used to pass the session ID.
        * `source` - Location where the session ID parameter is passed. Value: Header.
      * `switch` - Inference service affinity switch. Value: On/Off.
    * `containers` - Container configuration list of the inference service.
      * `environment_variables` - Container runtime environment variable list.
        * `key` - Variable name.
        * `value` - Variable value.
      * `image_type` - Image type. Value: TCR.
      * `startup_command` - Container startup command.
      * `tcr_repository_config` - TCR image repository configuration.
        * `image` - Image address.
        * `region_name` - Region name.
        * `registry_id` - Image repository instance ID. Required when TCRType = Enterprise.
        * `tcr_type` - TCR service type. Value: Personal/Enterprise.
    * `listen_port` - Port that the model service needs to listen on.
    * `request_paths` - Request path list of the inference service.
    * `resource_config` - Resource configuration of the inference service.
      * `auto_scaling_config` - Automatic scaling configuration.
        * `min_instance_count` - Minimum number of instances.
        * `scaling_policies` - Scaling policy list.
          * `policy_name` - Policy name.
          * `policy_type` - Policy type. Value: ScheduledScaling.
          * `scheduled_scaling_policy` - Scheduled scaling configuration.
            * `scheduled_actions` - Scheduled scaling action list.
              * `cron_expression` - Cron expression.
              * `min_instance_count` - Minimum instance count after the scheduled scaling action is triggered.
      * `concurrency` - Concurrency of a single instance.
      * `hardware_config` - Hardware configuration.
        * `cpu_num` - Number of CPU cores allocated to a single instance of the inference service.
        * `disk_size` - Temporary disk size allocated to a single instance of the inference service. Unit: MB.
        * `gpu_num` - Number of GPU cards allocated to a single instance of the inference service.
        * `mem_size` - Memory size allocated to a single instance of the inference service. Unit: MB.
      * `hardware_spec_id` - Hardware specification unique identifier ID.
      * `hardware_spec` - Hardware specification identifier. Deprecated.
      * `manual_instance_config` - Manual instance configuration.
        * `fixed_instance_count` - Fixed instance count.
      * `scaling_mode` - Scaling mode. Value: Auto/Manual.
  * `operation` - Deployment operation type. Value: create/update/resume/stop.
  * `record_id` - Deployment record ID.
  * `status` - Deployment status. Value: processing/succeeded/failed.


