Provides an action to modify (enable) an extension module of a Valkey (Redis) instance, via the Redis `ModifyValkeyModule` API.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

~> **NOTE:** This action executes a one-time cloud-side side effect. The enabled module persists after the action completes; destroying the action does not disable the module. Repeating the invocation repeats the cloud-side side effect. The `ModifyValkeyModule` API is asynchronous: the action polls `DescribeTaskInfo` with the returned task id until the task reaches a terminal state, subject to the provider write retry timeout.

Example Usage

```hcl
action "tencentcloud_redis_valkey_module" "example" {
  config {
    instance_id = "crs-1111222233334444"
    name        = "ValkeyBloom"
    status      = "ON"
  }
}
```

```shell
# Trigger explicitly (Terraform 1.14+); use `terraform plan -invoke=...` for a dry-run first
terraform apply -invoke=action.tencentcloud_redis_valkey_module.example
```

The action can also be triggered implicitly from a resource lifecycle block. It runs automatically during apply for the configured events, including on drift-driven `after_update` events, and Terraform does not guarantee the remote operation happens only once. A failed action does not roll back the cloud-side side effect that already happened; the `on_failure` meta argument accepts `halt` (default), `continue` and `taint`, and `halt` is not a destroy-protection guardrail (in a full destroy walk an action error is downgraded to a warning):

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
