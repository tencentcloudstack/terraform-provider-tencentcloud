# config-update-config-deliver Specification

## Purpose
TBD - created by archiving change add-config-update-config-deliver. Update Purpose after archive.
## Requirements
### Requirement: Resource schema for Config delivery settings
The `tencentcloud_config_update_config_deliver` resource SHALL expose the following arguments for managing Tencent Cloud Config delivery (投递设置):
- `status` (required, int): Delivery switch. `0` = disabled, `1` = enabled.
- `deliver_name` (optional, string): Delivery service name.
- `target_arn` (optional, string): Resource six-segment ARN. COS: `qcs::cos:$region:$account:prefix/$appid/$BucketName`; CLS: `qcs::cls:$region:$account:cls/topicId`.
- `deliver_prefix` (optional, string): Log prefix for stored delivery content.
- `deliver_type` (optional, string): Delivery type. Valid values: `COS`, `CLS`.
- `deliver_content_type` (optional, int): Content type. `1` = configuration change, `2` = resource list, `3` = all.
- `create_time` (computed, string): Creation time of the delivery configuration.

#### Scenario: Required field enforced
- **WHEN** a user applies a configuration without setting `status`
- **THEN** Terraform SHALL reject the plan with a missing-required-argument error

#### Scenario: Computed field populated on read
- **WHEN** the resource is created and `DescribeConfigDeliver` returns a non-nil `CreateTime`
- **THEN** the `create_time` field SHALL be populated in state from the API response

### Requirement: Create via Update API
Because the Config delivery setting is a global singleton with no dedicated Create API, the Create handler SHALL call `UpdateConfigDeliver` with all provided fields, then set the resource ID to a generated token via `helper.BuildToken()`. After a successful Update call the Create handler SHALL invoke Read to refresh state.

#### Scenario: First apply creates the config
- **WHEN** the resource is applied for the first time with `status = 1`
- **THEN** the system SHALL call `UpdateConfigDeliver` with the provided arguments, set `d.SetId(helper.BuildToken())`, and read the result back into state

### Requirement: Read from DescribeConfigDeliver
The Read handler SHALL call `DescribeConfigDeliver` (no request parameters) and populate each non-nil field of the response into state. If the response is nil the handler SHALL log a warning with the resource id and clear the id from state.

#### Scenario: Normal read
- **WHEN** Read is invoked and `DescribeConfigDeliver` returns a populated response
- **THEN** each present field (`status`, `deliver_name`, `target_arn`, `deliver_prefix`, `deliver_type`, `deliver_content_type`, `create_time`) SHALL be set into state

#### Scenario: Empty response clears state
- **WHEN** `DescribeConfigDeliver` returns `nil`
- **THEN** the handler SHALL log `[WARN] ... resource tencentcloud_config_update_config_deliver [<id>] not found` and call `d.SetId("")`

### Requirement: Update via UpdateConfigDeliver
The Update handler SHALL call `UpdateConfigDeliver` with the current state values for all fields, wrapped in a `resource.Retry(tccommon.WriteRetryTimeout, ...)` block using `tccommon.RetryError` for error wrapping. After a successful Update the handler SHALL invoke Read to refresh state.

#### Scenario: Field change triggers update
- **WHEN** the user changes `deliver_content_type` from `1` to `3` and applies
- **THEN** the system SHALL call `UpdateConfigDeliver` with the new values and refresh state from `DescribeConfigDeliver`

#### Scenario: API failure retries
- **WHEN** `UpdateConfigDeliver` returns a retriable error
- **THEN** the retry block SHALL re-attempt the call up to `tccommon.WriteRetryTimeout` before surfacing the error

### Requirement: Delete is a no-op
The Delete handler SHALL be a no-op (return nil) because the Config delivery setting has no dedicated Delete API and the configuration persists on the cloud side.

#### Scenario: Destroy does not delete remote config
- **WHEN** the user runs `terraform destroy` on the resource
- **THEN** the Delete handler SHALL return nil without calling any cloud API, leaving the remote delivery setting unchanged

### Requirement: Provider registration and documentation
The resource SHALL be registered in `tencentcloud/provider.go` ResourcesMap as `tencentcloud_config_update_config_deliver` and listed in the provider.go file comment index under the Config product. A documentation file `resource_tc_config_update_config_deliver.md` SHALL provide a one-line description (mentioning Tencent Cloud Config) and an Example Usage HCL block.

#### Scenario: Resource available to provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_config_update_config_deliver` SHALL be present in `ResourcesMap` pointing to `config.ResourceTencentCloudConfigUpdateConfigDeliver()`

