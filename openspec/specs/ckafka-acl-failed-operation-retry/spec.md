# ckafka-acl-failed-operation-retry Specification

## Purpose
TBD - created by archiving change add-ckafka-acl-failed-operation-retry. Update Purpose after archive.
## Requirements
### Requirement: CreateAcl 对 FailedOperation 错误进行固定次数重试

The system SHALL, in the `CkafkaService.CreateAcl` method (`tencentcloud/services/ckafka/service_tencentcloud_ckafka.go`), retry the `CreateAcl` cloud API call when it returns a `TencentCloudSDKError` with `Code == "FailedOperation"`. The retry SHALL be performed **up to 3 times** after the initial failed attempt, with a fixed **5-second** interval between each attempt. Retry-related parameters (retry times = 3, interval = 5s) SHALL be defined as named constants in `tencentcloud/services/ckafka/extension_ckafka.go` so they can be referenced by both implementation and tests.

#### Scenario: FailedOperation persists across all retries
- **WHEN** the `CreateAcl` API returns `Code == "FailedOperation"` on the initial call and on all 3 subsequent retries
- **THEN** the system SHALL make a total of 4 attempts (1 initial + 3 retries), sleeping 5 seconds between attempts
- **AND** the system SHALL return the final `FailedOperation` error to the caller

#### Scenario: FailedOperation occurs transiently then succeeds
- **WHEN** the `CreateAcl` API returns `Code == "FailedOperation"` on the initial call (or on some retries) and then succeeds on a later attempt within the 3-retry budget
- **THEN** the system SHALL return nil error and continue with the original success-path logic (OperateStatusCheck), without surfacing any error

#### Scenario: Non-FailedOperation error is returned
- **WHEN** the `CreateAcl` API returns an error whose code is NOT `FailedOperation` (e.g., `InvalidParameter`, `UnauthorizedOperation`)
- **THEN** the system SHALL NOT apply the fixed-count retry for that error
- **AND** the error SHALL be handled by the existing `resource.Retry` + `tccommon.RetryError` logic (other retryable codes like network errors / rate limits still retried per original behavior)

#### Scenario: CreateAcl succeeds on the first attempt
- **WHEN** the `CreateAcl` API succeeds on the first call
- **THEN** the system SHALL NOT sleep or retry
- **AND** the behavior SHALL be identical to the current implementation (no performance regression on the happy path)

### Requirement: 重试行为单元测试

The system SHALL provide unit tests for the `CreateAcl` FailedOperation retry behavior in `tencentcloud/services/ckafka/resource_tc_ckafka_acl_test.go`, using gomonkey to mock the `CreateAcl` cloud API. The tests SHALL verify the retry count, the 5-second interval, and both the eventual-success and eventual-failure paths.

#### Scenario: Mocked API returns FailedOperation 3 times then success
- **WHEN** a unit test mocks `CreateAcl` to return `Code == "FailedOperation"` for the first 3 calls and success on the 4th call
- **THEN** the test SHALL assert that `CkafkaService.CreateAcl` returns no error
- **AND** the test SHALL assert that the mocked API was invoked exactly 4 times

#### Scenario: Mocked API always returns FailedOperation
- **WHEN** a unit test mocks `CreateAcl` to always return `Code == "FailedOperation"`
- **THEN** the test SHALL assert that `CkafkaService.CreateAcl` returns the `FailedOperation` error
- **AND** the test SHALL assert that the mocked API was invoked exactly 4 times (1 initial + 3 retries)

#### Scenario: Mocked API returns a non-FailedOperation error
- **WHEN** a unit test mocks `CreateAcl` to return an error whose code is NOT `FailedOperation`
- **THEN** the test SHALL assert that `CkafkaService.CreateAcl` returns the error immediately
- **AND** the test SHALL assert that the mocked API was invoked only once (no fixed-count retry)

