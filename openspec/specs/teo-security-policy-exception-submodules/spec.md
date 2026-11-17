# teo-security-policy-exception-submodules Specification

## Purpose
TBD - created by syncing change add-teo-security-policy-exception-submodules. Update Purpose after archive.

## Requirements

### Requirement: web_security_submodules_for_exception schema definition
The `tencentcloud_teo_security_policy_config` resource SHALL include an Optional `web_security_submodules_for_exception` parameter (TypeSet, Elem of TypeString) under the `security_policy.exception_rules.rules` block, placed alongside the existing `web_security_modules_for_exception` field. This parameter specifies the security protection **submodules** that the exception rule skips, and is only valid when `skip_scope` is `WebSecuritySubmodules`.

#### Scenario: Field is optional and backward compatible
- **WHEN** a user creates or updates a `tencentcloud_teo_security_policy_config` resource without specifying `web_security_submodules_for_exception` in `exception_rules.rules`
- **THEN** the resource SHALL be created/updated successfully and the existing behavior SHALL be preserved without plan diff

#### Scenario: Field accepts a set of submodule strings
- **WHEN** a user provides `web_security_submodules_for_exception = ["websec-mod-managed-rules/managed-rule-groups", "websec-mod-rate-limiting-rules"]` in an `exception_rules.rules` block
- **THEN** the schema SHALL accept the values as a TypeSet of strings and mark them as candidates for the `ModifySecurityPolicy` request

#### Scenario: Field is a TypeSet (unordered, deduplicated)
- **WHEN** a user provides duplicate or differently-ordered values in `web_security_submodules_for_exception`
- **THEN** the schema SHALL treat the collection as a set (order-insensitive, deduplicated), consistent with the sibling `web_security_modules_for_exception` field

### Requirement: Read operation for web_security_submodules_for_exception
The resource Read operation SHALL populate `web_security_submodules_for_exception` from the `DescribeSecurityPolicy` API response, reading `SecurityPolicy.ExceptionRules.Rules[].WebSecuritySubmodulesForException` and flattening it into the Terraform state, with a nil check before setting.

#### Scenario: Read with WebSecuritySubmodulesForException present
- **WHEN** the `DescribeSecurityPolicy` response contains a non-nil `WebSecuritySubmodulesForException` field on an exception rule
- **THEN** the Read operation SHALL set `web_security_submodules_for_exception` in the rules map with the returned `[]*string` value

#### Scenario: Read with WebSecuritySubmodulesForException absent
- **WHEN** the `DescribeSecurityPolicy` response has a nil `WebSecuritySubmodulesForException` field on an exception rule
- **THEN** the Read operation SHALL NOT set `web_security_submodules_for_exception` in the rules map (leave it empty), consistent with the existing nil-check pattern

### Requirement: Create and Update operations populate web_security_submodules_for_exception
The resource Create and Update operations SHALL expand the `web_security_submodules_for_exception` Terraform configuration (a `*schema.Set`) into the `SecurityPolicy.ExceptionRules.Rules[].WebSecuritySubmodulesForException` field (a `[]*string`) of the `ModifySecurityPolicy` API request, iterating over the set and appending a pointer to each non-nil string element.

#### Scenario: Create with web_security_submodules_for_exception
- **WHEN** a user creates a resource with `web_security_submodules_for_exception` configured in an `exception_rules.rules` block
- **THEN** the Create operation SHALL expand the set and include `WebSecuritySubmodulesForException` in the `ExceptionRule` of the `ModifySecurityPolicy` request

#### Scenario: Update with web_security_submodules_for_exception changes
- **WHEN** a user updates the resource's `web_security_submodules_for_exception` configuration
- **THEN** the Update operation SHALL expand the new configuration and include `WebSecuritySubmodulesForException` in the `ModifySecurityPolicy` request

#### Scenario: Set elements are converted to []*string
- **WHEN** the `web_security_submodules_for_exception` set contains one or more string values
- **THEN** the provider SHALL iterate the set, take the address of each non-nil string, and append it to `exceptionRule.WebSecuritySubmodulesForException`

#### Scenario: Empty or absent set does not populate the field
- **WHEN** the `web_security_submodules_for_exception` set is empty or absent from the configuration
- **THEN** the provider SHALL NOT populate `exceptionRule.WebSecuritySubmodulesForException` (leave it nil), consistent with the sibling field's handling

### Requirement: Unit tests for web_security_submodules_for_exception
The system SHALL provide unit tests in `resource_tc_teo_security_policy_config_test.go` using gomonkey to mock the cloud API calls, covering the new `web_security_submodules_for_exception` parameter across Create, Read, and Update operations.

#### Scenario: Unit test covers Create with the new parameter
- **WHEN** a unit test simulates creating a resource with `web_security_submodules_for_exception` set
- **THEN** the mocked `ModifySecurityPolicyWithContext` SHALL be invoked with an `ExceptionRule.WebSecuritySubmodulesForException` matching the configured values

#### Scenario: Unit test covers Read populating the new parameter
- **WHEN** a unit test simulates reading a resource whose `DescribeSecurityPolicy` response includes `WebSecuritySubmodulesForException`
- **THEN** the test SHALL assert `web_security_submodules_for_exception` is populated in the Terraform state with the expected values

#### Scenario: Unit test covers Update changing the new parameter
- **WHEN** a unit test simulates updating `web_security_submodules_for_exception` to a new set of values
- **THEN** the mocked `ModifySecurityPolicyWithContext` SHALL be invoked with the updated `ExceptionRule.WebSecuritySubmodulesForException`

### Requirement: Resource documentation for web_security_submodules_for_exception
The system SHALL update the `resource_tc_teo_security_policy_config.md` example to demonstrate the `web_security_submodules_for_exception` field within an `exception_rules.rules` block, and the `website/docs/` documentation SHALL be generated via `make doc` in the finalization phase (no manual edits to `website/`).

#### Scenario: Documentation example includes the new field
- **WHEN** the resource markdown is updated
- **THEN** the Example Usage SHALL show `web_security_submodules_for_exception` inside an `exception_rules.rules` block with representative submodule values

#### Scenario: Website docs generated via make doc
- **WHEN** the finalization phase runs `make doc`
- **THEN** the `website/docs/` markdown SHALL be regenerated from the resource `.md` file without any manual edits to the `website/` directory
