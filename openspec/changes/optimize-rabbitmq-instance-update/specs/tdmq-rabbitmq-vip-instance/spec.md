## ADDED Requirements

### Requirement: Auto-Renew Flag Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL allow users to update the `auto_renew_flag` field through Terraform apply operations for prepaid instances.

#### Scenario: User enables auto-renewal on existing prepaid instance
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 1` (prepaid) and `auto_renew_flag = false`
- **WHEN** the user changes `auto_renew_flag` to `true` in the configuration
- **THEN** `terraform plan` shows the auto-renewal change
- **AND** `terraform apply` successfully updates the auto-renewal setting on the cloud resource
- **AND** the Terraform state reflects the updated `auto_renew_flag` value

#### Scenario: User disables auto-renewal on existing prepaid instance
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 1` (prepaid) and `auto_renew_flag = true`
- **WHEN** the user changes `auto_renew_flag` to `false` in the configuration
- **THEN** `terraform plan` shows the auto-renewal change
- **AND** `terraform apply` successfully updates the auto-renewal setting on the cloud resource
- **AND** the Terraform state reflects the updated `auto_renew_flag` value

#### Scenario: User attempts to update auto-renewal on postpaid instance
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 0` (postpaid)
- **WHEN** the user attempts to change `auto_renew_flag` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message indicating auto-renew flag can only be updated for prepaid instances
- **AND** the instance remains unchanged

#### Scenario: API integration for auto-renewal update
- **GIVEN** a user modifies `auto_renew_flag` in the configuration
- **WHEN** the Update operation detects the change via `d.HasChange("auto_renew_flag")`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `AutoRenewFlag` parameter with the new value
- **AND** the API call succeeds and the auto-renewal setting is updated on the cloud resource

### Requirement: Time Span Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL allow users to update the `time_span` field through Terraform apply operations for prepaid instances.

#### Scenario: User increases time span on existing prepaid instance
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 1` (prepaid) and `time_span = 1`
- **WHEN** the user changes `time_span` to `3` in the configuration
- **THEN** `terraform plan` shows the time span change
- **AND** `terraform apply` successfully updates the purchase duration on the cloud resource
- **AND** the Terraform state reflects the updated `time_span` value

#### Scenario: User decreases time span on existing prepaid instance
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 1` (prepaid) and `time_span = 6`
- **WHEN** the user changes `time_span` to `3` in the configuration
- **THEN** `terraform plan` shows the time span change
- **AND** `terraform apply` successfully updates the purchase duration on the cloud resource
- **AND** the Terraform state reflects the updated `time_span` value

#### Scenario: User attempts to update time span on postpaid instance
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 0` (postpaid)
- **WHEN** the user attempts to change `time_span` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message indicating time span can only be updated for prepaid instances
- **AND** the instance remains unchanged

#### Scenario: API integration for time span update
- **GIVEN** a user modifies `time_span` in the configuration
- **WHEN** the Update operation detects the change via `d.HasChange("time_span")`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `TimeSpan` parameter with the new value
- **AND** the API call succeeds and the purchase duration is updated on the cloud resource

### Requirement: Post-Update State Verification
The resource SHALL verify state synchronization after update operations to ensure changes were successfully applied.

#### Scenario: Successful update with state verification
- **GIVEN** a user performs an update operation that modifies `auto_renew_flag` or `time_span`
- **WHEN** the Update operation receives a successful response from `ModifyRabbitMQVipInstance`
- **THEN** the provider calls `resourceTencentCloudTdmqRabbitmqVipInstanceRead()` to refresh the state
- **AND** the read operation waits for the state to stabilize using retry logic with `ReadRetryTimeout*10`
- **AND** the Terraform state reflects the updated values from the cloud resource
- **AND** `terraform apply` completes successfully without showing pending changes

#### Scenario: State verification with eventual consistency
- **GIVEN** a user performs an update operation that modifies `auto_renew_flag`
- **WHEN** the Update operation receives a successful response but the read operation returns stale data
- **THEN** the provider retries the read operation up to `ReadRetryTimeout*10` duration
- **AND** once the read operation returns the updated value, the state is synchronized
- **AND** no error is reported to the user
- **AND** the Terraform state shows the correct updated value

#### Scenario: State verification timeout handling
- **GIVEN** a user performs an update operation
- **WHEN** the read operation continues to return stale data beyond `ReadRetryTimeout*10`
- **THEN** the provider logs a warning message indicating state verification timeout
- **AND** the provider proceeds to complete the Terraform operation
- **AND** a subsequent `terraform plan` may show pending changes
- **AND** users can retry `terraform apply` to complete the state synchronization

### Requirement: Retry Mechanism for Update Operations
The resource SHALL implement retry logic for update operations to handle transient API failures.

#### Scenario: Update succeeds after retry due to network issue
- **GIVEN** a user performs an update operation that modifies `auto_renew_flag`
- **WHEN** the `ModifyRabbitMQVipInstance` API call fails with a transient network error
- **THEN** the provider uses `resource.Retry()` with `WriteRetryTimeout` to retry the operation
- **AND** the API call succeeds on a retry attempt
- **AND** the update completes successfully
- **AND** the Terraform state reflects the updated value

#### Scenario: Update fails after all retries due to persistent error
- **GIVEN** a user performs an update operation with an invalid field value
- **WHEN** the `ModifyRabbitMQVipInstance` API call fails with a validation error
- **THEN** the provider retries the operation up to the `WriteRetryTimeout` limit
- **AND** all retry attempts fail with the same error
- **THEN** the provider returns a detailed error message to the user
- **AND** the Terraform state remains unchanged (previous field values)
- **AND** users can fix the configuration and retry `terraform apply`

