Provides an action to dispatch FWM (Firewall Manager) pending strategies for a product via the DispatchStrategy API (下发规则). This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources.

Example Usage

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

```hcl
action "tencentcloud_fwm_dispatch_strategy" "example" {
  config {
    product = "cfw_edge_acl"
    status  = 1

    timeouts {
      invoke = "3m"
    }
  }
}
```
