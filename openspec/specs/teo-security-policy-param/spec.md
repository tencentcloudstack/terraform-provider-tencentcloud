## ADDED Requirements

### Requirement: SecurityPolicy parameter support in ModifySecurityPolicy API

The `tencentcloud_teo_security_policy_config` resource SHALL support the `security_policy` parameter that maps to `request.SecurityPolicy` in the `ModifySecurityPolicy` API. When the user specifies the `security_policy` block in their terraform configuration, the resource MUST construct the `SecurityPolicy` struct and pass it to the `ModifySecurityPolicy` API request during create and update operations.

#### Scenario: User configures security_policy with custom rules
- **WHEN** user specifies a `security_policy` block with `custom_rules` configuration in their terraform resource
- **THEN** the resource MUST construct the `SecurityPolicy.CustomRules` struct and include it in the `ModifySecurityPolicy` API request

#### Scenario: User configures security_policy with managed rules
- **WHEN** user specifies a `security_policy` block with `managed_rules` configuration
- **THEN** the resource MUST construct the `SecurityPolicy.ManagedRules` struct and include it in the `ModifySecurityPolicy` API request

#### Scenario: User does not specify security_policy
- **WHEN** user does not include the `security_policy` block in their terraform configuration
- **THEN** the resource MUST NOT set `request.SecurityPolicy` on the `ModifySecurityPolicy` API request

### Requirement: DescribeSecurityPolicy reads SecurityPolicy response

The resource SHALL call the `DescribeSecurityPolicy` API with proper input parameters (`ZoneId`, `Entity`, `Host`, `TemplateId`) during read operations and MUST flatten the `response.Response.SecurityPolicy` output into the `security_policy` terraform attribute.

#### Scenario: Read operation retrieves security policy
- **WHEN** terraform performs a read/refresh operation on the `tencentcloud_teo_security_policy_config` resource
- **THEN** the resource MUST call `DescribeSecurityPolicy` with `ZoneId`, `Entity`, `Host`, and `TemplateId` extracted from the resource ID, and set the returned `SecurityPolicy` data into the `security_policy` attribute

#### Scenario: DescribeSecurityPolicy returns nil SecurityPolicy
- **WHEN** the `DescribeSecurityPolicy` API returns a nil `SecurityPolicy` in the response
- **THEN** the resource MUST NOT attempt to set the `security_policy` attribute and SHALL treat the resource as deleted (remove from state)

### Requirement: Input parameters for DescribeSecurityPolicy

The resource SHALL pass the following input parameters to the `DescribeSecurityPolicy` API:
- `ZoneId` from the `zone_id` schema field
- `Entity` from the `entity` schema field
- `Host` from the `host` schema field (when entity is "Host")
- `TemplateId` from the `template_id` schema field (when entity is "Template")

#### Scenario: Zone default policy query
- **WHEN** the resource ID contains entity "ZoneDefaultPolicy"
- **THEN** the resource MUST call `DescribeSecurityPolicy` with `ZoneId` and `Entity` set to "ZoneDefaultPolicy"

#### Scenario: Host-level policy query
- **WHEN** the resource ID contains entity "Host" with a host value
- **THEN** the resource MUST call `DescribeSecurityPolicy` with `ZoneId`, `Entity` set to "Host", and `Host` set to the domain name

#### Scenario: Template-level policy query
- **WHEN** the resource ID contains entity "Template" with a template_id value
- **THEN** the resource MUST call `DescribeSecurityPolicy` with `ZoneId`, `Entity` set to "Template", and `TemplateId` set to the template ID
