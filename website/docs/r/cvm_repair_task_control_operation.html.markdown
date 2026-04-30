---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_repair_task_control_operation"
sidebar_current: "docs-tencentcloud-resource-cvm_repair_task_control_operation"
description: |-
  Provides a resource to authorize a CVM repair task that is currently in the pending-authorization state, by invoking the `RepairTaskControl` cloud API.
---

# tencentcloud_cvm_repair_task_control_operation

Provides a resource to authorize a CVM repair task that is currently in the pending-authorization state, by invoking the `RepairTaskControl` cloud API.

This is an operation-style resource: only `Create` triggers the cloud-side action. `Read` and `Delete` are no-ops because the underlying API has no `Describe` or reverse counterpart. Any change to the input arguments forces the resource to be re-created (all fields are `ForceNew`). To re-trigger an authorization, use `terraform taint`.

Use the `tencentcloud_cvm_repair_tasks` data source to discover repair tasks (e.g., those with `task_status = 1`, pending authorization) before authorizing them with this resource.

## Example Usage

### Authorize a repair task immediately

```hcl
resource "tencentcloud_cvm_repair_task_control_operation" "demo" {
  product      = "CVM"
  instance_ids = ["ins-xxxxxxxx"]
  task_id      = "rep-xxxxxxxx"
  operate      = "AuthorizeRepair"
}
```

### Schedule the authorization at a future time

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

### Authorize with lossy local-disk migration

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

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, List: [`String`], ForceNew) List of instance IDs to operate on. Only repair tasks related to these instance IDs are authorized. Can be obtained from `InstanceId` in the `DescribeTaskInfo` API response.
* `operate` - (Required, String, ForceNew) Operation type. Currently only `AuthorizeRepair` is supported.
* `product` - (Required, String, ForceNew) Product type the pending-authorization task instance belongs to. Valid values: `CVM` (Cloud Virtual Machine), `CDH` (Cloud Dedicated Host), `CPM2.0` (Cloud Physical Machine 2.0).
* `task_id` - (Required, String, ForceNew) ID of the repair task to operate on. Can be obtained from `TaskId` in the `DescribeTaskInfo` API response.
* `order_auth_time` - (Optional, String, ForceNew) Scheduled authorization time, format `YYYY-MM-DD HH:MM:SS`. The scheduled time must be at least 5 minutes later than the current time and within 48 hours.
* `task_sub_method` - (Optional, String, ForceNew) Additional authorization handling strategy. When empty, the default authorization is used. For repair tasks supporting lossy migration, set to `LossyLocal` to allow lossy local-disk migration. WARNING: when `LossyLocal` is used on a local-disk instance, all local disk data will be wiped, equivalent to redeploying the local-disk instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

This resource is an operation-style resource and does not support `terraform import`. The cloud-side authorization is a one-shot side-effect with no `Describe` API to refresh state.

