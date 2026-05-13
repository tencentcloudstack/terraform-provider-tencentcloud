## ADDED Requirements

### Requirement: Alias Domain Resource Schema
The resource \`tencentcloud_teo_alias_domain\` SHALL define a Terraform schema with the following fields:
- \`zone_id\` (Required, String, ForceNew): 站点 ID
- \`alias_name\` (Required, String, ForceNew): 别称域名名称
- \`target_name\` (Required, String): 目标域名名称
- \`cert_type\` (Optional, String): 证书配置类型，创建时可选值为 "none" 或 "hosting"，修改时额外支持 "apply"
- \`cert_id\` (Optional, List of String): 证书 ID 列表，当 cert_type 为 "hosting" 时必填
- \`status\` (Computed, String): 别称域名状态
- \`forbid_mode\` (Computed, Int): 封禁模式
- \`created_on\` (Computed, String): 创建时间
- \`modified_on\` (Computed, String): 修改时间

#### Scenario: Schema definition with required fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include \`zone_id\`, \`alias_name\`, and \`target_name\` as required fields, with \`zone_id\` and \`alias_name\` marked as ForceNew

#### Scenario: Schema definition with optional fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include \`cert_type\` and \`cert_id\` as optional fields

#### Scenario: Schema definition with computed fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include \`status\`, \`forbid_mode\`, \`created_on\`, and \`modified_on\` as computed fields

### Requirement: Alias Domain Create
The resource SHALL support creating an alias domain by calling \`CreateAliasDomain\` API with zone_id, alias_name, target_name, and optionally cert_type and cert_id.

#### Scenario: Create with required fields only
- **WHEN** a user creates an alias domain resource with zone_id, alias_name, and target_name
- **THEN** the system SHALL call \`CreateAliasDomain\` API with ZoneId, AliasName, and TargetName
- **AND** set the resource ID to \`zone_id#alias_name\` using \`tccommon.FILED_SP\` as separator

#### Scenario: Create with certificate configuration
- **WHEN** a user creates an alias domain resource with cert_type="hosting" and cert_id specified
- **THEN** the system SHALL call \`CreateAliasDomain\` API with CertType and CertId parameters included

#### Scenario: Create API call uses retry
- **WHEN** the CreateAliasDomain API is called
- **THEN** the system SHALL wrap the API call with \`resource.Retry(tccommon.WriteRetryTimeout, ...)\` and use \`tccommon.RetryError(e)\` for error handling

#### Scenario: Create response validation
- **WHEN** the CreateAliasDomain API returns a response
- **THEN** the system SHALL check if the response is empty and return NonRetryableError if so

### Requirement: Alias Domain Read
The resource SHALL support reading an alias domain by parsing the composite ID and calling \`DescribeAliasDomains\` API with filters.

#### Scenario: Read with valid resource ID
- **WHEN** a user reads an alias domain resource with ID \`zone_id#alias_name\`
- **THEN** the system SHALL parse the ID to extract zone_id and alias_name
- **AND** call the service layer \`DescribeTeoAliasDomainById\` method with zone_id and alias_name

#### Scenario: Read sets computed fields from API response
- **WHEN** the DescribeAliasDomains API returns an AliasDomain object
- **THEN** the system SHALL set \`status\`, \`forbid_mode\`, \`created_on\`, \`modified_on\`, and \`target_name\` from the response
- **AND** the system SHALL check each field for nil before setting it

#### Scenario: Read handles resource not found
- **WHEN** the DescribeAliasDomains API returns no matching alias domain
- **THEN** the system SHALL set \`d.SetId("")\` and return nil

#### Scenario: Read does not set CertType and CertId
- **WHEN** the Read operation is performed
- **THEN** the system SHALL NOT attempt to set \`cert_type\` and \`cert_id\` from the API response since they are not returned by DescribeAliasDomains

### Requirement: Alias Domain Update
The resource SHALL support updating an alias domain by calling \`ModifyAliasDomain\` API when mutable fields change.

#### Scenario: Update target_name
- **WHEN** a user changes the \`target_name\` field
- **THEN** the system SHALL call \`ModifyAliasDomain\` API with ZoneId, AliasName, and TargetName

#### Scenario: Update cert_type and cert_id
- **WHEN** a user changes the \`cert_type\` or \`cert_id\` field
- **THEN** the system SHALL call \`ModifyAliasDomain\` API with ZoneId, AliasName, CertType, and CertId

#### Scenario: Update checks mutable args
- **WHEN** the Update function is called
- **THEN** the system SHALL check \`d.HasChange()\` for \`target_name\`, \`cert_type\`, and \`cert_id\` before making the API call

#### Scenario: Update API call uses retry
- **WHEN** the ModifyAliasDomain API is called
- **THEN** the system SHALL wrap the API call with \`resource.Retry(tccommon.WriteRetryTimeout, ...)\` and use \`tccommon.RetryError(e)\` for error handling

### Requirement: Alias Domain Delete
The resource SHALL support deleting an alias domain by calling \`DeleteAliasDomain\` API.

#### Scenario: Delete with valid resource ID
- **WHEN** a user deletes an alias domain resource with ID \`zone_id#alias_name\`
- **THEN** the system SHALL parse the ID to extract zone_id and alias_name
- **AND** call \`DeleteAliasDomain\` API with ZoneId and AliasNames containing the single alias_name

#### Scenario: Delete API call uses retry
- **WHEN** the DeleteAliasDomain API is called
- **THEN** the system SHALL wrap the API call with \`resource.Retry(tccommon.WriteRetryTimeout, ...)\` and use \`tccommon.RetryError(e)\` for error handling

### Requirement: Alias Domain Import
The resource SHALL support importing an existing alias domain.

#### Scenario: Import with composite ID
- **WHEN** a user imports an alias domain with ID \`zone_id#alias_name\`
- **THEN** the system SHALL parse the ID and call the Read function
- **AND** the Read function SHALL populate all available fields from the API response

#### Scenario: Import documentation
- **WHEN** the resource documentation is created
- **THEN** it SHALL include an import example showing the composite ID format \`zone_id#alias_name\`

### Requirement: Service Layer Method
The service layer SHALL provide a \`DescribeTeoAliasDomainById\` method in \`TeoService\`.

#### Scenario: Service method uses DescribeAliasDomains API
- **WHEN** \`DescribeTeoAliasDomainById\` is called with zone_id and alias_name
- **THEN** it SHALL call \`DescribeAliasDomains\` API with Filters containing alias-name exact match

#### Scenario: Service method handles pagination
- **WHEN** the DescribeAliasDomains API returns results
- **THEN** the service method SHALL use Limit=1000 (API maximum) and handle Offset-based pagination

#### Scenario: Service method uses retry
- **WHEN** the DescribeAliasDomains API is called
- **THEN** it SHALL be wrapped with \`resource.Retry(tccommon.ReadRetryTimeout, ...)\`

#### Scenario: Service method returns nil when not found
- **WHEN** no matching alias domain is found
- **THEN** the service method SHALL return nil without error

### Requirement: Provider Registration
The resource SHALL be registered in \`provider.go\`.

#### Scenario: Resource registration in provider
- **WHEN** the provider is initialized
- **THEN** \`tencentcloud_teo_alias_domain\` SHALL be mapped to \`teo.ResourceTencentCloudTeoAliasDomain()\` in the resource map

### Requirement: Unit Tests
The resource SHALL have unit tests using mock (gomonkey) approach.

#### Scenario: Create function unit test
- **WHEN** unit tests are executed
- **THEN** the Create function SHALL be tested with mock for the CreateAliasDomain API call

#### Scenario: Read function unit test
- **WHEN** unit tests are executed
- **THEN** the Read function SHALL be tested with mock for the DescribeAliasDomains API call

#### Scenario: Update function unit test
- **WHEN** unit tests are executed
- **THEN** the Update function SHALL be tested with mock for the ModifyAliasDomain API call

#### Scenario: Delete function unit test
- **WHEN** unit tests are executed
- **THEN** the Delete function SHALL be tested with mock for the DeleteAliasDomain API call
