## Context

The Tencent Cloud GA2 product line is a multi-tier object model:

```
GlobalAccelerator (parent instance, holds tags + cross-border config)
└── Listener                 ← THIS CHANGE
    └── EndpointGroup        ← already shipped (tencentcloud_ga2_endpoint_group)
        └── EndpointConfigurations
```

The provider already ships `tencentcloud_ga2_endpoint_group` and is shipping `tencentcloud_ga2_global_accelerator`. The middle tier — Listener — must be added so users can construct the full chain in HCL with no console hand-off.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` already exposes:
- `CreateListenerWithContext` → returns `{ TaskId, ListenerId }` (asynchronous)
- `DescribeListenersWithContext` → paged list with `Filters`, requires `GlobalAcceleratorId`
- `ModifyListenerWithContext` → returns `{ TaskId }` (asynchronous)
- `DeleteListenerWithContext` → returns `{ TaskId }` (asynchronous)
- `DescribeTaskResultWithContext` → returns `{ Status }`, used as the polling oracle

The existing `Ga2Service` already provides `WaitForGa2TaskFinish(ctx, taskId, timeout)`; we will reuse it verbatim.

## Goals / Non-Goals

**Goals:**
- Provide full lifecycle management of a GA2 Listener through Terraform.
- Schema fields exactly mirror `CreateListenerRequest` (no field renaming, no synthetic flags), per the user's explicit rule.
- All async writes wait for `Status == SUCCESS` on the returned `TaskId` before returning to Terraform, so dependent resources (`tencentcloud_ga2_endpoint_group`) can immediately reference the new listener.
- Code style matches `tencentcloud_igtm_monitor` and the previously-shipped GA2 resources (single-file resource layout, retry on every SDK call, defensive nil checks on response payloads).
- Filename conventions for the markdown doc and `_test.go` mirror `resource_tc_config_compliance_pack.md` / `_test.go`.

**Non-Goals:**
- Forwarding rules / HTTP-rule resource (separate change if needed).
- A `tencentcloud_ga2_listeners` datasource — resource-only here.
- Any field that exists only in the describe response but not in `CreateListener` is exposed Computed-only, never input.

## Decisions

### D1. Composite resource ID = `<GlobalAcceleratorId>#<ListenerId>`
Why: `DescribeListeners` / `Modify*` / `Delete*` all require **both** `GlobalAcceleratorId` and `ListenerId`. Persisting only the listener ID in `d.Id()` would force a re-lookup every time, which is wasteful and brittle.
Implementation: use the project-standard separator `tccommon.FILED_SP` (already used by `tencentcloud_ga2_endpoint_group`).
Alternative considered: keep ID as bare `ListenerId` and source `GlobalAcceleratorId` from the schema. Rejected because import then requires the user to also set `global_accelerator_id` post-import, which the standard `ImportStatePassthrough` cannot do.

### D2. Reuse `WaitForGa2TaskFinish` as-is
It already accepts a caller-supplied `timeout time.Duration`, polls `DescribeTaskResult`, and treats `SUCCESS` as terminal. No listener-specific behavior is needed.

### D3. Add `DescribeGa2ListenerById` to the existing `Ga2Service`
Pattern matches `DescribeGa2EndpointGroupById` and the new `DescribeGa2GlobalAcceleratorById`:
- Build the request **outside** the for-loop (only `Offset` / `Limit` mutate per page).
- `request.GlobalAcceleratorId` is set once outside the loop.
- `Filters` is set to `[{Name: "listener-id", Values: [listenerId]}]`.
- Page size = `100` (the documented maximum), passed as a literal — no new package-level constant.
- Strict-equals on `*item.ListenerId == listenerId` (and `*item.GlobalAcceleratorId == gaId` for symmetry) before returning.
- Returns `(nil, nil)` when not found; the resource layer treats this as "deleted out of band" and calls `d.SetId("")`.

### D4. Async retry topology
Every SDK call is wrapped in `resource.Retry(timeoutScope, func() *resource.RetryError { ... })`:
- Read paths: `tccommon.ReadRetryTimeout`.
- Write paths (Create / Modify / Delete): `tccommon.WriteRetryTimeout`.
- Async polling (after the write succeeds): `d.Timeout(schema.TimeoutCreate|Update|Delete)` passed into `WaitForGa2TaskFinish`.
- Resource-level `Timeouts` block defaults to **5 minutes** for Create/Update/Delete (matches the other GA2 resources).

### D5. Schema parity with `CreateListenerRequest`
Mapping (every CreateListener input field appears, no extras, no renames):

