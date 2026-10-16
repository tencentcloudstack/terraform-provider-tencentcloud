## ADDED Requirements

### Requirement: SSM secret encryption type is configurable
The `tencentcloud_ssm_secret` resource SHALL expose an `encrypt_type` argument that maps to the SSM `CreateSecret.EncryptType` input field, allowing users to choose the secret's encryption type at creation time. Accepted values SHALL be `0` (KMS key encryption) and `1` (soft key encryption). The field SHALL be optional and computed so that when omitted the SSM API default applies and is read back into state.

#### Scenario: Create secret with KMS encryption (encrypt_type = 0)
- **WHEN** a user creates `tencentcloud_ssm_secret` with `encrypt_type = 0`
- **THEN** the provider SHALL set `request.EncryptType = 0` in the `CreateSecret` API call
- **AND** after creation the resource state SHALL contain `encrypt_type = 0`

#### Scenario: Create secret with soft key encryption (encrypt_type = 1)
- **WHEN** a user creates `tencentcloud_ssm_secret` with `encrypt_type = 1`
- **THEN** the provider SHALL set `request.EncryptType = 1` in the `CreateSecret` API call
- **AND** after creation the resource state SHALL contain `encrypt_type = 1`

#### Scenario: Create secret without specifying encrypt_type
- **WHEN** a user creates `tencentcloud_ssm_secret` without setting `encrypt_type`
- **THEN** the provider SHALL NOT set `request.EncryptType` in the `CreateSecret` API call
- **AND** the SSM API default encryption type SHALL apply
- **AND** the resource state SHALL be populated with the value returned by `DescribeSecret`

### Requirement: SSM secret encryption type is read back from the API
The provider SHALL read the `EncryptType` field from the `DescribeSecret` response and set it on the `encrypt_type` state attribute during the read operation. The value SHALL be obtained via the existing `DescribeSecretByName` service method, which SHALL capture the `EncryptType` response field.

#### Scenario: Read secret reflects current encryption type
- **WHEN** the provider reads an existing `tencentcloud_ssm_secret`
- **THEN** it SHALL call `DescribeSecret` and extract `EncryptType` from the response
- **AND** set `encrypt_type` in the Terraform state to the returned value
- **AND** when the API returns `nil` for `EncryptType`, the field SHALL fall back to `0`

### Requirement: SSM secret encryption type is immutable
Because the SSM `UpdateSecret` API does not accept an `EncryptType` parameter, the `encrypt_type` argument SHALL be treated as immutable after creation. The `encrypt_type` argument SHALL NOT be `ForceNew`; instead the provider SHALL list `encrypt_type` in the `immutableArgs` array of the update method so that changing `encrypt_type` on an existing resource returns an explicit error rather than triggering re-creation.

#### Scenario: Changing encrypt_type is rejected with an error
- **WHEN** a user changes `encrypt_type` on an existing `tencentcloud_ssm_secret`
- **THEN** the provider SHALL return an error indicating that `encrypt_type` cannot be changed
- **AND** the resource SHALL NOT be destroyed or re-created

#### Scenario: Update without changing encrypt_type
- **WHEN** a user updates other attributes (e.g. `description`, `is_enabled`, `tags`) without touching `encrypt_type`
- **THEN** the provider SHALL NOT send `EncryptType` in any update API call
- **AND** the existing `encrypt_type` SHALL be preserved