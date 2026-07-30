# Implementation Tasks: Add ES Instance Destroy Protection

## 1. Service Layer

- [x] 1.1 Extend the `ElasticsearchService.UpdateInstance` wrapper signature in `tencentcloud/services/es/service_tencentcloud_elasticsearch.go` to add a new `enableDestroyProtection string` parameter (append after `multiZoneInfo []*es.ZoneDetail`).
- [x] 1.2 Inside the wrapper, set `request.EnableDestroyProtection = &enableDestroyProtection` only when `enableDestroyProtection != ""`, before the `ratelimit.Check` / API call.

## 2. Schema Definition

- [x] 2.1 In `tencentcloud/services/es/resource_tc_elasticsearch_instance.go`, add a new schema field `enable_destroy_protection` (`schema.TypeString`, `Optional: true`, `Computed: true`) with a `ValidateFunc` restricting values to `OPEN`/`CLOSE`, and a description noting it toggles cluster destroy protection. Place it alongside the other modifiable optional+computed fields (e.g. near `public_access`/`protocol`).

## 3. Create Logic

- [x] 3.1 In `resourceTencentCloudElasticsearchInstanceCreate`, after the instance reaches normal status and after the existing post-create update blocks (es_acl, kibana_public_access, kibana_private_access, public_access, cos_backup), add a block: when `enable_destroy_protection` is set (`d.GetOk`), call `elasticsearchService.UpdateInstance(...)` passing the value as the new `enableDestroyProtection` argument (other args empty/nil), wrapped in `resource.Retry(tccommon.WriteRetryTimeout*2, ...)` with `tccommon.RetryError`, then call `tencentCloudElasticsearchInstanceUpgradeWaiting(...)`.
- [x] 3.2 Update all existing `elasticsearchService.UpdateInstance(...)` call sites within the create function to pass `""` for the new trailing `enableDestroyProtection` argument (no behavior change).

## 4. Read Logic

- [x] 4.1 In `resourceTencentCloudElasticsearchInstanceRead`, after obtaining `instance *es.InstanceInfo`, add a nil-safe set: `if instance.EnableDestroyProtection != nil { _ = d.Set("enable_destroy_protection", instance.EnableDestroyProtection) }`.

## 5. Update Logic

- [x] 5.1 In `resourceTencentCloudElasticsearchInstanceUpdate`, update all existing `elasticsearchService.UpdateInstance(...)` call sites to pass `""` for the new trailing `enableDestroyProtection` argument (compiler-enforced).
- [x] 5.2 Add a new `d.HasChange("enable_destroy_protection")` block (placed among the other update blocks, e.g. after `cos_backup`) that reads the new value via `d.Get("enable_destroy_protection").(string)` and calls `elasticsearchService.UpdateInstance(...)` with that value as `enableDestroyProtection` (other args empty/nil), wrapped in `resource.Retry(tccommon.WriteRetryTimeout*2, ...)` with `tccommon.RetryError`, followed by `tencentCloudElasticsearchInstanceUpgradeWaiting(...)`.

## 6. Documentation

- [x] 6.1 Update the source documentation file `tencentcloud/services/es/resource_tc_elasticsearch_instance.md` to add an example (or extend an existing example) demonstrating `enable_destroy_protection = "OPEN"`. Keep the one-line description, Example Usage, and Import sections only (no manual `Argument Reference` / `Attribute Reference`).

## 7. Tests

- [x] 7.1 In `tencentcloud/services/es/resource_tc_elasticsearch_instance_test.go`, add mock-based (gomonkey) unit tests covering: schema field presence/defaults, create flow invoking the post-create `UpdateInstance` with `EnableDestroyProtection`, update flow invoking `UpdateInstance` on `d.HasChange`, and read flow nil-safe setting from `InstanceInfo.EnableDestroyProtection`. Run the relevant test file with `go test -gcflags=all=-l` to confirm it passes.

## 8. Finalization (run via tfpacer-finalize skill only)

- [ ] 8.1 Run `gofmt` on the modified Go files.
- [ ] 8.2 Run `make doc` to generate `website/docs/r/elasticsearch_instance.html.markdown`.
- [ ] 8.3 Create a changelog file under `.changelog/` describing the new `enable_destroy_protection` parameter for `tencentcloud_elasticsearch_instance`.

## Summary

- **Total Tasks**: 14
- **Critical Path**: Service Layer (1.x) → Schema (2.x) → Create/Read/Update (3.x–5.x) → Docs (6.x) → Tests (7.x) → Finalization (8.x)
- **Notes**: No SDK upgrade required (`EnableDestroyProtection` already present in vendored `es/v20180416`). The `CreateInstance` API does not support destroy protection, so it is applied as a post-create update. Build/lint verification is handled outside this task list by the build pipeline; gofmt, `make doc`, and changelog creation are deferred to the tfpacer-finalize skill.
