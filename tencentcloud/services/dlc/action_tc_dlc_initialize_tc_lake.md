Provides an action to activate (open) the DLC (Data Lake Compute) TCLake service via the InitializeTCLake API. This is a one-time, idempotent initialization operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_dlc_initialize_tc_lake" "example" {
  config {}
}
```
