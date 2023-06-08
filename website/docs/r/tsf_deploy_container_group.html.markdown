---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_deploy_container_group"
sidebar_current: "docs-tencentcloud-resource-tsf_deploy_container_group"
description: |-
  Provides a resource to create a tsf deploy_container_group
---

# tencentcloud_tsf_deploy_container_group

Provides a resource to create a tsf deploy_container_group

## Example Usage

```hcl
resource "tencentcloud_tsf_deploy_container_group" "deploy_container_group" {
  group_id          = "group-yqml6w3a"
  cpu_request       = "0.25"
  mem_request       = "640"
  server            = "ccr.ccs.tencentyun.com"
  reponame          = "tsf_100011913960/terraform"
  tag_name          = "terraform-only-1"
  do_not_start      = false
  instance_num      = 1
  update_type       = 1
  update_ivl        = 10
  mem_limit         = "1280"
  cpu_limit         = "0.5"
  agent_cpu_request = "0.1"
  agent_cpu_limit   = "0.2"
  agent_mem_request = "125"
  agent_mem_limit   = "400"
  max_surge         = "25%"
  max_unavailable   = "0"
  service_setting {
    access_type = 1
    protocol_ports {
      protocol    = "TCP"
      port        = 18081
      target_port = 18081
      node_port   = 30001
    }
    subnet_id                        = ""
    disable_service                  = false
    headless_service                 = false
    allow_delete_service             = true
    open_session_affinity            = false
    session_affinity_timeout_seconds = 10800

  }
  health_check_settings {
    readiness_probe {
      action_type           = "TCP"
      initial_delay_seconds = 0
      timeout_seconds       = 3
      period_seconds        = 30
      success_threshold     = 1
      failure_threshold     = 3
      scheme                = "HTTP"
      port                  = 80
      path                  = "/"
      type                  = "TSF_DEFAULT"
    }
  }
  scheduling_strategy {
    type = "NONE"
  }
  deploy_agent = true
  repo_type    = "personal"
  volume_clean = false
  jvm_opts     = "-Xms128m -Xmx512m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=512m"
  warmup_setting {
    enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) group Id.
* `instance_num` - (Required, Int) instance number.
* `tag_name` - (Required, String, ForceNew) image version name, v1.
* `agent_cpu_limit` - (Optional, String, ForceNew) The maximum number of CPU cores allocated to the agent container corresponds to the limit field in Kubernetes.
* `agent_cpu_request` - (Optional, String, ForceNew) The number of CPU cores allocated to the agent container corresponds to the request field in Kubernetes.
* `agent_mem_limit` - (Optional, String, ForceNew) The maximum amount of memory in MiB allocated to the agent container corresponds to the &amp;#39;limit&amp;#39; field in Kubernetes.
* `agent_mem_request` - (Optional, String, ForceNew) The amount of memory in MiB allocated to the agent container corresponds to the request field in Kubernetes.
* `agent_profile_list` - (Optional, List, ForceNew) javaagent info: SERVICE_AGENT/OT_AGENT.
* `cpu_limit` - (Optional, String, ForceNew) The maximum number of CPU cores for the business container, corresponding to the limit in K8S. If not specified, it defaults to twice the request.
* `cpu_request` - (Optional, String, ForceNew) The number of CPU cores allocated to the business container, corresponding to the request in K8S. The default value is 0.25.
* `deploy_agent` - (Optional, Bool, ForceNew) Whether to deploy the agent container. If this parameter is not specified, the agent container will not be deployed by default.
* `do_not_start` - (Optional, Bool, ForceNew) Not start right away.
* `envs` - (Optional, List, ForceNew) The environment variables that the application runs in the deployment group. If this parameter is not specified, no additional environment variables are set by default.
* `health_check_settings` - (Optional, List, ForceNew) The configuration information for health checks. If this parameter is not specified, the health check is not set by default.
* `incremental_deployment` - (Optional, Bool, ForceNew) Whether to perform incremental deployment. The default value is false, which means full update.
* `istio_cpu_limit` - (Optional, String, ForceNew) The maximum amount of CPU cores allocated to the istio proxy container corresponds to the &amp;#39;limit&amp;#39; field in Kubernetes.
* `istio_cpu_request` - (Optional, String, ForceNew) The number of CPU cores allocated to the istio proxy container corresponds to the &amp;#39;request&amp;#39; field in Kubernetes.
* `istio_mem_limit` - (Optional, String, ForceNew) The maximum amount of memory in MiB allocated to the agent container corresponds to the request field in Kubernetes.
* `istio_mem_request` - (Optional, String, ForceNew) The amount of memory in MiB allocated to the agent container corresponds to the request field in Kubernetes.
* `jvm_opts` - (Optional, String, ForceNew) jvm options.
* `max_surge` - (Optional, String, ForceNew) MaxSurge parameter in Kubernetes rolling update strategy.
* `max_unavailable` - (Optional, String, ForceNew) MaxUnavailable parameter in Kubernetes rolling update strategy.
* `mem_limit` - (Optional, String, ForceNew) The maximum memory size in MiB for the business container, corresponding to the limit in K8S. If not specified, it defaults to twice the request.
* `mem_request` - (Optional, String, ForceNew) The amount of memory in MiB allocated to the business container, corresponding to the request in K8S. The default value is 640 MiB.
* `repo_name` - (Optional, String, ForceNew) (Priority use) New image name, such as /tsf/nginx.
* `repo_type` - (Optional, String, ForceNew) repo type, tcr or leave it blank.
* `reponame` - (Optional, String, ForceNew) old image name, eg: /tsf/server.
* `scheduling_strategy` - (Optional, List, ForceNew) Node scheduling strategy. If this parameter is not specified, the node scheduling strategy will not be used by default.
* `server` - (Optional, String, ForceNew) image server.
* `service_setting` - (Optional, List, ForceNew) Network settings for container deployment groups.
* `update_ivl` - (Optional, Int, ForceNew) update Interval, is required when rolling update.
* `update_type` - (Optional, Int, ForceNew) Update method: 0 for fast update, 1 for rolling update.
* `volume_clean` - (Optional, Bool, ForceNew) Whether to clear the volume information. Default is false.
* `volume_info_list` - (Optional, List, ForceNew) Volume information, as a list.
* `volume_mount_info_list` - (Optional, List, ForceNew) Volume mount point information, list type.
* `warmup_setting` - (Optional, List, ForceNew) warmup setting.

The `agent_profile_list` object supports the following:

* `agent_type` - (Optional, String) Agent type.
* `agent_version` - (Optional, String) Agent version.

The `envs` object supports the following:

* `name` - (Required, String) env param name.
* `value_from` - (Optional, List) Kubernetes ValueFrom configuration. Note: This field may return null, indicating that no valid values can be obtained.
* `value` - (Optional, String) value of env.

The `field_ref` object supports the following:

* `field_path` - (Optional, String) The FieldPath configuration of Kubernetes. Note: This field may return null, indicating that no valid values can be obtained.

The `health_check_settings` object supports the following:

* `liveness_probe` - (Optional, List) Liveness probe. Note: This field may return null, indicating that no valid values can be obtained.
* `readiness_probe` - (Optional, List) Readiness health check. Note: This field may return null, indicating that no valid values can be obtained.

The `liveness_probe` object supports the following:

* `action_type` - (Required, String) The health check method. HTTP: checks through an HTTP interface; CMD: checks by executing a command; TCP: checks by establishing a TCP connection. Note: This field may return null, indicating that no valid values can be obtained.
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

The `protocol_ports` object supports the following:

* `port` - (Required, Int) port.
* `protocol` - (Required, String) TCP or UDP.
* `target_port` - (Required, Int) container port.
* `node_port` - (Optional, Int) node port.

The `readiness_probe` object supports the following:

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

The `resource_field_ref` object supports the following:

* `resource` - (Optional, String) The Resource configuration of Kubernetes. Note: This field may return null, indicating that no valid values can be obtained.

The `scheduling_strategy` object supports the following:

* `type` - (Required, String) NONE: Do not use scheduling strategy; CROSS_AZ: Deploy across availability zones. Note: This field may return null, indicating that no valid values can be obtained.

The `service_setting` object supports the following:

* `access_type` - (Required, Int) 0: Public network, 1: Access within the cluster, 2: NodePort, 3: Access within VPC. Note: This field may return null, indicating that no valid values can be obtained.
* `protocol_ports` - (Required, List) Container port mapping. Note: This field may return null, indicating that no valid values can be obtained.
* `subnet_id` - (Required, String) subnet Id.
* `allow_delete_service` - (Optional, Bool) When set to true and DisableService is also true, the previously created service will be deleted. Please use with caution. Note: This field may return null, indicating that no valid values can be obtained.
* `disable_service` - (Optional, Bool) Whether to create a Kubernetes service. The default value is false. Note: This field may return null, indicating that no valid values can be obtained.
* `headless_service` - (Optional, Bool) Whether the service is of headless type. Note: This field may return null, indicating that no valid values can be obtained.
* `open_session_affinity` - (Optional, Bool) Enable session affinity. true means enabled, false means disabled. The default value is false. Note: This field may return null, indicating that no valid values can be obtained.
* `session_affinity_timeout_seconds` - (Optional, Int) Session affinity session time. The default value is 10800. Note: This field may return null, indicating that no valid values can be obtained.

The `value_from` object supports the following:

* `field_ref` - (Optional, List) The FieldRef configuration of Kubernetes env. Note: This field may return null, indicating that no valid values can be obtained.
* `resource_field_ref` - (Optional, List) The ResourceFieldRef configuration of Kubernetes env. Note: This field may return null, indicating that no valid values can be obtained.

The `volume_info_list` object supports the following:

* `volume_name` - (Required, String) volume name.
* `volume_type` - (Required, String) volume type.
* `volume_config` - (Optional, String) volume config.

The `volume_mount_info_list` object supports the following:

* `volume_mount_name` - (Required, String) mount volume name.
* `volume_mount_path` - (Required, String) mount path.
* `read_or_write` - (Optional, String) Read and write access mode. 1: Read-only. 2: Read-write.
* `volume_mount_sub_path` - (Optional, String) mount subPath.

The `warmup_setting` object supports the following:

* `curvature` - (Optional, Int) Preheating curvature, with a value between 1 and 5.
* `enabled_protection` - (Optional, Bool) Whether to enable preheating protection. If protection is enabled and more than 50% of nodes are in preheating state, preheating will be aborted.
* `enabled` - (Optional, Bool) Whether to enable preheating.
* `warmup_time` - (Optional, Int) warmup time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



