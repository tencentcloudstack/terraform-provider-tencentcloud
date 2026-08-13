## 1. Service Layer

- [x] 1.1 Append `DescribeDlcWorkGroupPolicyAttachmentById(ctx, workGroupId, policyId string) (*dlc.Policy, error)` to `service_tencentcloud_dlc.go` — wraps `DescribeWorkGroupInfo` (`WorkGroupId`, `Type="DataAuth"`, `Limit=100`, `Offset` pagination) with `resource.Retry` + `ratelimit.Check`; iterates `WorkGroupInfo.DataPolicyInfo.PolicySet` to find the policy whose `PolicyId` matches `policyId`; returns the matched `*dlc.Policy` (nil if not found), nil/length-safe
- [x] 1.2 Append `DeleteDlcAttachWorkGroupPolicyAttachmentByPolicyId(ctx, workGroupId, policyId string) error` to `service_tencentcloud_dlc.go` — wraps `DetachWorkGroupPolicy` (`WorkGroupId`, `PolicyIds=[policyId]`) with `resource.Retry(tccommon.WriteRetryTimeout)` + `ratelimit.Check`; guards response/Response non-nil; logs on failure

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_dlc_attach_work_group_policy_attachment.go` with `ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment()` schema: top-level `work_group_id` (TypeInt, Required, ForceNew) and `policy_set` (TypeList, Required, ForceNew, MaxItems 1) reusing the same nested field set as `resource_tc_dlc_attach_work_group_policy_operation.go`; declare `Importer` with `schema.ImportStatePassthrough`; Create/Read/Delete only (no Update)
- [x] 2.2 Implement `resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentCreate`: build `AttachWorkGroupPolicyRequest` from schema (`WorkGroupId`, `PolicySet` from the single block), call `AttachWorkGroupPolicy` in `resource.Retry(WriteRetryTimeout)` with `tccommon.RetryError`; guard response/`Response.PolicySet` non-nil & non-empty; check first element `PolicyId` non-empty (else `NonRetryableError` after `log.Printf` of logId & d.Id()); `d.SetId(workGroupId + tccommon.FILED_SP + policyId)`; return Read
- [x] 2.3 Implement `resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentRead`: split `d.Id()` by `tccommon.FILED_SP` (require 2 parts); call `DescribeDlcWorkGroupPolicyAttachmentById`; if nil → `log.Printf("[CRUD] dlc_attach_work_group_policy_attachment id=%s", d.Id())` then `d.SetId("")`; else repopulate `work_group_id` and nested `policy_set` with nil-guarded `d.Set`
- [x] 2.4 Implement `resourceTencentCloudDlcAttachWorkGroupPolicyAttachmentDelete`: split `d.Id()`; call `DeleteDlcAttachWorkGroupPolicyAttachmentByPolicyId`; return nil on success

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_dlc_attach_work_group_policy_attachment` in `tencentcloud/provider.go` ResourcesMap (map key → `dlc.ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment()`)
- [x] 3.2 Add `tencentcloud_dlc_attach_work_group_policy_attachment` to the DLC Resource section of the `tencentcloud/provider.go` file comment (for gendoc index → `tencentcloud/provider.md`)

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/dlc/resource_tc_dlc_attach_work_group_policy_attachment.md`: one-line description ("Provides a resource to create a DLC attach work group policy attachment"), Example Usage (with `policy_set` block), and Import section documenting the composite ID `WorkGroupId#PolicyId` (no Argument/Attribute Reference sections)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/dlc/resource_tc_dlc_attach_work_group_policy_attachment_test.go` using gomonkey mocks (no terraform test suite): mock `UseDlcClient`, `AttachWorkGroupPolicy`, `DescribeWorkGroupInfo`, `DetachWorkGroupPolicy`; test Create (assert composite ID), Read (binding exists / drift-clear), Delete (success); runnable via `go test -gcflags=all=-l`
- [x] 5.2 Run `go test ./tencentcloud/services/dlc/ -run "TestDlcAttachWorkGroupPolicyAttachment" -v -count=1 -gcflags=all=-l` and ensure all cases pass

## 6. Verification

- [x] 6.1 Confirm `proposal.md`, `design.md`, `specs/dlc-attach-work-group-policy-attachment/spec.md`, and `tasks.md` are complete and consistent
- [x] 6.2 Confirm code correctness: `AttachWorkGroupPolicy` params exist on create, `DetachWorkGroupPolicy` params (`PolicyIds`) exist on delete, `DescribeWorkGroupInfo` params (`WorkGroupId`, `Type`) exist on read
