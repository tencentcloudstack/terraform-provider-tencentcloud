# dlc-attach-work-group-policy-attachment Specification

## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_dlc_attach_work_group_policy_attachment`

The provider SHALL register a new attachment resource named `tencentcloud_dlc_attach_work_group_policy_attachment` whose Create / Read / Delete callbacks manage the bind/unbind of authorization policies to a DLC work group. The resource MUST expose only Create / Read / Delete (no Update); every argument SHALL be `ForceNew`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_dlc_attach_work_group_policy_attachment"` mapped to `dlc.ResourceTencentCloudDlcAttachWorkGroupPolicyAttachment()`, alongside existing `tencentcloud_dlc_*` resource entries.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the DLC Resource section MUST include `tencentcloud_dlc_attach_work_group_policy_attachment` so that `website/docs/r/dlc_attach_work_group_policy_attachment.html.markdown` is generated.

### Requirement: Schema MUST mirror the AttachWorkGroupPolicy API input

The resource schema SHALL declare these top-level argument keys, with semantics matching the SDK request fields of `AttachWorkGroupPolicyRequestParams`:

| HCL key | SDK field | Type | Required | ForceNew |
|---|---|---|---|---|
| `work_group_id` | `WorkGroupId` | TypeInt | Yes | Yes |
| `policy_set` | `PolicySet` | TypeList (MaxItems 1) | Yes | Yes |

The `policy_set` block SHALL be limited to a single element (`MaxItems: 1`) so that one resource instance represents exactly one policy binding. The nested `policy_set` schema SHALL reuse the same field set as the existing `tencentcloud_dlc_attach_work_group_policy_operation` resource (`database`, `catalog`, `table`, `operation`, `policy_type`, `function`, `view`, `column`, `data_engine`, `re_auth`, `source`, `mode`, `operator`, `create_time`, `source_id`, `source_name`, `id`).

#### Scenario: Required fields enforce on plan

- **WHEN** the user writes a config that omits `work_group_id` or `policy_set`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: All arguments are ForceNew

- **WHEN** the user changes any argument (`work_group_id` or any nested `policy_set` field) after creation
- **THEN** Terraform SHALL plan to destroy and recreate the resource (no in-place update path exists).

### Requirement: Create SHALL call AttachWorkGroupPolicy and store a composite ID

The Create callback SHALL build an `AttachWorkGroupPolicyRequest` from the schema (`WorkGroupId`, `PolicySet` constructed from the single `policy_set` block), call `AttachWorkGroupPolicy` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError`, and guard that the response and `Response.PolicySet` are non-nil and non-empty. The callback SHALL then verify that the first (only) returned policy's `PolicyId` is non-empty; if it is empty, the callback SHALL return `NonRetryableError`. On success the resource ID SHALL be set to `WorkGroupId + "#" + PolicyId` (joined with `tccommon.FILED_SP`).

#### Scenario: Successful create stores composite ID

- **WHEN** `AttachWorkGroupPolicy` returns a `PolicySet` whose first element has `PolicyId = "policy-xxxx"`
- **AND** the configured `work_group_id` is `23184`
- **THEN** the resource `d.Id()` SHALL equal `"23184#policy-xxxx"`.

#### Scenario: Empty PolicyId is rejected

- **WHEN** `AttachWorkGroupPolicy` returns a `PolicySet` whose first element has an empty `PolicyId`
- **THEN** Create SHALL return a `NonRetryableError` and SHALL NOT write an ID to state.

#### Scenario: API error is retried

- **WHEN** `AttachWorkGroupPolicy` returns a retryable error
- **THEN** the `resource.Retry` loop SHALL re-invoke the API until success or `WriteRetryTimeout` is exhausted.

### Requirement: Read SHALL verify the binding via DescribeWorkGroupInfo

The Read callback SHALL split the composite ID (`WorkGroupId#PolicyId`), call `DescribeWorkGroupInfo` with `WorkGroupId` and `Type="DataAuth"` (paginating with `Limit=100` / `Offset` increments until the target `PolicyId` is found or pages are exhausted), and confirm the stored `PolicyId` exists in `WorkGroupInfo.DataPolicyInfo.PolicySet`. If the policy is found, the schema fields (`work_group_id` and the nested `policy_set`) SHALL be repopulated with nil-guarded `d.Set` calls. If the policy is NOT found, the callback SHALL first log `log.Printf("[CRUD] dlc_attach_work_group_policy_attachment id=%s", d.Id())` and then call `d.SetId("")`.

