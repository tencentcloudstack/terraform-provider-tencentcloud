## Requirements

### Requirement: Data source tencentcloud_teo_content_quota schema definition
The system SHALL provide a data source `tencentcloud_teo_content_quota` with the following schema:
- `zone_id` (TypeString, Required): 站点 ID，用于指定查询配额的站点
- `purge_quota` (TypeList, Computed): 刷新相关配额列表，每个元素包含 `batch`(TypeInt)、`daily`(TypeInt)、`daily_available`(TypeInt)、`type`(TypeString)
- `prefetch_quota` (TypeList, Computed): 预热相关配额列表，每个元素包含 `batch`(TypeInt)、`daily`(TypeInt)、`daily_available`(TypeInt)、`type`(TypeString)
- `result_output_file` (TypeString, Optional): 用于保存结果

#### Scenario: Schema defines all required and computed fields
- **WHEN** the data source schema is initialized
- **THEN** it SHALL contain `zone_id` as required, `purge_quota` and `prefetch_quota` as computed, and `result_output_file` as optional

### Requirement: Data source calls DescribeContentQuota API
The system SHALL call `DescribeContentQuota` API from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` with `ZoneId` parameter set from `zone_id` schema field.

#### Scenario: Successful API call with zone_id
- **WHEN** data source read is invoked with `zone_id` set
- **THEN** the system SHALL construct a `DescribeContentQuotaRequest` with `ZoneId` set and call the API
- **THEN** the API call SHALL be wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with error handling using `tccommon.RetryError()`

#### Scenario: API call failure with retry
- **WHEN** the `DescribeContentQuota` API call fails
- **THEN** the error SHALL be wrapped with `tccommon.RetryError()` for retry logic

### Requirement: Data source maps API response to schema fields
The system SHALL map the `DescribeContentQuota` API response fields to Terraform schema fields:
- `response.Response.PurgeQuota` → `purge_quota`
- `response.Response.PrefetchQuota` → `prefetch_quota`

For each Quota element in the lists:
- `Batch` → `batch`
- `Daily` → `daily`
- `DailyAvailable` → `daily_available`
- `Type` → `type`

#### Scenario: Response with non-nil quota fields
- **WHEN** the API returns non-nil `PurgeQuota` and `PrefetchQuota`
- **THEN** the system SHALL flatten each Quota element into a map and set `purge_quota` and `prefetch_quota` in the state

#### Scenario: Response with nil quota fields
- **WHEN** the API returns nil `PurgeQuota` or `PrefetchQuota`
- **THEN** the system SHALL skip setting the corresponding field (nil check before set)

#### Scenario: Quota element with nil sub-fields
- **WHEN** a Quota element has nil `Batch`, `Daily`, `DailyAvailable`, or `Type`
- **THEN** the system SHALL skip setting that sub-field (nil check before set)

### Requirement: Data source sets unique ID
The system SHALL use `helper.BuildToken()` to generate a unique ID for the data source after a successful read.

#### Scenario: ID generation after successful read
- **WHEN** data source read completes successfully
- **THEN** `d.SetId(helper.BuildToken())` SHALL be called

### Requirement: Provider registration
The system SHALL register `tencentcloud_teo_content_quota` data source in `provider.go` and document it in `provider.md`.

#### Scenario: Data source registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_content_quota` SHALL be available as a data source key mapping to `teo.DataSourceTencentCloudTeoContentQuota()`

### Requirement: Unit test with gomonkey mock
The system SHALL provide unit tests using gomonkey to mock the cloud API `DescribeContentQuota`, covering successful read scenarios.

#### Scenario: Mock successful API response
- **WHEN** unit test runs with mocked `DescribeContentQuota` returning valid response
- **THEN** the test SHALL verify that schema fields are correctly populated from the response
