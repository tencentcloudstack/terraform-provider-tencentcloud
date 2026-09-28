## 1. Connectivity & Service Layer

- [x] 1.1 Add BDRC SDK import (`bdrcv20260330`) and `bdrcv20260330Conn` field to `tencentcloud/connectivity/client.go`, implement `UseBdrcV20260330Client()` method (mirror `UseIgtmV20231024Client`)
- [x] 1.2 Create `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go` with `BdrcService` struct and `DescribeDisasterRecoverySitePairById(ctx, sitePairId, sitePairType)` method wrapping `DescribeDisasterRecoverySitePairs` (set `Limit` to 100, filter by `SitePairIds`, return first matching `SitePair`)

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_site_pair.go` with full schema definition (required ForceNew fields, optional fields, computed fields including nested `protected_resource_set`, `protected_resource_status_set`, `cross_cloud_details`)
- [x] 2.2 Implement `resourceTencentCloudBdrcDisasterRecoverySitePairCreate`: call `CreateDisasterRecoverySitePair` with retry, check `SitePairId` is non-empty (return `NonRetryableError` if empty), set `d.Id()`, then call Read
- [x] 2.3 Implement `resourceTencentCloudBdrcDisasterRecoverySitePairRead`: call service layer `DescribeDisasterRecoverySitePairById`, nil-check each field before `d.Set()`, log `[CRUD]` with `d.Id()` before `d.SetId("")` when not found
- [x] 2.4 Implement `resourceTencentCloudBdrcDisasterRecoverySitePairUpdate`: call `ModifySitePairAttribute` with `SitePairId` from `d.Id()` and new `SitePairName` when `site_pair_name` changed, then call Read
- [x] 2.5 Implement `resourceTencentCloudBdrcDisasterRecoverySitePairDelete`: call `DeleteDisasterRecoverySitePairs` with `SitePairIds` as single-element list containing `d.Id()`, with retry

## 3. Provider Registration

- [x] 3.1 Import `bdrc` service package and register `tencentcloud_bdrc_disaster_recovery_site_pair` in `provider.go` ResourcesMap
- [x] 3.2 Add BDRC resource entry to `provider.md` under a new BDRC section

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_site_pair.md` with one-line description, Example Usage, and Import section
- [x] 4.2 Create `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_site_pair_test.go` with gomonkey mocks (no TF_ACC test suite) covering Create, Read, Update, Delete business logic

## 5. Verification (finalization phase only)

- [ ] 5.1 Run `gofmt` on all new/modified Go files (via tfpacer-finalize skill)
- [ ] 5.2 Run `make doc` to generate website/docs documentation (via tfpacer-finalize skill)
- [ ] 5.3 Create changelog entry (via tfpacer-finalize skill)
