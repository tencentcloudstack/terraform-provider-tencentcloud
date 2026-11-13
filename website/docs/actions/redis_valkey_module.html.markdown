---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_valkey_module"
sidebar_current: "docs-tencentcloud-action-redis_valkey_module"
description: |-
  Provides an action to modify (enable) an extension module of a Valkey (Redis) instance, via the Redis `ModifyValkeyModule` API.
---

# tencentcloud_redis_valkey_module

Provides an action to modify (enable) an extension module of a Valkey (Redis) instance, via the Redis `ModifyValkeyModule` API.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

~> **NOTE:** This action executes a one-time cloud-side side effect. The enabled module persists after the action completes; destroying the action does not disable the module. Repeating the invocation repeats the cloud-side side effect. The `ModifyValkeyModule` API is asynchronous: the action polls `DescribeTaskInfo` with the returned task id until the task reaches a terminal state, subject to the provider write retry timeout.

## Example Usage

```hcl
action "tencentcloud_redis_valkey_module" "example" {
  config {
    instance_id = "crs-1111222233334444"
    name        = "ValkeyBloom"
    status      = "ON"
  }
}
```



```hcl
shell
# Trigger explicitly (Terraform 1.14+); use `terraform plan -invoke=...` for a dry-run first
terraform apply - invoke = action.tencentcloud_redis_valkey_module.example
```

### is not a destroy-protection guardrail (in a full destroy walk an action error is downgraded to a warning):

```hcl
resource "tencentcloud_redis_instance" "example" {
  # ...

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.tencentcloud_redis_valkey_module.example]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) ID of the instance to modify, e.g. `crs-xjhsdj****`. Please log in to the Redis console and copy the instance ID from the instance list.
* `name` - (Required, String) The extension module to modify. Valid values: `ValkeyBloom`, `ValkeyJSON`, `ValkeySearch`.
* `status` - (Required, String) Sets the enabled status of the extension module. Only `ON` (enable the module) is supported.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke: [ModifyValkeyModule](https://www.tencentcloud.com/document/api/239/85838)

