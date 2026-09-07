## MODIFIED Requirements

### Requirement: Resource schema SHALL expose binding parameters and limit policy_set to a single policy
The resource schema SHALL expose `user_id` (Required, ForceNew), `policy_set` (Required, ForceNew, MinItems 1, MaxItems 1), and `account_type` (Optional, ForceNew). The `policy_set` SHALL be a list containing exactly one policy object whose input-eligible fields (`database`, `catalog`, `table`, `operation`, `policy_type`, `function`, `view`, `column`, `data_engine`, `re_auth`, `engine_generation`, `model`, `policy_id`) can be set by the user. The `re_auth` field SHALL be `Optional` and `Computed` (TypeBool): the user MAY set it on create to control whether the grantee can further delegate the permissions, and it SHALL be populated from the cloud API response on read.

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

#### Scenario: re_auth is an optional, computed boolean input
- **WHEN** a user sets `policy_set.0.re_auth = true` in the configuration
- **THEN** Terraform SHALL accept the value as a valid boolean input (no validation error)
- **AND** the provider SHALL forward `ReAuth = true` to the `AttachUserPolicy` API request `PolicySet[0].ReAuth`

#### Scenario: re_auth may be omitted
- **WHEN** a user does not specify `policy_set.0.re_auth` in the configuration
- **THEN** Terraform SHALL accept the configuration without error (the field is optional)
- **AND** the cloud API default (`false`) SHALL apply

#### Scenario: re_auth is populated from the API response on read
- **WHEN** the provider reads an existing `tencentcloud_dlc_attach_user_policy_attachment` resource
- **THEN** the provider SHALL set `policy_set.0.re_auth` in Terraform state from the `ReAuth` field of the matched policy in the `DescribeUserInfo` response (when present)

### Requirement: Cloud API calls SHALL use retry and proper error handling
Create, Read, and Delete SHALL wrap their cloud API calls in `resource.Retry` using `tccommon.WriteRetryTimeout` (for Create/Delete) and `tccommon.ReadRetryTimeout` (for Read). Errors from the cloud API SHALL be wrapped with `tccommon.RetryError`. State mutations (setting the ID and fields) SHALL occur outside the retry block, after successful retry completion. The Create handler SHALL map every user-settable `policy_set` sub-field — including `re_auth` — from Terraform state into the corresponding `dlc.Policy` field of the `AttachUserPolicy` request `PolicySet`.

#### Scenario: Transient API failure is retried
- **WHEN** a DLC API call fails with a retryable error during Create, Read, or Delete
- **THEN** the provider SHALL retry the call within the configured timeout and wrap the error with `tccommon.RetryError`

#### Scenario: Create validates the response contains exactly one policy with a non-empty PolicyId
- **WHEN** the `AttachUserPolicy` API returns a nil response, or a `PolicySet` whose length is not 1, or a policy with an empty `PolicyId`
- **THEN** the provider SHALL return a `NonRetryableError` rather than writing an empty or invalid id to state

#### Scenario: Create forwards the configured re_auth to the API
- **WHEN** the Create handler builds the `AttachUserPolicy` request from a configuration that sets `policy_set.0.re_auth`
- **THEN** the provider SHALL set `PolicySet[0].ReAuth` on the request `dlc.Policy` from the configured `re_auth` value before calling the DLC API
