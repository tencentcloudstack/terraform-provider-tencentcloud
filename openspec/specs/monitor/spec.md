# Monitor Service Specification

## Requirement: External Cluster Resource Management
Terraform Provider SHALL support managing Prometheus external clusters through the `tencentcloud_monitor_external_cluster` resource.

### Scenario: Register external cluster to TMP instance
- **WHEN** user creates a `tencentcloud_monitor_external_cluster` resource with valid `instance_id` and `cluster_region`
- **THEN** the provider SHALL call `CreateExternalCluster` API
- **AND** the provider SHALL store the resource with composite ID format `{instanceId}#{clusterId}`
- **AND** the cluster SHALL be successfully registered to the TMP instance

### Scenario: Query external cluster information
- **WHEN** Terraform performs a read operation on the external cluster resource
- **THEN** the provider SHALL call `DescribePrometheusClusterAgents` API with `InstanceId` and `ClusterIds` parameters
- **AND** the `ClusterIds` parameter SHALL be formatted as a string array `[clusterId]`
- **AND** the provider SHALL correctly map the `ClusterType` field from API response to Terraform state as a computed attribute

### Scenario: Delete external cluster association
- **WHEN** user destroys a `tencentcloud_monitor_external_cluster` resource
- **THEN** the provider SHALL parse the composite resource ID to extract `instanceId` and `clusterId`
- **AND** the provider SHALL retrieve `cluster_type` value from Terraform state using `d.GetOk("cluster_type")`
- **AND** the provider SHALL call `DeletePrometheusClusterAgent` API with correct `Agents` parameter structure
- **AND** the cluster association SHALL be successfully removed from the TMP instance

### Scenario: Resource import support
- **WHEN** user imports an existing external cluster using `terraform import`
- **THEN** the provider SHALL accept resource ID in format `{instanceId}#{clusterId}`
- **AND** the provider SHALL retrieve cluster information via `DescribePrometheusClusterAgents` API
- **AND** all resource attributes SHALL be correctly populated in Terraform state

## Requirement: Resource Schema Compliance
The `tencentcloud_monitor_external_cluster` resource Schema SHALL comply with Terraform SDK v2 conventions and project standards.

### Scenario: Required fields validation
- **WHEN** user defines the resource without `instance_id` or `cluster_region`
- **THEN** Terraform validation SHALL fail with appropriate error message
- **AND** resource creation SHALL not proceed

### Scenario: ForceNew attribute behavior
- **WHEN** user modifies `instance_id` in existing resource configuration
- **THEN** Terraform SHALL plan to destroy and recreate the resource
- **AND** the provider SHALL not attempt in-place update

### Scenario: Computed fields population
- **WHEN** resource is created or refreshed
- **THEN** the `cluster_type` field SHALL be populated from API response
- **AND** the field SHALL be marked as `Computed: true` in schema
- **AND** users SHALL be able to reference this value in Terraform configurations

## Requirement: Error Handling and Retry Logic
The resource implementation SHALL handle API errors gracefully with appropriate retry mechanisms.

### Scenario: API transient error retry
- **WHEN** API call fails with retryable error (e.g., rate limit, temporary network issue)
- **THEN** the provider SHALL retry the operation using `tccommon.WriteRetryTimeout`
- **AND** the provider SHALL log each retry attempt with appropriate log level
- **AND** if retry succeeds, the operation SHALL complete normally

### Scenario: Non-retryable error handling
- **WHEN** API call fails with non-retryable error (e.g., invalid parameters, resource not found)
- **THEN** the provider SHALL return error immediately without retry
- **AND** error message SHALL be descriptive and user-friendly
- **AND** Terraform SHALL halt execution with clear error information

### Scenario: Resource not found during read
- **WHEN** `DescribePrometheusClusterAgents` returns empty result or resource not found
- **THEN** the provider SHALL log a warning message
- **AND** the provider SHALL set resource ID to empty string to mark for recreation
- **AND** Terraform SHALL detect drift and plan to recreate resource if it exists in state

## Requirement: API Parameter Mapping Correctness
The resource SHALL correctly map Terraform schema attributes to Tencent Cloud API parameters.

### Scenario: External labels array mapping
- **WHEN** user provides `external_labels` block in resource configuration
- **THEN** each label SHALL be mapped to API `Label` structure with `Name` and `Value` fields
- **AND** the labels array SHALL be passed to `CreateExternalCluster` API
- **AND** API SHALL accept the formatted labels without validation errors

### Scenario: ClusterIds parameter formatting
- **WHEN** provider calls `DescribePrometheusClusterAgents` API during read operation
- **THEN** the `ClusterIds` parameter SHALL be formatted as string array `[clusterId]`
- **AND** API SHALL successfully filter results by the provided cluster ID
- **AND** provider SHALL receive correct cluster information in response

### Scenario: Agents parameter construction for deletion
- **WHEN** provider calls `DeletePrometheusClusterAgent` API during delete operation
- **THEN** `Agents` parameter SHALL be constructed as array of `PrometheusAgentInfo` objects
- **AND** each agent object SHALL contain `ClusterId` parsed from resource ID
- **AND** each agent object SHALL contain `ClusterType` retrieved from Terraform state
- **AND** API SHALL successfully process the deletion request

## Requirement: Code Quality and Standards Compliance
The implementation SHALL follow project coding standards and reference implementation patterns.

### Scenario: Code structure follows igtm_strategy pattern
- **WHEN** reviewing the resource implementation code
- **THEN** the code structure SHALL closely match `resource_tc_igtm_strategy.go` patterns
- **AND** function naming SHALL follow project conventions (e.g., `resourceTencentCloudMonitorExternalClusterCreate`)
- **AND** error handling patterns SHALL be consistent with reference implementation

### Scenario: Logging standards compliance
- **WHEN** resource performs any operation (create, read, update, delete)
- **THEN** operations SHALL be logged using `tccommon.LogElapsed` defer pattern
- **AND** API request and response bodies SHALL be logged at DEBUG level
- **AND** critical errors SHALL be logged at CRITICAL level with context

### Scenario: Consistency check and validation
- **WHEN** resource operation completes
- **THEN** `defer tccommon.InconsistentCheck(d, meta)()` SHALL be called to validate state consistency
- **AND** any inconsistencies SHALL be detected and reported
- **AND** Terraform state SHALL remain valid after all operations
