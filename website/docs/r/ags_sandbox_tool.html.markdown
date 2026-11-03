---
subcategory: "Agent Sandbox(AGS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_tool"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_tool"
description: |-
  Provides a resource to create a AGS sandbox tool
---

# tencentcloud_ags_sandbox_tool

Provides a resource to create a AGS sandbox tool

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "tf-example"
  tool_type = "custom"
  network_configuration {
    network_mode = "PUBLIC"
  }
  description     = "example sandbox tool"
  default_timeout = "10m"
  tags = {
    "env" = "test"
  }
}

resource "tencentcloud_ags_sandbox_tool" "vpc" {
  tool_name = "tf-vpc-example"
  tool_type = "custom"
  network_configuration {
    network_mode = "VPC"
    vpc_config {
      subnet_ids         = ["subnet-xxxxxxxx"]
      security_group_ids = ["sg-xxxxxxxx"]
    }
  }
  description     = "VPC sandbox tool example"
  default_timeout = "30m"
  custom_configuration {
    image = "ccr.ccs.tencentyun.com/namespace/image:tag"
    resources {
      cpu    = "2"
      memory = "4Gi"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `tool_name` - (Required, String, ForceNew) Sandbox tool name, length 1-50 characters, supports Chinese, English, numbers, underscores. Must be unique under the same AppId.
* `tool_type` - (Required, String, ForceNew) Sandbox tool type. Currently supports: browser, code-interpreter, custom, mobile, osworld, swebench, aio, android-world, waa.
* `client_token` - (Optional, String, ForceNew) Idempotency token, max length 64 characters.
* `computer_configuration` - (Optional, List) Computer environment sandbox configuration.
* `custom_configuration` - (Optional, List) Sandbox tool custom configuration.
* `default_timeout` - (Optional, String, ForceNew) Default timeout, supports formats: 5m, 300s, 1h etc. If not specified, system default (5 minutes) is used. Max 24 hours.
* `description` - (Optional, String) Sandbox tool description, max length 200 characters.
* `log_configuration` - (Optional, List, ForceNew) Sandbox tool log push configuration.
* `network_configuration` - (Optional, List) Network configuration.
* `persistent` - (Optional, Bool, ForceNew) Persistent sandbox flag.
* `role_arn` - (Optional, String, ForceNew) Role ARN bound to the sandbox tool.
* `storage_mounts` - (Optional, List, ForceNew) Sandbox tool storage mount configurations.
* `tags` - (Optional, Map) Tags for the sandbox tool, up to 50 tags.

The `agent_bucket` object of `storage_source` supports the following:

* `library_id` - (Required, String) Library ID.
* `space_id` - (Required, String) Space ID.
* `access_domain` - (Optional, String) Access domain.

The `capabilities` object of `security_context` supports the following:

* `add` - (Optional, List) Capabilities to add.
* `drop` - (Optional, List) Capabilities to drop.

The `cfs` object of `storage_source` supports the following:

* `file_system_id` - (Required, String) CFS file system ID.
* `path` - (Required, String) CFS mount path.

The `cls_config` object of `log_configuration` supports the following:

* `topic_id` - (Optional, String) CLS topic ID.

The `computer_configuration` object supports the following:

* `os_world_configuration` - (Optional, List) OSWorld configuration.
* `waa_configuration` - (Optional, List) WAA configuration.

The `cos` object of `storage_source` supports the following:

* `bucket_name` - (Optional, String) COS bucket name.
* `bucket_path` - (Optional, String) COS bucket path.
* `endpoint` - (Optional, String) COS endpoint.

The `custom_configuration` object supports the following:

* `args` - (Optional, List) Startup arguments.
* `command` - (Optional, List) Startup command.
* `dns_config` - (Optional, List) DNS configuration.
* `env` - (Optional, Set) Environment variables.
* `image_registry_type` - (Optional, String) Image registry type: enterprise, personal, custom.
* `image` - (Optional, String) Image address.
* `ports` - (Optional, List) Port configurations.
* `probe` - (Optional, List) Probe configuration.
* `resources` - (Optional, List) Resource configuration.
* `security_context` - (Optional, List) Security context.

The `dns_config` object of `custom_configuration` supports the following:

* `options` - (Optional, List) DNS options.
* `searches` - (Optional, List) DNS search domains.
* `servers` - (Optional, List) DNS server addresses.

The `env` object of `custom_configuration` supports the following:

* `name` - (Required, String) Environment variable name.
* `value` - (Required, String) Environment variable value.

The `http_get` object of `probe` supports the following:

* `path` - (Required, String) HTTP path.
* `port` - (Required, Int) HTTP port.
* `scheme` - (Optional, String) HTTP scheme.

The `image` object of `storage_source` supports the following:

* `reference` - (Required, String) Image reference address.
* `digest` - (Optional, String) Image digest.
* `image_registry_type` - (Optional, String) Image registry type: enterprise, personal, custom.
* `sub_path` - (Optional, String) Sub path in the image.

The `log_configuration` object supports the following:

* `cls_config` - (Optional, List) CLS log configuration.
* `log_sources` - (Optional, List) Log source configuration.

The `log_sources` object of `log_configuration` supports the following:

* `files` - (Optional, List) Log file paths to collect, must be in /logs/ directory, max 10 files.

The `network_configuration` object supports the following:

* `network_mode` - (Required, String) Network mode. Valid values: PUBLIC, VPC.
* `vpc_config` - (Optional, List) VPC network configuration. Required when network_mode is VPC.

The `os_world_configuration` object of `computer_configuration` supports the following:

* `version` - (Optional, String) OSWorld version.

The `ports` object of `custom_configuration` supports the following:

* `port` - (Required, Int) Port number.
* `name` - (Optional, String) Port name.
* `protocol` - (Optional, String) Port protocol.

The `probe` object of `custom_configuration` supports the following:

* `failure_threshold` - (Optional, Int) Failure threshold.
* `http_get` - (Optional, List) HTTP GET probe configuration.
* `probe_period_ms` - (Optional, Int) Probe period in milliseconds.
* `probe_timeout_ms` - (Optional, Int) Probe timeout in milliseconds.
* `ready_timeout_ms` - (Optional, Int) Ready timeout in milliseconds.
* `success_threshold` - (Optional, Int) Success threshold.

The `resources` object of `custom_configuration` supports the following:

* `cpu` - (Required, String) CPU resource specification.
* `memory` - (Required, String) Memory resource specification.
* `storage` - (Optional, String) Storage resource specification.

The `security_context` object of `custom_configuration` supports the following:

* `capabilities` - (Optional, List) POSIX capabilities for the container.
* `privileged` - (Optional, Bool) Whether to enable privileged mode. Default: true.

The `storage_mounts` object supports the following:

* `mount_path` - (Required, String) Local mount path on the sandbox instance.
* `name` - (Required, String) Storage mount configuration name.
* `storage_source` - (Required, List) Storage source configuration.
* `read_only` - (Optional, Bool) Read-only mount permission, defaults to false.
* `volume_template` - (Optional, List) Volume template reference.
* `volume` - (Optional, List) Volume reference.

The `storage_source` object of `storage_mounts` supports the following:

* `agent_bucket` - (Optional, List) AgentBucket storage configuration.
* `cfs` - (Optional, List) CFS storage configuration.
* `cos` - (Optional, List) COS storage configuration.
* `image` - (Optional, List) Image storage configuration.

The `volume_template` object of `storage_mounts` supports the following:

* `reuse_key` - (Required, String) Reuse key for the volume template derivation.
* `volume_template_name` - (Required, String) Volume template name.

The `volume` object of `storage_mounts` supports the following:

* `volume_id` - (Required, String) Volume ID.
* `volume_name` - (Required, String) Volume name.

The `vpc_config` object of `network_configuration` supports the following:

* `security_group_ids` - (Optional, Set) Security group IDs.
* `subnet_ids` - (Optional, Set) VPC subnet IDs.

The `waa_configuration` object of `computer_configuration` supports the following:

* `image_id` - (Optional, String) WAA image ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Sandbox tool creation time, ISO8601 format.
* `default_timeout_seconds` - Default timeout in seconds.
* `image_digest` - Image digest of the custom configuration.
* `status_reason` - Reason for the current status.
* `status` - Sandbox tool status. Values: CREATING, ACTIVE, DELETING, FAILED.
* `tool_id` - Sandbox tool ID.
* `update_time` - Sandbox tool update time, ISO8601 format.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSandboxTool](https://www.tencentcloud.com/search/?wd=CreateSandboxTool&type=document&lang=en)
* Read: [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
* Update: [UpdateSandboxTool](https://www.tencentcloud.com/search/?wd=UpdateSandboxTool&type=document&lang=en)
* Delete: [DeleteSandboxTool](https://www.tencentcloud.com/search/?wd=DeleteSandboxTool&type=document&lang=en)

## Import

AGS sandbox tool can be imported using the toolId, e.g.

```
terraform import tencentcloud_ags_sandbox_tool.example sbt-xxxxxxxx
```