| Schema field | Type | Required? | ForceNew? | Source SDK field | Notes |
|---|---|---|---|---|---|
| `global_accelerator_id` | `TypeString` | Required | **Yes** | `GlobalAcceleratorId` | Cannot move a listener between accelerators |
| `name` | `TypeString` | Optional+Computed | No | `Name` | ≤60 bytes |
| `port_ranges` | `TypeList`, MaxItems=1 | Required | **Yes** | `PortRanges{FromPort, ToPort}` | `ModifyListener` does not accept it |
| `description` | `TypeString` | Optional+Computed | No | `Description` | ≤100 bytes |
| `listener_type` | `TypeString` | Optional+Computed | **Yes** | `ListenerType` | `ModifyListener` does not accept it |
| `protocol` | `TypeString` | Optional+Computed | **Yes** | `Protocol` | `ModifyListener` does not accept it; default TCP |
| `idle_timeout` | `TypeInt` | Optional+Computed | No | `IdleTimeout` | seconds |
| `get_real_ip_type` | `TypeString` | Optional+Computed | No | `GetRealIpType` | TOA / ProxyProtocol |
| `client_affinity` | `TypeString` | Optional+Computed | No | `ClientAffinity` | session stickiness toggle |
| `request_timeout` | `TypeInt` | Optional+Computed | No | `RequestTimeout` | |
| `x_forwarded_for_real_ip` | `TypeBool` | Optional+Computed | No | `XForwardedForRealIp` | L7 only |
| `certification_type` | `TypeString` | Optional+Computed | No | `CertificationType` | UNIDIRECTIONAL / MUTUAL |
| `cipher_policy_id` | `TypeString` | Optional+Computed | No | `CipherPolicyId` | TLS cipher pack |
| `server_certificates` | `TypeSet`, Elem=string | Optional+Computed | No | `ServerCertificates []*string` | TypeSet — order is not semantic |
| `client_ca_certificates` | `TypeSet`, Elem=string | Optional+Computed | No | `ClientCaCertificates []*string` | TypeSet — order is not semantic |
| `client_affinity_time` | `TypeInt` | Optional+Computed | No | `ClientAffinityTime` | **NOTE:** present only on `ModifyListener` (not on `CreateListener`); ignored on Create, applied on Update |

Computed-only fields (not in CreateListenerRequest, surfaced from `ListenerSet`):
- `listener_id` (string) — also stored as the second segment of `d.Id()`.
- `http_version` (string)
- `create_time` (string)
- `status` (string)
- `endpoint_group_counts` (int)

`port_ranges` is modeled as a `TypeList` with `MaxItems: 1` and a nested block exposing `from_port`/`to_port` (both `TypeInt`, Required) — this is the standard provider idiom for single-instance nested objects.

### D6. ForceNew choices justified by the API
- `global_accelerator_id`: a listener belongs to exactly one accelerator; moving it requires recreate.
- `port_ranges`: `ModifyListener` has no `PortRanges` field.
- `listener_type`: `ModifyListener` has no `ListenerType` field.
- `protocol`: `ModifyListener` has no `Protocol` field.

All other input fields are updatable in place via `ModifyListener`.

### D7. `client_affinity_time` asymmetry
`CreateListenerRequest` does **not** carry `ClientAffinityTime` (only `ModifyListenerRequest` does). The schema therefore exposes it as an Optional+Computed field; the value is silently ignored on Create (and the API will fall back to its own default, then surface it via Read), and forwarded only on Update if the user changes it. This is documented in the spec to avoid surprise.

### D8. Server / client CA certificate lists use `TypeSet`
Both `ServerCertificates` and `ClientCaCertificates` are unordered SDK string slices. Using `TypeList` would create spurious diffs on rotation. We model them as `TypeSet` with `Elem: schema.TypeString`, and convert via `(*schema.Set).List()` when building requests, mirroring the `endpoint_configurations` decision in `tencentcloud_ga2_endpoint_group`.

### D9. File layout
Single file: `resource_tc_ga2_listener.go` (schema + Create/Read/Update/Delete + build/flatten helpers). Service-level helper lives in the existing `service_tencentcloud_ga2.go`. This matches the user's strict feedback during the endpoint-group change to avoid `_crud.go` / `_helpers.go` splits.

### D10. `make doc` flow + `provider.md` registration
Per the previously-established workflow, the resource markdown lives at `tencentcloud/services/ga2/resource_tc_ga2_listener.md`. The website file at `website/docs/r/ga2_listener.html.markdown` is **never** hand-edited; it is regenerated by `make doc`. For `make doc` to discover the new resource, we must also append `tencentcloud_ga2_listener` to the `Global Accelerator(GA2)` Resources section in `tencentcloud/provider.md` (the catalog used by `gendoc`).

## Risks / Trade-offs

- **[Risk]** `client_affinity_time`'s asymmetric input availability (Modify-only) may confuse users who set it on Create. → **Mitigation**: ignore the value on Create; document the rule in the resource markdown; the value will reconcile on the first Update or via Computed reads.
- **[Risk]** Both `client_affinity_time` and the `protocol`/`listener_type` ForceNew constraints are inferred from the SDK schema, not from the user-supplied API doc text. → **Mitigation**: SDK Request/Response definitions are authoritative; this is the same approach used by every other resource in the provider.
- **[Risk]** Some L7-only fields (`request_timeout`, `x_forwarded_for_real_ip`, `certification_type`, `cipher_policy_id`, `server_certificates`, `client_ca_certificates`, `client_ca_certificates`, `http_version`) are silently ignored when the listener `protocol` is L4 (TCP/UDP). → **Mitigation**: all are Optional; rely on server-side validation to reject misuse, rather than encoding L4/L7 conditional logic in Schema (which historically led to over-engineered `CustomizeDiff` blocks).
- **[Trade-off]** `port_ranges` ForceNew means changing the listening port range destroys and recreates the listener. This matches the API's hard constraint and is explicitly documented.

## Migration Plan

This is purely additive. No state migration required:
1. Land the new resource + service helper + provider registration on the `feat/ga2_nr` branch.
2. After release, users opt in by adding `resource "tencentcloud_ga2_listener" "x" { ... }` to their config.
3. Existing `tencentcloud_ga2_endpoint_group` resources continue to reference `listener_id` exactly as before.

Rollback: pure revert of the new files + the `provider.go` and `provider.md` registration lines; no state mutations to undo.

## Open Questions

- None blocking. All required SDK surface area is vendored; the spec-level decisions (composite ID, ForceNew set, async polling reuse, TypeSet for cert lists) are determinate.
