# Capability: MQTT Message Enrichment Rule Management

## ADDED Requirements

### Requirement: Message Enrichment Rule Resource

Terraform Provider SHALL support managing Tencent Cloud MQTT message enrichment rules through the `tencentcloud_mqtt_message_enrichment_rule` resource, enabling users to create, read, update, delete, and import message attribute enrichment rules.

#### Scenario: Create message enrichment rule with required fields

- **WHEN** user provides valid `instance_id`, `rule_name`, `condition`, and `actions`
- **THEN** the system SHALL create a message enrichment rule with default priority 1
- **AND** return the rule ID assigned by the cloud service
- **AND** set resource ID in format `{instance_id}#{rule_id}`

#### Scenario: Create message enrichment rule with optional fields

- **WHEN** user provides `status` and `remark` in addition to required fields
- **THEN** the system SHALL create the rule with specified status and remark
- **AND** handle empty or missing optional fields gracefully

#### Scenario: Read message enrichment rule

- **WHEN** Terraform reads an existing rule
- **THEN** the system SHALL query the rule using DescribeMessageEnrichmentRules API
- **AND** populate all resource attributes including computed fields
- **AND** decode Base64-encoded `condition` and `actions` fields if needed

#### Scenario: Update message enrichment rule

- **WHEN** user modifies rule attributes (except `instance_id`)
- **THEN** the system SHALL update the rule using ModifyMessageEnrichmentRule API
- **AND** submit all current field values (full update semantics)
- **AND** preserve the existing `priority` value (computed field)
- **AND** refresh resource state after successful update

#### Scenario: Delete message enrichment rule

- **WHEN** user destroys the resource
- **THEN** the system SHALL delete the rule using DeleteMessageEnrichmentRule API
- **AND** use retry mechanism to handle eventual consistency
- **AND** return success when rule no longer exists

#### Scenario: Import existing message enrichment rule

- **WHEN** user runs `terraform import` with format `{instance_id}#{rule_id}`
- **THEN** the system SHALL parse the composite ID correctly
- **AND** fetch the rule state from cloud service
- **AND** populate Terraform state with current values

#### Scenario: Handle rule not found

- **WHEN** reading a rule that no longer exists in cloud
- **THEN** the system SHALL log a warning message
- **AND** set resource ID to empty string to remove from state
- **AND** not return an error to Terraform

### Requirement: Resource Schema Definition

The resource Schema SHALL define all fields according to the Tencent Cloud MQTT API specifications and Terraform best practices.

#### Scenario: Required fields validation

- **WHEN** user omits any required field (`instance_id`, `rule_name`, `condition`, `actions`)
- **THEN** Terraform SHALL reject the configuration at plan time
- **AND** display a clear error message indicating the missing field

#### Scenario: Priority field as computed-only

