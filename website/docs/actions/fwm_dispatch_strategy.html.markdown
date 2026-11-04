---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_dispatch_strategy"
sidebar_current: "docs-tencentcloud-action-fwm_dispatch_strategy"
description: |-
  Provides an action to dispatch FWM (Firewall Manager) pending strategies for a product via the DispatchStrategy API (下发规则). This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources.
---

# tencentcloud_fwm_dispatch_strategy

Provides an action to dispatch FWM (Firewall Manager) pending strategies for a product via the DispatchStrategy API (下发规则). This is a one-time operation; it does not create, update or destroy any Terraform-managed resource and exposes no attributes referenceable from other resources.

## Example Usage

### ~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

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

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) Product line to dispatch, e.g. `enterprise_sg` (enterprise security group), `cfw_vpc_acl` (VPC border firewall), `cfw_nat_acl` (NAT border firewall), `cfw_edge_acl` (internet border firewall). The value is forwarded to the cloud as-is.
* `status` - (Required, Int) Dispatch status of the DispatchStrategy API: `1` - dispatch, `2` - abort. The value is forwarded to the cloud as-is.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the invoke operation.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke:
  - [DescribeStrategyDispatchStatus](https://www.tencentcloud.com/search/?wd=DescribeStrategyDispatchStatus&type=document&lang=en)
  - [DispatchStrategy](https://www.tencentcloud.com/search/?wd=DispatchStrategy&type=document&lang=en)

