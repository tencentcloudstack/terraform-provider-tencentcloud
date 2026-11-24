---
subcategory: "Agent Sandbox(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_instance"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_instance"
description: |-
  Provides a resource to start an AGS sandbox instance.
---

# tencentcloud_ags_sandbox_instance

Provides a resource to start an AGS sandbox instance.

-> **Note:** `tool_id`, `tool_name`, `source_instance_id`, `snapshot_id`, `client_token`, `mount_options`, `custom_configuration`, `auth_mode`, `storage_mounts` and `start_mode` are create-only: the update API does not accept them, so changing any of them destroys the existing sandbox instance and starts a new one.

-> **Note:** `timeout`, `source_instance_id`, `snapshot_id`, `client_token`, `storage_mounts` and `start_mode` are write-only inputs. The describe API reports the timeout in seconds (`timeout_seconds`) instead of the original duration string and does not report the other fields at all, so the configured values are kept in state and drift on them cannot be detected.

-> **Note:** Destroying this resource only stops the sandbox instance itself. It does not touch the source instance, the sandbox tool or any storage (COS bucket, CFS file system or volume) referenced by `storage_mounts`. Whether the data of a stopped instance is retained afterwards is decided by the cloud side.

-> **Note:** The start and stop requests are asynchronous: the create path polls the instance status until it is `RUNNING`, and the destroy path polls it until it is `STOPPED` or no longer visible.

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_instance" "example" {
  tool_id = "tool-xxxxxxxx"
  timeout = "30m"

  metadata = [
    {
      name  = "run_id"
      value = "42"
    },
  ]

  custom_configuration = {
    image               = "ccr.ccs.tencentyun.com/test/go-httpbin:v2.16.1"
    image_registry_type = "personal"
    command             = ["/bin/go-httpbin"]
    args                = ["-port", "8080"]
    ports = [
      {
        name     = "http"
        port     = 8080
        protocol = "TCP"
      },
    ]
    resources = {
      cpu     = "1"
      memory  = "1Gi"
      storage = null
    }
  }

  auth_mode = "TOKEN"
}
```

## Argument Reference

The following arguments are supported:

* `auth_mode` - (Optional, String) Access authentication mode of the sandbox. Valid values: `DEFAULT`, `TOKEN`, `NONE`, `PUBLIC`. Default value: `DEFAULT`.
* `client_token` - (Optional, String) Idempotency token of the start request, up to 64 characters. When it is not set, a random token is generated for each apply so that provider level retries never create a duplicate instance.
* `custom_configuration` - (Optional, Object) Custom configuration of the sandbox instance.
* `metadata` - (Optional, List) Metadata of the sandbox instance.
* `mount_options` - (Optional, List) Storage mount options of the sandbox instance.
* `snapshot_id` - (Optional, String) Id of the sandbox snapshot to start from.
* `source_instance_id` - (Optional, String) Id of the source sandbox instance to clone from.
* `start_mode` - (Optional, String) Start mode of the sandbox instance. Valid values: `DEFAULT`, `WAIT_UNTIL_HEALTHY`. `WAIT_UNTIL_HEALTHY` uses the persistent background creation flow and the instance may be reported as `STARTING` right after the request is accepted; it is currently only supported by WAA sandboxes and GPU sandboxes. Default value: `DEFAULT`. It is not reported by the describe API, so the configured value is kept in state and drift on it cannot be detected.
* `storage_mounts` - (Optional, List) Additional storage mount rules applied when the sandbox instance starts. Only `volume`, `cfs`, `cos` and `image` sources are supported. It is not reported by the describe API, so the configured value is kept in state and drift on it cannot be detected.
* `timeout` - (Optional, String) Timeout of the sandbox instance. The instance is recycled automatically after it expires. Supported formats are `5m`, `300s`, `1h`, with a minimum of 30s and a maximum of 24h. It is not reported by the describe API as a duration string, so the configured value is kept in state and drift on it cannot be detected.
* `tool_id` - (Optional, String) Id of the sandbox tool. Either `tool_id` or `tool_name` must be set.
* `tool_name` - (Optional, String) Name of the sandbox tool. Either `tool_id` or `tool_name` must be set.

The `custom_configuration` object supports the following:

* `args` - (Optional, Set) Startup arguments.
* `command` - (Optional, Set) Startup command.
* `dns_config` - (Optional, Object) Sandbox DNS configuration.
* `env` - (Optional, List) Environment variables.
* `image_registry_type` - (Optional, String) Image registry type. Valid values: `enterprise`, `personal`, `custom`.
* `image` - (Optional, String) Image address.
* `ports` - (Optional, List) Port configuration.
* `probe` - (Optional, Object) Probe configuration.
* `resources` - (Optional, Object) Resource configuration.
* `security_context` - (Optional, Object) Sandbox security context.

The `dns_config` object of `custom_configuration` supports the following:

* `servers` - (Required, Set) DNS server addresses. Valid IP addresses are required. Default value: `10.0.0.1`.
* `options` - (Optional, Set) Options, corresponding to the options of resolv.conf.
* `searches` - (Optional, Set) Search domains, corresponding to the `search` directive of resolv.conf.

The `env` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Environment variable name.
* `value` - (Optional, String) Environment variable value.

The `ports` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Port name.
* `port` - (Optional, Int) Port number.
* `protocol` - (Optional, String) Protocol.

The `probe` object of `custom_configuration` supports the following:

* `failure_threshold` - (Optional, Int) Number of consecutive failures required to mark the probe as unhealthy.
* `http_get` - (Optional, Object) HTTP GET probe configuration.
* `probe_period_ms` - (Optional, Int) Interval in milliseconds between two probes.
* `probe_timeout_ms` - (Optional, Int) Timeout in milliseconds of a single probe.
* `ready_timeout_ms` - (Optional, Int) Timeout in milliseconds for the probe to become ready.
* `success_threshold` - (Optional, Int) Number of consecutive successes required to mark the probe as healthy.

The `http_get` object of `probe` supports the following:

* `path` - (Optional, String) Request path.
* `port` - (Optional, Int) Request port.
* `scheme` - (Optional, String) Request scheme.

The `resources` object of `custom_configuration` supports the following:

* `cpu` - (Optional, String) Amount of cpu resources.
* `gpu_resource_group_id` - (Optional, String) Id of the GPU resource group associated with the sandbox.
* `memory` - (Optional, String) Amount of memory resources.
* `storage` - (Optional, String) Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`.

