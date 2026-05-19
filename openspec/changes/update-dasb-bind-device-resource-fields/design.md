## Context

The resource `tencentcloud_dasb_bind_device_resource` wraps the DASB `BindDeviceResource` API and `DescribeDevices` API. The current Go SDK vendored in this repository (`dasb/v20191018`) only exposes three fields in `BindDeviceResourceRequest`: `DeviceIdSet`, `ResourceId`, and `DomainId`. The API documentation now defines six additional fields for K8S cluster managed account support. These fields are absent from the vendored SDK structs, so any code referencing them would fail to compile.

The `DescribeDevices` response `Device` struct already contains `DomainId` and `DomainName` fields in the SDK, but the Read function currently has two bugs:
1. `domain_id` is set inside the per-device loop (overwrites on each iteration) instead of being set once from the first applicable device.
2. `device_id_set` items are appended as `*uint64` pointers rather than dereferenced `uint64` / `int` values, causing Terraform state type mismatch.

## Goals / Non-Goals

**Goals:**
- Fix two existing bugs in the Read function (`domain_id` loop overwrite, `device_id_set` pointer dereference).
- Add `domain_name` as a Computed read-back field from `DescribeDevices`.
- Define schema stubs for the six new K8S fields (`manage_dimension`, `manage_account_id`, `manage_account`, `manage_kubeconfig`, `namespace`, `workload`) so they are ready to wire up after SDK upgrade.
- Document the SDK upgrade prerequisite clearly in code comments.

**Non-Goals:**
- Modifying the vendored SDK source code.
- Implementing the six new K8S fields in Create/Update until the SDK is upgraded.
- Adding new service-layer methods beyond what is required for the above.

## Decisions

### Decision 1: Schema stubs now, API wiring after SDK upgrade
Add the six new fields to the schema as `Optional` now (with a note in description that they require SDK upgrade), but guard their use in Create/Update behind a compile-time TODO comment. This keeps the schema forward-compatible and unblocks documentation/review, while preventing broken builds.

**Alternative considered**: Wait until SDK upgrade before any changes. Rejected because the bug fixes and `domain_name` mapping are independent of the SDK gap and should ship sooner.

### Decision 2: Fix `device_id_set` pointer dereference in Read
`item.Id` is `*uint64`. Appending it directly into `[]interface{}` stores a pointer, which causes state type mismatch when Terraform compares it against the schema's `TypeInt`. Must dereference: `int(*item.Id)`.

### Decision 3: `domain_id` set once outside loop
Because all devices returned by `DescribeDevicesByResourceId` are bound to the same resource, they share the same `DomainId`. Set it once from the first device that has a non-nil `DomainId` and break, rather than overwriting in every iteration.

### Decision 4: Code style follows `resource_tc_igtm_monitor.go`
Use `tccommon.NewResourceLifeCycleHandleFuncContext` for ctx, `resource.Retry` with `tccommon.WriteRetryTimeout` / `tccommon.ReadRetryTimeout` where applicable, and nil-guard every pointer read from API responses.

## Risks / Trade-offs

- [Risk] SDK fields missing → Mitigation: Clearly mark K8S fields with `// TODO: wire after SDK upgrade` comments; do not reference undefined struct fields.
- [Risk] Changing `device_id_set` Read behaviour could cause a one-time state diff on first plan after upgrade → Mitigation: This is a bug fix; document in CHANGELOG.
- [Risk] `DescribeDevices` with `ResourceIdSet` may return devices across multiple domain IDs if the API behaviour changes → Mitigation: Use first non-nil `DomainId` and log a warning if values differ.
