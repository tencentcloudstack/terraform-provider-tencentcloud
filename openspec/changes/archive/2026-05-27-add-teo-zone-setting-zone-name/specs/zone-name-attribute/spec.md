## ADDED Requirements

### Requirement: Expose zone_name as computed attribute
The `tencentcloud_teo_zone_setting` resource SHALL expose a `zone_name` attribute of type string that is Computed-only. The value SHALL be read from the `DescribeZoneSetting` API response field `ZoneSetting.ZoneName`. The attribute MUST NOT be settable by users in Terraform configuration.

#### Scenario: zone_name is populated on read
- **WHEN** the `DescribeZoneSetting` API returns a non-nil `ZoneSetting.ZoneName` value
- **THEN** the resource state SHALL contain `zone_name` set to that value

#### Scenario: zone_name is nil in API response
- **WHEN** the `DescribeZoneSetting` API returns a nil `ZoneSetting.ZoneName` value
- **THEN** the resource state SHALL NOT set `zone_name` (skip the `d.Set` call)

#### Scenario: zone_name is not sent on modify
- **WHEN** the user applies changes to the `tencentcloud_teo_zone_setting` resource
- **THEN** the `ModifyZoneSetting` API request SHALL NOT include any `ZoneName` parameter
