## ADDED Requirements

### Requirement: Resource schema for teo inference api token

The system SHALL provide a Terraform resource `tencentcloud_teo_inference_api_token` with the following schema fields:
- `zone_id` (string, Required, ForceNew): 站点 ID
- `name` (string, Required, ForceNew): 推理 API Token 名称
- `token_id` (string, Computed): 推理 API Token ID
- `content` (string, Computed, Sensitive): 推理 API Token 内容
- `create_time` (string, Computed): 创建时间

The resource SHALL use a composite ID in the format `zone_id#token_id` (using `tccommon.FILED_SP` as separator).

The resource SHALL support `terraform import` using the composite ID.

#### Scenario: Define resource with required fields
- **WHEN** a user defines a `tencentcloud_teo_inference_api_token` resource with `zone_id` and `name`
- **THEN** the resource is accepted as a valid configuration

#### Scenario: Resource supports import
- **WHEN** a user runs `terraform import tencentcloud_teo_inference_api_token.example zone-id#token-id`
- **THEN** the resource state is populated with `zone_id`, `token_id`, and all computed fields from the cloud API

### Requirement: Create teo inference api token

The system SHALL create a TEO inference API token by calling the `CreateInferenceAPIToken` cloud API with `ZoneId` and `Name` parameters.

Upon successful creation, the system SHALL check that the response is not nil and that `TokenId` is not empty. If the response or `TokenId` is empty, the system SHALL return a `NonRetryableError`.

After successful creation, the system SHALL set the resource ID to `zone_id#token_id` and read the resource to populate state.

#### Scenario: Successful token creation
- **WHEN** a user creates a `tencentcloud_teo_inference_api_token` with valid `zone_id` and `name`
- **THEN** the system calls `CreateInferenceAPIToken` and sets the resource ID to `zone_id#token_id`

#### Scenario: Creation returns empty token id
- **WHEN** `CreateInferenceAPIToken` returns a response with nil `TokenId`
- **THEN** the system returns a `NonRetryableError` and does not set an empty resource ID

#### Scenario: Creation API error
- **WHEN** `CreateInferenceAPIToken` returns an error
- **THEN** the system wraps the error with `tccommon.RetryError` and the creation fails after retries are exhausted

### Requirement: Read teo inference api token

The system SHALL read a TEO inference API token by calling `DescribeInferenceAPITokens` with `ZoneId` and `Limit` set to the API maximum (100), then matching by `token_id` in the returned `Tokens` array.

Inside the retry block, if the API returns an empty `Tokens` list or no matching `token_id`, the system SHALL return a `NonRetryableError`.

Outside the retry block, on retry failure, the system SHALL first log `[CRUD] teo inference_api_token id=<id>` and then call `d.SetId("")`.

When setting fields, the system SHALL check each response field for nil before calling `d.Set()`.

#### Scenario: Successful token read
- **WHEN** the system reads a token that exists in the cloud
- **THEN** `token_id`, `name`, `content`, and `create_time` are populated from the matched `Tokens` element

#### Scenario: Token not found
- **WHEN** `DescribeInferenceAPITokens` returns no matching token by `token_id`
- **THEN** the system logs `[CRUD] teo inference_api_token id=<id>` and sets `d.SetId("")` to remove the resource from state

#### Scenario: Read API error
- **WHEN** `DescribeInferenceAPITokens` returns an error
- **THEN** the system wraps the error with `tccommon.RetryError` and retries until timeout

### Requirement: Delete teo inference api token

The system SHALL delete a TEO inference API token by calling `DeleteInferenceAPIToken` with `ZoneId` and `TokenId` extracted from the composite resource ID.

#### Scenario: Successful token deletion
- **WHEN** a user deletes a `tencentcloud_teo_inference_api_token` resource
- **THEN** the system calls `DeleteInferenceAPIToken` with `ZoneId` and `TokenId` from the composite ID

#### Scenario: Delete API error
- **WHEN** `DeleteInferenceAPIToken` returns an error
- **THEN** the system wraps the error with `tccommon.RetryError` and the deletion fails after retries are exhausted

### Requirement: Update immutability for teo inference api token

The system SHALL NOT support updating a TEO inference API token because the cloud API provides no Update interface.

The Update function SHALL declare `immutableArgs` containing all top-level fields except `Id()`. If any field in `immutableArgs` has changed, the system SHALL return an error. If no immutable field has changed, the system SHALL call Read to refresh state.

#### Scenario: Attempt to update name
- **WHEN** a user changes `name` on an existing `tencentcloud_teo_inference_api_token`
- **THEN** the Update function detects the change in `immutableArgs` and returns an error

#### Scenario: No-op update
- **WHEN** a user applies a plan with no changes to immutable fields
- **THEN** the Update function calls Read to refresh state without error

### Requirement: Retry and timeout handling

All cloud API calls (Create/Read/Delete) SHALL use `resource.Retry` with the appropriate timeout:
- Create and Delete SHALL use `tccommon.WriteRetryTimeout`
- Read SHALL use `tccommon.ReadRetryTimeout`

Inside the retry block, the system SHALL only call the API and return results; setting the ID or other successful operations SHALL happen outside the retry block after error handling.

The system SHALL use `tccommon.RetryError()` to wrap errors returned by the cloud API.

#### Scenario: Transient API error triggers retry
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the system retries within the configured timeout period

#### Scenario: ID is set outside retry block
- **WHEN** `CreateInferenceAPIToken` succeeds
- **THEN** `d.SetId()` is called after the retry block completes, not inside it