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

### Requirement: SES domain resource supports tag_list parameter

The `tencentcloud_ses_domain` resource SHALL accept an optional `tag_list` nested block parameter that associates one or more tags with the domain. Each `tag_list` block SHALL contain `tag_key` and `tag_value` string sub-fields. The `tag_list` parameter SHALL be `ForceNew` — changing it requires resource recreation.

#### Scenario: Create domain with tags
- **WHEN** user creates a `tencentcloud_ses_domain` resource with a `tag_list` block containing `tag_key = "env"` and `tag_value = "prod"`
- **THEN** the provider SHALL call `CreateEmailIdentity` with a `TagList` containing one element with `TagKey = "env"` and `TagValue = "prod"`

#### Scenario: Create domain with multiple tags
- **WHEN** user creates a `tencentcloud_ses_domain` resource with multiple `tag_list` blocks
- **THEN** the provider SHALL call `CreateEmailIdentity` with a `TagList` containing one element per block, preserving all tag key/value pairs

#### Scenario: Create domain without tag
- **WHEN** user creates a `tencentcloud_ses_domain` resource without specifying `tag_list`
- **THEN** the provider SHALL call `CreateEmailIdentity` without setting `TagList`

#### Scenario: Read reflects tags
- **WHEN** the provider reads an existing SES domain via `GetEmailIdentity` and the response `TagList` contains at least one element
- **THEN** the provider SHALL populate `tag_list` in Terraform state with one block per element, each containing `tag_key` and `tag_value` from the response

#### Scenario: Read with no tag
- **WHEN** the provider reads an existing SES domain via `GetEmailIdentity` and the response `TagList` is empty or nil
- **THEN** the provider SHALL NOT set `tag_list` in Terraform state

### Requirement: Backward compatibility

The addition of `dkim_option` and `tag_list` parameters SHALL NOT break existing `tencentcloud_ses_domain` configurations. All new parameters SHALL be optional.

#### Scenario: Existing configuration without new parameters
- **WHEN** a user applies an existing `tencentcloud_ses_domain` configuration that does not include `dkim_option` or `tag_list`
- **THEN** the provider SHALL create and manage the domain as before, with no behavioral change
