## Context

The `tencentcloud_mongodb_sharding_instance` resource manages MongoDB sharding instances through the `ModifyDBInstanceSpec` API. The vendor SDK (`v20190725`) already supports `AddNodeList` and `RemoveNodeList` fields in the `ModifyDBInstanceSpecRequest` struct. The service layer (`MongodbService.UpgradeInstance`) already has handling code for these params via the `params map[string]interface{}` argument. The only missing piece is exposing these parameters in the Terraform resource schema and wiring them through the update function.

## Goals / Non-Goals

**Goals:**
- Add `add_node_list` (TypeList) and `remove_node_list` (TypeList) as Optional schema fields
- Wire these parameters to the existing `UpgradeInstance` service call in the update path
- Both fields are Optional and preserve full backward compatibility

**Non-Goals:**
- No changes to Create, Read, or Delete paths
- No new API calls or service methods
- No modification to the `UpgradeInstance` service method (it already handles these params)

## Decisions

1. **Schema field structure**: `add_node_list` uses `TypeList` with `MaxItems: 1` containing sub-fields `role` (Required, string) and `zone` (Required, string). `remove_node_list` uses `TypeList` with sub-fields `role` (Required, string), `node_name` (Required, string), and `zone` (Required, string). This mirrors the SDK struct layout exactly.

2. **Update trigger**: When `d.HasChange("add_node_list") || d.HasChange("remove_node_list")`, call the existing `UpgradeInstance` service method with these params in the `params` map. The `UpgradeInstance` method already handles constructing `AddNodeList` and `RemoveNodeList` from params.

3. **No state management in Read**: The `add_node_list` and `remove_node_list` parameters represent node addition/removal operations (not persistent state). They are not set back in the Read function. This follows the pattern of action-oriented parameters that drive changes but don't represent current state.

## Risks / Trade-offs

- **Risk**: After successful node addition/removal, the `d.SetId("")` won't be triggered but `add_node_list`/`remove_node_list` values remain in the Terraform config, causing a perpetual diff on next plan. → **Mitigation**: These fields should be documented as operational parameters that users should remove from config after the operation completes, or we can optionally clear them after successful execution.
- **Risk**: Adding/removing nodes is a long-running async operation that could timeout. → **Mitigation**: The `UpgradeInstance` service method already handles this via `DescribeDBInstanceDeal` polling with retry logic.