## Context

The `tencentcloud_scf_function` Terraform resource currently does not expose the `InstanceConcurrencyConfig` parameter. This parameter is already supported by the TencentCloud SCF API in three relevant endpoints:

- **CreateFunction**: Accepts `InstanceConcurrencyConfig` in `CreateFunctionRequest` (SDK: `models.go:697`)
- **UpdateFunctionConfiguration**: Accepts `InstanceConcurrencyConfig` in `UpdateFunctionConfigurationRequest` (SDK: `models.go:5503`)
- **GetFunction**: Returns `InstanceConcurrencyConfig` in `GetFunctionResponseParams` (SDK: `models.go:2524`)

The `UpdateFunctionCode` API does NOT accept this parameter, which aligns with expectations — it only deals with code changes.

The existing resource already follows a proven pattern for adding complex struct parameters: the `dns_cache` and `intranet_config` parameters were recently added following the same approach — storing values in `scfFunctionInfo`, passing them through `CreateFunctionRequest` and `UpdateFunctionConfigurationRequest`, and reading them back from `GetFunctionResponse`.

The `InstanceConcurrencyConfig` SDK struct contains:
```go
type InstanceConcurrencyConfig struct {
    DynamicEnabled            *string
    MaxConcurrency            *uint64
    InstanceIsolationEnabled  *string
    Type                      *string
    MixNodeConfig             []*MixNodeConfig
    SessionConfig             *SessionConfig
}
```

Where `MixNodeConfig` has `NodeSpec *string` and `Num *uint64`, and `SessionConfig` has `SessionSource *string`, `SessionName *string`, `MaximumConcurrencySessionPerInstance *uint64`, `MaximumTTLInSeconds *uint64`, `MaximumIdleTimeInSeconds *uint64`.

## Goals / Non-Goals

**Goals:**
- Add `instance_concurrency_config` (TypeList, Optional, Computed: false) to the `tencentcloud_scf_function` resource schema
- Include all SDK sub-fields: `dynamic_enabled`, `max_concurrency`, `instance_isolation_enabled`, `type`, `mix_node_config`, `session_config`
- Wire the parameter through Create (CreateFunction), Read (GetFunction), and Update (UpdateFunctionConfiguration) operations
- Maintain full backward compatibility — existing configurations continue working without changes

**Non-Goals:**
- No changes to `UpdateFunctionCode` API call (it does not support this parameter)
- No changes to trigger-related APIs (CreateTrigger, DeleteTrigger) — these are out of scope
- No changes to `DeleteFunction` API call
- No modification of any existing schema fields

## Decisions

### Decision 1: Use TypeList (not TypeSet) for the top-level parameter
**Rationale**: Follow the existing pattern used by `cfs_config` and `image_config` in the same resource. TypeList preserves order and is the Terraform convention for complex nested objects that represent a single configuration block (not a collection of independent items). The `instance_concurrency_config` is a single configuration block, not a set of independent items.

### Decision 2: Store as a single-element TypeList
**Rationale**: Consistent with the existing `cfs_config` pattern. The API returns a single `InstanceConcurrencyConfig` object, not an array. Using TypeList with MaxItems=1 would be ideal, but to stay consistent with the existing patterns in this resource (which don't enforce MaxItems), we follow the same approach.

### Decision 3: Wire `instance_concurrency_config` in Create and UpdateFunctionConfiguration only
**Rationale**: The `UpdateFunctionCode` API does not accept `InstanceConcurrencyConfig`. Configuration changes should go through `UpdateFunctionConfiguration`, while code changes go through `UpdateFunctionCode`. This is the same pattern used by `intranet_config` and `dns_cache`.

### Decision 4: All sub-fields are Optional with no Default
**Rationale**: The API treats all sub-fields as optional (`omitempty`). Users should be able to specify only the fields they need. No default values are set to avoid unintended API side effects.

## Risks / Trade-offs

- **Risk**: If user specifies `instance_concurrency_config` with invalid combinations of sub-fields, the API will reject the request
  - **Mitigation**: The API returns clear error messages; no additional client-side validation is needed beyond what the SDK type system provides

- **Risk**: The `mix_node_config` and `session_config` sub-structs are complex and rarely used
  - **Mitigation**: Making all sub-fields Optional means users can omit them entirely, keeping simple configurations clean
