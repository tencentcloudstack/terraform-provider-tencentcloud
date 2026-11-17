## Why

The `tencentcloud_dlc_data_engine` resource currently does not support binding tags to a data engine at creation time. Users who rely on tags for cost allocation, access control, or resource grouping must manage tags out-of-band, which breaks the infrastructure-as-code contract and causes state drift.

The DLC `CreateDataEngine` API already accepts a `Tags` input (`[]*TagInfo` with `TagKey`/`TagValue`), and the `DescribeDataEngine`/`DescribeDataEngines` APIs return the bound tags in `DataEngineInfo.TagList`. Exposing these as a Terraform `tags` block lets users declare tags alongside the rest of the engine configuration.

## What Changes

- Add a new `tags` block (TypeList) to the `tencentcloud_dlc_data_engine` resource schema, with `tag_key` and `tag_value` string sub-fields, modeled on the cloud API `TagInfo` structure.
- In `CreateDataEngine`, map the `tags` block to `request.Tags` (`[]*dlc.TagInfo`).
- In the resource `Read`, map `DataEngineInfo.TagList` (`[]*dlc.TagInfo`) back into the `tags` block. Note: the current `Read` uses the `DescribeDataEngine` API (singular) which returns the same `DataEngineInfo` type that carries `TagList`.
- Because `UpdateDataEngine` does **not** support a `Tags` field, mark the `tags` parameter as `ForceNew: true` so that changing tags triggers recreation (consistent with other immutable create-only parameters on this resource).

## Capabilities

### New Capabilities
- `dlc-data-engine-tags`: Tag binding support for the `tencentcloud_dlc_data_engine` resource — schema, create-time population, read-time synchronization, and ForceNew update constraint.

### Modified Capabilities
<!-- None — no existing spec for dlc-data-engine. -->

## Impact

- **Affected code:**
  - `tencentcloud/services/dlc/resource_tc_dlc_data_engine.go` — add `tags` schema field; populate `request.Tags` in `Create`; read `dataEngine.TagList` in `Read`; add `tags` to the `immutableArgs` list in `Update` (it is ForceNew, so Terraform handles recreation, but it is listed for consistency with the existing create-only arguments).
  - `tencentcloud/services/dlc/resource_tc_dlc_data_engine_test.go` — add unit tests covering tags marshalling and flattening using gomonkey mocks (per project rule for modified resources, extend the existing test suite).
  - `tencentcloud/services/dlc/resource_tc_dlc_data_engine.md` — document the new `tags` block and example usage.
- **APIs used:**
  - `CreateDataEngine` (input `Tags []*TagInfo`) — `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125`
  - `DescribeDataEngine` (output `DataEngine.TagList []*TagInfo`) — same package
- **Backward compatibility:** Fully backward compatible. `tags` is `Optional`; existing configurations and state are unaffected. `UpdateDataEngine` lacks a `Tags` field, so `tags` is `ForceNew` (recreation on change), consistent with the resource's existing create-only parameters.
- **Dependencies:** None new; the vendor SDK already defines `TagInfo`, `CreateDataEngineRequest.Tags`, and `DataEngineInfo.TagList`.