The `security_context` object of `custom_configuration` supports the following:

* `capabilities` - (Optional, Object) POSIX capabilities that can be added to or removed from the running container.
* `privileged` - (Optional, Bool) Whether the sandbox runs in privileged mode. Default value: `true`.

The `capabilities` object of `security_context` supports the following:

* `add` - (Optional, Set) Capabilities to add. Reference value: `["SYS_ADMIN"]`.
* `drop` - (Optional, Set) Capabilities to drop. Reference value: `["ALL"]`.

The `metadata` object supports the following:

* `name` - (Optional, String) Metadata name.
* `value` - (Optional, String) Metadata value.

The `mount_options` object supports the following:

* `mount_path` - (Optional, String) Local mount path inside the sandbox instance. It inherits the storage configuration of the tool when it is not set.
* `name` - (Optional, String) Name of the storage configuration in the sandbox tool.
* `read_only` - (Optional, Bool) Whether the mount is read only. It inherits the storage configuration of the tool when it is not set.
* `sub_path` - (Optional, String) Sub path of the storage mount inside the sandbox instance.

The `storage_mounts` object supports the following:

* `mount_path` - (Optional, String) Local mount path inside the sandbox instance.
* `name` - (Optional, String) Name of the storage mount configuration.
* `read_only` - (Optional, Bool) Whether the mount is read only. Default value: `false`.
* `storage_source` - (Optional, Object) Storage configuration.
* `volume_template` - (Optional, Object) Volume template reference and derivation rule.
* `volume` - (Optional, Object) Volume reference.

The `storage_source` object of `storage_mounts` supports the following:

