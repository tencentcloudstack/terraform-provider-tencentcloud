# teo-apply-free-certificate-action Specification

## Purpose
TBD - created by archiving change add-teo-apply-free-certificate. Update Purpose after archive.
## Requirement: Action Resource Schema Definition
The system SHALL define a Terraform Plugin Framework action `tencentcloud_teo_apply_free_certificate` with the following schema attributes:
- `zone_id` (Required, String): 站点 ID。
- `domain` (Required, String): 申请免费证书的目标域名。
- `verification_method` (Required, String): 申请免费证书的验证方式，取值 `http_challenge`（HTTP 访问文件验证）/ `dns_challenge`（DNS 委派验证）。

#### Scenario: Schema defines required input attributes
- **WHEN** the action schema is defined
- **THEN** it SHALL include `zone_id` as a required string attribute, `domain` as a required string attribute, and `verification_method` as a required string attribute

#### Scenario: Required attributes enforced
- **WHEN** a user omits `zone_id`, `domain`, or `verification_method` in the action config
- **THEN** the framework SHALL report a validation error before the action is invoked

## Requirement: Action Invoke Operation
The system SHALL implement the action `Invoke` handler to call the TEO `ApplyFreeCertificate` API via `ApplyFreeCertificateWithContext`, passing `ZoneId`, `Domain`, and `VerificationMethod` from the action config. After the apply call succeeds, the system SHALL poll `CheckFreeCertificateVerification` via `CheckFreeCertificateVerificationWithContext` (passing `ZoneId` and `Domain`) until the certificate application succeeds (response `CommonName` is non-empty), fails, or times out. The verification information (`DnsVerification` / `FileVerification`) returned by `ApplyFreeCertificate` SHALL be logged for observability. No cloud-side state SHALL be persisted (no id, no output attributes written).

#### Scenario: Successful apply and verification invocation
- **WHEN** the action is invoked with valid `zone_id`, `domain`, and `verification_method`
- **THEN** the system SHALL construct an `ApplyFreeCertificateRequest` with `ZoneId`, `Domain`, and `VerificationMethod`
- **AND** call `ApplyFreeCertificateWithContext` via the service layer with `resource.Retry(tccommon.WriteRetryTimeout, ...)`
- **AND** log the returned `DnsVerification` / `FileVerification` information
- **AND** poll `CheckFreeCertificateVerificationWithContext` via the service layer until the response `CommonName` is non-empty
- **AND** return no diagnostics on success

#### Scenario: Verification polling failure surfaced
- **WHEN** `CheckFreeCertificateVerification` reports verification not passed within the poll timeout (`3 * tccommon.ReadRetryTimeout`)
- **THEN** the system SHALL return an error diagnostic describing the verification polling timeout

#### Scenario: Missing required input rejected before API call
- **WHEN** `zone_id` is null/unknown/empty, or `domain` is null/unknown/empty, or `verification_method` is null/unknown/empty at invoke time
- **THEN** the system SHALL add a diagnostic error and return without calling the API

#### Scenario: Provider client not configured
- **WHEN** the provider client (`a.Client()`) is nil at invoke time
- **THEN** the system SHALL add a diagnostic error and return without calling the API

#### Scenario: Retryable API error on ApplyFreeCertificate
- **WHEN** `ApplyFreeCertificateWithContext` returns a retryable error
- **THEN** the service layer SHALL retry within `tccommon.WriteRetryTimeout` via `resource.Retry`, wrapping the error with `tccommon.RetryError`

#### Scenario: Non-retryable API error surfaced
- **WHEN** `ApplyFreeCertificateWithContext` returns a non-retryable error after retries are exhausted
- **THEN** the service layer SHALL return the error and the `Invoke` handler SHALL add it to `resp.Diagnostics`

#### Scenario: No state persisted after successful invoke
- **WHEN** the invoke completes successfully
- **THEN** the system SHALL NOT set any id, SHALL NOT write any output/computed attribute (including verification information), and SHALL NOT implement Read/Update/Delete lifecycle methods

## Requirement: Service Layer ApplyFreeCertificate Method
The TEO service layer SHALL provide an `ApplyFreeCertificate(ctx, zoneId, domain, verificationMethod string) (dnsVerification *teo.DnsVerification, fileVerification *teo.FileVerification, errRet error)` method that constructs an `ApplyFreeCertificateRequest`, fills `ZoneId`, `Domain`, and `VerificationMethod`, wraps the `ApplyFreeCertificateWithContext` call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `ratelimit.Check` and `tccommon.RetryError`, logs success/failure, and returns the `DnsVerification` / `FileVerification` from the response on success.

