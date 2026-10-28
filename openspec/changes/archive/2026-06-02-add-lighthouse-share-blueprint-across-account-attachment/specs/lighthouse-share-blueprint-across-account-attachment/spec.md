## ADDED Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_lighthouse_share_blueprint_across_account_attachment` SHALL define the following schema fields:

- `blueprint_id` (Required, ForceNew, TypeString): The Lighthouse blueprint ID to share. This field is immutable.
- `account_ids` (Required, TypeList of TypeString): List of target TencentCloud account IDs to share the blueprint with. This field is updatable.

`blueprint_id` SHALL be `ForceNew: true` since changing the blueprint requires a new resource.
`account_ids` SHALL NOT be `ForceNew` to allow incremental updates to the shared account list.

#### Scenario: Schema validation
- **WHEN** a user omits `blueprint_id` or `account_ids` in their Terraform configuration
- **THEN** Terraform SHALL report a validation error requiring both fields

#### Scenario: Schema ForceNew behavior for blueprint_id
- **WHEN** a user modifies `blueprint_id` on an existing resource
- **THEN** Terraform SHALL destroy the existing resource and create a new one

#### Scenario: Update behavior for account_ids
- **WHEN** a user modifies `account_ids` on an existing resource
- **THEN** Terraform SHALL call the Update function to incrementally add/remove accounts without destroying the resource

### Requirement: Create Operation (Share Blueprint)
The Create function SHALL call the `ShareBlueprintAcrossAccounts` API with `BlueprintId` and `AccountIds` from the resource configuration.

On successful API call, the resource ID SHALL be set as `blueprint_id`.

Before calling the Read function to refresh state, the Create function SHALL sort `account_ids` in ascending order and store them to state. This ensures consistent ordering regardless of the input order specified by the user.

The Create function SHALL use `resource.Retry` with `tccommon.WriteRetryTimeout` and `tccommon.RetryError` for error retry handling.

#### Scenario: Successful share creation
- **WHEN** a user creates a `tencentcloud_lighthouse_share_blueprint_across_account_attachment` resource with a valid `blueprint_id` and `account_ids`
- **THEN** the provider SHALL call `ShareBlueprintAcrossAccounts` API, set the resource ID to `blueprint_id`, and store sorted `account_ids` to state

#### Scenario: API error during creation
- **WHEN** the `ShareBlueprintAcrossAccounts` API returns an error
- **THEN** the provider SHALL retry and ultimately return the error to the user

#### Scenario: Empty API response during creation
- **WHEN** the `ShareBlueprintAcrossAccounts` API returns a nil or empty response
- **THEN** the provider SHALL return a `NonRetryableError`

#### Scenario: Account IDs sorting on create
- **WHEN** a user provides `account_ids` in any order (e.g., `["100087654321", "100012345678"]`)
- **THEN** the provider SHALL sort them in ascending order before storing to state (e.g., `["100012345678", "100087654321"]`)

### Requirement: Read Operation (Query Share Status)
The Read function SHALL call the `DescribeBlueprintsShareAcrossAccountInfos` API with `BlueprintIds` containing the `blueprint_id` from the resource ID.

The Read function SHALL use `d.Id()` directly as the `blueprintId` (no parsing needed).

If the API returns sharing information, the Read function SHALL sort the returned `account_ids` in ascending order and update the resource state. This ensures that the ordering in state matches what was stored during Create.

If no sharing information is found, the Read function SHALL clear the resource ID (`d.SetId("")`) to indicate the resource no longer exists.

#### Scenario: Successful read
- **WHEN** the Read function is called for an existing shared blueprint
- **THEN** the provider SHALL query `DescribeBlueprintsShareAcrossAccountInfos` and update the resource state with sorted current account IDs

#### Scenario: Resource deleted externally
- **WHEN** the Read function is called but the sharing relationship no longer exists
- **THEN** the provider SHALL clear the resource ID (`d.SetId("")`)

#### Scenario: No plan diff when account order differs
- **WHEN** a user has state with sorted `account_ids` and configures them in a different order
- **THEN** `terraform plan` SHALL show no changes (the order difference is normalized)

### Requirement: Update Operation (Incremental Modify Shared Accounts)
The Update function SHALL implement incremental updates when `account_ids` changes:

