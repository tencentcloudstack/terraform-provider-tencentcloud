## Context

The `tencentcloud_dbdc_db_custom_node` resource (`tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.go`) wraps the `CreateDBCustomNodes` API of the `dbdc` service (SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`).

**Current state:**
- The resource schema exposes `zone`, `image_id`, `vpc_id`, `subnet_id`, `node_type`, `period`, `node_name`, `login_settings`, `auto_voucher`, `voucher_ids`, `auto_renew`, `tags`, plus a set of Computed-only fields (`charge_type`, `system_disk`, `data_disks`, `node_id`, `cluster_id`, etc.).
- Create builds a `CreateDBCustomNodesRequest` from the schema and polls the returned `TaskId` via the existing `waitDBCustomTaskSucceeded` helper.
- Read calls the existing service helper `DescribeDBCustomNodeById` (wraps `DescribeDBCustomNodes`, returns `*dbdcv20201029.DBCustomNode`).
- Update handles `tags` (via `ModifyDBCustomNodeTags`) and `period`/`auto_renew` (via `RenewDBCustomNode`).
- Delete isolates then destroys the node.

**API behavior analysis** (verified against the vendored SDK structs):

| API | Field in Request | Field in Response (`DBCustomNode`) | Modify API available? |
|-----|------------------|------------------------------------|------------------------|
| `CreateDBCustomNodes` | `ChargeType` | `ChargeType` | No |
| `CreateDBCustomNodes` | `NetworkMode` | `NetworkMode` | No |
| `CreateDBCustomNodes` | `SystemDisk` (`DiskType`,`DiskSize`) | `SystemDisk` | No |
| `CreateDBCustomNodes` | `DataDisks` (`DiskType`,`DiskSize`,`DiskName`) | `DataDisks` | No |
| `CreateDBCustomNodes` | `HostName` | **(not returned)** | No |
| `CreateDBCustomNodes` | `SecurityGroupIds` | **(not in `DBCustomNode`)** | Yes: `ModifyDBCustomNodeSecurityGroups` |
| `DescribeDBCustomNodeSecurityGroups` | `NodeId` | `Groups[]` (`SecurityGroupId`) | (read API) |

**Key constraints:**
- `ChargeType`, `NetworkMode`, `SystemDisk`, `DataDisks`, `HostName` have no Modify API → they must be `ForceNew`.
- `HostName` is not returned by `DescribeDBCustomNodes` → it cannot be refreshed on Read (write-only).
- `SecurityGroupIds` is mutable: `ModifyDBCustomNodeSecurityGroups` accepts `NodeId` + `SecurityGroupIds[]` and performs a full overwrite (set semantic). Read must use the separate `DescribeDBCustomNodeSecurityGroups` API because `DBCustomNode` does not carry security groups.

## Goals / Non-Goals

**Goals:**
- Promote `charge_type` from Computed to Optional+Computed (ForceNew) and pass it to `CreateDBCustomNodesRequest.ChargeType`.
- Add `network_mode` (Optional, ForceNew) → `CreateDBCustomNodesRequest.NetworkMode`.
- Promote `system_disk` block (Optional+Computed, ForceNew) → `CreateDBCustomNodesRequest.SystemDisk` (`DiskType`, `DiskSize`).
- Promote `data_disks` block (Optional+Computed, ForceNew) → `CreateDBCustomNodesRequest.DataDisks` (`DiskType`, `DiskSize`, `DiskName`).
- Add `host_name` (Optional, ForceNew) → `CreateDBCustomNodesRequest.HostName`.
- Add `security_group_ids` (Optional, mutable) → `CreateDBCustomNodesRequest.SecurityGroupIds` on Create, `ModifyDBCustomNodeSecurityGroups` on Update, `DescribeDBCustomNodeSecurityGroups` on Read.
- Maintain full backward compatibility (existing configs/state unchanged).
- Update `.md` example and add test coverage.

**Non-Goals:**
- Adding an update path for `charge_type`, `network_mode`, `system_disk`, `data_disks`, `host_name` (no Modify API exists).
- Exposing `host_name` as a readable/computed field (API does not return it).
- Changes to the `tencentcloud_dbdc_db_custom_nodes` data source.

## Decisions

### Decision 1: ForceNew for parameters without a Modify API

`ChargeType`, `NetworkMode`, `SystemDisk`, `DataDisks`, `HostName` are accepted only by `CreateDBCustomNodes`. There is no Modify API to change them on an existing node. Setting `ForceNew: true` makes Terraform destroy+recreate the node when any of these change, which matches the API reality.

**Alternatives considered:** An `immutableArgs` array that returns a clear error (used by the cvm placement-group resource). That pattern is preferable when ForceNew would be silently destructive and an explicit error is friendlier. Here, however, these are genuine lifecycle-affecting inputs (disk size, charge type) where rebuild is the only valid response, so ForceNew is the conventional Terraform behavior and is consistent with the other already-ForceNew fields in this resource (`zone`, `image_id`, `node_type`, etc.).

### Decision 2: Promote Computed fields to Optional+Computed (not plain Optional)

`charge_type`, `system_disk`, `data_disks` are currently Computed (populated from the Describe response). To allow user input while preserving state refresh for resources created without these arguments (and for imports), they are promoted to `Optional: true, Computed: true`. This keeps backward compatibility: if the user omits them, the API defaults apply and the values are still read back into state.

### Decision 3: `host_name` is write-only (Optional, ForceNew, NOT Computed)

`DescribeDBCustomNodes` does not return `HostName`. If `host_name` were Computed, the Read would always overwrite the user's value with an empty string, causing a perpetual diff. Marking it Optional+ForceNew (without Computed) means Terraform keeps the configured value in state and never tries to reconcile it from the API. Import verification must ignore `host_name` (it cannot be repopulated).

### Decision 4: `security_group_ids` is mutable; Read via a dedicated service helper

`SecurityGroupIds` is not part of the `DBCustomNode` response. To refresh it, Read calls the new service helper `DescribeDBCustomNodeSecurityGroupsById(ctx, nodeId)`, which wraps `DescribeDBCustomNodeSecurityGroups` and returns `[]string` (the `SecurityGroupId` values from `Groups[]`). Update calls `ModifyDBCustomNodeSecurityGroups` with the full new list (set/overwrite semantic — the API replaces the node's security groups with the provided array). `d.HasChange("security_group_ids")` gates the Update call so it only fires when the list changes.

### Decision 5: Create reads new arguments with nil-guards, consistent with existing code

The existing Create builds the request field-by-field using `d.GetOk` / `helper.InterfacesHeadMap`. The new arguments follow the same pattern:
- `charge_type`, `network_mode`, `host_name`: `d.GetOk` → `helper.String`.
- `system_disk`: `helper.InterfacesHeadMap(d, "system_disk")` → build `dbdcv20201029.SystemDisk{DiskType, DiskSize}`.
- `data_disks`: iterate `d.Get("data_disks").([]interface{})` → build `[]*dbdcv20201029.DataDisk`.
- `security_group_ids`: iterate the list → `[]*string`.

### Decision 6: No SDK upgrade required

All required structs (`SystemDisk`, `DataDisk`), request fields (`ChargeType`, `NetworkMode`, `SystemDisk`, `DataDisks`, `HostName`, `SecurityGroupIds`), and client methods (`DescribeDBCustomNodeSecurityGroupsWithContext`, `ModifyDBCustomNodeSecurityGroupsWithContext`) already exist in the vendored SDK. No `go.mod`/vendor change is needed.

## Risks / Trade-offs

- **[Risk] ForceNew destroys nodes when an immutable field changes.** → Mitigation: This is the intended Terraform behavior for lifecycle-affecting inputs; descriptions will document that changes require recreation. Users are protected by `terraform plan` showing the destroy/create.
- **[Risk] `host_name` cannot be refreshed on Read/import.** → Mitigation: `host_name` is Optional+ForceNew (not Computed); the configured value stays in state. Import verification ignores `host_name`.
- **[Risk] `security_group_ids` Update overwrites the full list.** → Mitigation: This matches the `ModifyDBCustomNodeSecurityGroups` set semantic. Terraform computes the diff and sends the complete desired list, which is the correct behavior.
- **[Risk] Promoting Computed fields to Optional+Computed could surprise users who relied on the API default.** → Mitigation: Backward compatible — when omitted, nothing is sent to the API (gated by `d.GetOk`), so the API default still applies and the value is read back.
