## Context

The `tencentcloud_serverless_hbase_instance` resource (in `tencentcloud/services/emr/resource_tc_serverless_hbase_instance.go`) manages EMR Serverless HBase instances via the EMR v20190103 SDK.

Current `Update` behavior:
- An `immutableArgs` slice (`["pay_mode", "disk_type", "disk_size", "node_type", "time_span", "time_unit", "auto_renew_flag"]`) rejects any change to these fields.
- `instance_name` changes → `ModifySLInstanceBasic`.
- `tags` changes → tag service diff.
- `zone_settings.*.node_num` changes → `ModifySLInstanceWithContext` (uses `ModifySLInstance` API, sets `InstanceId`/`Zone`/`NodeNum`), then waits for status to reach `"2"` via `SLInstanceStateRefreshFunc`.

The cloud API `ModifySLInstance` (https://cloud.tencent.com/document/api/589/110869) actually accepts `InstanceId`, `Zone`, `NodeNum`, `ClientToken`, and `DiskSize`. Its `DiskSize` field is documented as: "存储空间大小，需要是100的倍数，且不允许比当前存储空间小（只允许扩容）, 单位：GB". So storage expansion is supported but currently blocked in Terraform by `immutableArgs`.

`ModifySLInstance` is asynchronous: after the call the instance enters the 变配中 (scaling) state, which the cloud API reports as `Status == 22`. The instance returns to 运行中 (running) with `Status == 2` once scaling completes. Status is queried through `DescribeSLInstance` / the `DescribeSLInstanceList`-based `GetSLInstanceStatus` helper that backs `SLInstanceStateRefreshFunc` (https://cloud.tencent.com/document/product/589/110871).

Constraints:
- Must remain backward compatible — `disk_size` stays `Required` `TypeInt`; only its lifecycle moves from immutable to updatable.
- Must follow project rules: retry cloud API calls with `tccommon.WriteRetryTimeout`, wrap errors with `tccommon.RetryError`, log with `logId`, use the resource name `serverless_hbase_instance` in messages.
- No `go build`/`go vet`/`golint` during generation; `gofmt`/`make doc`/`.changelog` only in the finalize phase.

## Goals / Non-Goals

**Goals:**
- Allow in-place expansion of `disk_size` on an existing `tencentcloud_serverless_hbase_instance` by calling `ModifySLInstance` with `InstanceId` + `DiskSize`.
- Wait for the async scaling to complete by polling instance status from `22` (变配中) to `2` (运行中) using the existing `SLInstanceStateRefreshFunc` + `tccommon.BuildStateChangeConf`.
- Keep the change isolated and non-conflicting with the existing `zone_settings.node_num` expansion path.

**Non-Goals:**
- Supporting disk shrink (the cloud API explicitly forbids it: 只允许扩容).
- Changing `disk_type`, `pay_mode`, `node_type`, `time_span`, `time_unit`, or `auto_renew_flag` mutability — they remain immutable.
- Modifying the resource schema shape (no new/removed fields).
- Adding a `Timeouts` schema block (the resource uses `tccommon.ReadRetryTimeout`-based state change conf and has no existing `Timeouts` block; introducing one is out of scope and not required for a backward-compatible change).

## Decisions

### Decision 1: Handle `disk_size` as a separate change branch before the `zone_settings` node_num block

**Choice:** Add an independent `if d.HasChange("disk_size")` branch in `Update` that calls `ModifySLInstance` with only `InstanceId` and `DiskSize`, then waits for status `2`. Remove `disk_size` from `immutableArgs`.

**Rationale:** The `ModifySLInstance` API fields `NodeNum`/`Zone` and `DiskSize` are independent optional parameters. The existing `node_num` loop reuses `ModifySLInstance` per-zone; mixing a global `DiskSize` into that per-zone loop would be incorrect (DiskSize is not zone-scoped). A dedicated branch keeps concerns separated and matches the API's actual semantics.

**Alternatives considered:**
- Reuse the existing `needChange`/`zone_settings` loop and append `DiskSize` to the request: rejected because `DiskSize` is instance-level, not per-zone, and the loop is driven by `zone_settings.%d.node_num` change detection which has no relation to `disk_size`.
- Add a new SDK wrapper method in the service file: rejected as unnecessary — `ModifySLInstanceWithContext` is already called directly in the resource file for `node_num`, so following the same inline pattern keeps consistency.

### Decision 2: Poll status with `22` as a pending state and `2` as the target

**Choice:** Reuse `SLInstanceStateRefreshFunc` (which returns the instance `Status` as a string) and `tccommon.BuildStateChangeConf` with:
- `pending`: `["22"]` (变配中)
- `target`: `["2"]` (运行中 / 变配完成)
- `timeout`: `10 * tccommon.ReadRetryTimeout` (same as the create and `node_num` wait)
- `minInterval`: `time.Second` (same as existing usage)

**Rationale:** This mirrors exactly how the resource already waits on `ModifySLInstance` for `node_num` expansion (lines 488–494) and for create (lines 260–263), minimizing risk and code divergence. The cloud API states Status `22` = 变配中, `2` = 运行中.

**Alternatives considered:**
- Adding `22` to a global `failStates`: rejected — `22` is a normal transient state, not a failure.
- Writing a new refresh function: rejected — the existing one returns the raw status string, which is sufficient.

### Decision 3: Retry semantics — retry only the API call, set state / wait outside the retry block

**Choice:** Per project rules, the `ModifySLInstance` call is wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; the state-change wait (`WaitForState`) runs **after** the retry block returns nil, matching the existing `node_num` pattern (lines 474–494).

**Rationale:** Project convention: do not set state or perform non-API operations inside the retry block; don't nest retries. The status wait uses its own `BuildStateChangeConf` polling, so no double retry.

### Decision 4: Validate Serverless HBase instance context naming in logs

**Choice:** Use `serverless_hbase_instance` in all log/error messages for this change (e.g., `[CRITAL]%s update serverless hbase instance disk_size failed`), consistent with existing log strings in the file.

**Rationale:** Project rule: use the resource's snake-case name in logs, not generic phrasing, for searchability.

## Risks / Trade-offs

- [Risk] The cloud API only allows disk expansion (新的 DiskSize 必须 >= 当前值). If a user lowers `disk_size`, `ModifySLInstance` will return a cloud-side error → Mitigation: surface the cloud API error via `tccommon.RetryError`; no Terraform-side pre-validation is added to avoid duplicating cloud constraints, but the error message will be clear. (Optionally a defensive `d.GetChange("disk_size")` check can be added, but matching the API's own validation is preferred.)
- [Risk] Concurrent change of `disk_size` and `zone_settings.node_num` in a single `terraform apply` → Mitigation: handle `disk_size` first (its own branch), then the `node_num` loop; both re-enter `WaitForState`, so ordering is safe. The two do not share request state.
- [Risk] Status polling timeout if scaling takes longer than `10 * ReadRetryTimeout` → Mitigation: same timeout budget already used for create and node_num expansion; acceptable and consistent with current behavior.
- [Trade-off] Not adding a `Timeouts` schema block: keeps backward compatibility and matches the existing resource design; users rely on the built-in `ReadRetryTimeout`-based wait.