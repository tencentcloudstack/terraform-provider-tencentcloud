## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform RESOURCE_KIND_CONFIG resource `tencentcloud_teo_dns_records_status` with the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID，同时用于 Read（`DescribeDnsRecords` 的 `request.ZoneId`）和 Update（`ModifyDnsRecordsStatus` 的 `request.ZoneId`）
- `records_id` (Required, ForceNew, TypeString): DNS 记录 ID，与 `zone_id` 组合作为资源唯一 ID（`zone_id#records_id`）
- `status` (Required, TypeString): DNS 记录状态，枚举值 `enable`（已生效）/`disable`（已停用），通过该字段控制记录的启用与停用

The resource ID SHALL be `zone_id#records_id` (using `tccommon.FILED_SP` as separator). The resource SHALL support import using the `zone_id#records_id` composite ID. The resource SHALL NOT expose `filters`, `sort_by`, `sort_order`, `match` (Describe query parameters), `dns_records` (response data), `limit`/`offset` pagination parameters, or `records_to_enable`/`records_to_disable`.

#### Scenario: Schema defines only required fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include only zone_id, records_id, status with correct types and constraints
- **AND** it SHALL NOT include filters, sort_by, sort_order, match, dns_records, limit, offset, records_to_enable, or records_to_disable

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** zone_id or records_id is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: status is updatable
- **WHEN** status is changed in the Terraform configuration
- **THEN** the resource SHALL update in-place by calling ModifyDnsRecordsStatus

#### Scenario: Resource supports import
- **WHEN** a user imports the resource using `zone_id#records_id` as composite ID
- **THEN** the system SHALL populate the resource state by calling the Read operation

### Requirement: Resource Create Operation
The resource Create method SHALL reuse the Update logic because RESOURCE_KIND_CONFIG has no independent creation API. The Create method SHALL set the resource ID to `zone_id + FILED_SP + records_id` via `d.SetId()` and then directly return `resourceTencentCloudTeoDnsRecordsStatusUpdate(d, meta)`. The Create method SHALL NOT independently construct the `ModifyDnsRecordsStatus` request or call the API; all API interaction is delegated to the Update method.

#### Scenario: Successful create delegates to update
- **WHEN** user applies a configuration with `status` set
- **THEN** the resource SHALL set the ID to `zone_id#records_id`
- **AND** delegate to the Update method which calls `ModifyDnsRecordsStatus` API when `d.HasChange("status")` is true
- **AND** poll `DescribeDnsRecords` until status reaches the target value
- **AND** finally call Read to populate state

#### Scenario: Create API returns error
- **WHEN** `ModifyDnsRecordsStatus` API returns a retryable error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`
- **AND** on non-retryable errors, the resource SHALL return the error directly

### Requirement: Resource Read Operation
The resource Read method SHALL call `DescribeDnsRecords` API with `zone_id` → `request.ZoneId` and `Filters` (Name="id", Values=[records_id]) to query the specific DNS record. The Read method SHALL use `resource.Retry(tccommon.ReadRetryTimeout, ...)` for retry logic, wrapping errors with `tccommon.RetryError()`. If the cloud API returns empty (response is nil, response.Response is nil, or len(response.Response.DnsRecords) == 0), the Read method SHALL first log `log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())` to preserve context, then call `d.SetId("")` to mark the resource as deleted. The Read method SHALL backfill the `status` field from `DnsRecord.Status` (if non-nil). The Read method SHALL NOT perform retry inside an existing retry block.

#### Scenario: Successful read
- **WHEN** the Read method queries DNS records with zone_id and records_id filter
- **THEN** the resource SHALL call `DescribeDnsRecords` with `request.ZoneId` and `request.Filters` populated
- **AND** if the record exists, the resource SHALL remain in state and backfill `status`

#### Scenario: Resource not found
- **WHEN** `DescribeDnsRecords` returns empty response or `DnsRecords` list is empty
- **THEN** the resource SHALL log `log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())` before clearing the ID
- **AND** call `d.SetId("")` to mark the resource as deleted

#### Scenario: Read API returns error
- **WHEN** `DescribeDnsRecords` API returns a retryable error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.ReadRetryTimeout`

