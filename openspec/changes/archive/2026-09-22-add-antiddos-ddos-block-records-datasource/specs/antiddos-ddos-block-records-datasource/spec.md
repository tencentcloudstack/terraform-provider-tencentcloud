## ADDED Requirements

### Requirement: Data source queries DDoS block records by time range
The `tencentcloud_antiddos_ddos_block_records` data source SHALL query AntiDDoS DDoS block/unblock records via the `DescribeDDoSBlockRecords` API (antiddos v20250903). The data source MUST accept required `start_time` and `end_time` parameters specifying the query time range, and SHALL pass them to the API request `StartTime` and `EndTime` fields.

#### Scenario: Query block records within a time range
- **WHEN** a user configures the data source with `start_time` and `end_time`
- **THEN** the data source calls `DescribeDDoSBlockRecords` with the given time range and returns the matching block records

#### Scenario: Missing required time parameters
- **WHEN** a user omits `start_time` or `end_time`
- **THEN** Terraform reports a validation error because both parameters are required

### Requirement: Data source supports filters for narrowing results
The data source SHALL accept an optional `filters` block list. Each filter item MUST contain a required `name` field and a required `values` set. The data source SHALL convert filters to the API request `Filters` field (`[]*v20250903.Filter`). Supported filter names include `Resource` (被封堵的 IP 或资源六段式) and `Status` (封堵状态: Blocked/Unblocking/Unblocked).

#### Scenario: Filter by status
- **WHEN** a user configures a filter with `name = "Status"` and `values = ["Blocked"]`
- **THEN** the data source returns only records whose status is Blocked

#### Scenario: Filter by resource IP
- **WHEN** a user configures a filter with `name = "Resource"` and `values = ["117.175.94.231"]`
- **THEN** the data source returns only records matching that resource

### Requirement: Data source returns block records list
The data source SHALL return a computed `block_records` list. Each element SHALL expose `resource` (被封堵资源的公网 IP), `block_time` (封堵时间), and `status` (封堵解封状态: Blocked/Unblocking/Unblocked) fields, mapped from `response.Response.BlockRecords[]` of type `DDoSBlockRecord`.

#### Scenario: Block records populated from API response
- **WHEN** the API returns non-empty `BlockRecords`
- **THEN** the data source populates `block_records` with one entry per record, each containing `resource`, `block_time`, and `status`

#### Scenario: Empty block records
- **WHEN** the API returns an empty `BlockRecords` list
- **THEN** the data source sets `block_records` to an empty list without error

### Requirement: Data source returns unblock quota info
The data source SHALL return a computed `unblock_quota_info` block exposing `total_quota`, `used_quota`, `quota_start_time`, and `quota_end_time` fields, mapped from `response.Response.UnblockQuotaInfo` of type `DDoSUnblockQuota`. The data source MUST handle the case where `UnblockQuotaInfo` is nil gracefully (no SetId clearing).

#### Scenario: Quota info populated
- **WHEN** the API returns a non-nil `UnblockQuotaInfo`
- **THEN** the data source populates `unblock_quota_info` with `total_quota`, `used_quota`, `quota_start_time`, and `quota_end_time`

#### Scenario: Quota info nil
- **WHEN** the API returns nil `UnblockQuotaInfo`
- **THEN** the data source leaves `unblock_quota_info` unset without error

### Requirement: Data source implements automatic pagination internally
The service layer method `DescribeAntiddosDDoSBlockRecordsByFilter` SHALL automatically paginate through all results using `Limit` (max 100, the API-documented maximum) and `Offset`, without exposing `limit`/`offset` to the Terraform user. The retry logic (using `tccommon.ReadRetryTimeout` and `tccommon.RetryError`) MUST be placed INSIDE the pagination loop.

#### Scenario: Large result set pagination
- **WHEN** more than 100 records match the query
- **THEN** the service layer fetches all pages and aggregates results into a single list

#### Scenario: API transient failure retried within page
- **WHEN** an API call fails transiently during a page fetch
- **THEN** the retry logic reattempts the same page before advancing the offset

### Requirement: Data source handles empty API response without clearing state
When the API Read returns an empty response (`response == nil` or `response.Response == nil` or `len(BlockRecords) == 0`), the data source MUST NOT call `d.SetId("")`. Instead it SHALL return a `NonRetryableError` from within the retry block so the outer retry continues, and on retry exhaustion SHALL log `[DATASOURCE] read empty, skip SetId`.

#### Scenario: API returns nil response
- **WHEN** the API returns a nil response or empty BlockRecords list during read
- **THEN** the data source returns a NonRetryableError from the retry block and does not clear the id

### Requirement: Data source supports result output file
The data source SHALL accept an optional `result_output_file` parameter. When provided, the data source SHALL write the query results to the specified file path in JSON format using `tccommon.WriteToFile`.

#### Scenario: Export results to file
- **WHEN** a user sets `result_output_file = "records.json"`
- **THEN** the data source writes the block records list to that file as JSON

### Requirement: Data source registered in provider
The provider SHALL register `tencentcloud_antiddos_ddos_block_records` in the `DataSourcesMap` of `tencentcloud/provider.go`, mapping to `antiddos.DataSourceTencentCloudAntiddosDDoSBlockRecords()`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_antiddos_ddos_block_records` is available as a data source

### Requirement: Connectivity layer exposes v20250903 antiddos client
The connectivity layer (`tencentcloud/connectivity/client.go`) SHALL provide a `UseAntiddosV20250903Client()` accessor returning `*v20250903.Client`, backed by a new `antiddosV20250903Conn` field. This is required because the existing `UseAntiddosClient()` returns the v20200309 client, which does not contain `DescribeDDoSBlockRecords`.

#### Scenario: Access v20250903 client
- **WHEN** service layer calls `me.client.UseAntiddosV20250903Client()`
- **THEN** a non-nil `*v20250903.Client` is returned and reused across calls