- **WHEN** resource is created or read
- **THEN** `priority` field SHALL be exposed as a computed (read-only) attribute
- **AND** internally use default value 1 during creation
- **AND** not allow user modification (omitted from Schema's Required/Optional)

#### Scenario: Instance ID forces recreation

- **WHEN** user changes `instance_id` value
- **THEN** Terraform SHALL mark the resource for ForceNew (recreate)
- **AND** destroy the old rule and create a new one

#### Scenario: Base64 encoding requirement

- **WHEN** user provides `condition` or `actions` fields
- **THEN** the system SHALL validate they are valid Base64-encoded JSON strings
- **AND** reject invalid encoding or malformed JSON at apply time

### Requirement: Service Layer Integration

The service layer (`service_tencentcloud_mqtt.go`) SHALL provide methods to interact with Tencent Cloud MQTT API for message enrichment rules.

#### Scenario: Query rule by ID

- **WHEN** service layer receives `instance_id` and `rule_id`
- **THEN** it SHALL call DescribeMessageEnrichmentRules API
- **AND** filter the results to find matching rule by ID
- **AND** return nil if rule not found (not an error)

#### Scenario: Delete rule by ID

- **WHEN** service layer receives delete request
- **THEN** it SHALL call DeleteMessageEnrichmentRule API with `instance_id` and `rule_id`
- **AND** use retry logic to handle transient failures
- **AND** return error only on permanent failures

#### Scenario: Handle API errors gracefully

- **WHEN** any API call fails
- **THEN** service layer SHALL log the request and response bodies
- **AND** return descriptive error to resource layer
- **AND** include API error code and message in error

### Requirement: Testing Coverage

The resource implementation SHALL include comprehensive acceptance tests covering all CRUD operations and edge cases.

#### Scenario: Basic CRUD acceptance test

- **WHEN** test creates, reads, updates, and deletes a rule
- **THEN** all operations SHALL succeed without errors
- **AND** resource state SHALL match expected values at each step
- **AND** rule SHALL not exist after deletion

#### Scenario: Import functionality test

- **WHEN** test imports an existing rule
- **THEN** imported state SHALL match the actual cloud resource
- **AND** subsequent plan SHALL show no changes

#### Scenario: Status field modification test

- **WHEN** test updates rule status from inactive to active
- **THEN** update SHALL succeed
- **AND** DescribeMessageEnrichmentRules SHALL reflect new status

#### Scenario: Field validation test

- **WHEN** test provides invalid Base64 encoding
- **THEN** apply operation SHALL fail with clear error
- **AND** no resource SHALL be created in cloud

### Requirement: Documentation

Resource documentation SHALL provide clear guidance on usage, parameters, and examples following Terraform Registry standards.

#### Scenario: Complete parameter documentation

- **WHEN** user views resource documentation
- **THEN** each parameter SHALL have description, type, and required/optional indicator
- **AND** constraints (e.g., length limits) SHALL be documented
- **AND** special encoding requirements SHALL be highlighted

#### Scenario: Usage example with Base64 encoding

- **WHEN** documentation includes usage example
- **THEN** example SHALL demonstrate proper Base64 encoding for `condition` and `actions`
- **AND** show JSON structure before encoding
- **AND** explain the purpose of each field

#### Scenario: Import example

- **WHEN** documentation covers import functionality
- **THEN** example SHALL show exact command syntax
- **AND** explain how to construct the composite ID
- **AND** mention what happens to computed fields after import

### Requirement: Error Handling and Logging

The implementation SHALL handle errors gracefully and provide detailed logging for troubleshooting.

#### Scenario: API failure during create

- **WHEN** CreateMessageEnrichmentRule API returns error
- **THEN** resource create operation SHALL fail
- **AND** error message SHALL include API error details
- **AND** no partial state SHALL be written to Terraform

#### Scenario: Eventual consistency handling

- **WHEN** rule is just created but not yet queryable
- **THEN** read operation SHALL use retry with backoff
- **AND** eventually succeed when rule becomes available
- **AND** log retry attempts for debugging

#### Scenario: Detailed operation logging

- **WHEN** any CRUD operation executes
- **THEN** system SHALL log operation start with resource ID
- **AND** log API request and response bodies at DEBUG level
- **AND** log operation completion time using LogElapsed
- **AND** log critical errors at CRITICAL level

### Requirement: Backward Compatibility

The new resource SHALL not break existing Terraform configurations or other MQTT resources.

#### Scenario: No impact on existing resources

- **WHEN** new resource is added to provider
- **THEN** existing MQTT resources SHALL continue to function unchanged
- **AND** provider version bump SHALL be minor (not major)
- **AND** no existing resource IDs SHALL conflict

#### Scenario: Shared service layer compatibility

- **WHEN** new methods added to `service_tencentcloud_mqtt.go`
- **THEN** existing service methods SHALL remain unchanged
- **AND** client initialization SHALL support new API calls
- **AND** rate limiting SHALL apply consistently