* `agent_bucket` - (Optional, Object) AgentBucket storage configuration.
* `cfs` - (Optional, Object) File storage configuration.
* `cos` - (Optional, Object) Object storage bucket configuration.
* `image` - (Optional, Object) Image volume configuration.

The `agent_bucket` object of `storage_source` supports the following:

* `access_domain` - (Optional, String) Access domain of the AgentBucket.
* `library_id` - (Optional, String) Library id of the AgentBucket.
* `space_id` - (Optional, String) Space id of the AgentBucket.

The `cfs` object of `storage_source` supports the following:

* `file_system_id` - (Optional, String) Id of the CFS file system.
* `path` - (Optional, String) CFS mount path.

The `cos` object of `storage_source` supports the following:

* `bucket_name` - (Optional, String) Object storage bucket name.
* `bucket_path` - (Optional, String) Object storage bucket path. It must be an absolute path starting with `/`.
* `endpoint` - (Optional, String) Object storage access endpoint.

The `image` object of `storage_source` supports the following:

* `image_registry_type` - (Optional, String) Image registry type. Valid values: `enterprise`, `personal`.
* `reference` - (Optional, String) Image address.
* `sub_path` - (Optional, String) Path inside the image.

The `volume` object of `storage_mounts` supports the following:

* `volume_id` - (Optional, String) Stable volume id in the form `vol-xxxxxxxx`, resolved by the cloud side.
* `volume_name` - (Optional, String) Volume name. The cloud side resolves it into a stable volume id under the current AppId.

The `volume_template` object of `storage_mounts` supports the following:

* `volume_template_name` - (Required, String) Name of the volume template.
* `reuse_key` - (Optional, String) Reuse key. It can be a literal value or a `${metadata.xxx}` expression.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the sandbox instance in ISO8601 format.
* `expires_at` - Expiry time of the sandbox instance in ISO8601 format. Null means no expiry time.
* `id` - Id of the resource, which is the same as `instance_id`.
* `instance_id` - Id of the sandbox instance.
* `network_mode` - Network mode of the sandbox instance. Valid values: `PUBLIC`, `SANDBOX`, `INTERNAL_SERVICE`. It overrides the network configuration of the tool.
* `persistent` - Whether the sandbox instance is persistent.
* `status` - Status of the sandbox instance. Valid values: `STARTING`, `RUNNING`, `STOPPING`, `STOPPED`, `STOP_FAILED`, `FAILED`, `RESIZING`.
* `stop_reason` - Reason why the sandbox instance was stopped. Valid values: `manual`, `timeout`, `error`, `system`. It is only reported when the status is `STOPPED`, `STOP_FAILED` or `FAILED`.
* `timeout_seconds` - Timeout of the sandbox instance in seconds. Null means no timeout is set.
* `update_time` - Last update time of the sandbox instance in ISO8601 format.

The `image` object of `storage_source` exports the following:

* `digest` - Image digest resolved by the cloud side. It is not required in requests and is only reported by the cloud.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeSandboxInstances](https://www.tencentcloud.com/search/?wd=DescribeSandboxInstances&type=document&lang=en)
  - [StartSandboxInstance](https://www.tencentcloud.com/search/?wd=StartSandboxInstance&type=document&lang=en)
* Read: [DescribeSandboxInstances](https://www.tencentcloud.com/search/?wd=DescribeSandboxInstances&type=document&lang=en)
* Update:
  - [DescribeSandboxInstances](https://www.tencentcloud.com/search/?wd=DescribeSandboxInstances&type=document&lang=en)
  - [UpdateSandboxInstance](https://www.tencentcloud.com/search/?wd=UpdateSandboxInstance&type=document&lang=en)
* Delete:
  - [DescribeSandboxInstances](https://www.tencentcloud.com/search/?wd=DescribeSandboxInstances&type=document&lang=en)
  - [StopSandboxInstance](https://www.tencentcloud.com/search/?wd=StopSandboxInstance&type=document&lang=en)

## Import

AGS sandbox instance can be imported using the sandbox instance id, e.g.

```
$ terraform import tencentcloud_ags_sandbox_instance.example instance-xxxxxxxx
```

