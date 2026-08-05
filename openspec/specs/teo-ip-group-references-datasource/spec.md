# teo-ip-group-references-datasource Specification

## Purpose
TBD - created by archiving change add-teo-ip-group-references-datasource. Update Purpose after archive.
## Requirements
### Requirement: Data source for querying TEO IP group references
The system SHALL provide a Terraform data source `tencentcloud_teo_ip_group_references` that queries the reference information of a specified IP group using the `DescribeIPGroupReferences` API.

#### Scenario: Query IP group references with required parameters
- **WHEN** user provides `zone_id` and `group_id` as required parameters
- **THEN** the data source SHALL call `DescribeIPGroupReferences` with the provided parameters and return `references`

#### Scenario: Automatic pagination for references
- **WHEN** the IP group is referenced by more entities than can be returned in a single API call
- **THEN** the data source SHALL automatically paginate through all results by using Limit=200 (the maximum value in the API annotation) and incrementing the offset until all references are retrieved

#### Scenario: Set computed fields from API response
- **WHEN** the API returns a successful response
- **THEN** the data source SHALL set `references` (TypeList) from `response.Response.References`, where each element maps the `IPGroupReference` fields (`zone_id`, `entity_type`, `entity_id`, `entity_name`, `sub_entity_type`, `sub_entity_id`, `sub_entity_name`), only if the response fields are not nil

#### Scenario: Generate data source ID
- **WHEN** the data source read completes successfully
- **THEN** the data source SHALL use `helper.BuildToken()` to generate the data source ID

### Requirement: Data source schema definition
The system SHALL define the following schema for `tencentcloud_teo_ip_group_references`:
- `zone_id` (TypeString, Required): Zone ID for the TEO site
- `group_id` (TypeInt, Required): IP group ID
- `references` (TypeList, Computed): List of references to the IP group. Each element is a schema.Resource with fields: `zone_id` (TypeString), `entity_type` (TypeString), `entity_id` (TypeString), `entity_name` (TypeString), `sub_entity_type` (TypeString), `sub_entity_id` (TypeString), `sub_entity_name` (TypeString)
- `result_output_file` (TypeString, Optional): Used to save results

#### Scenario: Schema fields match API parameters
- **WHEN** the data source schema is defined
- **THEN** `zone_id` SHALL map to API request parameter `ZoneId`, `group_id` SHALL map to API request parameter `GroupId`, and `references` SHALL map to API response field `References`

### Requirement: Error handling and retry
The system SHALL wrap the API call with `resource.Retry(tccommon.ReadRetryTimeout, ...)` and use `tccommon.RetryError()` for error handling. The retry block SHALL only call the cloud API; set operations SHALL be performed outside the retry block.

#### Scenario: API call fails with retryable error
- **WHEN** the `DescribeIPGroupReferences` API call fails
- **THEN** the error SHALL be wrapped with `tccommon.RetryError()` and retried up to `tccommon.ReadRetryTimeout`

#### Scenario: API returns empty response
- **WHEN** the API returns an empty response (`response == nil`, `response.Response == nil`, or `len(response.Response.References) == 0`)
- **THEN** the data source SHALL return a `NonRetryableError` (wrapped via `tccommon.RetryError`) inside the retry block instead of clearing the id, and the outer retry failure path SHALL log `[DATASOURCE] read empty, skip SetId` to aid troubleshooting

### Requirement: Provider registration
The system SHALL register the data source `tencentcloud_teo_ip_group_references` in `provider.go` and `provider.md`.

#### Scenario: Data source is available in Terraform
- **WHEN** the provider is initialized
- **THEN** the data source `tencentcloud_teo_ip_group_references` SHALL be available for use in Terraform configurations

### Requirement: Unit tests with gomonkey mock
The system SHALL provide unit tests using gomonkey mock for the data source, testing the Read operation without calling real cloud APIs.

#### Scenario: Unit test covers Read operation
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** the tests SHALL mock the `DescribeIPGroupReferences` API call and verify the data source correctly processes the response
