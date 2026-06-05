## ADDED Requirements

### Requirement: Resource schema definition

The resource `tencentcloud_teo_check_free_certificate_verification` SHALL define the following schema:

- `zone_id` (String, Required, ForceNew): 站点 ID。
- `domain` (String, Required, ForceNew): 需要验证的域名。
- `common_name` (String, Computed): 免费证书申请成功时，该证书颁发给的域名。
- `signature_algorithm` (String, Computed): 免费证书申请成功时，该证书使用的签名算法。
- `expire_time` (String, Computed): 免费证书申请成功时，该证书的过期时间。

#### Scenario: Schema contains all required and computed fields
- **WHEN** the resource schema is inspected
- **THEN** it SHALL contain `zone_id` as Required+ForceNew, `domain` as Required+ForceNew, and `common_name`, `signature_algorithm`, `expire_time` as Computed

### Requirement: Create operation calls CheckFreeCertificateVerification API

The resource Create method SHALL call the TEO `CheckFreeCertificateVerification` API with `zone_id` and `domain` from the Terraform configuration, and store the response fields (`CommonName`, `SignatureAlgorithm`, `ExpireTime`) into state.

#### Scenario: Successful certificate verification
- **WHEN** `terraform apply` is executed with valid `zone_id` and `domain`
- **THEN** the resource SHALL call `CheckFreeCertificateVerification` API, set the resource ID to `zone_id#domain`, and store `common_name`, `signature_algorithm`, `expire_time` from the API response into Terraform state

#### Scenario: API call fails
- **WHEN** the `CheckFreeCertificateVerification` API returns an error
- **THEN** the resource SHALL return the error to Terraform with appropriate error message

### Requirement: Read, Update, Delete are no-op

The resource SHALL implement Read, Update, and Delete as no-op methods since this is a RESOURCE_KIND_OPERATION type resource that does not maintain state.

#### Scenario: Read returns nil
- **WHEN** Terraform calls the Read method
- **THEN** the method SHALL return nil without performing any action

#### Scenario: Delete returns nil
- **WHEN** Terraform calls the Delete method
- **THEN** the method SHALL return nil without performing any action

### Requirement: Retry on API failure

The Create method SHALL wrap the API call with `resource.Retry` using `tccommon.ReadRetryTimeout` as the timeout, and use `tccommon.RetryError()` to wrap errors for retry.

#### Scenario: Transient error triggers retry
- **WHEN** the API call fails with a transient error
- **THEN** the resource SHALL retry the API call within the configured timeout

### Requirement: Response nil check

The Create method SHALL check if the API response and its fields are nil before setting state values.

#### Scenario: Response field is nil
- **WHEN** the API response contains nil fields (e.g., `CommonName` is nil)
- **THEN** the resource SHALL skip setting that field in state rather than causing a nil pointer error
