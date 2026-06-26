## ADDED Requirements

### Requirement: Resource supports advanced_origin_routing_parameters in branches.actions
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `advanced_origin_routing_parameters` parameter within `branches.actions`. This parameter defines the advanced origin routing optimization configuration and maps to the SDK `RuleEngineAction.AdvancedOriginRoutingParameters` field.

#### Scenario: Create resource with AdvancedOriginRouting action
- **WHEN** user defines `advanced_origin_routing_parameters` with `direction` in an `actions` block within `branches` of a `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the Create function SHALL map it to `RuleEngineAction.AdvancedOriginRoutingParameters` in the `CreateL7AccRules` API request

#### Scenario: Read resource returns advanced_origin_routing_parameters
- **WHEN** the `DescribeL7AccRules` API returns rules with `Branches[0].Actions[0].AdvancedOriginRoutingParameters` populated
- **THEN** the Read function SHALL set the `advanced_origin_routing_parameters` attribute in the Terraform state

#### Scenario: Update resource advanced_origin_routing_parameters
- **WHEN** user modifies the `advanced_origin_routing_parameters` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineAction.AdvancedOriginRoutingParameters` in the `ModifyL7AccRule` API request

### Requirement: Resource supports shield_parameters in branches.actions
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `shield_parameters` parameter within `branches.actions`. This parameter defines the origin offload (Shield) configuration and maps to the SDK `RuleEngineAction.ShieldParameters` field.

#### Scenario: Create resource with Shield action
- **WHEN** user defines `shield_parameters` with `shield_space_id` in an `actions` block within `branches`
- **THEN** the Create function SHALL map it to `RuleEngineAction.ShieldParameters` in the API request

#### Scenario: Read resource returns shield_parameters
- **WHEN** the API returns rules with `Branches[0].Actions[0].ShieldParameters` populated
- **THEN** the Read function SHALL set the `shield_parameters` attribute in the Terraform state

#### Scenario: Update resource shield_parameters
- **WHEN** user modifies the `shield_parameters` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineAction.ShieldParameters` in the API request

### Requirement: Resource supports site_failover_parameters in branches.actions
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `site_failover_parameters` parameter within `branches.actions`. This parameter defines the origin site failover configuration and maps to the SDK `RuleEngineAction.SiteFailoverParameters` field, which includes `site_failover_status_codes` and a list of `site_failover_params` (each being a `SiteFailover` structure).

#### Scenario: Create resource with SiteFailover action
- **WHEN** user defines `site_failover_parameters` with `site_failover_status_codes` and `site_failover_params` in an `actions` block within `branches`
- **THEN** the Create function SHALL map it to `RuleEngineAction.SiteFailoverParameters` in the API request

#### Scenario: Read resource returns site_failover_parameters
- **WHEN** the API returns rules with `Branches[0].Actions[0].SiteFailoverParameters` populated
- **THEN** the Read function SHALL set the `site_failover_parameters` attribute in the Terraform state

#### Scenario: Update resource site_failover_parameters
- **WHEN** user modifies the `site_failover_parameters` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineAction.SiteFailoverParameters` in the API request

### Requirement: Actions name description includes new action names
The `name` field description in `branches.actions` SHALL include `AdvancedOriginRouting`, `Shield`, and `SiteFailover` as valid operation names.

#### Scenario: Name field description is updated
- **WHEN** user reads the description of the `name` field in `branches.actions`
- **THEN** the description SHALL list `AdvancedOriginRouting`, `Shield`, and `SiteFailover` as valid operation names

### Requirement: All new parameters are optional
- **WHEN** user does not specify `advanced_origin_routing_parameters`, `shield_parameters`, or `site_failover_parameters` in the resource configuration
- **THEN** the resource SHALL behave exactly as before (no changes to existing behavior)
