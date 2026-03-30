## ADDED Requirements

### Requirement: NAT gateway datasource supports verbose level parameter
The tencentcloud_nats datasource SHALL accept a `verbose_level` parameter to control the level of detail returned by the DescribeNatGateways API.

#### Scenario: Valid verbose level value
- **WHEN** user sets `verbose_level` to "DETAIL", "COMPACT", or "SIMPLE"
- **THEN** the datasource SHALL pass the value to the VerboseLevel parameter of the DescribeNatGateways API request

#### Scenario: Invalid verbose level value
- **WHEN** user sets `verbose_level` to an invalid value (not "DETAIL", "COMPACT", or "SIMPLE")
- **THEN** the datasource SHALL return a validation error with a message listing the valid options

#### Scenario: Verbose level not specified
- **WHEN** user does not specify the `verbose_level` parameter
- **THEN** the datasource SHALL NOT set the VerboseLevel parameter in the API request and the API SHALL use its default behavior

### Requirement: Verbose level parameter validation
The verbose_level parameter MUST be validated to ensure it only accepts the values "DETAIL", "COMPACT", or "SIMPLE".

#### Scenario: DETAIL value validation
- **WHEN** user sets `verbose_level` to "DETAIL"
- **THEN** the validation SHALL pass and the parameter SHALL be accepted

#### Scenario: COMPACT value validation
- **WHEN** user sets `verbose_level` to "COMPACT"
- **THEN** the validation SHALL pass and the parameter SHALL be accepted

#### Scenario: SIMPLE value validation
- **WHEN** user sets `verbose_level` to "SIMPLE"
- **THEN** the validation SHALL pass and the parameter SHALL be accepted

#### Scenario: Case sensitivity validation
- **WHEN** user sets `verbose_level` to "detail" (lowercase)
- **THEN** the validation SHALL fail and return an error

### Requirement: Verbose level parameter is optional
The verbose_level parameter SHALL be an optional parameter in the datasource schema.

#### Scenario: Omitting verbose level
- **WHEN** user configures the tencentcloud_nats datasource without specifying `verbose_level`
- **THEN** the configuration SHALL be valid and the datasource SHALL execute successfully

#### Scenario: Including verbose level
- **WHEN** user configures the tencentcloud_nats datasource with `verbose_level` specified
- **THEN** the configuration SHALL be valid and the datasource SHALL execute successfully with the specified verbose level
