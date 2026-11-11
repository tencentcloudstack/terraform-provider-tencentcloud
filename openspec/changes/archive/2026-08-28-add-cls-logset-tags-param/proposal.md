## Why

The `tencentcloud_cls_logset` resource currently manages tags indirectly through the separate TencentCloud Tag Service API. The CLS `CreateLogset` and `ModifyLogset` APIs natively support a `Tags` array of `{Key, Value}` pairs, which allows tags to be created/updated atomically alongside the logset itself. Adding direct tag parameters that map to the CLS API's `Tags` structure gives users a native, cloud-API-aligned way to bind tags to a logset at creation and update time.

## What Changes

- Add a new optional `tags` (TypeList of nested `key`/`value`) schema block to `tencentcloud_cls_logset` that maps directly to the CLS `CreateLogset`/`ModifyLogset` API `Tags []*Tag` request parameter (where each `Tag` has `Key` and `Value` fields).
- Populate the `Tags` field of `CreateLogsetRequest` and `ModifyLogsetRequest` from the new schema block so tags are sent natively through the CLS API.
- Read tags back from the `DescribeLogsets` response (`LogsetInfo.Tags`) into the new schema block.
- **BREAKING**: The existing `tags` (TypeMap) field currently managed via the Tag Service will be replaced by the new native `tags` TypeList block to align with the CLS API contract. Existing configurations using the map-style `tags` will need to migrate to the list-style block.

## Capabilities

### New Capabilities
- `cls-logset-tags`: Native CLS API tag binding (`Key`/`Value` pairs) for the `tencentcloud_cls_logset` resource via `CreateLogset` and `ModifyLogset` APIs.

### Modified Capabilities
<!-- None -->

## Impact

- **Affected code**: `tencentcloud/services/cls/resource_tc_cls_logset.go` (schema, Create, Read, Update), `tencentcloud/services/cls/resource_tc_cls_logset_test.go` (test updates), `tencentcloud/services/cls/resource_tc_cls_logset.md` (doc updates).
- **APIs**: `CreateLogset` (cls/v20201016) — `Tags` request field; `ModifyLogset` (cls/v20201016) — `Tags` request field; `DescribeLogsets` (cls/v20201016) — `LogsetInfo.Tags` response field.
- **Dependencies**: No new vendor dependencies required; `cls.Tag` struct already available in the vendored SDK.
- **Backward compatibility**: The `tags` field changes from TypeMap (Tag Service managed) to TypeList (native CLS API managed). Users must migrate existing map-style tag configurations to the new list-style block.
