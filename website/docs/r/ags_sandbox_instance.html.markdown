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

~> **NOTE:** The `x-ags-` name prefix of `metadata` is reserved by the system. Such entries are added by the cloud side (for example `x-ags-precreate-match`), are not tracked in the state and must not be configured, otherwise a permanent diff is reported.

~> **NOTE:** The timeout of a sandbox instance that is not persistent is enforced by the cloud side: the instance is reclaimed automatically when `timeout` expires. Such an instance is removed from the state during the next refresh, so the next `terraform apply` creates a new sandbox instance instead of updating the existing one.

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "custom1" {
  custom_configuration = {
    args                = ["/nix/runtime/nexus/deployment/runtime/start.sh", "-H", "0.0.0.0", "-p", "8000", "-v", "3.12"]
    command             = ["/nix/runtime/nix-env/bin/bash"]
    dns_config          = null
    env                 = null
    image               = "ccr.ccs.tencentyun.com/finofliu/go-httpbin:2.19-alpine-with-envd"
    image_registry_type = "personal"
    ports = [
      {
        name     = "http"
        port     = 8000
        protocol = "TCP"
      },
    ]
    probe = {
      failure_threshold = 10
      http_get = {
        path   = "/health"
        port   = 8000
        scheme = "HTTP"
      }
      probe_period_ms   = 3000
      probe_timeout_ms  = 1000
      ready_timeout_ms  = 30000
      success_threshold = 1
    }
    resources = {
      cpu    = "2"
      memory = "2Gi"
    }
  }
  description = "image-volume regression baseline"
  network_configuration = {
    network_mode = "PUBLIC"
  }
  persistent = false
  role_arn   = "qcs::cam::uin/3321337994:roleName/cos-tcr-full-ags"
  storage_mounts = [
    {
      mount_path = "/nix"
      name       = "nix-runtime"
      read_only  = true
      storage_source = {
        image = {
          image_registry_type = "enterprise"
          reference           = "tencentags-guangzhou-994-test.tencentcloudcr.com/ags-e2e/nexus-runtime:0.2.2"
          sub_path            = "nix"
        }
      }
    },
    {
      mount_path = "/mnt/same-name-runtime"
      name       = "runtime"
      read_only  = true
      storage_source = {
        image = {
          image_registry_type = "enterprise"
          reference           = "tencentags-guangzhou-994-test.tencentcloudcr.com/ags-e2e/nexus-runtime:0.2.2"
          sub_path            = "nix/runtime/nix-env/bin"
        }
      }
    },
  ]
  tool_name = "tf-test"
  tool_type = "custom"
}

resource "tencentcloud_ags_sandbox_instance" "example1" {
  auth_mode = "DEFAULT"
  custom_configuration = {
    args                = ["/nix/runtime/nexus/deployment/runtime/start.sh", "-H", "0.0.0.0", "-p", "8000", "-v", "3.12"]
    command             = ["/nix/runtime/nix-env/bin/bash"]
    image               = "ccr.ccs.tencentyun.com/finofliu/go-httpbin:2.19-alpine-with-envd"
    image_registry_type = "personal"
    ports = [
      {
        name     = "http"
        port     = 8000
        protocol = "TCP"
      },
    ]
    probe = {
      failure_threshold = 10
      http_get = {
        path   = "/health"
        port   = 8000
        scheme = "HTTP"
      }
      probe_period_ms   = 3000
      probe_timeout_ms  = 1000
      ready_timeout_ms  = 30000
      success_threshold = 1
    }
    resources = {
      cpu    = "2"
      memory = "2Gi"
    }
  }
  mount_options = [
    {
      mount_path = "/nix"
      name       = "nix-runtime"
    },
    {
      mount_path = "/mnt/same-name-runtime"
      name       = "runtime"
    },
  ]
  timeout   = "5m"
  tool_name = tencentcloud_ags_sandbox_tool.custom1.tool_name
}
```

## Argument Reference

The following arguments are supported:

* `auth_mode` - (Optional, String) Authentication mode of the sandbox access. Valid values: `DEFAULT` (token authentication, the default value), `TOKEN` (token authentication, every port access requires a token), `NONE` (no authentication, no port access requires a token), `PUBLIC` (public mode, the ENVD management port 49983 requires a token while other ports do not).
* `custom_configuration` - (Optional, Object) Custom configuration of the sandbox instance.
* `metadata` - (Optional, List) Metadata of the sandbox instance. The `x-ags-` name prefix is reserved by the system: such entries are added by the cloud side, are not tracked in the state and must not be configured, otherwise a permanent diff is reported.
* `mount_options` - (Optional, List) Storage mount configuration of the sandbox instance.
* `snapshot_id` - (Optional, String) Id of the sandbox snapshot to start from.
* `source_instance_id` - (Optional, String) Id of the source sandbox instance to clone.
* `timeout` - (Optional, String) How long the sandbox instance keeps running before it is reclaimed automatically, for example `5m`, `300s` or `1h`. Default value: `5m`. Minimum `30s`, maximum `24h`. It is refreshed from the cloud side timeout in seconds on read.
* `tool_id` - (Optional, String) Id of the sandbox tool. At least one of `tool_id` and `tool_name` must be specified.
* `tool_name` - (Optional, String) Name of the sandbox tool. At least one of `tool_id` and `tool_name` must be specified.

The `custom_configuration` object supports the following:

* `args` - (Optional, List) Startup arguments. The order is preserved.
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
* `memory` - (Optional, String) Amount of memory resources.
* `storage` - (Optional, String) Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`.

