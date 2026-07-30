## ADDED Requirements

### Requirement: Resource calls STS AssumeRole API on create

The `tencentcloud_sts_assume_role_operation` resource SHALL call the STS `AssumeRole` API exactly once during the Create lifecycle. The resource SHALL pass all user-provided input parameters to the API request. The resource SHALL use `tccommon.WriteRetryTimeout` with retry logic via `resource.Retry`.

#### Scenario: Successful assume role with required parameters only
- **WHEN** user provides `role_arn` and `role_session_name`
- **THEN** the resource calls AssumeRole API with those parameters and stores the returned credentials, expired_time, and expiration in state

#### Scenario: Successful assume role with all parameters
- **WHEN** user provides `role_arn`, `role_session_name`, `duration_seconds`, `policy`, `external_id`, `tags`, `source_identity`, `serial_number`, and `token_code`
- **THEN** the resource calls AssumeRole API with all parameters mapped correctly and stores the response in state

#### Scenario: API call fails
- **WHEN** the AssumeRole API returns an error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and return it to Terraform

### Requirement: All input parameters are ForceNew

All input schema fields (`role_arn`, `role_session_name`, `duration_seconds`, `policy`, `external_id`, `tags`, `source_identity`, `serial_number`, `token_code`) SHALL be marked as `ForceNew: true`. This ensures any parameter change triggers resource recreation (a new API call).

#### Scenario: User changes role_arn
- **WHEN** user modifies the `role_arn` value in configuration
- **THEN** Terraform plans to destroy and recreate the resource (new AssumeRole call)

### Requirement: Read and Delete are no-ops

The Read and Delete lifecycle functions SHALL return nil immediately without making any API calls. This is the standard RESOURCE_KIND_OPERATION pattern.

#### Scenario: Terraform refresh
- **WHEN** Terraform calls Read on the resource
- **THEN** the function returns nil without any API interaction

#### Scenario: Terraform destroy
- **WHEN** Terraform calls Delete on the resource
- **THEN** the function returns nil without any API interaction

### Requirement: Response credentials are stored as Computed Sensitive fields

The resource SHALL expose the API response as Computed fields:
- `credentials`: A list (MaxItems: 1) containing `token`, `tmp_secret_id`, `tmp_secret_key` - all marked Sensitive
- `expired_time`: Integer Unix timestamp of credential expiry
- `expiration`: String ISO8601 UTC time of credential expiry

#### Scenario: Credentials stored in state after successful create
- **WHEN** AssumeRole API returns successfully with Credentials, ExpiredTime, and Expiration
- **THEN** the resource sets `credentials` list with token/tmp_secret_id/tmp_secret_key, sets `expired_time` and `expiration` in state

#### Scenario: Response field is nil
- **WHEN** AssumeRole API returns with a nil Credentials field
- **THEN** the resource SHALL NOT call d.Set for that field (skip nil fields)

### Requirement: Resource ID uses random token

The resource SHALL use `helper.BuildToken()` to generate its ID after a successful Create. The ID has no semantic meaning since this is a one-time operation.

#### Scenario: ID assignment after create
- **WHEN** the AssumeRole API call succeeds
- **THEN** the resource ID is set to a random token via `helper.BuildToken()`

### Requirement: Tags parameter maps to STS session tags

The `tags` schema field SHALL be a TypeList of objects, each containing `key` (Required, string) and `value` (Required, string). These map to the STS API's `Tags` field (`[]*Tag` with `Key` and `Value`).

#### Scenario: User provides session tags
- **WHEN** user specifies tags with key-value pairs
- **THEN** the resource constructs `[]*sts.Tag` objects and sets them on the request

### Requirement: Resource is registered in provider

The resource SHALL be registered in `tencentcloud/provider.go` under the STS service section and listed in `tencentcloud/provider.md`.

#### Scenario: Provider includes the resource
- **WHEN** Terraform initializes the tencentcloud provider
- **THEN** `tencentcloud_sts_assume_role_operation` is available as a resource type
