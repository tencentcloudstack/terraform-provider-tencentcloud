## ADDED Requirements

### Requirement: Resource shall manage DLC user-policy binding lifecycle
The Terraform provider SHALL provide a resource `tencentcloud_dlc_attach_user_policyr_attachment` of kind `RESOURCE_KIND_ATTACHMENT` that manages the binding between a DLC user and a set of authorization policies. The resource SHALL support Create (bind), Read (query), and Delete (unbind). The resource SHALL NOT support in-place Update; any change to a top-level argument SHALL trigger recreation.

#### Scenario: Bind policies to a user on create
- **WHEN** a user applies a configuration defining `tencentcloud_dlc_attach_user_policyr_attachment` with `user_id`, `policy_set`, and `account_type`
- **THEN** the provider SHALL call the DLC `AttachUserPolicy` API with the given `UserId`, `PolicySet`, and `AccountType`
- **AND** the provider SHALL set the resource ID to the composite of `user_id` and `account_type` joined by `tccommon.FILED_SP`
- **AND** the provider SHALL set `policy_set` in state from the API response

#### Scenario: Read reflects currently bound policies
- **WHEN** the provider reads an existing `tencentcloud_dlc_attach_user_policyr_attachment` resource
- **THEN** the provider SHALL call the DLC `DescribeUserInfo` API with the `UserId` and `AccountType` derived from the composite ID
- **AND** the provider SHALL verify the bound policies are still present
- **AND** if the binding is no longer present, the provider SHALL emit a `[CRUD]` log preserving the id and then set `d.SetId("")`

#### Scenario: Unbind policies on delete
- **WHEN** a user destroys a `tencentcloud_dlc_attach_user_policyr_attachment` resource
- **THEN** the provider SHALL call the DLC `DetachUserPolicy` API with the `UserId`, `AccountType`, and `PolicySet` derived from state
- **AND** the resource SHALL be removed from state

#### Scenario: Any change recreates the resource
- **WHEN** a user changes any top-level argument of an existing `tencentcloud_dlc_attach_user_policyr_attachment` resource
- **THEN** the provider SHALL reject the in-place update (the update method returns an error for immutable args) and Terraform SHALL recreate the resource

### Requirement: Resource SHALL use a composite identifier
Because the DLC `AttachUserPolicy` API returns no single resource identifier, the resource SHALL use a composite ID composed of `user_id` and `account_type` joined by `tccommon.FILED_SP`. The Read, Update, and Delete operations SHALL split the composite ID to recover `user_id` and `account_type` for use as API request parameters.

#### Scenario: Composite ID is split for API calls
- **WHEN** the provider performs Read, Update, or Delete on `tencentcloud_dlc_attach_user_policyr_attachment`
- **THEN** the provider SHALL split `d.Id()` by `tccommon.FILED_SP` into `user_id` and `account_type`
- **AND** the provider SHALL use these values as the `UserId` and `AccountType` request parameters

### Requirement: Resource schema SHALL expose binding parameters
The resource schema SHALL expose `user_id` (Required, ForceNew), `policy_set` (Required), and `account_type` (Optional, ForceNew). The `policy_set` SHALL be a list of policy objects whose input-eligible fields (`database`, `catalog`, `table`, `operation`, `policy_type`, `function`, `view`, `column`, `data_engine`, `re_auth`, `engine_generation`, `model`, `policy_id`) can be set by the user.

#### Scenario: Required arguments are enforced
- **WHEN** a user omits `user_id` or `policy_set` in the configuration
- **THEN** Terraform SHALL report a validation error before any API call is made

#### Scenario: account_type is optional
- **WHEN** a user does not specify `account_type`
- **THEN** the provider SHALL not set `AccountType` in the API request, allowing the cloud API to apply its default

### Requirement: Cloud API calls SHALL use retry and proper error handling
Create, Read, and Delete SHALL wrap their cloud API calls in `resource.Retry` using `tccommon.WriteRetryTimeout` (for Create/Delete) and `tccommon.ReadRetryTimeout` (for Read). Errors from the cloud API SHALL be wrapped with `tccommon.RetryError`. State mutations (setting the ID and fields) SHALL occur outside the retry block, after successful retry completion.

#### Scenario: Transient API failure is retried
- **WHEN** a DLC API call fails with a retryable error during Create, Read, or Delete
- **THEN** the provider SHALL retry the call within the configured timeout and wrap the error with `tccommon.RetryError`

#### Scenario: Create validates non-empty response
- **WHEN** the `AttachUserPolicy` API returns a nil response
- **THEN** the provider SHALL return a `NonRetryableError` rather than writing an empty id to state

### Requirement: Resource SHALL be registered in the provider
The provider SHALL register `tencentcloud_dlc_attach_user_policyr_attachment` in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md`. A documentation file `resource_tc_dlc_attach_user_policyr_attachment.md` SHALL exist under `tencentcloud/services/dlc/`.

#### Scenario: Resource is usable in a Terraform configuration
- **WHEN** a user references `tencentcloud_dlc_attach_user_policyr_attachment` in a `.tf` file
- **THEN** the provider SHALL recognize the resource type and execute its CRUD handlers

### Requirement: Unit tests SHALL use gomonkey mocks
The resource test file `resource_tc_dlc_attach_user_policyr_attachment_test.go` SHALL use gomonkey to mock the DLC cloud API client methods and test the business logic of Create/Read/Delete without relying on the Terraform acceptance test suite. The tests SHALL be runnable with `go test -gcflags=all=-l`.

#### Scenario: Create logic is tested with mocked API
- **WHEN** the unit test invokes the Create handler with a mocked `AttachUserPolicy` that returns a valid response
- **THEN** the test SHALL assert the composite ID is set and `policy_set` state is populated

#### Scenario: Delete logic is tested with mocked API
- **WHEN** the unit test invokes the Delete handler with a mocked `DetachUserPolicy`
- **THEN** the test SHALL assert `DetachUserPolicy` is called with the expected parameters
