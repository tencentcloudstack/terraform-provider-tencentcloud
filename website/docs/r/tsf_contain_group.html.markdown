---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_contain_group"
sidebar_current: "docs-tencentcloud-resource-tsf_contain_group"
description: |-
  Provides a resource to create a tsf contain_group
---

# tencentcloud_tsf_contain_group

Provides a resource to create a tsf contain_group

## Example Usage

```hcl
resource "tencentcloud_tsf_contain_group" "contain_group" {
  access_type         = 0
  application_id      = "application-y5r4nejv"
  cluster_id          = "cls-2yu5kxr8"
  cpu_limit           = "0.5"
  cpu_request         = "0.25"
  group_name          = "terraform-test"
  group_resource_type = "DEF"
  instance_num        = 1
  mem_limit           = "1280"
  mem_request         = "640"
  namespace_id        = "namespace-ydlezgxa"
  update_ivl          = 10
  update_type         = 1

  protocol_ports {
    node_port   = 0
    port        = 333
    protocol    = "TCP"
    target_port = 333
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_type` - (Required, Int) 0: public network 1: access within the cluster 2: NodePort.
* `application_id` - (Required, String) The application ID to which the group belongs.
* `cluster_id` - (Required, String) Cluster ID.
* `group_name` - (Required, String) Group name field, length 1~60, beginning with a letter or underscore, can contain alphanumeric underscore.
* `instance_num` - (Required, Int) number of instances.
* `namespace_id` - (Required, String) ID of the namespace to which the group belongs.
* `protocol_ports` - (Required, List) Protocol Ports array.
* `agent_cpu_limit` - (Optional, String) The maximum number of CPU cores for the agent container, corresponding to the limit of K8S.
* `agent_cpu_request` - (Optional, String) The number of CPU cores allocated by the agent container, corresponding to the K8S request.
* `agent_mem_limit` - (Optional, String) The maximum memory MiB of the agent container, corresponding to the limit of K8S.
* `agent_mem_request` - (Optional, String) The number of memory MiB allocated by the agent container, corresponding to the K8S request.
* `cpu_limit` - (Optional, String) The maximum number of allocated CPU cores, corresponding to the K8S limit.
* `cpu_request` - (Optional, String) Initially allocated CPU cores, corresponding to K8S request.
* `group_comment` - (Optional, String) Group remarks field, the length should not exceed 200 characters.
* `group_resource_type` - (Optional, String) Deployment Group Resource Type.
* `istio_cpu_limit` - (Optional, String) The maximum number of CPU cores for the istioproxy container corresponds to the limit of K8S.
* `istio_cpu_request` - (Optional, String) The number of CPU cores allocated by the istioproxy container, corresponding to the K8S request.
* `istio_mem_limit` - (Optional, String) The maximum memory MiB of the istioproxy container corresponds to the limit of K8S.
* `istio_mem_request` - (Optional, String) The number of memory MiB allocated by the istioproxy container, corresponding to the K8S request.
* `mem_limit` - (Optional, String) Maximum allocated memory MiB, corresponding to K8S limit.
* `mem_request` - (Optional, String) Initially allocated memory MiB, corresponding to K8S request.
* `subnet_id` - (Optional, String) subnet ID.
* `update_ivl` - (Optional, Int) Rolling update is required, update interval.
* `update_type` - (Optional, Int) Update method: 0: fast update 1: rolling update.

The `protocol_ports` object supports the following:

* `port` - (Required, Int) service port.
* `protocol` - (Required, String) TCP UDP.
* `target_port` - (Required, Int) container port.
* `node_port` - (Optional, Int) host port.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `application_name` - Application Name.
* `application_type` - App types.
* `cluster_ip` - Service ip.
* `cluster_name` - cluster name.
* `create_time` - creation time.
* `current_num` - Total number of instances launched.
* `envs` - environment variable array object.
  * `name` - environment variable name.
  * `value_from` - k8s ValueFrom.
    * `field_ref` - FieldRef for k8s env.
      * `field_path` - FieldPath of k8s.
    * `resource_field_ref` - ResourceFieldRef of k8s env.
      * `resource` - Resource of k8s.
  * `value` - environment variable value.
* `group_id` - Deployment group ID.
* `health_check_settings` - Deployment group health check settings.
  * `liveness_probe` - live health check.
    * `action_type` - health check method. HTTP: check by HTTP interface; CMD: check by executing command; TCP: check by establishing TCP connection.
    * `command` - Execute command check mode, the command to execute.
    * `failure_threshold` - Indicates the number of consecutive health check successes of the backend container from success to failure.
    * `initial_delay_seconds` - The time for the container to delay starting the health check.
    * `path` - The request path of the HTTP health check interface.
    * `period_seconds` - Interval between health checks.
    * `port` - Health check port, range 1~65535.
    * `scheme` - The inspection protocol used by the HTTP health check method. HTTP and HTTPS are supported.
    * `success_threshold` - Indicates the number of consecutive health check successes for the backend container from failure to success.
    * `timeout_seconds` - Maximum timeout for each health check response.
    * `type` - TSF_DEFAULT: tsf default readiness probe. K8S_NATIVE: k8s native probe. If not filled, it defaults to k8s native probe.
  * `readiness_probe` - readiness health check.
    * `action_type` - health check method. HTTP: check by HTTP interface; CMD: check by executing command; TCP: check by establishing TCP connection.
    * `command` - Execute command check mode, the command to execute.
    * `failure_threshold` - Indicates the number of consecutive health check successes for the backend container from success to failure.
    * `initial_delay_seconds` - The time for the container to delay starting the health check.
    * `path` - The request path of the HTTP health check interface.
    * `period_seconds` - The interval at which health checks are performed.
    * `port` - Health check port, range 1~65535.
    * `scheme` - The inspection protocol used by the HTTP health check method. HTTP and HTTPS are supported.
    * `success_threshold` - Indicates the number of consecutive health check successes for the backend container from failure to success.
    * `timeout_seconds` - The maximum timeout for each health check response.
    * `type` - TSF_DEFAULT: tsf default readiness probe. K8S_NATIVE: k8s native probe. If not filled, it defaults to k8s native probe.
* `instance_count` - Number of deployment group instances.
* `lb_ip` - load balancing ip.
* `max_surge` - The MaxSurge parameter of the kubernetes rolling update policy.
* `max_unavailable` - The MaxUnavailable parameter of the kubernetes rolling update policy.
* `message` - pod error message description.
* `microservice_type` - Service type.
* `namespace_name` - namespace name.
* `reponame` - Mirror name, such as /tsf/nginx.
* `server` - mirror server.
* `status` - Deployment group status.
* `tag_name` - Image version name.
* `updated_time` - Deployment group update timestamp.


## Import

tsf contain_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_contain_group.contain_group contain_group_id
```