#### Scenario: Update operation logging
- **GIVEN** a user performs an update operation
- **WHEN** the Update operation executes
- **THEN** each retry attempt is logged with `[DEBUG]` level messages
- **AND** the final success or failure is logged with appropriate level
- **AND** API request and response bodies are logged for troubleshooting

### Requirement: Differential Updates Implementation
The resource SHALL implement differential updates that only send changed fields to the API.

#### Scenario: Update only changed fields
- **GIVEN** an existing instance with `auto_renew_flag = false` and `cluster_name = "my-instance"`
- **WHEN** the user changes only `auto_renew_flag` to `true`
- **THEN** the Update operation detects only `auto_renew_flag` has changed
- **AND** the request to `ModifyRabbitMQVipInstance` includes only `AutoRenewFlag` parameter
- **AND** `ClusterName` is not included in the request
- **AND** the API call succeeds

#### Scenario: Update multiple changed fields
- **GIVEN** an existing instance with `auto_renew_flag = false` and `resource_tags` set
- **WHEN** the user changes both `auto_renew_flag` to `true` and adds new `resource_tags`
- **THEN** the Update operation detects both fields have changed
- **AND** the request to `ModifyRabbitMQVipInstance` includes both `AutoRenewFlag` and `Tags` parameters
- **AND** the API call succeeds
- **AND** both fields are updated on the cloud resource

#### Scenario: No update when no fields changed
- **GIVEN** an existing instance with no configuration changes
- **WHEN** the user runs `terraform apply`
- **THEN** `terraform plan` shows no changes
- **AND** no API call is made to `ModifyRabbitMQVipInstance`
- **AND** the Terraform state remains unchanged

### Requirement: Enhanced Error Handling
The resource SHALL provide detailed error messages for update operation failures.

#### Scenario: Invalid update on postpaid instance
- **GIVEN** an existing instance with `pay_mode = 0` (postpaid)
- **WHEN** the user attempts to update `auto_renew_flag`
- **THEN** `terraform apply` fails with error message: "auto_renew_flag can only be updated for prepaid instances (pay_mode = 1)"
- **AND** the error message suggests checking the pay_mode setting
- **AND** the Terraform state remains unchanged

#### Scenario: API validation error
- **GIVEN** a user attempts to update `time_span` to an invalid value (e.g., negative number)
- **WHEN** the `ModifyRabbitMQVipInstance` API returns a validation error
- **THEN** the provider returns the API error message to the user
- **AND** the error includes context about which field caused the validation failure
- **AND** the Terraform state remains unchanged

#### Scenario: Network timeout during update
- **GIVEN** a user performs an update operation
- **WHEN** the API call times out after all retry attempts
- **THEN** the provider returns an error message indicating network timeout
- **AND** the error message suggests checking network connectivity and retrying
- **AND** the Terraform state remains unchanged

## MODIFIED Requirements

### Requirement: Tag Update Support
The resource SHALL allow users to update tags through Terraform apply operations.

#### Scenario: User adds new tags to existing instance
- **GIVEN** an existing RabbitMQ VIP instance without tags in Terraform configuration
- **WHEN** the user adds resource_tags blocks (e.g., `tag_key = "cost-center", tag_value = "123"`) to the configuration
- **THEN** `terraform plan` shows the tag addition
- **AND** `terraform apply` successfully adds the tags to the instance using retry logic
- **AND** `terraform apply` verifies state synchronization by calling Read operation after update
- **AND** the Terraform state reflects the updated tags as a list

#### Scenario: User modifies existing tags
- **GIVEN** an existing instance with resource_tags blocks for `env = "dev"`
- **WHEN** the user changes the tag value to `env = "prod"`
- **THEN** `terraform plan` shows the tag value change
- **AND** `terraform apply` replaces all tags with the new configuration using retry logic
- **AND** `terraform apply` verifies state synchronization by calling Read operation after update
- **AND** the Terraform state reflects the updated tags

#### Scenario: User removes all tags
- **GIVEN** an existing instance with multiple resource_tags blocks
- **WHEN** the user removes all resource_tags blocks from the configuration
- **THEN** `terraform plan` shows tag removal
- **AND** `terraform apply` removes all tags from the instance using RemoveAllTags flag and retry logic
- **AND** `terraform apply` verifies state synchronization by calling Read operation after update
- **AND** the Terraform state shows an empty tags list

#### Scenario: User removes resource_tags field
- **GIVEN** an existing instance with resource_tags blocks
- **WHEN** the user removes the `resource_tags` field from the configuration
- **THEN** `terraform plan` shows no changes to tags
- **AND** tags on the cloud resource remain unchanged

### Requirement: Update API Integration
The resource SHALL correctly map tag changes between Terraform and Tencent Cloud APIs with retry logic.

#### Scenario: Update API integration with retry
- **GIVEN** a user modifies tags in an existing instance configuration
- **WHEN** the Update operation detects `resource_tags` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `Tags` array with all current tags
- **AND** the API call is wrapped in retry logic with `WriteRetryTimeout`
- **AND** the API call succeeds (after potential retries) and tags are updated on the cloud resource
- **AND** the Read operation is called to verify state synchronization

### Requirement: Error Handling for Tag Update
The resource SHALL handle tag-related update errors gracefully with retry logic.

#### Scenario: API error during tag update
- **GIVEN** a user updates tags on an existing instance
- **WHEN** the `ModifyRabbitMQVipInstance` API fails with a transient error
- **THEN** the provider retries the operation up to `WriteRetryTimeout`
- **AND** if all retry attempts fail, the error is logged and returned to the user with context
- **AND** Terraform state remains unchanged (previous tag values)
- **AND** subsequent `terraform apply` attempts the update again
