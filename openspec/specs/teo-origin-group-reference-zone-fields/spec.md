## ADDED Requirements

### Requirement: References block exposes zone information
The `references` nested block of `tencentcloud_teo_origin_group` resource SHALL include `zone_id`, `zone_name`, and `alias_zone_name` as computed attributes, sourced from the `DescribeOriginGroup` API response's `OriginGroupReference` struct.

#### Scenario: Read populates zone fields when API returns them
- **WHEN** the resource Read method calls `DescribeOriginGroup` and the response `OriginGroupReference` contains `ZoneId`, `ZoneName`, and `AliasZoneName` values
- **THEN** the Terraform state SHALL contain `zone_id`, `zone_name`, and `alias_zone_name` in each `references` block entry matching the API response values

#### Scenario: Read handles nil zone fields gracefully
- **WHEN** the resource Read method calls `DescribeOriginGroup` and the response `OriginGroupReference` has `ZoneId`, `ZoneName`, or `AliasZoneName` set to nil
- **THEN** the corresponding Terraform state fields SHALL NOT be set (nil fields are skipped, matching the existing pattern for other computed fields)

#### Scenario: Zone fields are computed and not user-settable
- **WHEN** a user creates or updates a `tencentcloud_teo_origin_group` resource
- **THEN** `zone_id`, `zone_name`, and `alias_zone_name` SHALL be computed-only fields that cannot be set by the user in the Terraform configuration

#### Scenario: Backward compatibility with existing state
- **WHEN** an existing `tencentcloud_teo_origin_group` resource state is refreshed
- **THEN** the existing `instance_type`, `instance_id`, and `instance_name` fields SHALL continue to work unchanged, and the new zone fields SHALL be populated without requiring any user configuration changes
