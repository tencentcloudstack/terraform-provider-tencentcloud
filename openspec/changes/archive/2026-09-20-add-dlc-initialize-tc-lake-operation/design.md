## Context

The DLC (Data Lake Compute) cloud API exposes an `InitializeTCLake` operation (`dlc/v20210125.InitializeTCLake`) that activates (opens) the TCLake service for the calling account. The API takes no request parameters and returns `response.Response.InstanceId` (`*string`) and `response.Response.IsSuccess` (`*bool`). This is a one-time, idempotent initialization step.

The Terraform provider already has a DLC service package (`tencentcloud/services/dlc/`) using the SDKv2 pattern and an SDK client accessor `UseDlcClient()` returning `*dlc.Client`. However, the user has explicitly requested that this resource be implemented with the **Terraform Plugin Framework action** pattern (not the SDKv2 operation resource pattern), following the existing reference implementations `NewTeoConfirmOriginAclUpdate` (`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`) and `NewBdrcRunCopyPairTasks` (`tencentcloud/services/bdrc/action_tc_bdrc_run_copy_pair_tasks.go`).

### Current State
- DLC service file: `tencentcloud/services/dlc/service_tencentcloud_dlc.go`
  - Existing `DlcService` struct holds `client *connectivity.TencentCloudClient`
  - Client accessor: `me.client.UseDlcClient()` returns `*dlc.Client`
- Framework registry: `tencentcloud/framework/registry.go`
  - `actionFactories` slice registers framework actions (currently `teo.NewTeoConfirmOriginAclUpdate`, `bdrc.NewBdrcRunCopyPairTasks`)
  - The `dlc` services package is NOT yet imported in `registry.go`
- Reference action implementations use `fw.ActionWithConfigure` (from `tencentcloud/framework/fw/action_with_configure.go`), which embeds `withMeta` and provides the `Client()` accessor plus a `Configure` method that reads `*sharedmeta.ProviderMeta`.

### Framework Action Constraint (Critical)
The Terraform Plugin Framework `action.Action` schema attributes **cannot be Computed** — `StringAttribute.IsComputed()` always returns `false` (see `vendor/.../action/schema/string_attribute.go`). The `action.InvokeResponse` struct only carries `Diagnostics` and a `SendProgress` function; it has no mechanism to write output state back to Terraform. Consequently, even though the `InitializeTCLake` response returns `InstanceId` and `IsSuccess`, these cannot be exposed as Terraform output attributes. The reference implementations confirm this: they invoke an API and surface only success/error diagnostics. The response fields will therefore be logged via `log.Printf` for observability but not persisted as Terraform state.

## Goals / Non-Goals

**Goals:**
- Add a framework action `tencentcloud_dlc_initialize_tc_lake` that calls the DLC `InitializeTCLake` API with an empty request.
- Follow the reference implementations (`NewTeoConfirmOriginAclUpdate`, `NewBdrcRunCopyPairTasks`) precisely for the action struct, factory, model, metadata, schema, and invoke flow.
- Wrap the API call with `tccommon.WriteRetryTimeout` retry and `tccommon.RetryError` error wrapping in the service layer.
- Register the action factory in `tencentcloud/framework/registry.go`'s `actionFactories` and add the `dlc` import.
- Provide gomonkey-based unit tests (no Terraform acceptance test suite) covering success, API error, client-not-configured, and schema/metadata validation — mirroring `action_tc_bdrc_run_copy_pair_tasks_test.go`.
- Provide a `.md` documentation file following the reference action docs format.

**Non-Goals:**
- Do NOT implement this as an SDKv2 operation resource (`resource_tc_dlc_initialize_tc_lake_operation.go`) — the user explicitly requested the framework action pattern.
- Do NOT expose `instance_id` / `is_success` as Terraform output attributes (framework actions do not support Computed/output attributes and `InvokeResponse` has no output mechanism).
- Do NOT register the resource in `tencentcloud/provider.go`'s `ResourcesMap` — framework actions are registered only in `tencentcloud/framework/registry.go`.
- Do NOT add a Read/Update/Delete lifecycle (actions only implement `Invoke`).
- Do NOT poll any Read interface afterwards — `InitializeTCLake` is a synchronous operation, not an asynchronous one.
- Do NOT modify any existing DLC resource, data source, or schema.

## Decisions

### 1. Implementation Pattern: Framework Action (not SDKv2 operation)
**Decision**: Implement `tencentcloud_dlc_initialize_tc_lake` as a Terraform Plugin Framework `action.Action`, following `NewTeoConfirmOriginAclUpdate` / `NewBdrcRunCopyPairTasks`.

**Rationale**:
- The user explicitly requested the framework action pattern and pointed at these two reference functions.
- Framework actions are the correct primitive for one-time operations that do not persist state, which matches the `InitializeTCLake` semantics exactly.

**Alternatives Considered**:
- SDKv2 operation resource (`resource_tc_dlc_initialize_tc_lake_operation.go` with Create/Read/Delete no-ops and a token id): This is the pattern used by `resource_tc_teo_confirm_origin_acl_update_operation.go`. Rejected because the user explicitly asked for the framework action pattern.

