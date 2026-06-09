## ADDED Requirements

### Requirement: Resource supports full ModifySecurityPolicy API parameters

The `tencentcloud_teo_security_policy_config` resource SHALL support the following parameters mapped to the `ModifySecurityPolicy` cloud API:
- `zone_id` (Required, ForceNew, String): Maps to `request.ZoneId`. The site ID.
- `entity` (Optional, ForceNew, String): Maps to `request.Entity`. Security policy type, valid values: `ZoneDefaultPolicy`, `Template`, `Host`.
- `host` (Optional, ForceNew, String): Maps to `request.Host`. Domain name, used when entity is `Host`.
- `template_id` (Optional, ForceNew, String): Maps to `request.TemplateId`. Template ID, used when entity is `Template`.
- `security_config` (Optional, Computed, List): Maps to `request.SecurityConfig`. Legacy security configuration.
- `security_policy` (Optional, List): Maps to `request.SecurityPolicy`. Expression-based security policy configuration.

#### Scenario: Create security policy with ZoneDefaultPolicy entity
- **WHEN** user specifies `zone_id` and `entity = "ZoneDefaultPolicy"` with `security_policy` configuration
- **THEN** the resource SHALL call `ModifySecurityPolicy` with the provided parameters and set the resource ID to `{zone_id}#ZoneDefaultPolicy`

#### Scenario: Create security policy with Host entity
- **WHEN** user specifies `zone_id`, `entity = "Host"`, and `host = "www.example.com"` with `security_policy` configuration
- **THEN** the resource SHALL call `ModifySecurityPolicy` with the provided parameters and set the resource ID to `{zone_id}#Host#{host}`

#### Scenario: Create security policy with Template entity
- **WHEN** user specifies `zone_id`, `entity = "Template"`, and `template_id = "temp-xxx"` with `security_policy` configuration
- **THEN** the resource SHALL call `ModifySecurityPolicy` with the provided parameters and set the resource ID to `{zone_id}#Template#{template_id}`

### Requirement: Resource supports DescribeSecurityPolicy API for reading state

The resource SHALL use the `DescribeSecurityPolicy` API to read back the current state, passing `zone_id`, `entity`, `host`, and `template_id` as request parameters, and mapping `response.Response.SecurityPolicy` to the `security_policy` attribute.

#### Scenario: Read security policy state
- **WHEN** the resource performs a Read operation
- **THEN** the resource SHALL call `DescribeSecurityPolicy` with `zone_id`, `entity`, `host`, and `template_id` extracted from the resource ID, and populate the `security_policy` attribute from the response

#### Scenario: Resource not found during read
- **WHEN** the `DescribeSecurityPolicy` API returns nil for `SecurityPolicy`
- **THEN** the resource SHALL remove itself from state by calling `d.SetId("")`

### Requirement: Resource supports import via composite ID

The resource SHALL support import using the composite ID format based on entity type.

#### Scenario: Import ZoneDefaultPolicy
- **WHEN** user imports with ID `{zone_id}#ZoneDefaultPolicy`
- **THEN** the resource SHALL parse the ID and read the zone-level security policy

#### Scenario: Import Host policy
- **WHEN** user imports with ID `{zone_id}#Host#{host}`
- **THEN** the resource SHALL parse the ID and read the domain-level security policy

#### Scenario: Import Template policy
- **WHEN** user imports with ID `{zone_id}#Template#{template_id}`
- **THEN** the resource SHALL parse the ID and read the template-level security policy
