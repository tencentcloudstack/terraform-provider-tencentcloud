---
subcategory: "Agent Sandbox(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_tool"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_tool"
description: |-
  Provides a resource to create an AGS sandbox tool.
---

# tencentcloud_ags_sandbox_tool

Provides a resource to create an AGS sandbox tool.

-> **Note:** `tool_name`, `tool_type`, `default_timeout`, `client_token`, `role_arn`, `storage_mounts`, `log_configuration` and `persistent` are create-only: the update API does not accept them, so changing any of them destroys the existing sandbox tool and creates a new one.

-> **Note:** `default_timeout` and `client_token` are only accepted by the create API. The describe API reports the timeout in seconds (`DefaultTimeoutSeconds`) instead of the original duration string, so the values configured here are kept as they are and drift on them cannot be detected.

-> **Note:** Removing `description`, `tags`, `custom_configuration`, `computer_configuration` or `persistent` from the configuration keeps the value that currently exists on the cloud side, because the update API does not transmit empty values. To clear those attributes, delete and recreate the sandbox tool.

-> **Note:** Destroying this resource deletes the sandbox tool itself only. It does not terminate sandbox instances that were created from the tool and it does not release any storage (COS bucket, CFS file system or volume) referenced by `storage_mounts`.

## Example Usage

### Create a custom sandbox tool with a public network

```hcl
resource "tencentcloud_ags_sandbox_tool" "example" {
  tool_name = "tf-example-tool"
  tool_type = "custom"

  network_configuration {
    network_mode = "PUBLIC"
  }

  description = "sandbox tool created by terraform"

  custom_configuration {
    image               = "ccr.ccs.tencentyun.com/tf-example/agent:latest"
    image_registry_type = "personal"
    command             = ["/bin/sh", "-c"]
    args                = ["python3 -m http.server 8080"]

    resources {
      cpu     = "1"
      memory  = "2Gi"
      storage = "5Gi"
    }

    ports {
      name     = "http"
      port     = 8080
      protocol = "TCP"
    }

    env {
      name  = "LOG_LEVEL"
      value = "info"
    }
  }
}
```

### Create a sandbox tool in a VPC with a storage mount and log push

```hcl
resource "tencentcloud_ags_sandbox_tool" "vpc" {
  tool_name  = "tf-example-vpc-tool"
  tool_type  = "browser"
  persistent = true

  network_configuration {
    network_mode = "VPC"

    vpc_config {
      subnet_ids         = ["subnet-xxxxxxxx"]
      security_group_ids = ["sg-xxxxxxxx"]
    }
  }

  storage_mounts {
    name       = "workspace"
    mount_path = "/workspace"
    read_only  = false

    storage_source {
      cos {
        endpoint    = "cos.ap-guangzhou.myqcloud.com"
        bucket_name = "tf-example-bucket"
        bucket_path = "/workspace"
      }
    }
  }

  log_configuration {
    cls_config {
      topic_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    }

    log_sources {
      files = ["/logs/app.log"]
    }
  }

  tags {
    key   = "created_by"
    value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_configuration` - (Required, Object) Network configuration of the sandbox tool.
* `tool_name` - (Required, String) Name of the sandbox tool. Length 1-50 characters and it must be unique under the same AppId. It cannot be changed after creation.
* `tool_type` - (Required, String) Type of the sandbox tool. Valid values: `browser`, `code-interpreter`, `mobile`, `osworld`, `custom`, `swebench`, `aio`, `android-world`, `waa`. It cannot be changed after creation.
* `client_token` - (Optional, String) Idempotency token used by the create API. Maximum length 64 characters. It cannot be changed after creation.
* `computer_configuration` - (Optional, Object) Configuration of desktop-like sandbox tools.
* `custom_configuration` - (Optional, Object) Custom configuration of the sandbox tool.
* `default_timeout` - (Optional, String) Default timeout of sandbox instances created from this tool, for example `5m`, `300s` or `1h`. Maximum 24 hours. When it is not specified, the system default of 5 minutes is used. It cannot be changed after creation and it is not reported by the describe API.
* `description` - (Optional, String) Description of the sandbox tool. Maximum length 200 characters.
* `log_configuration` - (Optional, Object) Log push configuration of the sandbox tool. It cannot be changed after creation.
* `persistent` - (Optional, Bool) Whether the sandbox tool is persistent. The cloud side applies a default value when it is not specified. It cannot be changed after creation.
* `role_arn` - (Optional, String) Arn of the role bound to the sandbox tool. It cannot be changed after creation.
* `storage_mounts` - (Optional, List) Storage mount configuration of sandbox instances. It cannot be changed after creation.
* `tags` - (Optional, List) Tags of the sandbox tool.