### Requirement: Resource Update Operation
The resource Update method SHALL call `ModifyDnsRecordsStatus` API when `status` changes (`d.HasChange("status")`), with the following parameter mapping:
- `zone_id` → `request.ZoneId`
- when `status` is `enable`: `records_id` → `request.RecordsToEnable`
- when `status` is `disable`: `records_id` → `request.RecordsToDisable`

The Update method SHALL use `resource.Retry(tccommon.WriteRetryTimeout, ...)` for retry logic, wrapping errors with `tccommon.RetryError()`. Setting the ID and other success operations SHALL be performed outside the retry block. Because `ModifyDnsRecordsStatus` is an async interface, the Update method SHALL poll `DescribeDnsRecords` using `resource.StateChangeConf` until the record `status` reaches the target value (`enable` or `disable`). The Update method SHALL end with a call to `resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)` to refresh state.

#### Scenario: Update status to enable
- **WHEN** `status` changes to `enable` in the Terraform configuration
- **THEN** the resource SHALL call `ModifyDnsRecordsStatus` with `request.ZoneId` and `request.RecordsToEnable` populated
- **AND** poll `DescribeDnsRecords` until record status becomes `enable`
- **AND** end with a call to Read to refresh state

#### Scenario: Update status to disable
- **WHEN** `status` changes to `disable` in the Terraform configuration
- **THEN** the resource SHALL call `ModifyDnsRecordsStatus` with `request.ZoneId` and `request.RecordsToDisable` populated
- **AND** poll `DescribeDnsRecords` until record status becomes `disable`
- **AND** end with a call to Read to refresh state

#### Scenario: Update API returns error
- **WHEN** `ModifyDnsRecordsStatus` API returns a retryable error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`
- **AND** on non-retryable errors, the resource SHALL return the error directly

#### Scenario: Update polls async status
- **WHEN** `ModifyDnsRecordsStatus` returns success but the record status has not yet reached the target value
- **THEN** the resource SHALL poll `DescribeDnsRecords` until the record `status` reaches the target value within `tccommon.ReadRetryTimeout`

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

#### Scenario: Unit test for Create with status enable
- **WHEN** the unit test for Create with `status=enable` is executed
- **THEN** it SHALL mock `ModifyDnsRecordsStatusWithContext` to return success
- **AND** mock `DescribeDnsRecordsWithContext` to return the DNS record with status `enable`
- **AND** verify that `ModifyDnsRecordsStatusWithContext` is called with `RecordsToEnable` populated
- **AND** verify the resource is created correctly with the `zone_id#records_id` as ID

#### Scenario: Unit test for Create with status disable
- **WHEN** the unit test for Create with `status=disable` is executed
- **THEN** it SHALL mock `ModifyDnsRecordsStatusWithContext` to return success
- **AND** mock `DescribeDnsRecordsWithContext` to return the DNS record with status `disable`
- **AND** verify that `ModifyDnsRecordsStatusWithContext` is called with `RecordsToDisable` populated
- **AND** verify the resource is created correctly with the `zone_id#records_id` as ID

#### Scenario: Unit test for Read
- **WHEN** the unit test for Read is executed
- **THEN** it SHALL mock `DescribeDnsRecordsWithContext` to return the DNS record
- **AND** verify the resource remains in state and `status` is backfilled

#### Scenario: Unit test for Update
- **WHEN** the unit test for Update with `status` change is executed
- **THEN** it SHALL mock `ModifyDnsRecordsStatusWithContext` to return success
- **AND** mock `DescribeDnsRecordsWithContext` to return the updated DNS record
- **AND** verify that `ModifyDnsRecordsStatusWithContext` is called
- **AND** verify the update operation completes

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_dns_records_status.md` with a one-line description mentioning TEO (EdgeOne), example usage, and import section. The documentation SHALL NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated by tooling).

#### Scenario: Documentation exists with required sections
- **WHEN** the resource documentation is created
- **THEN** the .md file SHALL exist with a one-line description mentioning TEO (EdgeOne)
- **AND** example usage SHALL demonstrate `zone_id`, `records_id`, `status` fields
- **AND** an import section SHALL be present showing the `zone_id#records_id` composite ID format
