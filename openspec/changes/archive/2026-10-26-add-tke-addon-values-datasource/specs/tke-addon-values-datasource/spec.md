# Capability: TKE Addon Values Data Source

## ADDED Requirements

### Requirement: Data Source Registration

The provider SHALL register a new data source named `tencentcloud_tke_addon_values` and document it in the provider's data source index.

#### Scenario: Data source is registered in DataSourcesMap

- **WHEN** the provider initializes its `DataSourcesMap`
- **THEN** the key `tencentcloud_tke_addon_values` SHALL map to `tke.DataSourceTencentCloudTkeAddonValues()`
- **AND** the same name SHALL be listed in `tencentcloud/provider.md`

#### Scenario: Data source file follows naming convention

- **WHEN** the data source is implemented
- **THEN** the implementation file SHALL be `tencentcloud/services/tke/data_source_tc_tke_addon_values.go` in package `tke`
- **AND** the exported constructor SHALL be `DataSourceTencentCloudTkeAddonValues`
- **AND** the read handler SHALL be `dataSourceTencentCloudTkeAddonValuesRead`
- **AND** the file SHALL NOT start with a generated-code comment

### Requirement: Input Arguments

The data source SHALL require `cluster_id` and `addon_name`, both of type string, and map them to the cloud API request fields `ClusterId` and `AddonName`.

#### Scenario: cluster_id is required

- **WHEN** defining the data source schema
- **THEN** `cluster_id` SHALL be `schema.TypeString` with `Required: true`
- **AND** its description SHALL identify it as the TKE cluster ID

#### Scenario: addon_name is required

- **WHEN** defining the data source schema
- **THEN** `addon_name` SHALL be `schema.TypeString` with `Required: true`
- **AND** its description SHALL identify it as the addon name

#### Scenario: Inputs are validated before calling the cloud API

- **WHEN** defining the data source schema
- **THEN** both `cluster_id` and `addon_name` SHALL have `ValidateFunc: validation.StringIsNotEmpty`
- **AND** an empty value SHALL be rejected at plan time, so the request never reaches the cloud API

#### Scenario: Inputs are passed to the cloud API request

- **WHEN** the read handler executes
- **THEN** the value of `cluster_id` SHALL be assigned to `DescribeAddonValuesRequest.ClusterId`
- **AND** the value of `addon_name` SHALL be assigned to `DescribeAddonValuesRequest.AddonName`

#### Scenario: result_output_file is optional

- **WHEN** defining the data source schema
- **THEN** `result_output_file` SHALL be `schema.TypeString` with `Optional: true`
- **AND** it SHALL be the last field declared in the schema map

### Requirement: Computed Outputs

The data source SHALL expose `values` and `default_values` as computed string fields taken verbatim from the cloud API response.

#### Scenario: values and default_values are computed strings

- **WHEN** defining the data source schema
- **THEN** `values` SHALL be `schema.TypeString` with `Computed: true`, mapped from `Response.Values`
- **AND** `default_values` SHALL be `schema.TypeString` with `Computed: true`, mapped from `Response.DefaultValues`
- **AND** neither field SHALL be decoded or restructured from JSON

#### Scenario: not-installed addon still returns a rendered payload

- **WHEN** the addon is not installed in the cluster but its chart exists in the addon repository
- **THEN** the cloud API still returns a rendered `Values` payload (the chart defaults) which is byte-identical to `DefaultValues`
- **AND** the data source SHALL pass it through unchanged (it SHALL NOT clear `values`, and it SHALL NOT pretend the addon is installed)
- **AND** the documentation SHALL state that `values` is NOT empty in this case, and that `tencentcloud_kubernetes_addons` is the way to check the installation status

#### Scenario: nil response fields are skipped

- **WHEN** the cloud API response contains a nil `Values` or nil `DefaultValues`
- **THEN** the read handler SHALL NOT call `d.Set` for that field
- **AND** the read handler SHALL NOT return an error for that reason
- **AND** the read handler SHALL still set the resource ID

### Requirement: Read Behavior

The read handler SHALL call `DescribeAddonValues` through the existing service method `TkeService.DescribeAddonValuesList`, wrapped in the standard retry and logging pattern.

#### Scenario: Successful read

