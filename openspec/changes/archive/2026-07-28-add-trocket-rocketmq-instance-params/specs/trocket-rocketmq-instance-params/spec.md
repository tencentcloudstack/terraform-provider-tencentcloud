## ADDED Requirements

### Requirement: TROcket Rocketmq Instance Billing and Deployment Schema

The system SHALL define the following additional Optional fields on the `tencentcloud_trocket_rocketmq_instance` resource Schema:
- `pay_mode` (Optional, TypeInt): 付费模式（0: 后付费；1: 预付费），默认值为 0。
- `renew_flag` (Optional, TypeInt): 预付费集群是否自动续费（0: 不自动续费；1: 自动续费），默认值为 0。
- `time_span` (Optional, TypeInt): 预付费集群的购买时长（单位：月），取值范围 1～60，默认值为 1。
- `max_topic_num` (Optional, TypeInt): 最大可创建主题数。
- `zone_ids` (Optional, TypeList with Elem TypeInt): 部署可用区列表。

All new fields SHALL be Optional and backward compatible; no existing Schema field SHALL be modified or removed.

#### Scenario: New fields are optional
- **WHEN** a user creates a `tencentcloud_trocket_rocketmq_instance` without specifying `pay_mode`, `renew_flag`, `time_span`, `max_topic_num`, or `zone_ids`
- **THEN** the system SHALL accept the configuration and rely on the cloud API defaults

#### Scenario: zone_ids accepts a list of integers
- **WHEN** a user sets `zone_ids = [100001, 100003]` in the Terraform configuration
- **THEN** the schema SHALL accept the list of integers and mark each element as TypeInt

### Requirement: TROcket Rocketmq Instance Create Operation

The system SHALL, in the `Create` operation, fill the following `CreateInstanceRequest` fields when the corresponding Terraform fields are set:
- `request.PayMode` from `pay_mode`
- `request.RenewFlag` from `renew_flag`
- `request.TimeSpan` from `time_span`
- `request.MaxTopicNum` from `max_topic_num`
- `request.ZoneIds` ([]*int64) from `zone_ids` list

#### Scenario: Create with prepaid billing fields
- **WHEN** the user sets `pay_mode = 1`, `renew_flag = 1`, `time_span = 12` in the Terraform configuration and creates the instance
- **THEN** the system SHALL fill `request.PayMode = 1`, `request.RenewFlag = 1`, `request.TimeSpan = 12` in the `CreateInstance` request

#### Scenario: Create without billing fields configured
- **WHEN** the user does not set `pay_mode`, `renew_flag`, or `time_span` in the Terraform configuration
- **THEN** the system SHALL NOT fill these fields in the `CreateInstance` request, preserving the cloud-side default

#### Scenario: Create with max_topic_num configured
- **WHEN** the user sets `max_topic_num = 500` in the Terraform configuration and creates the instance
- **THEN** the system SHALL fill `request.MaxTopicNum = 500` in the `CreateInstance` request

#### Scenario: Create with zone_ids configured
- **WHEN** the user sets `zone_ids = [100001, 100003]` in the Terraform configuration and creates the instance
- **THEN** the system SHALL convert each integer to `*int64` and fill `request.ZoneIds` in the `CreateInstance` request

### Requirement: TROcket Rocketmq Instance Read Operation

The system SHALL, in the `Read` operation, read the following fields from the `DescribeInstance` response (`rocketmqInstance` of type `Instance`):
- `pay_mode`: SHALL be mapped from `rocketmqInstance.PayMode` (string enum `POSTPAID`/`PREPAID`) to int (`POSTPAID`→0, `PREPAID`→1). Before calling `d.Set`, the system SHALL check that `rocketmqInstance.PayMode` is not nil; if nil, skip the set.
- `renew_flag`: SHALL be set from `rocketmqInstance.RenewFlag` (int64). Before calling `d.Set`, check non-nil; if nil, skip.
- `zone_ids`: SHALL be set from `rocketmqInstance.ZoneIds` ([]*int64). Before calling `d.Set`, check non-nil; if nil, skip.

The system SHALL NOT attempt to read `time_span` or `max_topic_num` from the `DescribeInstance` response (no corresponding field).

#### Scenario: Read pay_mode from DescribeInstance
- **WHEN** the Read operation calls `DescribeInstance` and the response `PayMode` is `"PREPAID"`
- **THEN** the system SHALL set `pay_mode = 1` in Terraform state

#### Scenario: Read pay_mode POSTPAID
- **WHEN** the response `PayMode` is `"POSTPAID"`
- **THEN** the system SHALL set `pay_mode = 0` in Terraform state

#### Scenario: Read with nil response fields
- **WHEN** the response `PayMode`, `RenewFlag`, or `ZoneIds` is nil
- **THEN** the system SHALL skip the corresponding `d.Set` call to avoid nil pointer issues

### Requirement: TROcket Rocketmq Instance Update Operation

The system SHALL treat `pay_mode`, `renew_flag`, and `time_span` as immutable by adding them to the `immutableArgs` array in the Update operation. When any of these fields change, the system SHALL return an error `argument <name> cannot be changed`.

The system SHALL support in-place update of `max_topic_num` and `zone_ids`:
- When `max_topic_num` changes, the system SHALL fill `request1.MaxTopicNum` (int64) and call `ModifyInstance`.
- When `zone_ids` changes, the system SHALL convert each integer to `*string` and fill `request1.ZoneIds` ([]*string), then call `ModifyInstance`.

#### Scenario: Update pay_mode is rejected
- **WHEN** `pay_mode` changes in the Terraform configuration
- **THEN** the system SHALL return an error and NOT call any API

#### Scenario: Update renew_flag is rejected
- **WHEN** `renew_flag` changes in the Terraform configuration
- **THEN** the system SHALL return an error and NOT call any API

#### Scenario: Update time_span is rejected
- **WHEN** `time_span` changes in the Terraform configuration
- **THEN** the system SHALL return an error and NOT call any API

#### Scenario: Update max_topic_num in-place
- **WHEN** `max_topic_num` changes from 500 to 1000 in the Terraform configuration
- **THEN** the system SHALL call `ModifyInstance` with `MaxTopicNum = 1000` and the update SHALL be performed in-place without resource recreation

#### Scenario: Update zone_ids in-place
- **WHEN** `zone_ids` changes in the Terraform configuration
- **THEN** the system SHALL call `ModifyInstance` with `ZoneIds` ([]*string) converted from the new list and the update SHALL be performed in-place

### Requirement: TROcket Rocketmq Instance Unit Tests

The system SHALL provide unit tests in `resource_tc_trocket_rocketmq_instance_test.go` using the Terraform test suite (TF_ACC), covering the new parameters `pay_mode`, `renew_flag`, `time_span`, `max_topic_num`, and `zone_ids` in create and update scenarios.

#### Scenario: Test covers new parameters
- **WHEN** the test suite is run
- **THEN** test cases SHALL cover creating an instance with the new fields and updating `max_topic_num`/`zone_ids`

### Requirement: TROcket Rocketmq Instance Resource Documentation

The system SHALL update the markdown documentation file `resource_tc_trocket_rocketmq_instance.md` to include example usage demonstrating the new `pay_mode`, `renew_flag`, `time_span`, `max_topic_num`, and `zone_ids` fields. The documentation SHALL NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated).

#### Scenario: Documentation updated
- **WHEN** the resource parameters are added
- **THEN** the `.md` file SHALL include example HCL demonstrating the new fields
