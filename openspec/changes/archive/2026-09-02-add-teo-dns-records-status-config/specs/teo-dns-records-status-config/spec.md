## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform RESOURCE_KIND_CONFIG resource `tencentcloud_teo_dns_records_status` with the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID，同时用于 Read（`DescribeDnsRecords` 的 `request.ZoneId`）和 Update（`ModifyDnsRecordsStatus` 的 `request.ZoneId`）
- `filters` (Optional, TypeList): 过滤条件列表，元素为 AdvancedFilter，对应 `request.Filters`，用于 Read 时按条件查询 DNS 记录
  - `name` (Required, TypeString): 需要过滤的字段，对应 `request.Filters.Name`
  - `values` (Required, TypeList of TypeString): 字段的过滤值，对应 `request.Filters.Values`
  - `fuzzy` (Optional, TypeBool): 是否启用模糊查询，对应 `request.Filters.Fuzzy`
- `sort_by` (Optional, TypeString): 排序依据，对应 `request.SortBy`
- `sort_order` (Optional, TypeString): 排序方式，对应 `request.SortOrder`
- `match` (Optional, TypeString): 匹配方式，对应 `request.Match`
- `records_to_enable` (Optional, TypeList of TypeString): 待启用的 DNS 记录 ID 列表，对应 `request.RecordsToEnable`，只管理单个资源，传入单个记录 ID
- `records_to_disable` (Optional, TypeList of TypeString): 待停用的 DNS 记录 ID 列表，对应 `request.RecordsToDisable`，只管理单个资源，传入单个记录 ID
- `dns_records` (Computed, TypeList): DNS 记录列表，对应 `response.Response.DnsRecords`，列表展开平铺，每个元素包含以下字段：
  - `zone_id` (Computed, TypeString): 站点 ID，对应 `response.Response.DnsRecords.ZoneId`
  - `record_id` (Computed, TypeString): DNS 记录 ID，对应 `response.Response.DnsRecords.RecordId`
  - `name` (Computed, TypeString): DNS 记录名，对应 `response.Response.DnsRecords.Name`
  - `type` (Computed, TypeString): DNS 记录类型，对应 `response.Response.DnsRecords.Type`
  - `location` (Computed, TypeString): DNS 记录解析线路，对应 `response.Response.DnsRecords.Location`
  - `content` (Computed, TypeString): DNS 记录内容，对应 `response.Response.DnsRecords.Content`
  - `ttl` (Computed, TypeInt): 缓存时间，对应 `response.Response.DnsRecords.TTL`
  - `weight` (Computed, TypeInt): DNS 记录权重，对应 `response.Response.DnsRecords.Weight`
  - `priority` (Computed, TypeInt): MX 记录优先级，对应 `response.Response.DnsRecords.Priority`
  - `status` (Computed, TypeString): DNS 记录解析状态，对应 `response.Response.DnsRecords.Status`
  - `created_on` (Computed, TypeString): 创建时间，对应 `response.Response.DnsRecords.CreatedOn`
  - `modified_on` (Computed, TypeString): 修改时间，对应 `response.Response.DnsRecords.ModifiedOn`

The resource ID SHALL be `zone_id`. The resource SHALL support import using the `zone_id` as ID. The resource SHALL NOT expose `limit`/`offset` pagination parameters.

#### Scenario: Schema defines all fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include zone_id, filters (with name/values/fuzzy sub-fields), sort_by, sort_order, match, records_to_enable, records_to_disable, and dns_records (with zone_id/record_id/name/type/location/content/ttl/weight/priority/status/created_on/modified_on sub-fields) with correct types and constraints

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** zone_id is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: Resource supports import
- **WHEN** a user imports the resource using `zone_id` as ID
- **THEN** the system SHALL populate the resource state by calling the Read operation

#### Scenario: Pagination parameters not exposed
- **WHEN** the resource schema is defined
- **THEN** it SHALL NOT include `limit` or `offset` fields

### Requirement: Resource Create Operation
The resource Create method SHALL reuse the Update logic because RESOURCE_KIND_CONFIG has no independent creation API. The Create method SHALL call `ModifyDnsRecordsStatus` API with `zone_id` → `request.ZoneId`, `records_to_enable` → `request.RecordsToEnable` (if set), `records_to_disable` → `request.RecordsToDisable` (if set). If both `records_to_enable` and `records_to_disable` are empty, the Create method SHALL skip the `ModifyDnsRecordsStatus` call and directly call Read. The Create method SHALL use `resource.Retry(tccommon.WriteRetryTimeout, ...)` for retry logic, wrapping errors with `tccommon.RetryError()`. After the API call succeeds, the Create method SHALL set the resource ID to `zone_id` and call Read to populate state. Setting the ID and other success operations SHALL be performed outside the retry block.