- **WHEN** the read handler is invoked with valid `cluster_id` and `addon_name`
- **THEN** it SHALL call `TkeService.DescribeAddonValuesList(ctx, clusterId, addonName)` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`
- **AND** it SHALL convert retryable errors via `tccommon.RetryError(e)`
- **AND** it SHALL return the retry error to Terraform when retries are exhausted

#### Scenario: deterministic errors are not retried

- **WHEN** the cloud API returns a deterministic failure, i.e. `Code=ResourceUnavailable` with message `not found chart` or `resource name may not be empty`
- **THEN** the read handler SHALL return `resource.NonRetryableError` immediately (a single API call)
- **AND** it SHALL NOT consume the whole `ReadRetryTimeout` (3 minutes by default)
- **AND** genuinely retryable codes (e.g. `RequestLimitExceeded`) SHALL still be retried

#### Scenario: Standard deferred helpers

- **WHEN** the read handler starts
- **THEN** it SHALL register `defer tccommon.LogElapsed("data_source.tencentcloud_tke_addon_values.read")()`
- **AND** it SHALL register `defer tccommon.InconsistentCheck(d, meta)()`

#### Scenario: Resource ID is a composite of cluster and addon

- **WHEN** the read handler completes successfully
- **THEN** it SHALL set the ID to `cluster_id` and `addon_name` joined by `tccommon.FILED_SP`
- **AND** it SHALL NOT set an empty ID

#### Scenario: Empty cloud API response is an error

- **WHEN** the cloud API returns a nil `Response` or nil `Response.Response`
- **THEN** the read handler SHALL surface an error to Terraform
- **AND** it SHALL NOT call `d.SetId("")`

#### Scenario: Result written to file when requested

- **WHEN** `result_output_file` is configured with a non-empty path
- **THEN** the read handler SHALL call `tccommon.WriteToFile(output, <payload>)` with the actual payload, i.e. a `map[string]interface{}` containing `values` and `default_values`
- **AND** it SHALL NOT pass `d` (`*schema.ResourceData`), which marshals to `{}`
- **AND** it SHALL return the error if writing fails

### Requirement: Documentation

The data source SHALL ship a Markdown template used by `make doc` to generate the website documentation.

#### Scenario: Documentation template structure

- **WHEN** creating `tencentcloud/services/tke/data_source_tc_tke_addon_values.md`
- **THEN** the first line SHALL be a single-sentence description mentioning TKE, in the form `Use this data source to query ...`
- **AND** it SHALL include an `Example Usage` section in HCL that references `data.tencentcloud_tke_addon_values`
- **AND** the example SHALL use `jsondecode()` when consuming `values` / `default_values`
- **AND** it SHALL NOT contain `Import`, `Argument Reference`, or `Attribute Reference` sections

### Requirement: Unit Test Coverage

The data source SHALL be covered by mock-based unit tests that do not require real cloud credentials.

#### Scenario: Cloud API is mocked with gomonkey

- **WHEN** writing `data_source_tc_tke_addon_values_test.go`
- **THEN** it SHALL use `gomonkey` to patch `TkeService.DescribeAddonValuesList`
- **AND** it SHALL provide a mock implementing `tccommon.ProviderMeta` whose `GetAPIV3Conn()` returns a `*connectivity.TencentCloudClient`
- **AND** it SHALL NOT use the Terraform acceptance test suite

#### Scenario: Covered cases

- **WHEN** the unit tests run
- **THEN** they SHALL assert both outputs are populated on a normal response
- **AND** they SHALL assert nil `Values` is handled without error
- **AND** they SHALL assert a service-layer error is propagated
- **AND** they SHALL assert the schema declares `cluster_id`, `addon_name`, `values`, `default_values`, and `result_output_file` with the expected Required/Computed flags
- **AND** they SHALL assert a not-installed addon keeps its rendered `values` (equal to `default_values`) instead of being cleared
- **AND** they SHALL assert a deterministic error is propagated after a single call (no retry)
- **AND** they SHALL assert a retryable error is still retried
- **AND** they SHALL assert `result_output_file` contains the real payload (not `{}`)
- **AND** they SHALL assert empty `cluster_id` / `addon_name` are rejected by `ValidateFunc`
