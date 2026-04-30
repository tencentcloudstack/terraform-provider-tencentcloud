## ADDED Requirements

### Requirement: Data source schema for content quota
The data source `tencentcloud_teo_content_quota` SHALL provide a schema with the following fields:
- `zone_id` (Required, string): The site ID to query content quota for
- `purge_quota` (Computed, list of objects): Cache purge quota list, each element containing:
  - `type` (Computed, string): Quota type (e.g., purge_prefix, purge_url, purge_host, purge_all, purge_cache_tag)
  - `batch` (Computed, int): Single batch submission quota limit
  - `daily` (Computed, int): Daily submission quota limit
  - `daily_available` (Computed, int): Daily remaining available quota
- `prefetch_quota` (Computed, list of objects): Cache prefetch quota list, each element containing:
  - `type` (Computed, string): Quota type (e.g., prefetch_url)
  - `batch` (Computed, int): Single batch submission quota limit
  - `daily` (Computed, int): Daily submission quota limit
  - `daily_available` (Computed, int): Daily remaining available quota
- `result_output_file` (Optional, string): Used to save results

#### Scenario: Query content quota with zone_id
- **WHEN** user provides a valid `zone_id` in the data source configuration
- **THEN** the data source SHALL call `DescribeContentQuota` API with the provided `zone_id` and return the `purge_quota` and `prefetch_quota` lists

#### Scenario: API returns null for quota fields
- **WHEN** the API response contains null for `PurgeQuota` or `PrefetchQuota`
- **THEN** the corresponding schema field SHALL not be set (remain empty), and no error SHALL be returned

### Requirement: Service layer API call with retry
The service layer method `DescribeTeoContentQuotaByFilter` SHALL wrap the `DescribeContentQuota` API call with `resource.Retry(tccommon.ReadRetryTimeout, ...)` for transient error handling. On API failure, the error SHALL be wrapped with `tccommon.RetryError()` and returned.

#### Scenario: Successful API call
- **WHEN** the `DescribeContentQuota` API call succeeds
- **THEN** the method SHALL return the `PurgeQuota` and `PrefetchQuota` lists from the response

#### Scenario: Transient API failure
- **WHEN** the `DescribeContentQuota` API call fails with a transient error
- **THEN** the method SHALL retry up to `ReadRetryTimeout` duration before returning the error

### Requirement: Data source ID generation
The data source SHALL use `helper.DataResourceIdsHash([]string{zoneId})` to generate a deterministic ID based on the `zone_id` input.

#### Scenario: ID generation from zone_id
- **WHEN** the data source read completes successfully
- **THEN** the resource ID SHALL be a hash of the `zone_id` value

### Requirement: Provider registration
The data source SHALL be registered in `provider.go` with key `tencentcloud_teo_content_quota` mapped to `teo.DataSourceTencentCloudTeoContentQuota()`, and a corresponding entry SHALL be added to `provider.md`.

#### Scenario: Data source available in provider
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_teo_content_quota` data source SHALL be available for use in Terraform configurations

### Requirement: Documentation
A markdown documentation file `data_source_tc_teo_content_quota.md` SHALL be created following the project documentation guidelines, including:
- A one-line description mentioning TEO product name
- Example Usage section with `jsonencode()` if JSON fields are involved
- No `Argument Reference` or `Attribute Reference` sections (auto-generated)

#### Scenario: Documentation file exists
- **WHEN** the data source implementation is complete
- **THEN** a corresponding `.md` documentation file SHALL exist in the teo service directory

### Requirement: Unit test with mock
The test file `data_source_tc_teo_content_quota_test.go` SHALL use gomonkey mock approach (not Terraform test suite) to test the business logic of the data source, and the test SHALL pass with `go test -gcflags=all=-l`.

#### Scenario: Unit test passes
- **WHEN** running `go test -gcflags=all=-l` on the test file
- **THEN** all test cases SHALL pass, verifying the data source read logic with mocked API responses