### 2. Output Attribute Handling
**Decision**: Do NOT expose `instance_id` / `is_success` as schema attributes. Log the response via `log.Printf` inside the service method (consistent with existing DLC service methods) and return the response to the Invoke handler, which only acts on errors.

**Rationale**:
- Framework action schema attributes cannot be `Computed` (`StringAttribute.IsComputed()` is hardcoded `false`).
- `action.InvokeResponse` only contains `Diagnostics` and `SendProgress` — there is no output/state write path.
- Both reference implementations invoke their API and surface only diagnostics; neither exposes response fields as attributes.
- The mapping in the requirement lists `instance_id`/`is_success` as cloud API output fields; these are captured in the Go response struct and logged, fulfilling observability without violating framework constraints.

**Alternatives Considered**:
- Define `instance_id`/`is_success` as `Optional` attributes and attempt to write them: Rejected — there is no response state-writing API on `InvokeResponse`, and `Optional`-only attributes without a write path would always be null, misleading users.

### 3. Schema Definition
**Decision**: The action schema defines NO input attributes (the `InitializeTCLakeRequest` is empty). The schema `Attributes` map is empty `map[string]schema.Attribute{}` with a `Description`.

**Rationale**:
- The cloud API request `InitializeTCLakeRequestParams` has no fields.
- An empty-attributes schema is valid for framework actions (the action simply takes no configuration).

### 4. Service Layer Method Signature
**Decision**: Add `func (me *DlcService) InitializeTCLake(ctx context.Context) (response *dlc.InitializeTCLakeResponse, errRet error)` to `service_tencentcloud_dlc.go`.

**Rationale**:
- No request parameters are needed, so the method takes only `ctx`.
- Returns the response pointer so the Invoke handler could (if ever needed) inspect it; primarily the response is logged inside the method's retry block (matching the existing `log.Printf("[DEBUG]...")` pattern).
- Uses `tccommon.WriteRetryTimeout` + `resource.Retry` + `tccommon.RetryError(e)` and `ratelimit.Check(request.GetAction())`, matching `ConfirmOriginACLUpdate` / `RunCopyPairTasks`.

### 5. Action Factory & Registration
**Decision**: Factory `NewDlcInitializeTCLake() action.Action` returns `&DlcInitializeTCLake{}`. Register `dlc.NewDlcInitializeTCLake` in `tencentcloud/framework/registry.go`'s `actionFactories` and add `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"` to the import block.

**Rationale**:
- Mirrors the exact registration pattern of the two reference actions.
- Framework actions are NOT registered in SDKv2 `provider.go`.

### 6. Naming
**Decision**:
- Type name: `DlcInitializeTCLake`
- Factory: `NewDlcInitializeTCLake`
- Model type: `DlcInitializeTCLakeModel` (empty struct, kept for symmetry with references)
- Action type name (Metadata): `tencentcloud_dlc_initialize_tc_lake`
- Files: `action_tc_dlc_initialize_tc_lake.go`, `action_tc_dlc_initialize_tc_lake_test.go`, `action_tc_dlc_initialize_tc_lake.md`

**Rationale**:
- Follows the `NewTeoConfirmOriginAclUpdate` / `NewBdrcRunCopyPairTasks` naming convention (action files prefixed with `action_tc_`).
- The action type name matches the requested resource name `tencentcloud_dlc_initialize_tc_lake`.

### 7. Test Strategy
**Decision**: Use gomonkey mocks (no Terraform acceptance test suite), mirroring `action_tc_bdrc_run_copy_pair_tasks_test.go`.

**Rationale**:
- Project guidelines require gomonkey-based unit tests for new framework actions (the reference test file uses gomonkey to mock `RunCopyPairTasksWithContext`).
- Tests cover: successful invoke (mock returns response, assert no diagnostics), API error (mock returns error, assert diagnostic contains API error), client-not-configured (nil client guard), and metadata/schema validation (TypeName + empty attributes).

## Risks / Trade-offs

- **[Risk] Users expect `instance_id`/`is_success` as outputs** → Mitigation: Document clearly in the `.md` and design that framework actions do not expose output attributes; the API response is logged for observability. The action's value is triggering the one-time initialization, not reading back state.
- **[Risk] Empty schema may confuse practitioners** → Mitigation: The `.md` example shows the `action` block with an empty `config {}` and a description explaining it takes no parameters. This is acceptable for a parameterless one-time operation.
- **[Risk] `InitializeTCLake` may fail if TCLake is already initialized** → Mitigation: The service layer retry handles transient errors via `tccommon.RetryError`; a non-retryable "already initialized" error surfaces as a diagnostic, which is the correct behavior for an idempotent one-time action.
- **[Risk] Forgetting to add the `dlc` import in `registry.go`** → Mitigation: The tasks list explicitly calls out editing `registry.go` to add both the import and the factory entry.