Provides an action to dispatch FWM (Firewall Manager) pending strategies for a product via the DispatchStrategy API (下发规则). This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources.

The dispatch submits every pending (new / pending-dispatch) strategy of the current account x product line to the receiving accounts; it cannot target a single rule ID. The cloud API has no idempotency key and does not block an already-running task, so the provider first queries the latest dispatch status and reuses an in-progress task instead of inserting another one. Dispatches of the same product should still be serialized to avoid concurrent triggers. Create success does NOT mean dispatch success: the action polls `DescribeStrategyDispatchStatus` to a terminal state (succeeded / failed) before returning. To re-dispatch, taint the action or change a trigger field; destroy does not roll back the dispatched rules and the far-side effect persists.

Example Usage

Dispatch pending strategies after creating strategy resources

```hcl
# Strategy resources must exist before dispatching them; establish the
# ordering by referencing them directly.
resource "tencentcloud_fwm_strategy" "example" {
  product = "cfw_edge_acl"

  pre_strategy {
    group_id = "fwmrg-qww1luux7z"
  }
}

action "tencentcloud_fwm_dispatch_strategy" "example" {
  config {
    product = "cfw_edge_acl"
    uin     = "100011949846"
  }

  # Ensures the strategies above are created before this dispatch runs.
  depends_on = [tencentcloud_fwm_strategy.example]
}
```

```shell
# Trigger explicitly (Terraform 1.14+); use `terraform plan -invoke=...` for a dry-run first.
terraform apply -invoke=action.tencentcloud_fwm_dispatch_strategy.example
```

Trigger automatically during apply through `lifecycle.action_trigger`

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product = "cfw_edge_acl"

  pre_strategy {
    group_id = "fwmrg-qww1luux7z"
  }

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.tencentcloud_fwm_dispatch_strategy.example]
      on_failure = halt
    }
  }
}

action "tencentcloud_fwm_dispatch_strategy" "example" {
  config {
    product = "cfw_edge_acl"
    uin     = "100011949846"
  }
}
```

Re-dispatch by tainting the action or changing a trigger field

```hcl
action "tencentcloud_fwm_dispatch_strategy" "example" {
  config {
    product = "cfw_edge_acl"
    uin     = "100011949846"
    # Changing sub_account_uin (or re-tainting the action) forces a new invocation.
    sub_account_uin = "100049315301"
  }
}
```
