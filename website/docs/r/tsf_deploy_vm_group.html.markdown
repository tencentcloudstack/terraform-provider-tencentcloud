---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_deploy_vm_group"
sidebar_current: "docs-tencentcloud-resource-tsf_deploy_vm_group"
description: |-
  Provides a resource to create a tsf deploy_vm_group
---

# tencentcloud_tsf_deploy_vm_group

Provides a resource to create a tsf deploy_vm_group

## Example Usage

```hcl
resource "tencentcloud_tsf_deploy_vm_group" "deploy_vm_group" {
  group_id            = "group-vzd97zpy"
  pkg_id              = "pkg-131bc1d3"
  startup_parameters  = "-Xms128m -Xmx512m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=512m"
  deploy_desc         = "deploy test"
  force_start         = false
  enable_health_check = true
  health_check_settings {
    readiness_probe {
      action_type           = "HTTP"
      initial_delay_seconds = 10
      timeout_seconds       = 2
      period_seconds        = 10
      success_threshold     = 1
      failure_threshold     = 3
      scheme                = "HTTP"
      port                  = "80"
      path                  = "/"
    }
  }
  update_type = 0
  jdk_name    = "konaJDK"
  jdk_version = "8"

}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) group id.
* `pkg_id` - (Required, String, ForceNew) program package ID.
* `agent_profile_list` - (Optional, List, ForceNew) javaagent info: SERVICE_AGENT/OT_AGENT.
* `deploy_batch` - (Optional, Set: [`Float64`], ForceNew) The ratio of instances participating in each batch during rolling release.
* `deploy_beta_enable` - (Optional, Bool, ForceNew) Whether to enable beta batch.
* `deploy_desc` - (Optional, String, ForceNew) group description.
* `deploy_exe_mode` - (Optional, String, ForceNew) The execution method of rolling release.
* `deploy_wait_time` - (Optional, Int, ForceNew) The time interval for each batch during rolling release.
* `enable_health_check` - (Optional, Bool, ForceNew) Whether to enable health check.
* `force_start` - (Optional, Bool, ForceNew) Whether to allow forced start.
* `health_check_settings` - (Optional, List, ForceNew) When enabling health check, configure the health check settings.
* `incremental_deployment` - (Optional, Bool, ForceNew) Whether to perform incremental deployment. The default value is false, which means full update.
* `jdk_name` - (Optional, String, ForceNew) JDK name: konaJDK or openJDK.
* `jdk_version` - (Optional, String, ForceNew) JDK version: 8 or 11(openJDK only support 8).
* `start_script` - (Optional, String, ForceNew) The base64-encoded startup script.
* `startup_parameters` - (Optional, String, ForceNew) start args of group.
* `stop_script` - (Optional, String, ForceNew) The base64-encoded stop script.
* `update_type` - (Optional, Int, ForceNew) Update method: 0 for fast update, 1 for rolling update.
* `warmup_setting` - (Optional, List, ForceNew) warmup setting.

The `agent_profile_list` object supports the following:

* `agent_type` - (Optional, String) Agent type.
* `agent_version` - (Optional, String) Agent version.

The `health_check_settings` object supports the following:

* `liveness_probe` - (Optional, List) Survival health check. Note: This field may return null, indicating that no valid value was found.
* `readiness_probe` - (Optional, List) Readiness health check. Note: This field may return null, indicating that no valid values can be obtained.

The `liveness_probe` object of `health_check_settings` supports the following:

* `action_type` - (Required, String) Health check method. HTTP: check through HTTP interface; CMD: check through executing command; TCP: check through establishing TCP connection. Note: This field may return null, indicating that no valid value was found.
* `command` - (Optional, Set) The command to be executed for command health checks. Note: This field may return null, indicating that no valid values can be obtained.
* `failure_threshold` - (Optional, Int) The number of consecutive successful health checks required for the backend container to transition from success to failure. Note: This field may return null, indicating that no valid values can be obtained.
* `initial_delay_seconds` - (Optional, Int) The time delay for the container to start the health check. Note: This field may return null, indicating that no valid values can be obtained.
* `path` - (Optional, String) The request path for HTTP health checks. Note: This field may return null, indicating that no valid values can be obtained.
* `period_seconds` - (Optional, Int) The time interval for performing health checks. Note: This field may return null, indicating that no valid values can be obtained.
* `port` - (Optional, Int) The port used for health checks, ranging from 1 to 65535. Note: This field may return null, indicating that no valid values can be obtained.
* `scheme` - (Optional, String) The protocol used for HTTP health checks. HTTP and HTTPS are supported. Note: This field may return null, indicating that no valid values can be obtained.
* `success_threshold` - (Optional, Int) The number of consecutive successful health checks required for the backend container to transition from failure to success. Note: This field may return null, indicating that no valid values can be obtained.
* `timeout_seconds` - (Optional, Int) The maximum timeout period for each health check response. Note: This field may return null, indicating that no valid values can be obtained.
* `type` - (Optional, String) The type of readiness probe. TSF_DEFAULT represents the default readiness probe of TSF, while K8S_NATIVE represents the native readiness probe of Kubernetes. If this field is not specified, the native readiness probe of Kubernetes is used by default. Note: This field may return null, indicating that no valid values can be obtained.

The `readiness_probe` object of `health_check_settings` supports the following:

* `action_type` - (Required, String) The health check method. HTTP indicates checking through an HTTP interface, CMD indicates checking through executing a command, and TCP indicates checking through establishing a TCP connection. Note: This field may return null, indicating that no valid values can be obtained.
* `command` - (Optional, Set) The command to be executed for command check. Note: This field may return null, indicating that no valid values can be obtained.
* `failure_threshold` - (Optional, Int) The number of consecutive successful health checks required for the backend container to transition from success to failure. Note: This field may return null, indicating that no valid values can be obtained.
* `initial_delay_seconds` - (Optional, Int) The time to delay the start of the container health check. Note: This field may return null, indicating that no valid values can be obtained.
* `path` - (Optional, String) The request path for HTTP health checks. Note: This field may return null, indicating that no valid values can be obtained.
* `period_seconds` - (Optional, Int) The time interval for performing health checks. Note: This field may return null, indicating that no valid values can be obtained.
* `port` - (Optional, Int) The port used for health checks, ranging from 1 to 65535. Note: This field may return null, indicating that no valid values can be obtained.
* `scheme` - (Optional, String) The protocol used for HTTP health checks. HTTP and HTTPS are supported. Note: This field may return null, indicating that no valid values can be obtained.
* `success_threshold` - (Optional, Int) The number of consecutive successful health checks required for the backend container to transition from failure to success. Note: This field may return null, indicating that no valid values can be obtained.
* `timeout_seconds` - (Optional, Int) The maximum timeout period for each health check response. Note: This field may return null, indicating that no valid values can be obtained.
* `type` - (Optional, String) The type of readiness probe. TSF_DEFAULT represents the default readiness probe of TSF, while K8S_NATIVE represents the native readiness probe of Kubernetes. If this field is not specified, the native readiness probe of Kubernetes is used by default. Note: This field may return null, indicating that no valid values can be obtained.

The `warmup_setting` object supports the following:

* `curvature` - (Optional, Int) Preheating curvature, with a value between 1 and 5.
* `enabled_protection` - (Optional, Bool) Whether to enable preheating protection. If protection is enabled and more than 50% of nodes are in preheating state, preheating will be aborted.
* `enabled` - (Optional, Bool) Whether to enable preheating.
* `warmup_time` - (Optional, Int) warmup time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



