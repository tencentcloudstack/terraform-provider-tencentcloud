---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_unblock_resources"
sidebar_current: "docs-tencentcloud-action-antiddos_unblock_resources"
description: |-
  Provides an action to apply for unblocking blocked AntiDDoS (DDoS protection) resources (public IP list) via the UnblockResources API. This is a one-time operation; no cloud-side state is persisted after the action completes.
---

# tencentcloud_antiddos_unblock_resources

Provides an action to apply for unblocking blocked AntiDDoS (DDoS protection) resources (public IP list) via the UnblockResources API. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

## Example Usage

```hcl
action "tencentcloud_antiddos_unblock_resources" "example" {
  config {
    resources = [
      "117.175.94.230",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `resources` - (Required, List) List of public IPs to unblock. The list length is limited to 10.


