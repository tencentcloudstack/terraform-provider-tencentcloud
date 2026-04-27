## ADDED Requirements

### Requirement: zone_setting computed attribute
The `tencentcloud_teo_l7_acc_setting` resource SHALL expose a computed attribute `zone_setting` of type `TypeList` with `MaxItems: 1` and `Computed: true`, mapping to the `DescribeL7AccSetting` API response's `ZoneSetting` field (type `ZoneConfigParameters`).

The `zone_setting` attribute SHALL contain the following sub-attributes:
- `zone_name` (TypeString, Computed): The site name, mapping from `ZoneConfigParameters.ZoneName`
- `zone_config` (TypeList, MaxItems:1, Computed): The site acceleration configuration, mapping from `ZoneConfigParameters.ZoneConfig`. The `zone_config` sub-attribute SHALL have the same schema structure as the existing top-level `zone_config` attribute in this resource.

#### Scenario: Read populates zone_setting from DescribeL7AccSetting response
- **WHEN** the `resourceTencentCloudTeoL7AccSettingRead` function is called and the `DescribeL7AccSetting` API returns a `ZoneSetting` with `ZoneName` and `ZoneConfig`
- **THEN** the `zone_setting` attribute SHALL be populated with `zone_name` set from `ZoneSetting.ZoneName` and `zone_config` set from `ZoneSetting.ZoneConfig`

#### Scenario: zone_setting is computed and not settable by user
- **WHEN** a user creates or updates the `tencentcloud_teo_l7_acc_setting` resource
- **THEN** the `zone_setting` attribute SHALL NOT be settable by the user; it SHALL only be populated from the API response during read operations

#### Scenario: Existing attributes remain unchanged
- **WHEN** the `zone_setting` attribute is added to the resource schema
- **THEN** the existing `zone_name` (computed) and `zone_config` (required) attributes SHALL continue to function as before with no changes to their behavior
