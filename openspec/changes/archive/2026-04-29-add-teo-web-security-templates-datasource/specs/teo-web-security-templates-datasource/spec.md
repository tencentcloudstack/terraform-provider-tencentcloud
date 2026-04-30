## ADDED Requirements

### Requirement: Data source queries web security templates by zone IDs
The data source `tencentcloud_teo_web_security_templates` SHALL accept `zone_ids` as a required list of strings parameter and call the `DescribeWebSecurityTemplates` API with those zone IDs.

#### Scenario: Query templates with valid zone IDs
- **WHEN** user provides a list of zone IDs in the `zone_ids` parameter
- **THEN** the data source SHALL call `DescribeWebSecurityTemplates` with the provided zone IDs and return the matching security policy templates

#### Scenario: API call retry on failure
- **WHEN** the `DescribeWebSecurityTemplates` API call fails with a retryable error
- **THEN** the data source SHALL retry the call with `tccommon.ReadRetryTimeout` duration using `resource.Retry`

### Requirement: Data source returns security policy templates list
The data source SHALL return a `security_policy_templates` computed attribute containing a list of security policy template objects, each with `zone_id`, `template_id`, `template_name`, and `bind_domains` fields.

#### Scenario: Templates returned with all fields populated
- **WHEN** the API returns SecurityPolicyTemplateInfo objects with ZoneId, TemplateId, TemplateName, and BindDomains fields
- **THEN** each template in the output SHALL include `zone_id`, `template_id`, `template_name`, and `bind_domains` with values from the corresponding API fields

#### Scenario: Template with nil fields
- **WHEN** a SecurityPolicyTemplateInfo object has nil fields
- **THEN** the data source SHALL skip setting those nil fields (not set empty values)

### Requirement: Bind domains nested structure
Each security policy template SHALL include a `bind_domains` computed attribute as a list of objects, each with `domain`, `zone_id`, and `status` fields mapping from the `BindDomainInfo` API structure.

#### Scenario: Template with bound domains
- **WHEN** a SecurityPolicyTemplateInfo has BindDomains populated
- **THEN** each bind domain SHALL include `domain` (from BindDomainInfo.Domain), `zone_id` (from BindDomainInfo.ZoneId), and `status` (from BindDomainInfo.Status)

#### Scenario: Template with no bound domains
- **WHEN** a SecurityPolicyTemplateInfo has nil BindDomains
- **THEN** the `bind_domains` attribute SHALL not be set for that template

### Requirement: Data source result output file
The data source SHALL support an optional `result_output_file` parameter to save results to a file, consistent with other data sources.

#### Scenario: Result output file specified
- **WHEN** user provides a `result_output_file` path
- **THEN** the data source SHALL write the results to that file using `tccommon.WriteToFile`

### Requirement: Data source registration
The data source SHALL be registered in `tencentcloud/provider.go` with key `tencentcloud_teo_web_security_templates` and listed in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** the data source `tencentcloud_teo_web_security_templates` SHALL be available for use in Terraform configurations

### Requirement: Data source documentation
The data source SHALL have a corresponding `.md` documentation file following the project documentation conventions, with a one-line description mentioning TEO, example usage, and no Argument/Attribute Reference sections.

#### Scenario: Documentation file exists
- **WHEN** the data source is implemented
- **THEN** a `data_source_tc_teo_web_security_templates.md` file SHALL exist in `tencentcloud/services/teo/` with proper description, example usage, and import sections as applicable

### Requirement: Unit tests with gomonkey mock
The data source SHALL have unit tests using gomonkey mock approach for the cloud API, running with `go test -gcflags=all=-l`.

#### Scenario: Unit test covers read operation
- **WHEN** unit tests are executed
- **THEN** the read operation SHALL be tested with mocked API responses covering normal cases and empty results
