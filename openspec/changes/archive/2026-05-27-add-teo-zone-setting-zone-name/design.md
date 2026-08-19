## Context

The `tencentcloud_teo_zone_setting` resource manages TEO zone settings via the `ModifyZoneSetting` and `DescribeZoneSetting` APIs. The `DescribeZoneSetting` response includes a `ZoneSetting` struct that contains a `ZoneName` field (站点名称), but this field is not currently exposed in the Terraform resource schema.

The `ModifyZoneSettingRequestParams` struct does NOT include `ZoneName`, confirming it is a read-only attribute from the API perspective.

Current resource file: `tencentcloud/services/teo/resource_tc_teo_zone_setting.go`

## Goals / Non-Goals

**Goals:**
- Expose `zone_name` as a Computed-only attribute in the `tencentcloud_teo_zone_setting` resource schema.
- Read and set the value from `DescribeZoneSetting` response's `ZoneSetting.ZoneName` field in the Read function.
- Maintain backward compatibility with existing Terraform state files.

**Non-Goals:**
- Making `zone_name` writable (it is not supported by `ModifyZoneSetting`).
- Changing any existing schema fields or behavior.
- Adding `zone_name` to Create or Update logic.

## Decisions

1. **Schema type: Computed-only string**
   - Rationale: `ZoneName` is only available in the read response, not in the modify request. It is a server-determined value that users cannot set.
   - Alternative: Optional+Computed was considered but rejected since the API does not accept this field on write.

2. **Placement in schema: top-level attribute**
   - Rationale: `ZoneName` is a top-level field in the `ZoneSetting` struct, similar to the existing `area` field which is also Computed-only at the top level.

3. **Nil check before setting**
   - Rationale: The API documentation notes the field may return null. A nil check must be performed before calling `d.Set("zone_name", ...)` to avoid setting nil values.

## Risks / Trade-offs

- [Risk] Field returns nil for some zones → Mitigation: Nil check before `d.Set`; field simply won't appear in state if nil.
- [Risk] Backward compatibility with existing state → Mitigation: Adding a Computed-only field is always backward compatible; existing states will simply gain the new attribute on next refresh.
