## 1. Schema & Service Layer

- [x] 1.1 Add `enable_cerebro`, `cerebro_public_access`, `cerebro_private_access`, `cerebro_private_domain` fields to resource schema in `resource_tc_elasticsearch_instance.go`
- [x] 1.2 Add `ES_CEREBRO_PUBLIC_ACCESS` and `ES_CEREBRO_PRIVATE_ACCESS` constants to `resource_tc_elasticsearch_instance.go` (reuse `ES_KIBANA_PUBLIC_ACCESS` / `ES_PRIVATE_ACCESS` pattern)
- [x] 1.3 Extend `UpdateInstance` service function signature in `service_tencentcloud_elasticsearch.go` to accept `enableCerebro *bool`, `cerebroPublicAccess string`, `cerebroPrivateAccess string`, `cerebroPrivateDomain string`
- [x] 1.4 Update `UpdateInstance` function body to set `request.EnableCerebro`, `request.CerebroPublicAccess`, `request.CerebroPrivateAccess`, `request.CerebroPrivateDomain` when non-empty/non-nil
- [x] 1.5 Update all 26 existing call sites of `UpdateInstance` to pass empty Cerebro parameters (`nil`, `""`, `""`, `""`)

## 2. Resource Update Logic

- [x] 2.1 Add `d.HasChange("enable_cerebro")` block in update function to call `UpdateInstance` with `EnableCerebro` field
- [x] 2.2 Add `d.HasChange("cerebro_public_access")` block in update function
- [x] 2.3 Add `d.HasChange("cerebro_private_access")` block in update function
- [x] 2.4 Add `d.HasChange("cerebro_private_domain")` block in update function
- [x] 2.5 Each update block SHALL use retry logic with `tccommon.WriteRetryTimeout` and wait for upgrade via `tencentCloudElasticsearchInstanceUpgradeWaiting`

## 3. Documentation

- [x] 3.1 Update `resource_tc_elasticsearch_instance.md` with example usage demonstrating Cerebro parameters
- [x] 3.2 Add Cerebro field descriptions to the resource schema (in Go code)

## 4. Verification

- [x] 4.1 Verify the code compiles correctly (all call sites updated, no missing imports)
- [x] 4.2 Verify existing tests still pass