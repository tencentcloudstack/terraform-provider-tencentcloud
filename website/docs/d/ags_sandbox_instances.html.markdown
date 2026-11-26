---
subcategory: "Agent Sandbox(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_instances"
sidebar_current: "docs-tencentcloud-data_source-ags_sandbox_instances"
description: |-
  Use this data source to query the list of AGS sandbox instances.
---

# tencentcloud_ags_sandbox_instances

Use this data source to query the list of AGS sandbox instances.

## Example Usage

```hcl
data "tencentcloud_ags_sandbox_instances" "example" {
  instance_ids = ["xw5xweta6cdpoi4kg6hlkvycc3roqo4d3j6fp5ou"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Optional, List) List of sandbox instance ids to query. If empty, all sandbox instances are queried. At most 100 ids are supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

        * `add` - Capabilities to add. Reference value: `["SYS_ADMIN"]`.
        * `drop` - Capabilities to drop. Reference value: `["ALL"]`.
        * `path` - Request path.
        * `port` - Request port.
        * `scheme` - Request scheme.
      * `capabilities` - POSIX capabilities that can be added to or removed from the running container.
      * `cpu` - Amount of cpu resources.
      * `failure_threshold` - Failure threshold of the health check.
      * `http_get` - HTTP GET probe configuration.
      * `image_id` - Id of the custom WAA image.
      * `memory` - Amount of memory resources.
      * `name` - Environment variable name.
      * `name` - Port name.
      * `options` - Configuration options, corresponding to the options of resolv.conf.
      * `port` - Port number.
      * `privileged` - Whether privileged mode is enabled. Default value: `true`.
      * `probe_period_ms` - Interval in milliseconds between two health check probes.
      * `probe_timeout_ms` - Timeout in milliseconds of a single health check probe.
      * `protocol` - Protocol.
      * `ready_timeout_ms` - Timeout in milliseconds for the health check to become ready.
      * `searches` - Search domains, corresponding to the `search` directive of resolv.conf.
      * `servers` - DNS server addresses. Valid IP addresses are required. Default value: `10.0.0.1`.
      * `storage` - Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`. Note: this field may return null, which indicates no valid value.
      * `success_threshold` - Success threshold of the health check.
      * `value` - Environment variable value.
      * `version` - Built-in OSWorld version. Valid values: `osworld1` (osworld v1), `osworld2` (osworld v2). Default value: `osworld1`.
    * `args` - Startup arguments.
    * `command` - Startup command.
    * `dns_config` - Sandbox DNS configuration.
    * `env` - Environment variables.
    * `image_digest` - Image digest.
    * `image_registry_type` - Image registry type. Valid values: `enterprise`, `personal`, `custom`. `enterprise`: TCR enterprise container image service; `personal`: CCR personal container image service.
    * `image` - Image address.
    * `mount_path` - Local mount path of the sandbox instance (optional). By default it inherits the storage configuration of the tool.
    * `name` - Metadata name.
    * `name` - Name of the storage configuration in the sandbox tool.
    * `os_world_configuration` - Built-in OSWorld configuration.
    * `ports` - Port configuration.
    * `probe` - Probe configuration.
    * `read_only` - Read and write permission of the sandbox instance storage mount (optional). By default it inherits the storage configuration of the tool.
    * `resources` - Resource configuration.
    * `security_context` - Sandbox security context.
    * `sub_path` - Sub path of the sandbox instance storage mount (optional).
    * `value` - Metadata value.
    * `waa_configuration` - WAA sandbox tool configuration.
  * `auth_mode` - Sandbox access authentication mode. Valid values: `DEFAULT` (default, i.e. TOKEN authentication), `TOKEN` (Token authentication, all port access requires a TOKEN), `NONE` (no authentication, all port access does not require a TOKEN), `PUBLIC` (public mode, the ENVD management port (49983) access requires a TOKEN while other ports do not). Default value: `DEFAULT`.
  * `computer_configuration` - Sandbox configuration of desktop computer environments.
  * `create_time` - Creation time in ISO 8601 format.
  * `custom_configuration` - Custom configuration of the sandbox instance.
  * `expires_at` - Expiration time in ISO 8601 format. Null means no expiration time.
  * `instance_id` - Unique identifier of the sandbox instance.
  * `metadata` - Sandbox instance metadata.
  * `mount_options` - Storage mount options.
  * `network_mode` - Network mode. Valid values: `PUBLIC` (public network access), `SANDBOX` (no network), `INTERNAL_SERVICE` (Tencent Cloud internal public services). It can override the network configuration at the tool level, but it is ineffective if the tool itself does not support VPC networks.
  * `persistent` - Whether the instance is persistent.
  * `status` - Instance status. Valid values: `STARTING` (starting), `RUNNING` (running), `STOPPING` (stopping), `STOPPED` (stopped), `STOP_FAILED` (stop failed), `FAILED` (failed), `RESIZING` (updating instance specifications).
  * `stop_reason` - Stop reason: `manual` (manual), `timeout` (timeout), `error` (error), `system` (system). It has a value only when the status is `STOPPED`, `STOP_FAILED` or `FAILED`. When the provider fails to stop the instance, the status is `STOP_FAILED` and the reason is `error`.
  * `timeout_seconds` - Timeout in seconds. Null means no timeout is set.
  * `tool_id` - Id of the sandbox tool the instance belongs to.
  * `tool_name` - Name of the sandbox tool the instance belongs to.
  * `update_time` - Last updated time in ISO 8601 format.
* `id` - ID of the data source, a random token.
* `instance_set` - Sandbox instance list.


