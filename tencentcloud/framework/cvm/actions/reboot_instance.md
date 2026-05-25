Provides an action that requests a CVM instance reboot. Actions are
declarative side-effecting operations triggered by `lifecycle.action_trigger`
on a host resource. This reference action is currently a stub and only
validates the input shape and emits a tflog record; it does not call the
real CVM API. Replace the body of `Invoke` with a real CVM client call
when promoting it from L0 to a production action.

Example Usage

```hcl
action "tencentcloud_reboot_instance" "demo" {
  config {
    instance_id = "ins-abcdef01"
    force       = true
  }
}

resource "terraform_data" "trigger" {
  input = "v1"

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.tencentcloud_reboot_instance.demo]
    }
  }
}
```
