# dlc-attach-user-policy-attachment Specification

## Purpose
TBD - created by archiving change add-dlc-attach-user-policy-attachment. Update Purpose after archive.
## Requirements
### Requirement: Resource shall manage DLC user-policy binding lifecycle
The Terraform provider SHALL provide a resource `tencentcloud_dlc_attach_user_policy_attachment` of kind `RESOURCE_KIND_ATTACHMENT` that manages the binding between a DLC user and exactly one authorization policy. The resource SHALL support Create (bind), Read (query), and Delete (unbind). The resource SHALL NOT support in-place Update; any change to a top-level argument SHALL trigger recreation.

#### Scenario: Bind a single policy to a user on create
- **WHEN** a user applies a configuration defining `tencentcloud_dlc_attach_user_policy_attachment` with `user_id`, `policy_set` (containing exactly one policy block), and `account_type`
- **THEN** the provider SHALL call the DLC `AttachUserPolicy` API with the given `UserId`, `PolicySet`, and `AccountType`
- **AND** the provider SHALL read the single returned `PolicyId` from the API response
- **AND** the provider SHALL set the resource ID to the composite of `user_id` and the returned `policy_id` joined by `tccommon.FILED_SP`
- **AND** the provider SHALL set `policy_set` in state from the API response

#### Scenario: Read reflects the currently bound policy
- **WHEN** the provider reads an existing `tencentcloud_dlc_attach_user_policy_attachment` resource
- **THEN** the provider SHALL call the DLC `DescribeUserInfo` API with the `UserId` derived from the composite ID and the `AccountType` from state
- **AND** the provider SHALL locate the policy in the response whose `PolicyId` matches the `policy_id` derived from the composite ID
- **AND** if the matching policy is no longer present, the provider SHALL emit a `[CRUD]` log preserving the id and then set `d.SetId("")`

#### Scenario: Unbind the policy on delete
- **WHEN** a user destroys a `tencentcloud_dlc_attach_user_policy_attachment` resource
- **THEN** the provider SHALL call the DLC `DetachUserPolicy` API with the `UserId` derived from the composite ID, the `AccountType` from state, and `PolicyIds` set to the single `policy_id` derived from the composite ID
- **AND** the resource SHALL be removed from state

#### Scenario: Any change recreates the resource
- **WHEN** a user changes any top-level argument of an existing `tencentcloud_dlc_attach_user_policy_attachment` resource
- **THEN** the provider SHALL reject the in-place update (the update method returns an error for immutable args) and Terraform SHALL recreate the resource

### Requirement: Resource SHALL use a composite identifier of user_id and policy_id
Because the DLC `AttachUserPolicy` API is bound to a single policy per call and returns a deterministic `PolicyId` for that policy, the resource SHALL use a composite ID composed of `user_id` and the returned `policy_id`, joined by `tccommon.FILED_SP`. The Read and Delete operations SHALL split the composite ID to recover `user_id` and `policy_id` for use as API request parameters.

#### Scenario: Composite ID is split for API calls
- **WHEN** the provider performs Read or Delete on `tencentcloud_dlc_attach_user_policy_attachment`
- **THEN** the provider SHALL split `d.Id()` by `tccommon.FILED_SP` into `user_id` and `policy_id`
- **AND** the provider SHALL use `user_id` as the `UserId` request parameter and `policy_id` to identify the target policy (via `PolicyIds` for Delete, or by matching `PolicyId` in the Read response)

### Requirement: Resource schema SHALL expose binding parameters and limit policy_set to a single policy
The resource schema SHALL expose `user_id` (Required, ForceNew), `policy_set` (Required, ForceNew, MinItems 1, MaxItems 1), and `account_type` (Optional, ForceNew). The `policy_set` SHALL be a list containing exactly one policy object whose input-eligible fields (`database`, `catalog`, `table`, `operation`, `policy_type`, `function`, `view`, `column`, `data_engine`, `re_auth`, `engine_generation`, `model`, `policy_id`) can be set by the user.

#### Scenario: Required arguments are enforced
- **WHEN** a user omits `user_id` or `policy_set` in the configuration
- **THEN** Terraform SHALL report a validation error before any API call is made

#### Scenario: policy_set is limited to exactly one policy
- **WHEN** a user configures `policy_set` with more than one policy block
- **THEN** Terraform SHALL report a schema validation error (`MaxItems` exceeded) before any API call is made
- **AND** as a defensive measure, the provider's Create handler SHALL also reject a `PolicySet` whose length is not exactly 1 before calling the DLC API

#### Scenario: account_type is optional
- **WHEN** a user does not specify `account_type`
- **THEN** the provider SHALL not set `AccountType` in the API request, allowing the cloud API to apply its default

### Requirement: Cloud API calls SHALL use retry and proper error handling
Create, Read, and Delete SHALL wrap their cloud API calls in `resource.Retry` using `tccommon.WriteRetryTimeout` (for Create/Delete) and `tccommon.ReadRetryTimeout` (for Read). Errors from the cloud API SHALL be wrapped with `tccommon.RetryError`. State mutations (setting the ID and fields) SHALL occur outside the retry block, after successful retry completion.

#### Scenario: Transient API failure is retried
- **WHEN** a DLC API call fails with a retryable error during Create, Read, or Delete
- **THEN** the provider SHALL retry the call within the configured timeout and wrap the error with `tccommon.RetryError`

#### Scenario: Create validates the response contains exactly one policy with a non-empty PolicyId
- **WHEN** the `AttachUserPolicy` API returns a nil response, or a `PolicySet` whose length is not 1, or a policy with an empty `PolicyId`
- **THEN** the provider SHALL return a `NonRetryableError` rather than writing an empty or invalid id to state

### Requirement: Resource SHALL be registered in the provider
The provider SHALL register `tencentcloud_dlc_attach_user_policy_attachment` in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md`. A documentation file `resource_tc_dlc_attach_user_policy_attachment.md` SHALL exist under `tencentcloud/services/dlc/`.

#### Scenario: Resource is usable in a Terraform configuration
- **WHEN** a user references `tencentcloud_dlc_attach_user_policy_attachment` in a `.tf` file
- **THEN** the provider SHALL recognize the resource type and execute its CRUD handlers

### Requirement: Unit tests SHALL use gomonkey mocks
The resource test file `resource_tc_dlc_attach_user_policy_attachment_test.go` SHALL use gomonkey to mock the DLC cloud API client methods and test the business logic of Create/Read/Delete without relying on the Terraform acceptance test suite. The tests SHALL be runnable with `go test -gcflags=all=-l`.

#### Scenario: Create logic is tested with mocked API
- **WHEN** the unit test invokes the Create handler with a mocked `AttachUserPolicy` that returns a valid response containing one policy with a `PolicyId`
- **THEN** the test SHALL assert the composite ID is set to `user_id#policy_id` and `policy_set` state is populated

#### Scenario: Read logic is tested with mocked API returning multiple policies
- **WHEN** the unit test invokes the Read handler with a mocked `DescribeUserInfo` that returns multiple policies for the user
- **THEN** the test SHALL assert only the policy matching the `policy_id` from the composite ID is kept in `policy_set` state

#### Scenario: Delete logic is tested with mocked API
- **WHEN** the unit test invokes the Delete handler with a mocked `DetachUserPolicy`
- **THEN** the test SHALL assert `DetachUserPolicy` is called with `PolicyIds` containing the single `policy_id` derived from the composite ID
