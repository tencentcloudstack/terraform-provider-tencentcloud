## ADDED Requirements

### Requirement: Resource manages TEO L7 acceleration rule lifecycle

The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL manage the full CRUD lifecycle of a TEO L7 acceleration rule using the following APIs:
- Create: `CreateL7AccRules` with `ZoneId` and `Rules` parameters
- Read: `DescribeL7AccRules` with `ZoneId` and `Filters` (rule-id filter using `Values`)
- Update: `ModifyL7AccRule` with `ZoneId` and `Rule` (containing `RuleId`, `Status`, `RuleName`, `Description`, `Branches`)
- Delete: `DeleteL7AccRules` with `ZoneId` and `RuleIds`

The resource SHALL use a composite ID format `{zone_id}#{rule_id}` with `tccommon.FILED_SP` as separator.

#### Scenario: Create a new L7 acceleration rule
- **WHEN** user applies a Terraform configuration with `tencentcloud_teo_l7_acc_rule_v2` resource specifying `zone_id`, `status`, `rule_name`, `description`, and `branches`
- **THEN** the resource SHALL call `CreateL7AccRules` API with the zone ID and a single `RuleEngineItem` containing the specified fields, and store the returned `rule_id` in state with composite ID format

#### Scenario: Read an existing L7 acceleration rule
- **WHEN** Terraform performs a refresh on the resource
- **THEN** the resource SHALL call `DescribeL7AccRules` API with `zone_id` and a filter for `rule-id` with the stored `rule_id` as the filter value, and update state with the returned `status`, `rule_name`, `description`, `branches`, and `rule_priority`

#### Scenario: Update an existing L7 acceleration rule
- **WHEN** user modifies `status`, `rule_name`, `description`, or `branches` in the Terraform configuration
- **THEN** the resource SHALL call `ModifyL7AccRule` API with `zone_id` and a `RuleEngineItem` containing the `rule_id` and updated fields

#### Scenario: Delete an L7 acceleration rule
- **WHEN** user destroys the resource
- **THEN** the resource SHALL call `DeleteL7AccRules` API with `zone_id` and `rule_id`

#### Scenario: Import an existing rule
- **WHEN** user imports a resource using `terraform import` with ID format `{zone_id}#{rule_id}`
- **THEN** the resource SHALL parse the composite ID and read the rule details from the API

### Requirement: Resource schema defines correct field types and constraints

The resource schema SHALL define the following fields:
- `zone_id`: TypeString, Required, ForceNew - Zone ID for the TEO site
- `status`: TypeString, Optional - Rule status (enable/disable)
- `rule_name`: TypeString, Optional - Rule name (max 255 characters)
- `description`: TypeList of TypeString, Optional - Rule annotations
- `branches`: TypeList of Resource, Optional - Sub-rule branches with conditions and actions
- `rule_id`: TypeString, Computed - Rule ID returned by the API
- `rule_priority`: TypeInt, Computed - Rule priority (output only)

#### Scenario: Zone ID is immutable after creation
- **WHEN** user attempts to change `zone_id` after resource creation
- **THEN** Terraform SHALL plan to destroy and recreate the resource (ForceNew behavior)

#### Scenario: Optional fields can be omitted
- **WHEN** user creates a resource without specifying `status`, `rule_name`, `description`, or `branches`
- **THEN** the resource SHALL be created successfully with API defaults for those fields

### Requirement: Error handling follows provider patterns

The resource SHALL implement proper error handling:
- All API calls MUST be wrapped with retry logic using `tccommon.RetryError()`
- Create MUST verify the response is not nil and `RuleIds` is not empty before storing the ID
- Read MUST check for nil response and empty rules list, logging the resource ID before clearing state

#### Scenario: Create API returns empty response
- **WHEN** the `CreateL7AccRules` API returns a nil response or empty `RuleIds`
- **THEN** the resource SHALL return a `NonRetryableError` with a descriptive message

#### Scenario: Read API returns no matching rule
- **WHEN** the `DescribeL7AccRules` API returns an empty rules list
- **THEN** the resource SHALL log the resource ID for debugging and set the ID to empty string to remove from state
