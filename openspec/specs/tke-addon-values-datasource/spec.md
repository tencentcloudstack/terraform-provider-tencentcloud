# tke-addon-values-datasource Specification

## Purpose

Provide a Terraform data source `tencentcloud_tke_addon_values` that queries the TKE `DescribeAddonValues` API for a given cluster and addon, exposing the addon's currently effective parameters (`values`) and the chart's default parameters (`default_values`) so they can be referenced by other Terraform resources and data sources.

提供 Terraform 数据源 `tencentcloud_tke_addon_values`，用于查询 TKE 指定集群下指定 addon 的当前生效参数（`values`）与 chart 默认参数（`default_values`），以便其他 Terraform 资源与数据源通过 `data.tencentcloud_tke_addon_values.xxx` 的方式引用。

## Requirements
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
- **AND** it SHALL convert errors via `tccommon.RetryError(e)`
- **AND** it SHALL return the retry error to Terraform when retries are exhausted

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
- **THEN** the read handler SHALL call `tccommon.WriteToFile(output, d)`
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

