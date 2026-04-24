## Context

The `tencentcloud_teo_zone` resource manages TEO (TencentCloud EdgeOne) sites. Currently, the `DeleteZone` API call in the delete function uses `d.Id()` to populate the `ZoneId` request parameter. The resource ID is set during create from the `CreateZone` response's `ZoneId` field. The `zone_id` is not exposed as a schema field, making it inaccessible for direct reference in Terraform configurations.

The `DeleteZone` API (from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) accepts a `ZoneId` parameter in its request struct. The resource already uses this field via `d.Id()`, but it should also be available as an explicit schema parameter `zone_id`.

## Goals / Non-Goals

**Goals:**
- Add `zone_id` as a Computed schema parameter to the `tencentcloud_teo_zone` resource
- Populate `zone_id` in the read function from the `DescribeZones` API response (`Zone.ZoneId`)
- Update the delete function to use `d.Get("zone_id")` for the `DeleteZone` API's `ZoneId` parameter instead of `d.Id()`
- Maintain backward compatibility - existing Terraform configurations and state must not break

**Non-Goals:**
- Changing the resource ID mechanism (d.Id() still returns zone_id)
- Modifying any other API operations (CreateZone, ModifyZone, ModifyZoneStatus, ModifyZoneWorkMode)
- Adding any other new parameters to this resource

## Decisions

### Decision 1: `zone_id` as Computed-only field
**Choice**: Add `zone_id` as `Computed: true` (not Optional+Computed)
**Rationale**: The `zone_id` is generated server-side by the CreateZone API and is equivalent to the resource ID. Users cannot specify it - it's always set from the API response. Making it Computed-only avoids confusion and maintains consistency with the resource ID.

### Decision 2: Delete function uses `d.Get("zone_id")` instead of `d.Id()`
**Choice**: In `resourceTencentCloudTeoZoneDelete`, read `zone_id` from `d.Get("zone_id")` and use it as the `ZoneId` in the `DeleteZone` request.
**Rationale**: This makes the parameter flow explicit and consistent with the schema definition. The `d.Id()` and `d.Get("zone_id")` will have the same value since `zone_id` is set from the create response which is the same value stored as the resource ID.

### Decision 3: Read function sets `zone_id` from API response
**Choice**: In `resourceTencentCloudTeoZoneRead`, set `zone_id` from `respData.ZoneId`.
**Rationale**: The `Zone` struct returned by `DescribeZones` contains `ZoneId` which should be persisted in the state as `zone_id`.

## Risks / Trade-offs

- [Risk] State migration for existing resources: existing state files won't have `zone_id` set → Mitigation: After the next `terraform refresh` or `terraform plan`, the read function will populate `zone_id` from the API response. Since it's Computed-only, no user action is needed.
- [Risk] `d.Id()` and `d.Get("zone_id")` could theoretically differ → Mitigation: Both are set from the same source (create response ZoneId / read response ZoneId), so they will always be consistent.
