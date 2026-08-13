## Context

The `tencentcloud_elasticsearch_instance` resource (`tencentcloud/services/es/resource_tc_elasticsearch_instance.go`) manages the full lifecycle of a Tencent Cloud Elasticsearch Service (ES) instance via the ES cloud APIs in package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416`.

The ES cloud API exposes a cluster destroy-protection switch through two endpoints:

1. **`UpdateInstance`** — accepts an input field `EnableDestroyProtection` (`*string`, values `OPEN`/`CLOSE`) on `UpdateInstanceRequest` (and `UpdateInstanceRequestParams`). This is a modifiable attribute.
2. **`DescribeInstances`** — returns `EnableDestroyProtection` (`*string`) on `InstanceInfo`, so the current protection state can be read back.

Verification against the vendored SDK confirms:
- `EnableDestroyProtection` is **NOT** present on `CreateInstanceRequest`/`CreateInstanceRequestParams`. Destroy protection cannot be set at creation time; it can only be toggled afterwards via `UpdateInstance`.
- The field exists on `InstanceInfo` (read output) and on `UpdateInstanceRequest`/`UpdateInstanceRequestParams` (update input). No other ES request type carries this field.

Today the Terraform resource schema has no equivalent field, so users cannot manage destroy protection through Terraform. This design adds a single new optional+computed schema field wired to those two APIs.

The existing ES service-layer `UpdateInstance` wrapper (`service_tencentcloud_elasticsearch.go`) uses a long positional parameter list. Adding the new parameter therefore requires extending that signature and updating every existing call site within the resource CRUD code (which is fully in-repo, so this is a mechanical change).

## Goals / Non-Goals

**Goals:**
- Expose destroy protection management for ES instances via a new `enable_destroy_protection` schema field on `tencentcloud_elasticsearch_instance`.
- Support updating the value (OPEN/CLOSE) through the existing `UpdateInstance` API with retry and upgrade-wait behavior consistent with the resource's other update paths.
- Read the value back from `DescribeInstances` so state stays accurate.
- Remain fully backward compatible (the new field is optional + computed).

**Non-Goals:**
- Do not add destroy protection to the create flow (the create API does not support it). If a user sets `enable_destroy_protection` on create, it will be applied as an immediate post-create update step (see Decisions).
- Do not change the delete flow. Note: when destroy protection is `OPEN`, a `terraform destroy` will fail at the cloud API `DeleteInstance` call until the user sets the field to `CLOSE` (or disables it out-of-band). This is the expected/secure behavior and is documented.
- Do not modify any other existing schema field.
- Do not introduce a separate resource or `_extension.go` file.

## Decisions

### Decision 1: Field name `enable_destroy_protection`, type string, Optional + Computed

**Choice**: Add `"enable_destroy_protection"` as `schema.TypeString`, `Optional: true`, `Computed: true`, with `ValidateFunc` restricting values to `OPEN`/`CLOSE`.

**Rationale**: The cloud API field is a string enum (`OPEN`/`CLOSE`). `Computed` is required because the value is read back from `DescribeInstances` (and may be unset/empty in API responses, represented as the field being absent from state). Using `Optional + Computed` matches the existing pattern used by sibling modifiable fields on this resource (e.g. `public_access`, `kibana_public_access`, `protocol`).

**Naming note**: The user request specified the SchemaName as `EnableDestroyProtection` (PascalCase). The codebase convention for new top-level resource fields is snake_case, and sibling fields on this resource (`kibana_public_access`, `public_access`, `es_public_acl`, `protocol`) are snake_case. We adopt `enable_destroy_protection` for consistency with the resource's existing schema and Terraform naming conventions; the mapping to the cloud API field `EnableDestroyProtection` is handled in code.

**Alternatives considered**:
- *PascalCase `EnableDestroyProtection`*: rejected for inconsistency with the rest of this resource's schema (would be the only PascalCase top-level field here), even though a few PascalCase fields exist elsewhere in the provider. Snake_case is preferable for a clean, consistent resource.
- *Boolean `enable_destroy_protection`*: rejected because the API is string-enum-based (`OPEN`/`CLOSE`), matching the pattern of `public_access`/`kibana_public_access`. A bool would require lossy translation.

### Decision 2: Apply on create via a post-create UpdateInstance call

**Choice**: Because `CreateInstanceRequest` has no `EnableDestroyProtection` field, the create handler will, after the instance reaches `NORMAL` status, perform a single `UpdateInstance` call passing `enable_destroy_protection` (when the user set it) — mirroring how the resource already applies `es_acl`, `kibana_public_access`, `public_access`, and `cos_backup` as post-create update steps.

**Rationale**: This keeps behavior intuitive: setting the field in a create config still results in the protection being enabled. The resource already uses this post-create-update pattern for several attributes, so this is consistent and low-risk.

**Alternatives considered**:
- *Ignore the field on create, only honor on update*: rejected — surprising UX (a user setting `enable_destroy_protection = "OPEN"` on create would see it silently ignored until a later apply).

### Decision 3: Extend the `UpdateInstance` service-layer wrapper signature

**Choice**: Add a new `enableDestroyProtection string` parameter to `ElasticsearchService.UpdateInstance(...)`. Inside the wrapper, set `request.EnableDestroyProtection = &enableDestroyProtection` when the value is non-empty. Update all existing call sites in `resource_tc_elasticsearch_instance.go` to pass `""` for the new argument (no behavior change) except the new destroy-protection update path which passes the configured value.

**Rationale**: The existing wrapper is the single chokepoint that builds `UpdateInstanceRequest`. Extending it (rather than introducing a parallel method) keeps all `UpdateInstance` request construction in one place and reuses the existing retry/error-logging logic. The positional-parameter list is already long; adding one more parameter is acceptable and avoids a larger refactor that is out of scope.

**Alternatives considered**:
- *New dedicated service method `UpdateInstanceDestroyProtection`*: rejected — would duplicate request construction and retry logic, and `UpdateInstance` supports combining multiple attribute changes in one call which the resource relies on.
- *Refactor wrapper to accept a `*es.UpdateInstanceRequest`*: larger change, out of scope for a single-field addition.

### Decision 4: Update flow uses `d.HasChange` + retry + upgrade-wait

**Choice**: In `resourceTencentCloudElasticsearchInstanceUpdate`, add a `d.HasChange("enable_destroy_protection")` block that reads the new value and calls `elasticsearchService.UpdateInstance(...)` passing `enableDestroyProtection`, wrapped in `resource.Retry(tccommon.WriteRetryTimeout*2, ...)` with `tccommon.RetryError`, followed by `tencentCloudElasticsearchInstanceUpgradeWaiting(...)` — identical to the existing `public_access`, `es_public_acl`, `cos_backup`, etc. update blocks.

**Rationale**: Matches the established, proven pattern for every other updatable attribute on this resource.

### Decision 5: Read flow is nil-safe

**Choice**: In `resourceTencentCloudElasticsearchInstanceRead`, after obtaining `instance *es.InstanceInfo`, set state only when `instance.EnableDestroyProtection != nil`:
```go
if instance.EnableDestroyProtection != nil {
    _ = d.Set("enable_destroy_protection", instance.EnableDestroyProtection)
}
```

**Rationale**: The SDK marks this field as possibly-null (`注意：此字段可能返回 null`). Setting only when non-nil avoids writing an empty string over a user-managed value and follows the resource's existing nil-safe read convention for optional API fields.

## Risks / Trade-offs

- **[Risk] `terraform destroy` fails when protection is `OPEN`** → This is the intended security behavior. Mitigation: document clearly that users must set `enable_destroy_protection = "CLOSE"` (and apply) before destroying, or disable protection out-of-band. The error surfaced will be the cloud API error from `DeleteInstance`.
- **[Risk] API returns null for the field** → Mitigation: nil-safe read (Decision 5); when null, state simply is not updated for this attribute, preserving the user's configured value.
- **[Risk] Forgetting to update a call site of `UpdateInstance`** → Mitigation: The new positional argument causes a compile error at every existing call site if missed, so the compiler enforces completeness. All call sites are within the single resource file and the create-flow post-create calls, all in-repo.
- **[Trade-off] Positional parameter list grows by one** → Accepted; a full refactor to a request-struct signature is explicitly a non-goal here.

## Migration Plan

- No state migration required. The new field is optional + computed.
- Deploy: add the field, ship new provider version. Existing configs without `enable_destroy_protection` continue to work unchanged.
- Rollback: removing the field reverts to prior behavior; no persisted state for the new field affects other attributes.

## Open Questions

- None outstanding. SDK support and create-vs-update API capabilities have been verified against the vendored `es/v20180416` package.
