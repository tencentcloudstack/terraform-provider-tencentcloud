## ADDED Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_lighthouse_share_blueprint_across_account_attachment` SHALL define the following schema fields:

- `blueprint_id` (Required, ForceNew, TypeString): The Lighthouse blueprint ID to share.
- `account_ids` (Required, ForceNew, TypeList of TypeString): List of target TencentCloud account IDs to share the blueprint with.

All fields SHALL be `ForceNew: true` since the binding relationship is immutable (no update API available).

#### Scenario: Schema validation
- **WHEN** a user omits `blueprint_id` or `account_ids` in their Terraform configuration
- **THEN** Terraform SHALL report a validation error requiring both fields

#### Scenario: Schema ForceNew behavior
- **WHEN** a user modifies `blueprint_id` or `account_ids` on an existing resource
- **THEN** Terraform SHALL destroy the existing resource and create a new one

### Requirement: Create Operation (Share Blueprint)
The Create function SHALL call the `ShareBlueprintAcrossAccounts` API with `BlueprintId` and `AccountIds` from the resource configuration.

On successful API call, the resource ID SHALL be set as the composite of `blueprint_id` and sorted `account_ids` joined by `tccommon.FILED_SP`.

The Create function SHALL use `resource.Retry` with `tccommon.WriteRetryTimeout` and `tccommon.RetryError` for error retry handling.

#### Scenario: Successful share creation
- **WHEN** a user creates a `tencentcloud_lighthouse_share_blueprint_across_account_attachment` resource with a valid `blueprint_id` and `account_ids`
- **THEN** the provider SHALL call `ShareBlueprintAcrossAccounts` API and set the resource ID to `blueprint_id#<account_ids>`

#### Scenario: API error during creation
- **WHEN** the `ShareBlueprintAcrossAccounts` API returns an error
- **THEN** the provider SHALL retry and ultimately return the error to the user

#### Scenario: Empty API response during creation
- **WHEN** the `ShareBlueprintAcrossAccounts` API returns a nil or empty response
- **THEN** the provider SHALL return a `NonRetryableError`

### Requirement: Read Operation (Query Share Status)
The Read function SHALL call the `DescribeBlueprintsShareAcrossAccountInfos` API with `BlueprintIds` containing the `blueprint_id` from the composite resource ID.

The Read function SHALL parse the composite ID from `d.Id()`, extract `blueprint_id` and `account_ids`.

If the API returns sharing information that matches, the Read function SHALL update the resource state with the current values.

If no matching sharing information is found, the Read function SHALL clear the resource ID (`d.SetId("")`) to indicate the resource no longer exists.

#### Scenario: Successful read
- **WHEN** the Read function is called for an existing shared blueprint
- **THEN** the provider SHALL query `DescribeBlueprintsShareAcrossAccountInfos` and update the resource state

#### Scenario: Resource deleted externally
- **WHEN** the Read function is called but the sharing relationship no longer exists
- **THEN** the provider SHALL clear the resource ID (`d.SetId("")`)

#### Scenario: Invalid composite ID
- **WHEN** the Read function is called with a malformed composite ID
- **THEN** the provider SHALL return an error indicating the ID format is invalid

### Requirement: Delete Operation (Cancel Share)
The Delete function SHALL call the `CancelShareBlueprintAcrossAccounts` API with `BlueprintId` and `AccountIds` extracted from the composite resource ID.

The Delete function SHALL use `resource.Retry` with `tccommon.WriteRetryTimeout` and `tccommon.RetryError` for error retry handling.

#### Scenario: Successful share cancellation
- **WHEN** a user destroys a `tencentcloud_lighthouse_share_blueprint_across_account_attachment` resource
- **THEN** the provider SHALL call `CancelShareBlueprintAcrossAccounts` API and remove the resource from state

#### Scenario: API error during deletion
- **WHEN** the `CancelShareBlueprintAcrossAccounts` API returns an error
- **THEN** the provider SHALL retry and ultimately return the error to the user

### Requirement: Resource Import Support
The resource SHALL support Terraform import via `schema.ImportStatePassthrough`.

The import ID SHALL be the composite format: `blueprint_id#<account_ids>`, where `<account_ids>` is the sorted, FILED_SP-joined list of account IDs.

#### Scenario: Import existing sharing relationship
- **WHEN** a user imports a resource using `terraform import tencentcloud_lighthouse_share_blueprint_across_account_attachment.example lhbp-xxx#100012345678`
- **THEN** the provider SHALL read the sharing information and populate the resource state

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
- Import section explaining the composite ID format

The final website documentation SHALL be generated via `make doc`.

#### Scenario: Documentation generation
- **WHEN** `make doc` is executed
- **THEN** `website/docs/r/lighthouse_share_blueprint_across_account_attachment.html.markdown` SHALL be generated with complete documentation

### Requirement: Unit Testing
The resource SHALL have unit tests at `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_attachment_test.go`.

Tests SHALL use gomonkey for mocking cloud API calls and cover Create, Read, Delete, and Import scenarios.

Tests SHALL be runnable with `go test -gcflags=all=-l` and pass successfully.

#### Scenario: Unit tests pass
- **WHEN** `go test -gcflags=all=-l` is run on the test file
- **THEN** all unit tests SHALL pass without errors