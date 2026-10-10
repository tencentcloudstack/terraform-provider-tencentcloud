## 1. Service Layer

- [x] 1.1 Add `DescribeConfigAggregatorById(ctx, accountGroupId, ownerUin)` helper to `tencentcloud/services/config/service_tencentcloud_config.go` wrapping `DescribeAggregatorWithContext` with `tccommon.ReadRetryTimeout` retry and `tccommon.RetryError` error wrapping

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/config/resource_tc_config_aggregator.go` with `ResourceTencentCloudConfigAggregator()` schema: `name`(required string), `description`(required string), `type`(required string, ForceNew), `owner_uin`(required string, ForceNew), `aggregator_accounts`(optional list of {`member_uin` int, `member_name` string}), `account_group_id`(computed string), `aggregator_status`(computed int); Importer passthrough
- [x] 2.2 Implement Create: build `CreateAggregatorRequest` (Name/Description/Type/AggregatorAccounts); `resource.Retry(WriteRetryTimeout)` → `CreateAggregatorWithContext`; nil-response guard; check `AccountGroupId` nil/""→ NonRetryableError; `d.SetId(accountGroupId#ownerUin)`; call Read
- [x] 2.3 Implement Read: parse composite id (2 parts); call `DescribeConfigAggregatorById`; on empty → `log.Printf("[CRUD] tencentcloud_config_aggregator id=%s", d.Id())` then `d.SetId("")`; nil-check each field before set; flatten `AggregatorAccounts`; set `account_group_id`, `aggregator_status`, `owner_uin`(from id)
- [x] 2.4 Implement Update: parse id; `needChange` over `["name","description","aggregator_accounts"]`; build `UpdateAggregatorRequest` (Name/Description/AccountGroupId/OwnerUin/AggregatorAccounts); `resource.Retry(WriteRetryTimeout)` → `UpdateAggregatorWithContext`; nil-response guard; call Read
- [x] 2.5 Implement Delete: parse id; build `DeleteAggregatorsRequest` (AccountGroupId/OwnerUin); `resource.Retry(WriteRetryTimeout)` → `DeleteAggregatorsWithContext`; nil-response guard

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_config_aggregator` → `config.ResourceTencentCloudConfigAggregator()` in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Register doc entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/config/resource_tc_config_aggregator.md` (one-line desc with "Config" product name; Example Usage; Import section explaining composite id `account_group_id#owner_uin`; no Argument/Attribute Reference)

## 5. Tests

- [x] 5.1 Create `tencentcloud/services/config/resource_tc_config_aggregator_test.go` using gomonkey (no terraform test suite): mock `UseConfigV20220802Client` + `CreateAggregatorWithContext`/`DescribeAggregatorWithContext`/`UpdateAggregatorWithContext`/`DeleteAggregatorsWithContext`; cover Create→Read→Update→Delete business logic
- [x] 5.2 Verify all error returns checked; functions that cannot fail use `_ =` for the err

## 6. Verification (post-implementation)

- [ ] 6.1 Run `gofmt` and `make doc` via tfpacer-finalize (do NOT run go build/vet/test manually)