## Context

The `tencentcloud_elasticsearch_instance` resource already manages Kibana and other optional web services via the `UpdateInstance` API. The Cerebro service (an alternative to Kibana for monitoring) is also configurable through the same `UpdateInstance` API, with fields:

- `EnableCerebro` (*bool) - enable/disable Cerebro
- `CerebroPublicAccess` (*string) - public network access (`OPEN`/`CLOSE`)
- `CerebroPrivateAccess` (*string) - private network access (`OPEN`/`CLOSE`)
- `CerebroPrivateDomain` (*string) - custom private domain

These fields exist in the SDK's `UpdateInstanceRequestParams` struct but are not yet exposed in the Terraform resource. The current `UpdateInstance` service function signature is already long (14 parameters) and will need to be extended.

**Key constraint**: The `DescribeInstances` API response (`InstanceInfo` struct) does **not** return Cerebro fields. This means these parameters are write-only from the descriptor's perspective, similar to how `password` is handled (it's present in the schema but never read back from the API).

## Goals / Non-Goals

**Goals:**
- Add `enable_cerebro`, `cerebro_public_access`, `cerebro_private_access`, and `cerebro_private_domain` as optional computed fields to the resource schema
- Extend the `UpdateInstance` service function to accept and pass Cerebro parameters
- Handle Cerebro updates in the resource Update flow with retry logic and upgrade-wait pattern
- Maintain backward compatibility — existing configurations without Cerebro fields must continue to work

**Non-Goals:**
- Reading Cerebro state from the API (not available in `DescribeInstances` response)
- Adding Cerebro to the Create flow (the API does not support setting Cerebro at creation time)
- Adding Cerebro to the Delete flow (not applicable)
- Adding Cerebro-related computed output fields (e.g., Cerebro URL) — the API does not return these

## Decisions

### Decision 1: Schema fields as Optional + Computed

**Choice**: All four Cerebro fields are defined as `Optional: true, Computed: true`.

**Rationale**: Since `DescribeInstances` does not return Cerebro fields, the Read function cannot populate them. Setting `Computed: true` allows Terraform to preserve the user's configured value in state even when the Read function does not set it. This is the same pattern used for `password` in this resource.

**Alternatives considered**:
- `Optional: true, Computed: false` — would cause Terraform to detect drift on every plan, since the Read function cannot set the value.
- `Optional: true, Computed: true, ForceNew: true` — unnecessary, Cerebro can be updated on an existing instance.

### Decision 2: Extend the existing UpdateInstance service function

**Choice**: Add `enableCerebro *bool`, `cerebroPublicAccess string`, `cerebroPrivateAccess string`, `cerebroPrivateDomain string` parameters to the existing `UpdateInstance` function signature.

**Rationale**: The existing pattern is to pass all updatable fields as parameters to `UpdateInstance`, which internally constructs the `UpdateInstanceRequest` and fills in non-empty values. This keeps the update logic centralized and avoids creating many small service functions.

**Alternatives considered**:
- Create a separate `UpdateInstanceCerebro` function — adds unnecessary indirection and code duplication. The `UpdateInstance` API is a single endpoint that handles all attributes.
- Pass a struct — would require defining a new type and refactoring all existing callers, which is higher risk.

### Decision 3: No read-back of Cerebro fields

**Choice**: The Read function does NOT attempt to set Cerebro fields from the API response, since `InstanceInfo` has no such fields.

**Rationale**: This is consistent with the `password` field pattern. The user's configured values are preserved in state because the fields are `Computed`. If the user changes the value outside of Terraform (e.g., via console), Terraform will not detect the drift — this is an accepted limitation.

**Risk**: If Cerebro state is changed outside of Terraform, the next `terraform apply` will not detect or correct the drift. This is the same trade-off as `password` and is documented in the field descriptions.

### Decision 4: Update flow follows existing HasChange pattern

**Choice**: Add a new `if d.HasChange(...)` block in the Update function for the Cerebro fields, following the same pattern as `kibana_public_access`, `kibana_private_access`, etc.

**Rationale**: Consistency with existing code reduces maintenance burden and makes the code reviewable. Each Cerebro field change is handled independently with its own retry and upgrade-wait logic.

## Risks / Trade-offs

- **[Risk] No read-back**: Cerebro state cannot be reconciled from the API. If changed outside Terraform, state will be stale. → **Mitigation**: Document this limitation in the field descriptions. Same pattern as `password`.

- **[Risk] Function signature bloat**: The `UpdateInstance` function already has 14 parameters. Adding 4 more makes it 18, which is unwieldy. → **Mitigation**: Accept this for now. A future refactoring could switch to a struct-based approach, but that's out of scope for this change.

- **[Risk] Terraform detects drift on Cerebro access fields when `enable_cerebro` is false**: If Cerebro is disabled, the access fields (`cerebro_public_access`, `cerebro_private_access`) are meaningless. → **Mitigation**: The `Computed` nature of these fields means Terraform won't flag them as changed if they're not in the user's config. If they ARE in the user's config, the update will be attempted (and the API may reject or ignore it).