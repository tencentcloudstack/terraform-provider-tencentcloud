## ADDED Requirements

### Requirement: Resource schema definition
The resource `tencentcloud_teo_security_js_injection_rule` SHALL define the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID
- `js_injection_rules` (Required, TypeList): JavaScript 注入规则列表，包含嵌套结构:
  - `rule_id` (Computed, TypeString): 规则 ID，由 API 返回
  - `name` (Required, TypeString): 规则名称
  - `priority` (Optional, Computed, TypeInt): 规则优先级，范围 0-100，默认 0
  - `condition` (Required, TypeString): 匹配条件，需符合表达式语法
  - `inject_js` (Optional, Computed, TypeString): JS 注入选项，取值 `no-injection` 或 `inject-sdk-only`
- `js_injection_rule_ids` (Computed, TypeList): 规则 ID 列表，由 Create API 返回

#### Scenario: Valid resource schema with required fields only
- **WHEN** user defines a `tencentcloud_teo_security_js_injection_rule` resource with `zone_id` and `js_injection_rules` containing `name` and `condition`
- **THEN** the resource SHALL accept the configuration and proceed with creation

#### Scenario: Missing required field zone_id
- **WHEN** user defines the resource without `zone_id`
- **THEN** Terraform SHALL report a validation error indicating `zone_id` is required

#### Scenario: Missing required nested field name or condition
- **WHEN** user defines `js_injection_rules` without `name` or `condition`
- **THEN** Terraform SHALL report a validation error indicating the required nested field is missing

### Requirement: Resource Create operation
The resource SHALL call `CreateSecurityJSInjectionRule` API to create JS injection rules when `terraform apply` is executed for a new resource.

#### Scenario: Successful creation
- **WHEN** creating a new resource with valid `zone_id` and `js_injection_rules`
- **THEN** the system SHALL call `CreateSecurityJSInjectionRule` with `ZoneId` and `JSInjectionRules` parameters
- **AND** set the resource ID to `zone_id`
- **AND** store `js_injection_rule_ids` from the API response
- **AND** call Read to refresh the full resource state

#### Scenario: Creation API failure with retryable error
- **WHEN** the `CreateSecurityJSInjectionRule` API returns a retryable error
- **THEN** the system SHALL retry the operation using `tccommon.WriteRetryTimeout`
- **AND** wrap the error using `tccommon.RetryError()`

### Requirement: Resource Read operation
The resource SHALL call `DescribeSecurityJSInjectionRule` API to read the current state of JS injection rules.

#### Scenario: Successful read with existing rules
- **WHEN** reading an existing resource
- **THEN** the system SHALL call `DescribeSecurityJSInjectionRule` with `ZoneId` parameter
- **AND** set `Limit` to 100 (API maximum) and handle pagination
- **AND** populate `js_injection_rules` from the API response including `rule_id`, `name`, `priority`, `condition`, and `inject_js`

#### Scenario: Resource not found
- **WHEN** the `DescribeSecurityJSInjectionRule` API returns no rules for the given `zone_id`
- **THEN** the system SHALL set `d.SetId("")` to indicate the resource no longer exists

#### Scenario: Read API failure with retryable error
- **WHEN** the `DescribeSecurityJSInjectionRule` API returns a retryable error
- **THEN** the system SHALL retry the operation using `tccommon.ReadRetryTimeout`
- **AND** wrap the error using `tccommon.RetryError()`

### Requirement: Resource Update operation
The resource SHALL call `ModifySecurityJSInjectionRule` API when `js_injection_rules` changes.

#### Scenario: Successful update
- **WHEN** `js_injection_rules` is changed and `terraform apply` is executed
- **THEN** the system SHALL call `ModifySecurityJSInjectionRule` with `ZoneId` and the full `JSInjectionRules` list
- **AND** call Read to refresh the resource state after update

#### Scenario: No changes detected
- **WHEN** no mutable fields have changed
- **THEN** the system SHALL NOT call the Modify API

#### Scenario: Update API failure with retryable error
- **WHEN** the `ModifySecurityJSInjectionRule` API returns a retryable error
- **THEN** the system SHALL retry the operation using `tccommon.WriteRetryTimeout`
- **AND** wrap the error using `tccommon.RetryError()`

### Requirement: Resource Delete operation
The resource SHALL call `DeleteSecurityJSInjectionRule` API to delete all JS injection rules when the resource is destroyed.

#### Scenario: Successful deletion
- **WHEN** destroying the resource
- **THEN** the system SHALL read `js_injection_rule_ids` from the state
- **AND** call `DeleteSecurityJSInjectionRule` with `ZoneId` and `JSInjectionRuleIds` parameters

#### Scenario: Delete API failure with retryable error
- **WHEN** the `DeleteSecurityJSInjectionRule` API returns a retryable error
- **THEN** the system SHALL retry the operation using `tccommon.WriteRetryTimeout`
- **AND** wrap the error using `tccommon.RetryError()`

### Requirement: Resource Import
The resource SHALL support Terraform import using `zone_id` as the import ID.

#### Scenario: Import existing resource
- **WHEN** user runs `terraform import tencentcloud_teo_security_js_injection_rule.example <zone_id>`
- **THEN** the system SHALL import the resource by calling Read with the provided `zone_id`

### Requirement: Provider registration
The resource SHALL be registered in `tencentcloud/provider.go` ResourcesMap as `"tencentcloud_teo_security_js_injection_rule"` and listed in `tencentcloud/provider.md` under TEO category.

#### Scenario: Resource available in provider
- **WHEN** the provider is loaded
- **THEN** `tencentcloud_teo_security_js_injection_rule` SHALL be available as a resource type

### Requirement: Unit tests with gomonkey mock
The resource SHALL have unit tests using gomonkey mock to verify business logic without calling real cloud APIs.

#### Scenario: Create function unit test
- **WHEN** running unit tests for the Create function
- **THEN** the test SHALL mock `CreateSecurityJSInjectionRule` and `DescribeSecurityJSInjectionRule` API calls
- **AND** verify the resource ID is set correctly and state is populated

#### Scenario: Read function unit test
- **WHEN** running unit tests for the Read function
- **THEN** the test SHALL mock `DescribeSecurityJSInjectionRule` API call
- **AND** verify the state is correctly populated from the API response

#### Scenario: Update function unit test
- **WHEN** running unit tests for the Update function
- **THEN** the test SHALL mock `ModifySecurityJSInjectionRule` and `DescribeSecurityJSInjectionRule` API calls
- **AND** verify the update request contains the correct parameters

#### Scenario: Delete function unit test
- **WHEN** running unit tests for the Delete function
- **THEN** the test SHALL mock `DeleteSecurityJSInjectionRule` API call
- **AND** verify the delete request contains the correct ZoneId and JSInjectionRuleIds
