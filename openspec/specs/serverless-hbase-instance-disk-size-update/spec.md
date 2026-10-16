# serverless-hbase-instance-disk-size-update Specification

## Purpose
TBD - created by archiving change modify-serverless-hbase-disk-size. Update Purpose after archive.
## Requirements
### Requirement: disk_size is updatable in place

The `tencentcloud_serverless_hbase_instance` resource SHALL allow the `disk_size` field to be changed on an existing instance without forcing replacement. The `Update` function SHALL NOT include `disk_size` in its immutable-arguments check.

#### Scenario: Increasing disk_size triggers ModifySLInstance

- **WHEN** the user changes only `disk_size` to a larger value (a multiple of 100 GB) on an existing `tencentcloud_serverless_hbase_instance`
- **THEN** the `Update` function constructs a `ModifySLInstanceRequest` with `InstanceId` set from `d.Id()`, `Zone` set from `zone_settings.0.zone` (required by the API) and `DiskSize` set to the new value, and calls `ModifySLInstanceWithContext`
- **AND** if `zone_settings.0.zone` is empty, `Update` returns an error before calling the API
- **AND** the call is wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with errors wrapped via `tccommon.RetryError`

#### Scenario: disk_size change is no longer rejected as immutable

- **WHEN** the user changes `disk_size` and applies
- **THEN** the resource does not return the `argument 'disk_size' cannot be changed` error
- **AND** the update proceeds to call the cloud API

### Requirement: disk_size update waits for async scaling to complete

Because `ModifySLInstance` is asynchronous, the `Update` function SHALL wait for the instance to finish scaling before returning. The instance status `22` (变配中) SHALL be treated as a pending state, and status `2` (运行中, 变配完成) SHALL be the success target. Status polling SHALL use the existing `SLInstanceStateRefreshFunc` and `tccommon.BuildStateChangeConf`.

#### Scenario: Polls until status transitions from 22 to 2

- **WHEN** `ModifySLInstance` for `disk_size` returns success
- **THEN** the `Update` function builds a state change conf with pending `["22"]` and target `["2"]` using `10 * tccommon.ReadRetryTimeout` and `time.Second` min interval
- **AND** calls `WaitForState` to block until the instance reaches status `2`
- **AND** if the wait fails, the error is returned to the caller

#### Scenario: Wait happens outside the API retry block

- **WHEN** the `ModifySLInstance` retry block completes successfully
- **THEN** the status-change wait is executed after the retry block
- **AND** no retry or state mutation is performed inside the retry block

### Requirement: disk_size update is independent of node_num expansion

The `disk_size` update branch SHALL be independent from the existing `zone_settings.node_num` expansion logic. Changing `disk_size` SHALL only read `Zone` from `zone_settings.0.zone` and SHALL NOT modify `NodeNum`, and changing `node_num` SHALL NOT modify `DiskSize`. Both changes MAY occur in the same `terraform apply` and SHALL be applied sequentially without conflict.

#### Scenario: Only disk_size changed, zone_settings unchanged

- **WHEN** the user changes `disk_size` but does not change `zone_settings`
- **THEN** the `disk_size` branch calls `ModifySLInstance` with `InstanceId`, `Zone` (from `zone_settings.0.zone`) and `DiskSize` (no `NodeNum`)
- **AND** the `node_num` loop is not entered

#### Scenario: Both disk_size and node_num changed together

- **WHEN** the user changes both `disk_size` and `zone_settings.0.node_num` in one apply
- **THEN** the `disk_size` branch runs `ModifySLInstance` and waits for status `2`
- **AND** subsequently the `node_num` loop runs `ModifySLInstance` per changed zone and waits for status `2`
- **AND** neither call corrupts the other's request parameters

### Requirement: disk_size update errors are logged with resource name

Any failure in the `disk_size` update path SHALL be logged using `log.Printf` with the `logId` and the resource name `serverless hbase instance`, then returned to the caller.

#### Scenario: ModifySLInstance failure is logged

- **WHEN** the `ModifySLInstance` call for `disk_size` fails after retries are exhausted
- **THEN** a log line is emitted in the form `[CRITAL]%s update serverless hbase instance disk_size failed, reason:%+v` with the `logId` and the error
- **AND** the error is returned from `Update`