#### Scenario: Successful create with records_to_enable
- **WHEN** user applies a configuration with `records_to_enable` set to a single record ID
- **THEN** the resource SHALL call `ModifyDnsRecordsStatus` API with `request.ZoneId` and `request.RecordsToEnable` populated
- **AND** after success, the resource SHALL set the ID to `zone_id` and call Read to populate state

#### Scenario: Successful create with records_to_disable
- **WHEN** user applies a configuration with `records_to_disable` set to a single record ID
- **THEN** the resource SHALL call `ModifyDnsRecordsStatus` API with `request.ZoneId` and `request.RecordsToDisable` populated
- **AND** after success, the resource SHALL set the ID to `zone_id` and call Read to populate state

#### Scenario: Create with no records_to_enable and no records_to_disable
- **WHEN** user applies a configuration with both `records_to_enable` and `records_to_disable` empty
- **THEN** the resource SHALL skip the `ModifyDnsRecordsStatus` call
- **AND** set the ID to `zone_id` and call Read to populate state

#### Scenario: Create API returns error
- **WHEN** `ModifyDnsRecordsStatus` API returns a retryable error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`
- **AND** on non-retryable errors, the resource SHALL return the error directly

### Requirement: Resource Read Operation
The resource Read method SHALL call `DescribeDnsRecords` API with the following parameter mapping:
- `zone_id` → `request.ZoneId`
- `filters` → `request.Filters` (each filter element: `name` → `AdvancedFilter.Name`, `values` → `AdvancedFilter.Values`, `fuzzy` → `AdvancedFilter.Fuzzy`)
- `sort_by` → `request.SortBy` (if set)
- `sort_order` → `request.SortOrder` (if set)
- `match` → `request.Match` (if set)

The Read method SHALL use `resource.Retry(tccommon.ReadRetryTimeout, ...)` for retry logic, wrapping errors with `tccommon.RetryError()`. The Read method SHALL populate `dns_records` from `response.Response.DnsRecords`, flattening the list elements. Before each `d.Set` call, the Read method SHALL check that the Response field is not nil; if nil, the Read method SHALL NOT call `d.Set` for that field. If the cloud API returns empty (response is nil, response.Response is nil, or len(response.Response.DnsRecords) == 0), the Read method SHALL first log `log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())` to preserve context, then call `d.SetId("")` to mark the resource as deleted. The Read method SHALL NOT perform retry inside an existing retry block.

#### Scenario: Successful read with filters
- **WHEN** the Read method queries DNS records with filters set
- **THEN** the resource SHALL call `DescribeDnsRecords` with `request.Filters` populated from the filters schema
- **AND** populate `dns_records` from `response.Response.DnsRecords`

#### Scenario: Read with nil DnsRecords fields
- **WHEN** the Read method queries DNS records where some `DnsRecords` element fields are nil
- **THEN** the resource SHALL check each field is not nil before calling `d.Set`
- **AND** no error SHALL be returned for nil fields

#### Scenario: Resource not found
- **WHEN** `DescribeDnsRecords` returns empty response or `DnsRecords` list is empty
- **THEN** the resource SHALL log `log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())` before clearing the ID
- **AND** call `d.SetId("")` to mark the resource as deleted

#### Scenario: Read API returns error
- **WHEN** `DescribeDnsRecords` API returns a retryable error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.ReadRetryTimeout`

### Requirement: Resource Update Operation
The resource Update method SHALL call `ModifyDnsRecordsStatus` API when `records_to_enable` or `records_to_disable` changes, with the following parameter mapping:
- `zone_id` → `request.ZoneId`
- `records_to_enable` → `request.RecordsToEnable` (if set)
- `records_to_disable` → `request.RecordsToDisable` (if set)

The Update method SHALL use `resource.Retry(tccommon.WriteRetryTimeout, ...)` for retry logic, wrapping errors with `tccommon.RetryError()`. Setting the ID and other success operations SHALL be performed outside the retry block. Because `ModifyDnsRecordsStatus` is an asynchronous API, after the API call succeeds, the Update method SHALL call `DescribeDnsRecords` (via the Read operation) to poll until the target record's `status` reaches the expected value (`enable` or `disable`) or the update timeout is reached. The Update method SHALL end with a call to `resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)` to refresh state.

