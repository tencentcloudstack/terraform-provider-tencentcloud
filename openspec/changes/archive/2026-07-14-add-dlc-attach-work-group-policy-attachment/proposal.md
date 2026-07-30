## Why

The `dlc` (Data Lake Compute) service exposes the ability to bind/unbind
authorization policies to work groups via the cloud APIs `AttachWorkGroupPolicy`
and `DetachWorkGroupPolicy`, but the Terraform provider only has thin
"operation"-style resources (`tencentcloud_dlc_attach_work_group_policy_operation`
and `tencentcloud_dlc_detach_work_group_policy_operation`) whose Read/Delete
callbacks are no-ops, so the bound state drifts unmanaged. There is no proper
`RESOURCE_KIND_ATTACHMENT` resource that tracks the binding lifecycle
(create = bind, read = verify membership, delete = unbind).

This change adds `tencentcloud_dlc_attach_work_group_policy_attachment` to
declaratively manage the bind/unbind of one or more policies to a work group.

## What Changes

- Add a new Terraform resource
  `tencentcloud_dlc_attach_work_group_policy_attachment` of kind
  `RESOURCE_KIND_ATTACHMENT` (Create / Read / Delete, all args `ForceNew`).
  - Create: calls `AttachWorkGroupPolicy` (`WorkGroupId` + `PolicySet`).
  - Read: calls `DescribeWorkGroupInfo` (`Type=DataAuth`) and verifies that the
    attached policies (by `PolicyId`) are present in the work group's
    `DataPolicyInfo.PolicySet`.
  - Delete: calls `DetachWorkGroupPolicy` (`WorkGroupId` + `PolicySet` and/or
    `PolicyIds`).
- The resource file is named `resource_tc_dlc_attach_work_group_policy_attachment.go`
  (the conventional `resource_tc_<Product>_<name>_attachment.go` format).
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add a gomonkey-based unit test file
  `resource_tc_dlc_attach_work_group_policy_attachment_test.go`.
- Add a doc stub `resource_tc_dlc_attach_work_group_policy_attachment.md`.

## Capabilities

### New Capabilities
- `dlc-attach-work-group-policy-attachment`: Manage the bind/unbind of one or
  more authorization policies (`PolicySet`) to a DLC work group (`WorkGroupId`)
  as a Terraform attachment resource, verified through
  `DescribeWorkGroupInfo`.

### Modified Capabilities
<!-- None. -->

## Impact

- **New files**:
  - `tencentcloud/services/dlc/resource_tc_dlc_attach_work_group_policy_attachment.go`
  - `tencentcloud/services/dlc/resource_tc_dlc_attach_work_group_policy_attachment_test.go`
  - `tencentcloud/services/dlc/resource_tc_dlc_attach_work_group_policy_attachment.md`
- **Modified files**:
  - `tencentcloud/provider.go` (register the new resource key)
  - `tencentcloud/provider.md` (add the resource to the DLC Resource section for
    gendoc)
- **SDK**: reuses the vendored
  `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125`
  client (`AttachWorkGroupPolicy`, `DescribeWorkGroupInfo`,
  `DetachWorkGroupPolicy`). No SDK upgrade required.
- **Backward compatibility**: purely additive; existing
  `tencentcloud_dlc_attach_work_group_policy_operation` resource is untouched.
