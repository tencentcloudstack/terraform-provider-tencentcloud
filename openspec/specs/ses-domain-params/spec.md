## ADDED Requirements

### Requirement: SES domain resource supports DKIM option parameter

The `tencentcloud_ses_domain` resource SHALL accept an optional `dkim_option` parameter (integer) that specifies the DKIM key length. Value `0` means 1024-bit and value `1` means 2048-bit. The parameter SHALL be `ForceNew` — changing it requires resource recreation.

#### Scenario: Create domain with DKIM option
- **WHEN** user creates a `tencentcloud_ses_domain` resource with `dkim_option = 1`
- **THEN** the provider SHALL call `CreateEmailIdentity` with `DKIMOption` set to `1`
- **AND** the resource SHALL be recreated if `dkim_option` is changed

#### Scenario: Create domain without DKIM option
- **WHEN** user creates a `tencentcloud_ses_domain` resource without specifying `dkim_option`
- **THEN** the provider SHALL call `CreateEmailIdentity` without setting `DKIMOption`
- **AND** the cloud API SHALL apply its default DKIM option

#### Scenario: Read reflects DKIM option
- **WHEN** the provider reads an existing SES domain via `GetEmailIdentity`
- **THEN** the provider SHALL populate `dkim_option` in Terraform state from the response `DKIMOption` field if it is not nil

### Requirement: SES domain resource supports tag parameters

The `tencentcloud_ses_domain` resource SHALL accept optional `tag_key` and `tag_value` parameters (string) that associate a single tag with the domain. Both parameters SHALL be `ForceNew` — changing either requires resource recreation.

#### Scenario: Create domain with tag
- **WHEN** user creates a `tencentcloud_ses_domain` resource with `tag_key = "env"` and `tag_value = "prod"`
- **THEN** the provider SHALL call `CreateEmailIdentity` with a `TagList` containing one element with `TagKey = "env"` and `TagValue = "prod"`

#### Scenario: Create domain without tag
- **WHEN** user creates a `tencentcloud_ses_domain` resource without specifying `tag_key` or `tag_value`
- **THEN** the provider SHALL call `CreateEmailIdentity` without setting `TagList`

#### Scenario: Read reflects tag
- **WHEN** the provider reads an existing SES domain via `GetEmailIdentity` and the response `TagList` contains at least one element
- **THEN** the provider SHALL populate `tag_key` and `tag_value` in Terraform state from the first element of the response `TagList`

#### Scenario: Read with no tag
- **WHEN** the provider reads an existing SES domain via `GetEmailIdentity` and the response `TagList` is empty or nil
- **THEN** the provider SHALL NOT set `tag_key` or `tag_value` in Terraform state

### Requirement: Backward compatibility

The addition of `dkim_option`, `tag_key`, and `tag_value` parameters SHALL NOT break existing `tencentcloud_ses_domain` configurations. All new parameters SHALL be optional.

#### Scenario: Existing configuration without new parameters
- **WHEN** a user applies an existing `tencentcloud_ses_domain` configuration that does not include `dkim_option`, `tag_key`, or `tag_value`
- **THEN** the provider SHALL create and manage the domain as before, with no behavioral change
