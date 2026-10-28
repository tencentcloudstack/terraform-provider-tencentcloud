## Context

The `tencentcloud_dbdc_db_custom_node` resource (RESOURCE_KIND_GENERAL) currently manages the full CRUD lifecycle of a DB Custom node via the `dbdc` v20201029 SDK. It already exposes create-time parameters such as `zone`, `image_id`, `vpc_id`, `subnet_id`, `node_type`, `login_settings`, `security_group_ids`, etc., and computed fields refreshed from `DescribeDBCustomNodes`.

**Current state:**
- Resource file: `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.go`
- Service layer: `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` (`DescribeDBCustomNodeById` wraps `DescribeDBCustomNodes`, returns `*dbdcv20201029.DBCustomNode`)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` (already vendored, no upgrade needed)

**API behavior analysis (vendored SDK):**

| API | `DisasterRecoverGroupIds` in Request | `DisasterRecoverGroupId` in Response |
|-----|--------------------------------------|--------------------------------------|
| `CreateDBCustomNodes` | Yes (`DisasterRecoverGroupIds []*string`, "仅支持指定一个") | No |
| `DescribeDBCustomNodes` | N/A | Yes (`DBCustomNode.DisasterRecoverGroupId *string`, inside `Response.NodeSet[]`) |
| `RenewDBCustomNode` / `ModifyDBCustomNodeTags` / `ModifyDBCustomNodeSecurityGroups` | No | N/A |
| `IsolateDBCustomNode` / `DestroyDBCustomNode` | No | N/A |

**Key constraint:** `DisasterRecoverGroupIds` is only accepted by `CreateDBCustomNodes`. There is no update API that modifies the placement group of an existing node, so the input parameter is `ForceNew`. The `DescribeDBCustomNodes` response includes `DisasterRecoverGroupId`, so it can be refreshed on Read as a computed field.

## Goals / Non-Goals

**Goals:**
- Add `disaster_recover_group_ids` (Optional, `ForceNew`, `TypeList` of `TypeString`, `MaxItems: 1`) input parameter to `tencentcloud_dbdc_db_custom_node`, mapped to `request.DisasterRecoverGroupIds` in `CreateDBCustomNodes`. The API documents that only one placement group ID is supported.
- Add `disaster_recover_group_id` (Computed, `TypeString`) output parameter to `tencentcloud_dbdc_db_custom_node`, mapped to `DBCustomNode.DisasterRecoverGroupId` from the `DescribeDBCustomNodes` response, refreshed during Read.
- Maintain full backward compatibility — existing configurations continue to work unchanged.

**Non-Goals:**
- Supporting modification of the placement group after creation (the API has no such update path).
- Adding placement-group-related parameters to any `dbdc` data source (out of scope).

## Decisions

### Decision 1: `disaster_recover_group_ids` is `ForceNew` (not immutable-args pattern)

**Rationale:** The `CreateDBCustomNodes` API is the only API that accepts `DisasterRecoverGroupIds`; there is no update API that can change a node's placement group. Because changing this value requires recreating the node, `ForceNew: true` is the correct Terraform semantics. Unlike resources that only have CRD-style APIs (where the project convention is to use an `immutableArgs` array), this resource has a real Update flow for other fields (tags, period, security groups), so `ForceNew` on the single immutable input is the idiomatic choice.

### Decision 2: Input is a list (`disaster_recover_group_ids`) with `MaxItems: 1`; output is a single string (`disaster_recover_group_id`)

**Rationale:** The API request field `DisasterRecoverGroupIds` is `[]*string` (an array), so the Terraform input mirrors it as a `TypeList` to preserve the natural mapping and forward compatibility if the API ever relaxes the single-item limit. The API documents "仅支持指定一个", so `MaxItems: 1` is enforced to give users a clear early error. The API response field `DisasterRecoverGroupId` is a single `*string` on the node, so the computed output is a single `TypeString` named `disaster_recover_group_id`. This mirrors the existing pattern in the codebase where list inputs and singular computed outputs coexist.

### Decision 3: Read `disaster_recover_group_id` from the existing `DescribeDBCustomNodeById` result

**Rationale:** The service-layer helper `DescribeDBCustomNodeById` already returns `*dbdcv20201029.DBCustomNode`, whose struct now includes `DisasterRecoverGroupId`. No service-layer change is required — the Read function simply adds a nil-guarded `d.Set("disaster_recover_group_id", respData.DisasterRecoverGroupId)` alongside the existing field reads.

## Risks / Trade-offs

- **[Risk] Changing `disaster_recover_group_ids` destroys and recreates the node**: Using `ForceNew: true` means changing the placement group forces resource replacement.
  - **Mitigation:** This is the only correct behavior since no update API exists; the recreation is explicit and visible to the user in the plan.

- **[Risk] `disaster_recover_group_id` may be empty for nodes not in a placement group**: The API may return an empty/nil `DisasterRecoverGroupId` for nodes created without a placement group.
  - **Mitigation:** The Read function nil-checks `respData.DisasterRecoverGroupId` before calling `d.Set`, consistent with the existing Read pattern for all other computed fields.