#### Scenario: Binding exists repopulates state

- **WHEN** `DescribeWorkGroupInfo` returns a `DataPolicyInfo.PolicySet` containing a policy with `PolicyId = "policy-xxxx"`
- **AND** the resource ID is `"23184#policy-xxxx"`
- **THEN** Read SHALL set `work_group_id` to `23184` and SHALL populate the `policy_set` block from the matched policy, and SHALL keep `d.Id()` unchanged.

#### Scenario: Binding missing clears state with logging

- **WHEN** `DescribeWorkGroupInfo` returns a `DataPolicyInfo.PolicySet` that does NOT contain the stored `PolicyId`
- **THEN** Read SHALL emit `[CRUD] dlc_attach_work_group_policy_attachment id=23184#policy-xxxx` to the log and SHALL call `d.SetId("")`.

### Requirement: Delete SHALL call DetachWorkGroupPolicy using PolicyIds

The Delete callback SHALL split the composite ID, build a `DetachWorkGroupPolicyRequest` with `WorkGroupId` and `PolicyIds=[policyId]`, and call `DetachWorkGroupPolicy` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError`. The callback SHALL guard that the response and `Response` are non-nil.

#### Scenario: Successful delete unbinds the policy

- **WHEN** `DetachWorkGroupPolicy` succeeds for `WorkGroupId=23184`, `PolicyIds=["policy-xxxx"]`
- **THEN** Delete SHALL return nil without error.

#### Scenario: Delete API error is retried

- **WHEN** `DetachWorkGroupPolicy` returns a retryable error
- **THEN** the `resource.Retry` loop SHALL re-invoke the API until success or `WriteRetryTimeout` is exhausted.

### Requirement: Resource SHALL support import by composite ID

The resource SHALL declare `Importer` with `schema.ImportStatePassthrough`, so that `terraform import tencentcloud_dlc_attach_work_group_policy_attachment.example 23184#policy-xxxx` imports the binding. The `.md` doc MUST document that the composite ID (`WorkGroupId#PolicyId`) is required for import.

#### Scenario: Import via composite ID

- **WHEN** the user runs `terraform import tencentcloud_dlc_attach_work_group_policy_attachment.example 23184#policy-xxxx`
- **THEN** the provider SHALL call Read with `d.Id() = "23184#policy-xxxx"` and populate state from the live binding.

### Requirement: Unit tests SHALL use gomonkey mocks (no terraform test suite)

The test file `resource_tc_dlc_attach_work_group_policy_attachment_test.go` SHALL mock the DLC client methods (`AttachWorkGroupPolicy`, `DescribeWorkGroupInfo`, `DetachWorkGroupPolicy`) using `gomonkey` and assert the business logic of Create / Read / Delete (including the composite-ID and drift-clear paths) without requiring real cloud credentials. The tests SHALL be runnable via `go test -gcflags=all=-l`.

#### Scenario: Mocked Create test

- **WHEN** `AttachWorkGroupPolicy` is mocked to return a `PolicySet` with `PolicyId = "policy-xxxx"`
- **THEN** the Create test SHALL assert `d.Id() == "23184#policy-xxxx"` and `assert.NoError` on the returned error.

#### Scenario: Mocked drift-clear Read test

- **WHEN** `DescribeWorkGroupInfo` is mocked to return a `DataPolicyInfo.PolicySet` without the stored `PolicyId`
- **THEN** the Read test SHALL assert `d.Id() == ""` after Read returns.