The `computer_configuration` object supports the following:

* `os_world_configuration` - (Optional, Object) Built-in OSWorld configuration.
* `waa_configuration` - (Optional, Object) WAA sandbox configuration.

The `os_world_configuration` object of `computer_configuration` supports the following:

* `version` - (Optional, String) Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. Default value: `osworld1`.

The `waa_configuration` object of `computer_configuration` supports the following:

* `image_id` - (Optional, String) Custom WAA image id.

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
* `memory` - (Optional, String) Amount of memory resources.
* `storage` - (Optional, String) Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`.

The `security_context` object of `custom_configuration` supports the following:

* `capabilities` - (Optional, Object) POSIX capabilities that can be added to or removed from the running container.
* `privileged` - (Optional, Bool) Whether the sandbox runs in privileged mode. Default value: `true`.

The `capabilities` object of `security_context` supports the following:

* `add` - (Optional, Set) Capabilities to add. Reference value: `["SYS_ADMIN"]`.
* `drop` - (Optional, Set) Capabilities to drop. Reference value: `["ALL"]`.

The `log_configuration` object supports the following:

* `cls_config` - (Optional, Object) Log push configuration of CLS.
* `log_sources` - (Optional, Object) Log source configuration.

The `cls_config` object of `log_configuration` supports the following:

* `topic_id` - (Optional, String) Id of the CLS log topic used to push sandbox tool logs.

The `log_sources` object of `log_configuration` supports the following:

* `files` - (Optional, Set) Log files to collect. They must live under the `/logs/` directory, sub directories are not supported and at most 10 files are allowed.

The `network_configuration` object supports the following:

* `network_mode` - (Required, String) Network mode. Valid values: `PUBLIC`, `VPC`, `SANDBOX`.
* `vpc_config` - (Optional, Object) VPC related configuration. It is used when `network_mode` is `VPC`.

The `vpc_config` object of `network_configuration` supports the following:

* `security_group_ids` - (Optional, Set) List of security group ids.
* `subnet_ids` - (Optional, Set) List of VPC subnet ids.

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

The `tags` object supports the following:

* `key` - (Optional, String) Tag key.
* `value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of the sandbox tool in ISO8601 format.
* `id` - Id of the resource, which is the same as `tool_id`.
* `status_reason` - Reason why the sandbox tool is in its current status.
* `status` - Status of the sandbox tool. Valid values: `CREATING`, `ACTIVE`, `DELETING`, `FAILED`.
* `tool_id` - Id of the sandbox tool.
* `update_time` - Last update time of the sandbox tool in ISO8601 format.

The `image` object of `storage_source` exports the following:

* `digest` - Image digest resolved by the cloud side. It is not required in requests and is only reported by the cloud.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateSandboxTool](https://www.tencentcloud.com/search/?wd=CreateSandboxTool&type=document&lang=en)
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
* Read: [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
* Update:
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
  - [UpdateSandboxTool](https://www.tencentcloud.com/search/?wd=UpdateSandboxTool&type=document&lang=en)
* Delete:
  - [DeleteSandboxTool](https://www.tencentcloud.com/search/?wd=DeleteSandboxTool&type=document&lang=en)
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)

## Import

AGS sandbox tool can be imported using the sandbox tool id, e.g.

```
$ terraform import tencentcloud_ags_sandbox_tool.example tool-xxxxxxxx
```

