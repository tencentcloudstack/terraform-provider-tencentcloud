---
subcategory: "AGS(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_tool"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_tool"
description: |-
  Provides a resource to create and manage an AGS sandbox tool.
---

# tencentcloud_ags_sandbox_tool

Provides a resource to create and manage an AGS sandbox tool.

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "my-sandbox-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "PUBLIC"
  }

  description     = "An example sandbox tool"
  default_timeout = "10m"

  tags {
    key   = "env"
    value = "test"
  }

  custom_configuration {
    image               = "nginx:latest"
    image_registry_type = "personal"
    command             = ["/bin/bash"]
    args                = ["-c", "echo hello"]

    env {
      name  = "MY_VAR"
      value = "my_value"
    }

    resources {
      cpu     = "2"
      memory  = "4Gi"
      storage = "5Gi"
    }
  }
}
```

### ## Import

```hcl
terraform import tencentcloud_ags_sandbox_tool.example < tool-id >
```

## Argument Reference

The following arguments are supported:

* `network_configuration` - (Required, Object) Framework attribute.
* `tool_name` - (Required, String) Sandbox tool name, 1-50 characters, supports English letters, digits, underscores and hyphens. Must be unique under the same AppId.
* `tool_type` - (Required, String) Sandbox tool type. Supported values: browser, code-interpreter, mobile, osworld, custom, swebench, aio, android-world, waa.
* `client_token` - (Optional, String) Idempotency token, maximum 64 characters.
* `computer_configuration` - (Optional, Object) Desktop/computer environment sandbox configuration.
* `custom_configuration` - (Optional, Object) Custom configuration for sandbox tool.
* `default_timeout` - (Optional, String) Default timeout, supports formats like 5m, 300s, 1h. Maximum 24 hours. If not specified, the system default (5 minutes) is used.
* `description` - (Optional, String) Sandbox tool description, maximum 200 characters.
* `log_configuration` - (Optional, Object) Log push configuration for sandbox tool.
* `persistent` - (Optional, Bool) Persistent sandbox flag.
* `role_arn` - (Optional, String) Role ARN.
* `storage_mounts` - (Optional, List) Storage mount configuration for sandbox tool instances.
* `tags` - (Optional, List) Tags for the sandbox tool.

The `computer_configuration` object supports the following:

* `os_world_configuration` - (Optional, Object) Built-in OSWorld configuration.
* `waa_configuration` - (Optional, Object) WAA sandbox tool configuration.

The `os_world_configuration` object of `computer_configuration` supports the following:

* `version` - (Optional, String) Built-in OSWorld version. Supported: osworld1, osworld2. Default: osworld1.

The `waa_configuration` object of `computer_configuration` supports the following:

* `image_id` - (Optional, String) Custom WAA image ID.

The `custom_configuration` object supports the following:

* `args` - (Optional, Set) Startup arguments.
* `command` - (Optional, Set) Startup command.
* `dns_config` - (Optional, Object) Sandbox DNS configuration.
* `env` - (Optional, List) Environment variables.
* `image_registry_type` - (Optional, String) Image registry type: enterprise, personal, custom.
* `image` - (Optional, String) Image address.
* `ports` - (Optional, List) Port configuration.
* `probe` - (Optional, Object) Probe configuration.
* `resources` - (Optional, Object) Resource configuration.
* `security_context` - (Optional, Object) Sandbox security context.

The `dns_config` object of `custom_configuration` supports the following:

* `servers` - (Required, Set) DNS server addresses (valid IP addresses). Default: 10.0.0.1.
* `options` - (Optional, Set) Configuration items (corresponds to resolv.conf options).
* `searches` - (Optional, Set) Search domains (corresponds to resolv.conf search directive).

The `env` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Environment variable name.
* `value` - (Optional, String) Environment variable value.

The `ports` object of `custom_configuration` supports the following:

* `name` - (Optional, String) Port name.
* `port` - (Optional, Int) Port number.
* `protocol` - (Optional, String) Protocol.

The `probe` object of `custom_configuration` supports the following:

* `failure_threshold` - (Optional, Int) Health check failure threshold.
* `http_get` - (Optional, Object) HTTP GET probe configuration.
* `probe_period_ms` - (Optional, Int) Probe interval.
* `probe_timeout_ms` - (Optional, Int) Single probe timeout.
* `ready_timeout_ms` - (Optional, Int) Health check readiness timeout.
* `success_threshold` - (Optional, Int) Health check success threshold.

The `http_get` object of `probe` supports the following:

* `path` - (Optional, String) Path.
* `port` - (Optional, Int) Port.
* `scheme` - (Optional, String) Scheme.

The `resources` object of `custom_configuration` supports the following:

* `cpu` - (Optional, String) CPU resource amount.
* `memory` - (Optional, String) Memory resource amount.
* `storage` - (Optional, String) Custom disk size. Supported values: 1Gi, 5Gi, 10Gi, 20Gi.

The `security_context` object of `custom_configuration` supports the following:

* `capabilities` - (Optional, Object) POSIX capabilities to add or remove from the container.
* `privileged` - (Optional, Bool) Whether to enable privileged mode. Default: true.

The `capabilities` object of `security_context` supports the following:

* `add` - (Optional, Set) Capabilities to add. Example: ["SYS_ADMIN"].
* `drop` - (Optional, Set) Capabilities to drop. Example: ["ALL"].

The `log_configuration` object supports the following:

* `cls_config` - (Optional, Object) CLS log push configuration.
* `log_sources` - (Optional, Object) Log source configuration.

The `cls_config` object of `log_configuration` supports the following:

* `topic_id` - (Optional, String) CLS log topic ID for sandbox tool log push.

The `log_sources` object of `log_configuration` supports the following:

* `files` - (Optional, Set) Log file paths to collect, must be under /logs/ directory, no subdirectories, max 10 files.

The `network_configuration` object supports the following:

* `network_mode` - (Required, String) Network mode. Supported values: PUBLIC, VPC, SANDBOX.
* `vpc_config` - (Optional, Object) VPC network configuration.

The `vpc_config` object of `network_configuration` supports the following:

* `security_group_ids` - (Optional, Set) List of security group IDs.
* `subnet_ids` - (Optional, Set) List of VPC subnet IDs.

The `storage_mounts` object supports the following:

* `mount_path` - (Optional, String) Local mount path in the sandbox instance.
* `name` - (Optional, String) Storage mount configuration name.
* `read_only` - (Optional, Bool) Read-only permission configuration, defaults to false.
* `storage_source` - (Optional, Object) Storage source configuration.
* `volume_template` - (Optional, Object) VolumeTemplate reference and derivation rules.
* `volume` - (Optional, Object) Volume reference.

The `storage_source` object of `storage_mounts` supports the following:

* `agent_bucket` - (Optional, Object) AgentBucket storage configuration.
* `cfs` - (Optional, Object) CFS storage configuration.
* `cos` - (Optional, Object) COS storage configuration.
* `image` - (Optional, Object) Image volume configuration.

The `agent_bucket` object of `storage_source` supports the following:

* `access_domain` - (Optional, String) AgentBucket AccessDomain.
* `library_id` - (Optional, String) AgentBucket LibraryID.
* `space_id` - (Optional, String) AgentBucket spaceId.

The `cfs` object of `storage_source` supports the following:

* `file_system_id` - (Optional, String) CFS resource ID.
* `path` - (Optional, String) CFS mount path.

The `cos` object of `storage_source` supports the following:

* `bucket_name` - (Optional, String) COS bucket name.
* `bucket_path` - (Optional, String) COS bucket path, must be an absolute path starting with /.
* `endpoint` - (Optional, String) COS access domain.

The `image` object of `storage_source` supports the following:

* `digest` - (Optional, String) Image digest, not required in requests.
* `image_registry_type` - (Optional, String) Image registry type: enterprise, personal.
* `reference` - (Optional, String) Image address.
* `sub_path` - (Optional, String) Path inside the image.

The `volume` object of `storage_mounts` supports the following:

* `volume_id` - (Optional, String) Stable volume ID, format vol-xxxxxxxx.
* `volume_name` - (Optional, String) Volume name; resolved to a stable VolumeId under the current AppId.

The `volume_template` object of `storage_mounts` supports the following:

* `volume_template_name` - (Required, String) VolumeTemplate name.
* `reuse_key` - (Optional, String) Reuse key, can be a literal or ${metadata.xxx}.

The `tags` object supports the following:

* `key` - (Optional, String) Tag key.
* `value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time in ISO8601 format.
* `id` - Resource identifier (same as tool_id).
* `status` - Sandbox tool status: CREATING, ACTIVE, DELETING, FAILED.
* `tool_id` - Sandbox tool ID.
* `update_time` - Update time in ISO8601 format.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateSandboxTool](https://www.tencentcloud.com/search/?wd=CreateSandboxTool&type=document&lang=en)
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
* Read: [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
* Update:
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
  - [UpdateSandboxTool](https://www.tencentcloud.com/search/?wd=UpdateSandboxTool&type=document&lang=en)
* Delete: [DeleteSandboxTool](https://www.tencentcloud.com/search/?wd=DeleteSandboxTool&type=document&lang=en)

