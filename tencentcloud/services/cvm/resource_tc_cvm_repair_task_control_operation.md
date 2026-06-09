Provides a resource to authorize a CVM repair task that is currently in the pending-authorization state, by invoking the `RepairTaskControl` cloud API.

This is an operation-style resource: only `Create` triggers the cloud-side action. `Read` and `Delete` are no-ops because the underlying API has no `Describe` or reverse counterpart. Any change to the input arguments forces the resource to be re-created (all fields are `ForceNew`). To re-trigger an authorization, use `terraform taint`.

Use the `tencentcloud_cvm_repair_tasks` data source to discover repair tasks (e.g., those with `task_status = 1`, pending authorization) before authorizing them with this resource.

Example Usage

Authorize a repair task immediately

```hcl
resource "tencentcloud_cvm_repair_task_control_operation" "demo" {
  product      = "CVM"
  instance_ids = ["ins-xxxxxxxx"]
  task_id      = "rep-xxxxxxxx"
  operate      = "AuthorizeRepair"
}
```

Schedule the authorization at a future time

The scheduled time must be at least 5 minutes later than the current time and within 48 hours.

```hcl
resource "tencentcloud_cvm_repair_task_control_operation" "scheduled" {
  product         = "CVM"
  instance_ids    = ["ins-xxxxxxxx"]
  task_id         = "rep-xxxxxxxx"
  operate         = "AuthorizeRepair"
  order_auth_time = "2030-01-01 12:00:00"
}
```

Authorize with lossy local-disk migration

WARNING: setting task_sub_method to LossyLocal on a local-disk instance will WIPE ALL LOCAL DISK DATA, equivalent to redeploying the local-disk instance. Make sure to back up critical data and review your /etc/fstab mounts (consider adding nofail) before proceeding.

```hcl
resource "tencentcloud_cvm_repair_task_control_operation" "lossy_local" {
  product         = "CVM"
  instance_ids    = ["ins-xxxxxxxx"]
  task_id         = "rep-xxxxxxxx"
  operate         = "AuthorizeRepair"
  task_sub_method = "LossyLocal"
}
```

Import

This resource is an operation-style resource and does not support `terraform import`. The cloud-side authorization is a one-shot side-effect with no `Describe` API to refresh state.
