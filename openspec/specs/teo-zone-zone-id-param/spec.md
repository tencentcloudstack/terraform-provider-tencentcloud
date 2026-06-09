## ADDED Requirements

### Requirement: zone_id schema parameter
The `tencentcloud_teo_zone` resource SHALL include a `zone_id` schema field of type `schema.TypeString` with `Computed: true`. This field represents the site ID returned by the TEO API.

#### Scenario: zone_id is populated after resource creation
- **WHEN** a `tencentcloud_teo_zone` resource is created successfully
- **THEN** the `zone_id` field SHALL be set to the `ZoneId` value from the `CreateZone` API response

#### Scenario: zone_id is populated during resource read
- **WHEN** a `tencentcloud_teo_zone` resource is read (refresh/plan)
- **THEN** the `zone_id` field SHALL be set to the `ZoneId` value from the `DescribeZones` API response (`Zone.ZoneId`)

### Requirement: DeleteZone uses zone_id parameter
The `resourceTencentCloudTeoZoneDelete` function SHALL use `d.Get("zone_id")` to obtain the `ZoneId` value for the `DeleteZone` API request, instead of using `d.Id()`.

#### Scenario: DeleteZone request uses zone_id from schema
- **WHEN** a `tencentcloud_teo_zone` resource is deleted
- **THEN** the `DeleteZone` API request's `ZoneId` field SHALL be populated from `d.Get("zone_id")`

#### Scenario: zone_id is empty during delete
- **WHEN** a `tencentcloud_teo_zone` resource is deleted and `zone_id` is not set in state
- **THEN** the delete function SHALL fall back to `d.Id()` as the `ZoneId` value for the `DeleteZone` API request