#### Scenario: Update records_to_enable
- **WHEN** `records_to_enable` changes in the Terraform configuration
- **THEN** the resource SHALL call `ModifyDnsRecordsStatus` with `request.ZoneId` and `request.RecordsToEnable` populated
- **AND** poll `DescribeDnsRecords` until the target record's `status` reaches `enable` or the update timeout is reached

#### Scenario: Update records_to_disable
- **WHEN** `records_to_disable` changes in the Terraform configuration
- **THEN** the resource SHALL call `ModifyDnsRecordsStatus` with `request.ZoneId` and `request.RecordsToDisable` populated
- **AND** poll `DescribeDnsRecords` until the target record's `status` reaches `disable` or the update timeout is reached

#### Scenario: Update API returns error
- **WHEN** `ModifyDnsRecordsStatus` API returns a retryable error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`
- **AND** on non-retryable errors, the resource SHALL return the error directly

#### Scenario: Update finalizes with Read
- **WHEN** the Update operation completes
- **THEN** the Update method SHALL call `resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)` to refresh state

### Requirement: Resource Delete Operation
The resource Delete method SHALL be a no-op (the CONFIG resource does not support cloud-side deletion via Terraform). The Delete method SHALL return nil without calling any cloud API and without resetting DNS record status.

#### Scenario: Delete is a no-op
- **WHEN** user destroys the `tencentcloud_teo_dns_records_status` resource
- **THEN** the Delete method SHALL return nil without calling any cloud API

### Requirement: Provider Registration
The resource `tencentcloud_teo_dns_records_status` SHALL be registered in `tencentcloud/provider.go` mapping `"tencentcloud_teo_dns_records_status"` to `teo.ResourceTencentCloudTeoDnsRecordsStatus()`. The resource SHALL also be registered in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider
- **WHEN** user references `tencentcloud_teo_dns_records_status` in their configuration
- **THEN** Terraform SHALL recognize it as a valid resource type

### Requirement: Unit Tests with gomonkey mock
The resource SHALL have unit tests in `resource_tc_teo_dns_records_status_config_test.go` that use gomonkey to mock the cloud API calls (`DescribeDnsRecordsWithContext` and `ModifyDnsRecordsStatusWithContext`), covering Create, Read, and Update operations. The tests SHALL NOT use the Terraform test suite; instead they SHALL use gomonkey mock for business logic unit testing.

#### Scenario: Unit test for Create with records_to_enable
- **WHEN** the unit test for Create with `records_to_enable` set is executed
- **THEN** it SHALL mock `ModifyDnsRecordsStatusWithContext` to return success
- **AND** mock `DescribeDnsRecordsWithContext` to return the DNS records
- **AND** verify the resource is created correctly with `dns_records` populated in state

#### Scenario: Unit test for Create with records_to_disable
- **WHEN** the unit test for Create with `records_to_disable` set is executed
- **THEN** it SHALL mock `ModifyDnsRecordsStatusWithContext` to return success
- **AND** mock `DescribeDnsRecordsWithContext` to return the DNS records
- **AND** verify the resource is created correctly

#### Scenario: Unit test for Read
- **WHEN** the unit test for Read is executed
- **THEN** it SHALL mock `DescribeDnsRecordsWithContext` to return the DNS records
- **AND** verify the resource state is populated with `dns_records` fields

#### Scenario: Unit test for Update
- **WHEN** the unit test for Update with `records_to_enable` change is executed
- **THEN** it SHALL mock `ModifyDnsRecordsStatusWithContext` to return success
- **AND** mock `DescribeDnsRecordsWithContext` to return the updated DNS records
- **AND** verify the update operation completes and state is refreshed

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_dns_records_status.md` with a one-line description mentioning TEO (EdgeOne), example usage, and import section. The documentation SHALL NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated by tooling).

#### Scenario: Documentation exists with required sections
- **WHEN** the resource documentation is created
- **THEN** the .md file SHALL exist with a one-line description mentioning TEO (EdgeOne)
- **AND** example usage SHALL demonstrate `zone_id`, `filters`, `records_to_enable`/`records_to_disable` fields
- **AND** an import section SHALL be present showing the `zone_id` ID format