The `security_context` object of `custom_configuration` supports the following:

* `capabilities` - (Optional, Object) POSIX capabilities that can be added to or removed from the running container.
* `privileged` - (Optional, Bool) Whether the sandbox runs in privileged mode. Default value: `true`.

The `capabilities` object of `security_context` supports the following:

* `add` - (Optional, Set) Capabilities to add. Reference value: `["SYS_ADMIN"]`.
* `drop` - (Optional, Set) Capabilities to drop. Reference value: `["ALL"]`.

The `metadata` object supports the following:

* `name` - (Optional, String) Metadata name. The `x-ags-` prefix is reserved by the system and must not be used here, otherwise a permanent diff is reported.
* `value` - (Optional, String) Metadata value.

The `mount_options` object supports the following:

* `mount_path` - (Optional, String) Local mount path inside the sandbox instance. By default it inherits the mount path of the storage configuration in the tool.
* `name` - (Optional, String) Name of the storage configuration in the sandbox tool.
* `read_only` - (Optional, Bool) Read and write permission of the mount. By default it inherits the storage configuration in the tool.
* `sub_path` - (Optional, String) Sub path of the storage mounted into the sandbox instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `image_id` - Custom WAA image id.
    * `version` - Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. Default value: `osworld1`.
  * `os_world_configuration` - Built-in OSWorld configuration.
  * `waa_configuration` - WAA sandbox configuration.
* `computer_configuration` - Desktop environment configuration of the sandbox instance.
* `create_time` - Creation time in ISO 8601 format.
* `expires_at` - Expiration time in ISO 8601 format. `null` means the instance never expires.
* `id` - Id of the resource, which is the same as `instance_id`.
* `instance_id` - Unique identifier of the sandbox instance.
* `network_mode` - Network mode. Valid values: `PUBLIC` (public network access), `SANDBOX` (no network), `INTERNAL_SERVICE` (Tencent Cloud internal public service). It can override the network configuration of the tool level, but when a tool itself does not support VPC networks, the VPC mode chosen for the instance is invalid.
* `persistent` - Whether the sandbox instance is a persistent instance.
* `status` - Status of the instance. Valid values: `STARTING` (starting), `RUNNING` (running), `STOPPING` (stopping), `STOPPED` (stopped), `STOP_FAILED` (stop failed), `FAILED` (failed), `RESIZING` (the specification is being updated).
* `stop_reason` - Stop reason: `manual` (manual), `timeout` (timeout), `error` (error) or `system` (system). It only has a value when the status is `STOPPED`, `STOP_FAILED` or `FAILED`. When the provider fails to stop the instance, the status is `STOP_FAILED` and the reason is `error`.
* `update_time` - Last update time in ISO 8601 format.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeSandboxInstances](https://www.tencentcloud.com/search/?wd=DescribeSandboxInstances&type=document&lang=en)
  - [StartSandboxInstance](https://www.tencentcloud.com/search/?wd=StartSandboxInstance&type=document&lang=en)
* Read: [DescribeSandboxInstances](https://www.tencentcloud.com/search/?wd=DescribeSandboxInstances&type=document&lang=en)
* Update:
  - [DescribeSandboxInstances](https://www.tencentcloud.com/search/?wd=DescribeSandboxInstances&type=document&lang=en)
  - [UpdateSandboxInstance](https://www.tencentcloud.com/search/?wd=UpdateSandboxInstance&type=document&lang=en)
* Delete: [StopSandboxInstance](https://www.tencentcloud.com/search/?wd=StopSandboxInstance&type=document&lang=en)

## Import

AGS sandbox instance can be imported using the id, e.g.

```
terraform import tencentcloud_ags_sandbox_instance.example instance-xxxxxxxx
```

