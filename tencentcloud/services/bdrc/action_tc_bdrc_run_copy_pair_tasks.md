Provides an action to launch a group of BDRC (Business Data Resilient Disaster Recovery) copy pair tasks via the RunCopyPairTasks API. This is a one-time operation; no cloud-side state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

```hcl
action "tencentcloud_bdrc_run_copy_pair_tasks" "example" {
  config {
    copy_pair_ids = [
      "cvmcopypair-avp8kqtj",
    ]
    copy_pair_type = "INSTANCE"
  }
}
```
