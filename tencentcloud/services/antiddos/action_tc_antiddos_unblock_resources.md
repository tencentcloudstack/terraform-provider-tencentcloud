Provides an action to apply for unblocking blocked AntiDDoS (DDoS protection) resources (public IP list) via the UnblockResources API. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_antiddos_unblock_resources" "example" {
  config {
    resources = [
      "117.175.94.230",
    ]
  }
}
```