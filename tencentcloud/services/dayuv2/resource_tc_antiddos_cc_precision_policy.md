Provides a resource to create a antiddos cc_precision_policy

Example Usage

```hcl
resource "tencentcloud_antiddos_cc_precision_policy" "cc_precision_policy" {
  instance_id   = "bgpip-0000078h"
  ip            = "212.64.62.191"
  protocol      = "http"
  domain        = "t.baidu.com"
  policy_action = "drop"
  policy_list {
    field_type     = "value"
    field_name     = "cgi"
    value          = "a.com"
    value_operator = "equal"
  }

  policy_list {
    field_type     = "value"
    field_name     = "ua"
    value          = "test"
    value_operator = "equal"
  }
}
```

Import

antiddos cc_precision_policy can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_cc_precision_policy.cc_precision_policy ${instanceId}#${policyId}#${instanceIp}#${domain}#${protocol}
```