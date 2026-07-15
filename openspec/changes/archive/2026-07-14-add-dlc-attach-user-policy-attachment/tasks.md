## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policyr_attachment.go` with the `ResourceTencentCloudDlcAttachUserPolicyrAttachment` schema function, defining `user_id` (Required, ForceNew), `policy_set` (Required, list of Policy objects with input-eligible fields), and `account_type` (Optional, ForceNew). Follow the code style of `tencentcloud_igtm_strategy`.
- [x] 1.2 Implement `resourceTencentCloudDlcAttachUserPolicyrAttachmentCreate`: build `AttachUserPolicyRequest` from schema, call `AttachUserPolicyWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError`, validate the response is non-nil, set the composite ID (`user_id#account_type` via `tccommon.FILED_SP`) outside the retry block, and set `policy_set` from the response.
- [x] 1.3 Implement `resourceTencentCloudDlcAttachUserPolicyrAttachmentRead`: split the composite ID into `user_id` and `account_type`, call `DescribeUserInfoWithContext` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`, verify the bound policies are present in the returned `UserDetailInfo` (Type=DataAuth), set fields only when response fields are non-nil, and on empty result emit `log.Printf("[CRUD] ...")` preserving the id before `d.SetId("")`.
- [x] 1.4 Implement `resourceTencentCloudDlcAttachUserPolicyrAttachmentUpdate` as an immutable guard: since this is a CRD-only resource, set only `Id()` as ForceNew; collect all other top-level args into `immutableArgs` and return an error if any changed.
- [x] 1.5 Implement `resourceTencentCloudDlcAttachUserPolicyrAttachmentDelete`: split the composite ID, build `DetachUserPolicyRequest` with `UserId`, `AccountType`, and `PolicySet` from state, call `DetachUserPolicyWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError`.
- [x] 1.6 Ensure all returned errors are checked; for functions that cannot fail, assign the error to `_ = func()` to avoid unused-variable errors. Do not add file-header comments.

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_dlc_attach_user_policyr_attachment` in `tencentcloud/provider.go` (add the resource map entry pointing to `dlc.ResourceTencentCloudDlcAttachUserPolicyrAttachment`), referencing how `tencentcloud_igtm_strategy` is registered.
- [x] 2.2 Add the resource entry to `tencentcloud/provider.md`.

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policyr_attachment.md` following `gendoc/README.md` and other resources' `.md` format: a one-line description mentioning DLC, Example Usage (use `jsonencode()` for any json-string field values), and an Import section (RESOURCE_KIND_ATTACHMENT includes import; document the composite id requirement). Do NOT add `Argument Reference` / `Attribute Reference` sections (auto-generated).

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policyr_attachment_test.go` using gomonkey mocks (no Terraform acceptance test suite) to test the business logic of Create, Read, and Delete handlers.
- [x] 4.2 Run the unit tests with `go test -gcflags=all=-l` against the test file and ensure they pass.

## 5. Verification

- [x] 5.1 Verify the generated Go code compiles conceptually against the vendor DLC SDK models (`AttachUserPolicyRequest`, `DescribeUserInfoRequest`/`Response`, `DetachUserPolicyRequest`, `Policy`, `UserDetailInfo`, `Policys`, `Filter`) and that every schema field used in Create exists in `AttachUserPolicyRequest`, every field used in Read exists in `DescribeUserInfoRequest`/`DescribeUserInfoResponse`, and every field used in Delete exists in `DetachUserPolicyRequest`.
- [x] 5.2 Confirm no `website/` files are edited directly (documentation is generated via `make doc` in the finalize phase) and no `.changelog/` files are created outside the finalize phase.
