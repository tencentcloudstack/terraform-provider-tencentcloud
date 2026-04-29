## ADDED Requirements

### Requirement: Data source for querying TEO security IP group content
The system SHALL provide a Terraform data source `tencentcloud_teo_security_ip_group_content` that queries the IP list within a specified security IP group using the `DescribeSecurityIPGroupContent` API.

#### Scenario: Query security IP group content with required parameters
- **WHEN** user provides `zone_id` and `group_id` as required parameters
- **THEN** the data source SHALL call `DescribeSecurityIPGroupContent` with the provided parameters and return `ip_total_count` and `ip_list`

#### Scenario: Automatic pagination for IP list
- **WHEN** the IP group contains more IPs than can be returned in a single API call
- **THEN** the data source SHALL automatically paginate through all results by incrementing the offset until all IPs are retrieved

#### Scenario: Set computed fields from API response
- **WHEN** the API returns a successful response
- **THEN** the data source SHALL set `ip_total_count` (TypeInt) from `response.Response.IPTotalCount` and `ip_list` (TypeList of TypeString) from `response.Response.IPList`, only if the response fields are not nil

#### Scenario: Generate data source ID
- **WHEN** the data source read completes successfully
- **THEN** the data source SHALL use `helper.BuildToken()` to generate the data source ID

### Requirement: Data source schema definition
The system SHALL define the following schema for `tencentcloud_teo_security_ip_group_content`:
- `zone_id` (TypeString, Required): Zone ID for the TEO site
- `group_id` (TypeInt, Required): IP group ID
- `ip_total_count` (TypeInt, Computed): Total count of IPs in the group
- `ip_list` (TypeList of TypeString, Computed): List of IPs or CIDR blocks in the group
- `result_output_file` (TypeString, Optional): Used to save results

#### Scenario: Schema fields match API parameters
- **WHEN** the data source schema is defined
- **THEN** `zone_id` SHALL map to API request parameter `ZoneId`, `group_id` SHALL map to API request parameter `GroupId`, `ip_total_count` SHALL map to API response field `IPTotalCount`, and `ip_list` SHALL map to API response field `IPList`

### Requirement: Error handling and retry
The system SHALL wrap the API call with `resource.Retry(tccommon.ReadRetryTimeout, ...)` and use `tccommon.RetryError()` for error handling.

#### Scenario: API call fails with retryable error
- **WHEN** the `DescribeSecurityIPGroupContent` API call fails
- **THEN** the error SHALL be wrapped with `tccommon.RetryError()` and retried up to `tccommon.ReadRetryTimeout`

### Requirement: Provider registration
The system SHALL register the data source `tencentcloud_teo_security_ip_group_content` in `provider.go` and `provider.md`.

#### Scenario: Data source is available in Terraform
- **WHEN** the provider is initialized
- **THEN** the data source `tencentcloud_teo_security_ip_group_content` SHALL be available for use in Terraform configurations

### Requirement: Unit tests with gomonkey mock
The system SHALL provide unit tests using gomonkey mock for the data source, testing the Read operation without calling real cloud APIs.

#### Scenario: Unit test covers Read operation
- **WHEN** unit tests are executed with `go test -gcflags=all=-l`
- **THEN** the tests SHALL mock the `DescribeSecurityIPGroupContent` API call and verify the data source correctly processes the response