1. Calculate the difference between old and new account lists:
   - Accounts in old but not in new → to remove
   - Accounts in new but not in old → to add
2. If there are accounts to remove, call `CancelShareBlueprintAcrossAccounts` with those accounts
3. If there are accounts to add, call `ShareBlueprintAcrossAccounts` with those accounts
4. Call Read to refresh state with actual values

Both API calls SHALL use `resource.Retry` with `tccommon.WriteRetryTimeout` and `tccommon.RetryError` for error retry handling.

#### Scenario: Add accounts only
- **WHEN** a user adds new account(s) to `account_ids`
- **THEN** the provider SHALL call `ShareBlueprintAcrossAccounts` only for the newly added accounts

#### Scenario: Remove accounts only
- **WHEN** a user removes account(s) from `account_ids`
- **THEN** the provider SHALL call `CancelShareBlueprintAcrossAccounts` only for the removed accounts

#### Scenario: Both add and remove accounts
- **WHEN** a user both adds and removes accounts from `account_ids`
- **THEN** the provider SHALL first call `CancelShareBlueprintAcrossAccounts` for removed accounts, then call `ShareBlueprintAcrossAccounts` for added accounts

### Requirement: Delete Operation (Cancel Share)
The Delete function SHALL call the `CancelShareBlueprintAcrossAccounts` API with `BlueprintId` from the resource ID and `AccountIds` from the resource state (`d.Get("account_ids")`).

The Delete function SHALL use `resource.Retry` with `tccommon.WriteRetryTimeout` and `tccommon.RetryError` for error retry handling.

#### Scenario: Successful share cancellation
- **WHEN** a user destroys a `tencentcloud_lighthouse_share_blueprint_across_account_attachment` resource
- **THEN** the provider SHALL call `CancelShareBlueprintAcrossAccounts` API with all current account IDs and remove the resource from state

#### Scenario: API error during deletion
- **WHEN** the `CancelShareBlueprintAcrossAccounts` API returns an error
- **THEN** the provider SHALL retry and ultimately return the error to the user

### Requirement: Resource Import Support
The resource SHALL support Terraform import via `schema.ImportStatePassthrough`.

The import ID SHALL be simply the `blueprint_id` (e.g., `lhbp-xxx`).

On import, the Read function SHALL query the API and populate the resource state with the actual shared account IDs.

#### Scenario: Import existing sharing relationship
- **WHEN** a user imports a resource using `terraform import tencentcloud_lighthouse_share_blueprint_across_account_attachment.example lhbp-xxx`
- **THEN** the provider SHALL read the sharing information using the blueprint ID and populate the resource state with all shared account IDs

### Requirement: Provider Registration
The resource SHALL be registered in `tencentcloud/provider.go` under the Lighthouse resources section with the key `tencentcloud_lighthouse_share_blueprint_across_account_attachment`.

A corresponding entry SHALL be added to `tencentcloud/provider.md`.

#### Scenario: Resource accessible after registration
- **WHEN** the provider is built with the new resource registered
- **THEN** users SHALL be able to use `tencentcloud_lighthouse_share_blueprint_across_account_attachment` in their Terraform configurations

### Requirement: Documentation
The resource SHALL have a documentation file at `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment.md`.

The documentation SHALL include:
- A one-sentence description with the cloud product name (TEO → Lighthouse)
- Example Usage section with valid HCL examples
- Import section explaining the simple blueprint_id format

The final website documentation SHALL be generated via `make doc`.

#### Scenario: Documentation generation
- **WHEN** `make doc` is executed
- **THEN** `website/docs/r/lighthouse_share_blueprint_across_account_attachment.html.markdown` SHALL be generated with complete documentation

### Requirement: Unit Testing
The resource SHALL have unit tests at `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment_test.go`.

Tests SHALL use gomonkey for mocking cloud API calls and cover Create, Read, Update, Delete, and Import scenarios.

Tests SHALL be runnable with `go test -gcflags=all=-l` and pass successfully.

#### Scenario: Unit tests pass
- **WHEN** `go test -gcflags=all=-l` is run on the test file
- **THEN** all unit tests SHALL pass without errors