#### Scenario: Service method constructs and sends request
- **WHEN** `ApplyFreeCertificate` service method is called with zone ID, domain, and verification method
- **THEN** the system SHALL create `teov20220901.NewApplyFreeCertificateRequest()`, set `ZoneId`, `Domain`, and `VerificationMethod`
- **AND** call `ApplyFreeCertificateWithContext` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`, invoking `ratelimit.Check(request.GetAction())` before the call
- **AND** on a retryable error return `tccommon.RetryError(e)`

#### Scenario: Service method returns verification info on success
- **WHEN** `ApplyFreeCertificateWithContext` succeeds
- **THEN** the service method SHALL verify `result.Response != nil`
- **AND** return the `DnsVerification` and `FileVerification` from the response (either may be nil depending on `VerificationMethod`)

#### Scenario: Service method records failure log
- **WHEN** the service method returns an error
- **THEN** a deferred function SHALL log `[CRITAL]` with the log ID, action name, request body, and error reason

## Requirement: Service Layer CheckFreeCertificateVerification Polling Method
The TEO service layer SHALL provide a `CheckFreeCertificateVerification(ctx, zoneId, domain string, timeout time.Duration) (errRet error)` method that polls `CheckFreeCertificateVerificationWithContext` via `resource.Retry(timeout, ...)`. The method SHALL return `nil` when the response `CommonName` is non-empty (certificate application succeeded). When the API returns a verification-not-passed error, the method SHALL return `resource.RetryableError` to continue polling. When the API returns a permission/authorization error, the method SHALL return `resource.NonRetryableError`. The method SHALL log success/failure.

#### Scenario: Polling succeeds when CommonName is non-empty
- **WHEN** `CheckFreeCertificateVerificationWithContext` returns a response with non-empty `CommonName`
- **THEN** the service method SHALL return `nil`

#### Scenario: Polling retries on verification-not-passed error
- **WHEN** `CheckFreeCertificateVerificationWithContext` returns a verification-not-passed error
- **THEN** the service method SHALL return `resource.RetryableError` to allow the outer retry loop to continue polling

#### Scenario: Polling fails on permission error
- **WHEN** `CheckFreeCertificateVerificationWithContext` returns a permission/authorization error (e.g. `UnauthorizedOperation.NoPermission`, `UnauthorizedOperation.CamUnauthorized`)
- **THEN** the service method SHALL return `resource.NonRetryableError` so the polling stops immediately

#### Scenario: Polling timeout
- **WHEN** the polling does not succeed within the `timeout` duration
- **THEN** the service method SHALL return the timeout error and the `Invoke` handler SHALL add it to `resp.Diagnostics`

## Requirement: Action Registration
The system SHALL register the `tencentcloud_teo_apply_free_certificate` action factory `teo.NewTeoApplyFreeCertificate` in `tencentcloud/framework/registry.go`'s `actionFactories` slice, and SHALL NOT register it in the SDKv2 `provider.go` `ResourcesMap`.

#### Scenario: Action factory registered
- **WHEN** the framework provider collects action factories
- **THEN** `frameworkActions()` SHALL include `teo.NewTeoApplyFreeCertificate`

#### Scenario: Not registered in SDKv2 provider
- **WHEN** the SDKv2 provider resources map is built
- **THEN** it SHALL NOT contain an entry keyed `tencentcloud_teo_apply_free_certificate_operation`

## Requirement: Unit Tests
The system SHALL provide unit tests in `tencentcloud/services/teo/action_tc_teo_apply_free_certificate_test.go` using gomonkey to mock the TEO cloud API, testing only business logic (no Terraform acceptance test suite).

#### Scenario: Successful invoke test
- **WHEN** a test invokes the action with valid inputs and mocks `UseTeoV20220901Client`, `ApplyFreeCertificateWithContext` to return `DnsVerification`/`FileVerification`, and `CheckFreeCertificateVerificationWithContext` to return a non-empty `CommonName`
- **THEN** the test SHALL assert the mocked `ApplyFreeCertificateWithContext` is called with the expected `ZoneId`/`Domain`/`VerificationMethod`, and that no error diagnostic is produced

#### Scenario: ApplyFreeCertificate API error test
- **WHEN** a test invokes the action and mocks `ApplyFreeCertificateWithContext` to return an error
- **THEN** the test SHALL assert an error diagnostic is produced containing the API error message

#### Scenario: Verification polling timeout test
- **WHEN** a test mocks `ApplyFreeCertificateWithContext` to succeed but `CheckFreeCertificateVerificationWithContext` to always return verification-not-passed
- **THEN** the test SHALL assert an error diagnostic is produced

#### Scenario: Missing input test
- **WHEN** a test invokes the action with `zone_id`, `domain`, or `verification_method` null/empty
- **THEN** the test SHALL assert an error diagnostic is produced and the API is not called

#### Scenario: Client not configured test
- **WHEN** the provider client is not set on the action
- **THEN** the test SHALL assert an error diagnostic mentioning "Provider not configured"

#### Scenario: Metadata and schema validation test
- **WHEN** a test calls `Metadata` and `Schema`
- **THEN** the test SHALL assert `TypeName == "tencentcloud_teo_apply_free_certificate"` and that all required attributes exist with correct types

## Requirement: Action Documentation
The system SHALL provide a markdown documentation file `action_tc_teo_apply_free_certificate.md` with a one-line description mentioning TEO, an Example Usage section using the framework action `action` block syntax, and a NOTE about Terraform version support. It SHALL NOT contain an Import section or manually-written Argument/Attribute Reference sections.

#### Scenario: Documentation file exists with required sections
- **WHEN** the action is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning TEO, an Example Usage block using `action "tencentcloud_teo_apply_free_certificate" "example" { config { ... } }` syntax, and a NOTE about Terraform 1.14+ support
