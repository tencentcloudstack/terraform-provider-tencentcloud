## Context

The `tencentcloud_teo_l7_acc_setting` resource manages TEO (TencentCloud EdgeOne) site-level L7 acceleration settings. It previously had a `zone_setting` computed attribute that duplicated data from `zone_name` and `zone_config`. The SDK's `ZoneConfig` struct also includes `NetworkErrorLogging *NetworkErrorLoggingParameters` which was not exposed in the Terraform resource.

## Goals / Non-Goals

**Goals:**
- Remove the redundant `zone_setting` computed attribute to simplify the resource schema
- Add `network_error_logging` configuration support under `zone_config`
- Maintain backward compatibility for existing `zone_config` configurations

**Non-Goals:**
- Changing existing `zone_config` sub-fields
- Adding any other missing `ZoneConfig` fields beyond `network_error_logging`

## Decisions

1. **Remove `zone_setting` entirely**: The attribute was a computed-only duplicate of `zone_name` + `zone_config`. Removing it simplifies the schema and eliminates ~500 lines of redundant code in both schema definition and Read function.

2. **Add `network_error_logging` as Optional under `zone_config`**: Following the pattern of all other `zone_config` sub-fields (e.g., `grpc`, `quic`). The `NetworkErrorLoggingParameters` struct has a single `Switch` field (on/off).

3. **Place `network_error_logging` between `grpc` and `accelerate_mainland`**: This matches the order in the SDK's `ZoneConfig` struct where `NetworkErrorLogging` appears between `Grpc` and `AccelerateMainland`.

4. **Remove unit tests for `zone_setting`**: The 3 gomonkey mock tests (`TestL7AccSettingZoneSetting_Read_Success`, `_Read_APIError`, `_Read_EmptyResponse`) were specifically for validating `zone_setting` population. The acceptance tests remain and cover the overall resource functionality.

## Risks / Trade-offs

- **Breaking change for `zone_setting` users**: Users who referenced `zone_setting` in their configurations will need to switch to `zone_name`/`zone_config`. This is acceptable since `zone_setting` was a recently added computed attribute with limited adoption.
- **Adding `network_error_logging`**: No risk, purely additive Optional field.
