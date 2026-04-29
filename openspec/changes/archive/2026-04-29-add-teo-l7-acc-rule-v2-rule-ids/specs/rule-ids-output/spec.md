## MODIFIED Requirements

### Requirement: Remove rule_ids computed attribute
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL NOT have a `rule_ids` computed attribute. The `rule_id` attribute provides the single rule ID managed by this resource.

#### Scenario: rule_ids is not in schema
- **WHEN** a user inspects the `tencentcloud_teo_l7_acc_rule_v2` resource schema
- **THEN** there SHALL NOT be a `rule_ids` attribute

## ADDED Requirements

### Requirement: Vary action support in branches.actions
The `tencentcloud_teo_l7_acc_rule_v2` resource's `branches.actions` SHALL support a `Vary` action with `vary_parameters` block (TypeList, Optional, MaxItems:1) containing a `switch` field (TypeString, Required, values: "on"/"off"), mapping to `RuleEngineAction.VaryParameters` in the cloud API.

#### Scenario: Vary action is set via Create/Update
- **WHEN** a user configures an action with `name = "Vary"` and `vary_parameters { switch = "on" }`
- **THEN** the `CreateL7AccRules`/`ModifyL7AccRule` API SHALL be called with `RuleEngineAction.Name = "Vary"` and `RuleEngineAction.VaryParameters.Switch = "on"`

#### Scenario: Vary action is read from API
- **WHEN** the Read function processes an action with `VaryParameters` set
- **THEN** the `vary_parameters.switch` attribute SHALL be populated from `actions.VaryParameters.Switch`

### Requirement: OriginAuthentication action support in branches.actions
The `tencentcloud_teo_l7_acc_rule_v2` resource's `branches.actions` SHALL support an `OriginAuthentication` action with `origin_authentication_parameters` block (TypeList, Optional, MaxItems:1) containing a `request_properties` list (TypeList, Required) where each item has `type` (TypeString, Required), `name` (TypeString, Required), and `value` (TypeString, Required), mapping to `RuleEngineAction.OriginAuthenticationParameters` in the cloud API.

#### Scenario: OriginAuthentication action is set via Create/Update
- **WHEN** a user configures an action with `name = "OriginAuthentication"` and `origin_authentication_parameters { request_properties { type = "Header" name = "Authorization" value = "Bearer token" } }`
- **THEN** the API SHALL be called with `RuleEngineAction.Name = "OriginAuthentication"` and `RuleEngineAction.OriginAuthenticationParameters.RequestProperties` containing the corresponding entry

#### Scenario: OriginAuthentication action is read from API
- **WHEN** the Read function processes an action with `OriginAuthenticationParameters` set
- **THEN** the `origin_authentication_parameters.request_properties` attribute SHALL be populated from `actions.OriginAuthenticationParameters.RequestProperties`

### Requirement: Correct action name descriptions
The `actions.name` field description SHALL list all valid action names matching the SDK's `RuleEngineAction.Name` definition, with correct naming (e.g., `SetContentIdentifier` not `SetContentIdentifierParameters`).

#### Scenario: name description matches SDK
- **WHEN** a user reads the `actions.name` field description
- **THEN** it SHALL include all supported actions from `RuleEngineAction.Name` in the SDK, including `SetContentIdentifier`, `Vary`, `ContentCompression`, and `OriginAuthentication`
