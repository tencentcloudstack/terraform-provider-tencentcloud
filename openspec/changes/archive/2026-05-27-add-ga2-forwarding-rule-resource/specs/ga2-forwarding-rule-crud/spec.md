## ADDED Requirements

### Requirement: Create forwarding rule
The system SHALL create a GA2 Layer-7 forwarding rule by calling the CreateForwardingRule API with the specified parameters (global_accelerator_id, listener_id, forwarding_policy_id, rule_conditions, rule_actions, origin_headers, enable_origin_sni, origin_sni, origin_host). After creation, the system SHALL poll DescribeTaskResult until the task succeeds, then store the composite ID (global_accelerator_id#listener_id#forwarding_policy_id#forwarding_rule_id) in state.

#### Scenario: Successful creation with all parameters
- **WHEN** user provides global_accelerator_id, listener_id, forwarding_policy_id, rule_conditions, rule_actions, origin_headers, enable_origin_sni, origin_sni, and origin_host
- **THEN** the system calls CreateForwardingRule API, waits for the async task to complete, and stores the composite ID in Terraform state

#### Scenario: Successful creation with minimal parameters
- **WHEN** user provides only global_accelerator_id, listener_id, forwarding_policy_id, rule_conditions, and rule_actions (required fields)
- **THEN** the system calls CreateForwardingRule API with only the provided fields, waits for the async task to complete, and stores the composite ID in Terraform state

#### Scenario: Create API returns nil response
- **WHEN** the CreateForwardingRule API returns a nil response or nil ForwardingRuleId
- **THEN** the system returns a NonRetryableError indicating the creation failed

### Requirement: Read forwarding rule
The system SHALL read a GA2 Layer-7 forwarding rule by calling the DescribeForwardingRule API with pagination (Limit=100) and filtering the results by forwarding_rule_id to find the specific rule. The system SHALL set all attributes from the ForwardingRuleSet response into Terraform state.

#### Scenario: Successful read of existing rule
- **WHEN** the resource exists and DescribeForwardingRule returns a ForwardingRuleSet containing the target forwarding_rule_id
- **THEN** the system sets rule_conditions, rule_actions, origin_headers, enable_origin_sni, origin_sni, and origin_host from the response into state

#### Scenario: Rule not found (deleted externally)
- **WHEN** the DescribeForwardingRule API does not return a ForwardingRuleSet entry matching the forwarding_rule_id
- **THEN** the system removes the resource from state (sets ID to empty string)

### Requirement: Update forwarding rule
The system SHALL update a GA2 Layer-7 forwarding rule by calling the ModifyForwardingRule API with the changed parameters. The system SHALL poll DescribeTaskResult until the task succeeds, then call Read to refresh state.

#### Scenario: Successful update of rule conditions and actions
- **WHEN** user modifies rule_conditions, rule_actions, origin_headers, enable_origin_sni, origin_sni, or origin_host
- **THEN** the system calls ModifyForwardingRule API with all current values, waits for the async task to complete, and refreshes state via Read

#### Scenario: Update with HasChange detection
- **WHEN** user applies a configuration change
- **THEN** the system only calls ModifyForwardingRule if at least one mutable field has changed

### Requirement: Delete forwarding rule
The system SHALL delete a GA2 Layer-7 forwarding rule by calling the DeleteForwardingRule API with the four identifying IDs. The system SHALL poll DescribeTaskResult until the task succeeds.

#### Scenario: Successful deletion
- **WHEN** the resource exists and DeleteForwardingRule is called
- **THEN** the system calls DeleteForwardingRule API, waits for the async task to complete, and removes the resource from state

### Requirement: Import forwarding rule
The system SHALL support importing an existing GA2 forwarding rule using the composite ID format: `global_accelerator_id#listener_id#forwarding_policy_id#forwarding_rule_id`.

#### Scenario: Successful import with composite ID
- **WHEN** user runs terraform import with the composite ID format
- **THEN** the system parses the four components from the ID and calls Read to populate state

#### Scenario: Invalid import ID format
- **WHEN** user provides an import ID that does not contain exactly 4 parts separated by #
- **THEN** the system returns an error indicating the expected format

### Requirement: Resource schema definition
The system SHALL define the Terraform schema with the following fields:
- `global_accelerator_id` (Required, ForceNew, String): Global accelerator instance ID
- `listener_id` (Required, ForceNew, String): Listener ID
- `forwarding_policy_id` (Required, ForceNew, String): Forwarding policy ID
- `rule_conditions` (Required, List of objects with rule_condition_type and rule_condition_value): Rule conditions
- `rule_actions` (Required, List of objects with rule_action_type and rule_action_value): Rule actions
- `origin_headers` (Optional, List of objects with key and value): Origin headers
- `enable_origin_sni` (Optional, Bool): Whether to enable origin SNI
- `origin_sni` (Optional, String): Origin SNI value
- `origin_host` (Optional, String): Origin host value
- `forwarding_rule_id` (Computed, String): The forwarding rule ID returned by the API
- `task_id` (Computed, String): The async task ID from the last mutation operation

#### Scenario: Schema validation
- **WHEN** the resource schema is registered in the provider
- **THEN** all Required fields MUST be provided by the user, ForceNew fields trigger recreation on change, and Computed fields are set by the system

### Requirement: Retry and timeout handling
The system SHALL wrap all API calls with retry logic using tccommon.ReadRetryTimeout for read operations and tccommon.WriteRetryTimeout for write operations. The resource SHALL declare Timeouts for Create, Update, and Delete operations (default 20 minutes each).

#### Scenario: Transient API error during create
- **WHEN** the CreateForwardingRule API returns a transient error
- **THEN** the system retries the call within the WriteRetryTimeout period

#### Scenario: Task polling timeout
- **WHEN** the async task does not reach SUCCESS status within the configured timeout
- **THEN** the system returns a timeout error to the user
