# Monitor Service - Spec Delta

## ADDED Requirements

### Requirement: External Cluster Register Command DataSource
Terraform Provider SHALL support querying Prometheus external cluster register commands through the `tencentcloud_monitor_external_cluster_register_command` data source.

#### Scenario: Query register command with valid parameters
- **WHEN** user queries `tencentcloud_monitor_external_cluster_register_command` data source with valid `instance_id` and `cluster_id`
- **THEN** the provider SHALL call `DescribeExternalClusterRegisterCommand` API
- **AND** the provider SHALL return the complete register command information
- **AND** the data source SHALL set ID to `{instanceId}#{clusterId}` format

#### Scenario: Missing required parameters
- **WHEN** user defines the data source without `instance_id` or `cluster_id`
- **THEN** Terraform validation SHALL fail with appropriate error message
- **AND** the API call SHALL not be executed

#### Scenario: Register command data mapping
- **WHEN** API successfully returns register command data
- **THEN** all returned fields SHALL be correctly mapped to data source computed attributes
- **AND** users SHALL be able to reference these values in Terraform configurations

#### Scenario: Export to file support
- **WHEN** user specifies `result_output_file` parameter
- **THEN** the provider SHALL write the query result to the specified file
- **AND** the file SHALL contain valid JSON representation of the data

### Requirement: DataSource Schema Compliance
The `tencentcloud_monitor_external_cluster_register_command` data source Schema SHALL comply with Terraform SDK v2 conventions and project standards.

#### Scenario: Required fields validation
- **WHEN** user defines the data source without required parameters
- **THEN** Terraform validation SHALL fail during plan phase
- **AND** a clear error message SHALL indicate which required fields are missing

#### Scenario: Computed fields population
- **WHEN** data source is read
- **THEN** all output fields SHALL be populated from API response
- **AND** the fields SHALL be marked as `Computed: true` in schema
- **AND** users SHALL be able to reference these computed values

#### Scenario: Data source ID generation
- **WHEN** data source read operation succeeds
- **THEN** the provider SHALL generate unique ID in format `{instanceId}#{clusterId}`
- **AND** this ID SHALL remain consistent across multiple reads with same parameters

### Requirement: Error Handling and Retry Logic
The data source implementation SHALL handle API errors gracefully with appropriate retry mechanisms.

#### Scenario: API transient error retry
- **WHEN** API call fails with retryable error (e.g., rate limit, temporary network issue)
- **THEN** the provider SHALL retry the operation using `tccommon.ReadRetryTimeout`
- **AND** the provider SHALL log each retry attempt with appropriate log level
- **AND** if retry succeeds, the operation SHALL complete normally

#### Scenario: Non-retryable error handling
- **WHEN** API call fails with non-retryable error (e.g., invalid parameters, resource not found)
- **THEN** the provider SHALL return error immediately without retry
- **AND** error message SHALL be descriptive and user-friendly
- **AND** Terraform SHALL halt execution with clear error information

#### Scenario: Resource not found handling
- **WHEN** `DescribeExternalClusterRegisterCommand` returns empty result or resource not found
- **THEN** the provider SHALL log a warning message
- **AND** the provider SHALL return an error indicating the cluster or instance does not exist
- **AND** Terraform SHALL report the error to user

### Requirement: API Parameter Mapping Correctness
The data source SHALL correctly map Terraform schema attributes to Tencent Cloud API parameters.

#### Scenario: API request parameter construction
- **WHEN** data source read operation is triggered
- **THEN** `InstanceId` parameter SHALL be set from `instance_id` input
- **AND** `ClusterId` parameter SHALL be set from `cluster_id` input
- **AND** both parameters SHALL be passed to `DescribeExternalClusterRegisterCommand` API
- **AND** API SHALL accept the parameters without validation errors

#### Scenario: API response field mapping
- **WHEN** API returns register command data
- **THEN** each response field SHALL be mapped to corresponding data source attribute
- **AND** nested structures SHALL be properly flattened or represented as nested blocks
- **AND** all field types SHALL match their schema definitions

### Requirement: Code Quality and Standards Compliance
The implementation SHALL follow project coding standards and reference implementation patterns.

#### Scenario: Code structure follows igtm_instance_list pattern
- **WHEN** reviewing the data source implementation code
- **THEN** the code structure SHALL closely match `data_source_tc_igtm_instance_list.go` patterns
- **AND** function naming SHALL follow project conventions (e.g., `dataSourceTencentCloudMonitorExternalClusterRegisterCommandRead`)
- **AND** error handling patterns SHALL be consistent with reference implementation

#### Scenario: Logging standards compliance
- **WHEN** data source performs read operation
- **THEN** operation SHALL be logged using `tccommon.LogElapsed` defer pattern
- **AND** API request and response bodies SHALL be logged at DEBUG level
- **AND** critical errors SHALL be logged at CRITICAL level with context

#### Scenario: Consistency check and validation
- **WHEN** data source operation completes
- **THEN** `defer tccommon.InconsistentCheck(d, meta)()` SHALL be called to validate state consistency
- **AND** any inconsistencies SHALL be detected and reported

#### Scenario: Result output file support
- **WHEN** data source implements result export functionality
- **THEN** it SHALL support `result_output_file` parameter
- **AND** use `tccommon.WriteToFile()` helper method
- **AND** follow the same pattern as reference data source implementation
