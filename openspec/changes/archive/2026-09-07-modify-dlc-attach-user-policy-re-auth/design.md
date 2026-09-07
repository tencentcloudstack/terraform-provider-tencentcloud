## Context

The `tencentcloud_dlc_attach_user_policy_attachment` resource binds a single DLC authorization policy to a user via the `AttachUserPolicy` API and reads it back via `DescribeUserInfo`. Today the `policy_set.re_auth` schema field is declared `Computed: true` only, so users cannot set it; the value is solely populated from the API response on read. The DLC `Policy` struct in the vendored SDK (`dlc/v20210125/models.go`) defines `ReAuth *bool` with the semantics "whether the grantee is allowed to further grant the permissions (default `false`)". Because `AttachUserPolicyRequest.PolicySet` accepts full `Policy` objects, `ReAuth` is already an accepted input parameter on the create API — the provider simply never forwards it.

This change makes `re_auth` user-settable on create while preserving the read-back behavior.

## Goals / Non-Goals

**Goals:**
- Allow users to configure `policy_set.re_auth` (bool) when creating the resource, so the grantee can be allowed to re-grant permissions.
- Forward the configured `re_auth` value to the `AttachUserPolicy` API request `PolicySet[0].ReAuth`.
- Keep the field readable from the API response (`DescribeUserInfo`) so state stays accurate.
- Maintain full backward compatibility: existing configs without `re_auth` continue to work with the API default (`false`).

**Non-Goals:**
- No in-place Update support is added. The resource remains `RESOURCE_KIND_ATTACHMENT` (Create/Read/Delete only); any change to `re_auth` triggers recreation, consistent with the existing `ForceNew` behavior on `policy_set`.
- No new API endpoints or SDK upgrades.
- No changes to other `policy_set` sub-fields.

## Decisions

### Decision 1: Make `re_auth` `Optional: true, Computed: true` rather than `Required`

**Rationale**: The API treats `ReAuth` as optional with a default of `false`. Marking it `Optional` + `Computed` lets users opt in to delegation while keeping existing configurations valid (no `re_auth` → API default applies). `Computed` is retained so the read handler still populates it from the API response, avoiding drift when the API omits or defaults the value.

**Alternatives considered**:
- `Required`: Would break every existing configuration that does not set `re_auth`. Rejected.
- `Optional` only (drop `Computed`): Would cause perpetual diffs because the read handler sets the value from the API response but Terraform would treat any non-configured value as "should be zero". Keeping `Computed` preserves the current read-back semantics.

### Decision 2: Forward `re_auth` in the Create handler only

**Rationale**: The resource has no Update handler (it is an attachment resource where any top-level change recreates). The Create handler iterates the `policy_set` list and builds `dlc.Policy` objects; adding a `re_auth` mapping there is the single place that needs to change. The Delete handler (`DetachUserPolicy`) does not take a `Policy` body, and the Read handler already flattens `ReAuth` from the response via `flattenDlcAttachUserPolicyAttachmentPolicySet`, so no change is needed on read.

**Alternatives considered**:
- Mirror into an Update handler: There is no Update handler, so not applicable.

### Decision 3: Use `d.GetOkExists` semantics via `dMap["re_auth"]` consistent with sibling fields

**Rationale**: The existing Create handler reads each sub-field from the `dMap` (`item.(map[string]interface{})`) using `if v, ok := dMap["<field>"]; ok`. For booleans, Terraform's SDK returns the value through the map regardless of whether it was set (defaults to `false`). To match the existing pattern in this handler and to send an explicit `ReAuth` when the user sets it, the create handler maps `dMap["re_auth"]` to `policy.ReAuth = helper.Bool(v.(bool))`. This is consistent with how other boolean-capable fields are handled across the provider's attachment resources and ensures the user's intent (`true`) is sent. When the user omits it, the SDK default `false` is forwarded, which matches the API default — acceptable and non-breaking.

**Alternatives considered**:
- Skip sending `ReAuth` when false: Would require distinguishing "unset" from "explicitly false", which Terraform SDK v2 does not expose cleanly via the map for booleans. Given the API default is also `false`, sending `false` is harmless.

## Risks / Trade-offs

- **[Risk] Sending `ReAuth=false` for users who omitted it could differ from API default in edge cases** → Mitigation: The API default is documented as `false`, so sending `false` is equivalent to the default; behavior is unchanged for existing configurations.
- **[Risk] State drift if API returns a different `ReAuth` than configured** → Mitigation: The field remains `Computed`, so the read handler overrides state from the API response, reconciling any difference. Because `policy_set` is `ForceNew`, a detected difference after read would surface as a recreate on the next plan, which is the existing contract for this attachment resource.
- **[Trade-off] No in-place update for `re_auth`** → Accepted; consistent with the resource's attachment lifecycle (recreate on change).
