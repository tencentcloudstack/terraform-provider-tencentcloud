## ADDED Requirements

### Requirement: Data source returns accelerate regions list

The data source `tencentcloud_ga2_accelerate_regions` SHALL call the `DescribeAccelerateRegions` API and return the full list of available accelerate regions in the `accelerator_region_set` attribute.

#### Scenario: Successful query of accelerate regions

- **WHEN** user declares a `data "tencentcloud_ga2_accelerate_regions"` block in Terraform configuration
- **THEN** the data source SHALL call `DescribeAccelerateRegions` API and populate `accelerator_region_set` with all returned regions

#### Scenario: API returns empty list

- **WHEN** the `DescribeAccelerateRegions` API returns an empty `AcceleratorRegionSet`
- **THEN** the data source SHALL set `accelerator_region_set` to an empty list without error

### Requirement: Each region entry contains complete region information

Each element in `accelerator_region_set` SHALL contain the following computed attributes mapped from the API response:
- `name` (String): Region Chinese name, mapped from `AcceleratorRegionSet[].Name`
- `is_available` (Int): Availability status (0: unavailable, 1: available), mapped from `AcceleratorRegionSet[].IsAvailable`
- `region` (String): Region identifier, mapped from `AcceleratorRegionSet[].Region`
- `area_name` (String): Area name, mapped from `AcceleratorRegionSet[].AreaName`
- `is_china_mainland` (Int): Whether it is a China mainland region, mapped from `AcceleratorRegionSet[].IsChinaMainland`
- `support_isp_type` (List of String): Supported ISP types, mapped from `AcceleratorRegionSet[].SupportIspType`
- `is_tencent_region` (Int): Whether it is a Tencent region, mapped from `AcceleratorRegionSet[].IsTencentRegion`

#### Scenario: Region entry with all fields populated

- **WHEN** the API returns a region with all fields non-nil
- **THEN** the data source SHALL map all fields to the corresponding Terraform attributes

#### Scenario: Region entry with nil fields

- **WHEN** the API returns a region with some fields as nil
- **THEN** the data source SHALL only set non-nil fields and skip nil fields without error

### Requirement: API call uses retry mechanism

The data source SHALL use `tccommon.ReadRetryTimeout` as the timeout and wrap the API call with retry logic using `resource.RetryContext`.

#### Scenario: API call fails with retryable error

- **WHEN** the `DescribeAccelerateRegions` API returns a retryable error
- **THEN** the data source SHALL retry the call within the configured timeout

#### Scenario: API call fails with non-retryable error

- **WHEN** the `DescribeAccelerateRegions` API returns a non-retryable error
- **THEN** the data source SHALL return the error immediately without retry

### Requirement: Data source supports result output file

The data source SHALL support an optional `result_output_file` attribute that, when specified, writes the query results to the given file path.

#### Scenario: Result output file specified

- **WHEN** user specifies `result_output_file` in the data source configuration
- **THEN** the data source SHALL write the accelerator region set data to the specified file

#### Scenario: Result output file not specified

- **WHEN** user does not specify `result_output_file`
- **THEN** the data source SHALL not write any output file

### Requirement: Data source is registered in provider

The data source `tencentcloud_ga2_accelerate_regions` SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Provider includes the data source

- **WHEN** the Terraform provider is initialized
- **THEN** the data source `tencentcloud_ga2_accelerate_regions` SHALL be available for use
