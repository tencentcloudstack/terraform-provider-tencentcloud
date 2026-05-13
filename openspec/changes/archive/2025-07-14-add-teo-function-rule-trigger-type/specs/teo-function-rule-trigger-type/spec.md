## ADDED Requirements

### Requirement: trigger_type parameter in tencentcloud_teo_function_rule resource
The `tencentcloud_teo_function_rule` resource SHALL include an optional `trigger_type` parameter of type string. Valid values SHALL be `direct`, `weight`, and `region`. When not specified, the API defaults to `direct`. The parameter SHALL NOT have ForceNew set, allowing it to be updated in-place.

#### Scenario: Create function rule with trigger_type set to direct
- **WHEN** a user creates a `tencentcloud_teo_function_rule` resource with `trigger_type = "direct"`
- **THEN** the Create handler SHALL set `request.TriggerType` to `"direct"` in the `CreateFunctionRule` API call

#### Scenario: Create function rule without trigger_type
- **WHEN** a user creates a `tencentcloud_teo_function_rule` resource without specifying `trigger_type`
- **THEN** the Create handler SHALL NOT set `request.TriggerType` in the API call, and the API SHALL default to `"direct"`

#### Scenario: Read function rule with trigger_type
- **WHEN** the Read handler calls `DescribeFunctionRules` and the response `FunctionRule.TriggerType` is not nil
- **THEN** the handler SHALL set `trigger_type` in the Terraform state to the value of `respData.TriggerType`

#### Scenario: Read function rule with nil TriggerType
- **WHEN** the Read handler calls `DescribeFunctionRules` and the response `FunctionRule.TriggerType` is nil
- **THEN** the handler SHALL NOT set `trigger_type` in the Terraform state

#### Scenario: Update function rule trigger_type
- **WHEN** a user updates the `trigger_type` parameter of an existing `tencentcloud_teo_function_rule` resource
- **THEN** the Update handler SHALL detect the change via `d.HasChange("trigger_type")`, include `trigger_type` in the mutable arguments, and set `request.TriggerType` in the `ModifyFunctionRule` API call

#### Scenario: trigger_type validation
- **WHEN** a user provides an invalid value for `trigger_type` (e.g., "invalid")
- **THEN** Terraform SHALL produce a validation error at plan time
