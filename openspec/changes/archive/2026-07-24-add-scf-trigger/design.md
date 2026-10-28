## Context

The Terraform Provider for TencentCloud manages resources under `tencentcloud/services/<service>/`. For the SCF (Serverless Cloud Function) product, the SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416` is already vendored and provides four native trigger lifecycle APIs:

- `CreateTrigger` — creates a trigger bound to a function.
- `ListTriggers` — lists triggers of a function (supports `Filters` by `TriggerName`).
- `UpdateTrigger` — updates mutable trigger attributes (type, enable, qualifier, namespace, trigger_desc, description, custom_argument).
- `DeleteTrigger` — deletes a trigger.

An existing resource, `tencentcloud_scf_trigger_config`, exists but is not a full lifecycle resource: its Create path calls `UpdateTrigger` (upsert-style) and its Delete is a no-op. It also only exposes a subset of fields. This change introduces a **new, separate** resource `tencentcloud_scf_trigger` that performs true create/read/update/delete against the native APIs, following the RESOURCE_KIND_GENERAL pattern used by `tencentcloud_igtm_strategy`.

The resource does not have a single server-issued ID. A trigger is uniquely identified by the combination of `function_name`, `namespace`, and `trigger_name`, so a composite id is required.

## Goals / Non-Goals

**Goals:**
- Provide a fully-managed Terraform resource `tencentcloud_scf_trigger` with real CRUD backed by SCF APIs.
- Map all documented create/update parameters to the schema.
- Support import using the composite id (`function_name#namespace#trigger_name`).
- Follow provider conventions: `tccommon` retry/timeout helpers, nil checks before `d.Set`, composite id via `tccommon.FILED_SP`, gomonkey-based unit tests.
- Keep the change fully additive (no breaking changes to existing resources).

**Non-Goals:**
- Do not modify or deprecate the existing `tencentcloud_scf_trigger_config` resource (out of scope; backward compatibility must be preserved).
- Do not expose `Offset`/`Limit` pagination as user-facing schema fields; pagination is internal to the read/query helper.
- Do not model `ListTriggers`-only fields (`order_by`, `order`, `filters`) as top-level resource arguments (they are query helpers, not trigger attributes). They are used internally for the read path only.

## Decisions

### 1. Composite ID format: `function_name#namespace#trigger_name`

The SCF trigger APIs do not return a single unique ID. `CreateTrigger` returns a `Trigger` object whose identity is the combination of function name, namespace, and trigger name. `DeleteTrigger`/`UpdateTrigger` require `FunctionName` + `TriggerName` + `Type` + `Namespace` (+ optional `Qualifier`/`TriggerDesc`).

Decision: use `function_name + tccommon.FILED_SP + namespace + tccommon.FILED_SP + trigger_name` as the Terraform id, consistent with the existing `trigger_config` resource. The parts are parsed in Read/Update/Delete to rebuild API requests.

Alternative considered: include `type` in the id. Rejected because `DeleteTrigger` requires `Type`, but `Type` is set as `ForceNew` (cannot change after create), so it is always recoverable from state during Delete. Keeping the id format identical to the existing resource reduces import friction and inconsistency risk.

### 2. `Type` is Required + ForceNew; `FunctionName`/`TriggerName`/`Namespace` are Required + ForceNew

`CreateTrigger` requires `FunctionName`, `TriggerName`, `Type`. `UpdateTrigger` cannot change `Type` (it is part of the trigger identity). Changing any of these would mean a different trigger, so they are `ForceNew`.

`Namespace` defaults to `"default"` in the API but is part of the composite identity, so it is `Required` + `ForceNew` with a default of `"default"` (Optional with Default + ForceNew) to match the existing convention.

### 3. Mutable fields on Update

`UpdateTrigger` accepts: `FunctionName`, `TriggerName`, `Type`, `Enable`, `Qualifier`, `Namespace`, `TriggerDesc`, `Description`, `CustomArgument`. Since `FunctionName`/`TriggerName`/`Type`/`Namespace` are ForceNew, the mutable fields exercised in Update are: `enable`, `qualifier`, `trigger_desc`, `description`, `custom_argument`. The Update function checks `d.HasChange` for these mutable args and calls `UpdateTrigger` only when needed.

### 4. `enable` representation: string "OPEN"/"CLOSE"

`CreateTrigger`/`UpdateTrigger` take `Enable` as a string (`OPEN`/`CLOSE`). The read-back `TriggerInfo.Enable` is an int64 (`1`=open, `0`=close). The Read function converts the int64 to the string representation ("OPEN"/"CLOSE") to keep the schema a single `schema.TypeString`.

### 5. Read path uses `ListTriggers` with a `TriggerName` filter

There is no single GetTrigger API. Read uses `ListTriggers` filtered by `FunctionName` + `Namespace` + `Filters(TriggerName=...)`, matching the existing `DescribeScfTriggerConfigById` helper. A new service helper `DescribeScfTriggerById` mirrors this. If the returned list is empty, the resource is treated as gone (log + `d.SetId("")`).

### 6. Computed read-back fields

`TriggerInfo` returns additional metadata not settable by the user: `AvailableStatus`, `AddTime`, `ModTime`. These are exposed as `Computed` schema fields (`available_status`, `add_time`, `mod_time`) for observability. The create-time `Trigger` response (`CreateTriggerResponse.Response.TriggerInfo`) is a `Trigger` struct and is used only to confirm creation; the authoritative read is done via `ListTriggers`.

### 7. Create-time response nil check

Per provider rules, after `CreateTrigger` the code MUST verify the response is non-nil and that the `TriggerInfo` is present; if empty, return `NonRetryableError` to avoid writing a broken id.

### 8. Unit testing with gomonkey

Per project rules for newly added Terraform resources, tests use gomonkey to mock the SCF client methods (not the Terraform acceptance test framework). Run with `go test -gcflags=all=-l`.

## Risks / Trade-offs

- [Risk] Two resources (`tencentcloud_scf_trigger` and `tencentcloud_scf_trigger_config`) can manage the same underlying SCF trigger, which could cause state drift if both are used for the same trigger. → Mitigation: they have different resource type names; users should pick one. The new resource is the recommended full-lifecycle option. Documented as separate resources.
- [Risk] The `enable` field is a string ("OPEN"/"CLOSE") on write but int64 (1/0) on read, requiring conversion. → Mitigation: a single conversion point in Read maps 1→"OPEN", 0→"CLOSE"; any other value is left unset.
- [Risk] `ListTriggers` is eventually consistent; a Read immediately after Create might not find the trigger. → Mitigation: the Read helper is wrapped in `tccommon.ReadRetryTimeout` retry so transient absence is retried.
- [Trade-off] `order_by`, `order`, and `filters` from `ListTriggers` are not exposed as resource arguments because they are query controls, not trigger attributes. They are used internally only.
- [Trade-off] We do not add the deprecated fields (`ResourceId`, `BindStatus`, `TriggerAttribute`) from `TriggerInfo` to the schema.
