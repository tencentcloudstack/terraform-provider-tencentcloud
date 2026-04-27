## Context

The `tencentcloud_teo_l7_acc_setting` resource manages TEO (TencentCloud EdgeOne) site-level L7 acceleration settings. Currently, the `DescribeL7AccSetting` API response's `ZoneSetting` field (of type `ZoneConfigParameters`) is partially mapped: `ZoneName` is exposed as the computed `zone_name` attribute, and `ZoneConfig` is exposed as the required `zone_config` attribute. However, the full `ZoneSetting` response structure is not available as a single computed attribute, which limits users who want to access the complete configuration as a unified block.

The cloud API `DescribeL7AccSetting` returns:
- `ZoneSetting` (*ZoneConfigParameters): contains `ZoneName` (string) and `ZoneConfig` (*ZoneConfig)

The `ZoneConfig` already has all its sub-fields mapped in the existing `zone_config` schema attribute.

## Goals / Non-Goals

**Goals:**
- Add a `zone_setting` computed attribute to `tencentcloud_teo_l7_acc_setting` that represents the full `ZoneConfigParameters` response from `DescribeL7AccSetting`
- Maintain full backward compatibility with existing Terraform configurations and state
- Follow existing code patterns in the resource

**Non-Goals:**
- Modifying or removing existing `zone_name` and `zone_config` attributes
- Changing the Create/Update/Delete logic (the new attribute is computed-only)
- Adding new input parameters to the resource (zone_id already exists)

## Decisions

1. **`zone_setting` as TypeList with MaxItems:1 and Computed:true**: Following the existing pattern in this resource where complex objects are represented as `TypeList` with `MaxItems: 1`. The attribute is computed-only since it represents API response data.

2. **Reuse existing ZoneConfig schema structure for sub-fields**: The `zone_setting.zone_config` sub-attribute will have the same schema structure as the existing top-level `zone_config` attribute, ensuring consistency. Similarly, `zone_setting.zone_name` will mirror the top-level `zone_name`.

3. **Populate in Read function only**: Since `zone_setting` is computed, it only needs to be populated in the `resourceTencentCloudTeoL7AccSettingRead` function from the API response.

## Risks / Trade-offs

- **Data duplication**: The `zone_setting` attribute will contain data already available via `zone_name` and `zone_config`. This is acceptable as it provides a structured view of the complete API response → Users can choose which attribute to reference based on their needs.
- **Schema complexity**: Adding a deeply nested computed attribute increases schema complexity → The structure follows the API response format, which is the natural representation.
