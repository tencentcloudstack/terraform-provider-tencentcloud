## ADDED Requirements

### Requirement: Data source schema definition
The system SHALL provide a Terraform data source named `tencentcloud_ga2_cross_border_settlement` with the following schema:
- Required input parameters:
  - `global_accelerator_id` (String): The GA2 global accelerator instance ID
  - `accelerate_region` (String): The acceleration region
  - `endpoint_group_region` (String): The endpoint group region
  - `settlement_month` (Int): The billing year-month time
- Computed output attributes:
  - `traffic` (Float): Traffic usage in GB with 6 decimal places precision
- Optional parameter:
  - `result_output_file` (String): Path to output file for saving results

#### Scenario: Schema is correctly registered
- **WHEN** Terraform initializes the provider
- **THEN** the data source `tencentcloud_ga2_cross_border_settlement` SHALL be available with all required and computed fields defined

### Requirement: Query cross-border settlement traffic
The system SHALL call the `DescribeCrossBorderSettlement` API with the user-provided parameters and return the traffic usage value.

#### Scenario: Successful query with valid parameters
- **WHEN** user provides valid `global_accelerator_id`, `accelerate_region`, `endpoint_group_region`, and `settlement_month`
- **THEN** the system SHALL call `DescribeCrossBorderSettlement` API and set the `traffic` attribute to the returned `Traffic` value

#### Scenario: API returns nil traffic
- **WHEN** the API response `Traffic` field is nil
- **THEN** the system SHALL NOT set the `traffic` attribute (skip setting nil values)

### Requirement: Retry on API failure
The system SHALL wrap the API call with retry logic using `tccommon.ReadRetryTimeout` as the timeout duration.

#### Scenario: Transient API error with retry
- **WHEN** the `DescribeCrossBorderSettlement` API call fails with a transient error
- **THEN** the system SHALL retry the call within the `tccommon.ReadRetryTimeout` duration using `tccommon.RetryError()` to wrap the error

#### Scenario: Permanent API failure
- **WHEN** the API call fails with a non-retryable error after exhausting retries
- **THEN** the system SHALL return the error to Terraform

### Requirement: Resource ID composition
The system SHALL compose the data source ID from the four input parameters joined by `tccommon.FILED_SP` separator.

#### Scenario: ID is set after successful read
- **WHEN** the API call succeeds
- **THEN** the system SHALL set the resource ID to `global_accelerator_id#accelerate_region#endpoint_group_region#settlement_month` (using `tccommon.FILED_SP` as separator)

### Requirement: Provider registration
The system SHALL register the data source in `provider.go` and document it in `provider.md` under the GA2 service section.

#### Scenario: Data source is registered in provider
- **WHEN** the Terraform provider is built
- **THEN** `tencentcloud_ga2_cross_border_settlement` SHALL appear in the provider's data source map

### Requirement: Result output file support
The system SHALL support writing query results to a file when `result_output_file` is specified.

#### Scenario: Output to file
- **WHEN** user specifies `result_output_file` parameter
- **THEN** the system SHALL write the query results to the specified file path
