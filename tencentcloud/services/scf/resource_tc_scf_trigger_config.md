Provides a resource to create a scf trigger_config

~> **NOTE:** Use of the current resource is no longer recommended; `tencentcloud_scf_trigger` is recommended instead.

Example Usage

```hcl

resource "tencentcloud_scf_trigger_config" "trigger_config" {
  enable        = "OPEN"
  function_name = "keep-1676351130"
  trigger_name  = "SCF-timer-1685540160"
  type          = "timer"
  qualifier     = "$DEFAULT"
  namespace     = "default"
  trigger_desc = "* 1 2 * * * *"
  description = "func"
  custom_argument = "Information"
}

```

Import

scf trigger_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_trigger_config.trigger_config functionName#namespace#triggerName
```