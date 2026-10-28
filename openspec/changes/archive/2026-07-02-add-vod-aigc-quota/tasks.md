## 1. Service Layer Implementation

- [x] 1.1 Add `CreateVodAigcQuota` method to `service_tencentcloud_vod.go` — wraps `CreateAigcQuota` API with WriteRetryTimeout retry and ratelimit check
- [x] 1.2 Add `DescribeVodAigcQuotaById` method to `service_tencentcloud_vod.go` — wraps `DescribeAigcQuotas` API with ReadRetryTimeout retry, filters by (SubAppId, QuotaType, ApiToken), returns `*vod.AigcQuotaItem` or nil
- [x] 1.3 Add `ModifyVodAigcQuota` method to `service_tencentcloud_vod.go` — wraps `ModifyAigcQuota` API with WriteRetryTimeout retry and ratelimit check
- [x] 1.4 Add `DeleteVodAigcQuota` method to `service_tencentcloud_vod.go` — wraps `DeleteAigcQuota` API with WriteRetryTimeout retry, idempotent on ResourceNotFound errors

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_vod_aigc_quota.go` — define `ResourceTencentCloudVodAigcQuota()` with schema fields (sub_app_id, quota_type, quota_limit, api_token, usage)
- [x] 2.2 Implement `parseVodAigcQuotaId()` — parse composite ID `{sub_app_id}#{quota_type}#{api_token}`
- [x] 2.3 Implement `resourceTencentCloudVodAigcQuotaCreate()` — call CreateAigcQuota, poll DescribeAigcQuotas until visible, set composite ID
- [x] 2.4 Implement `resourceTencentCloudVodAigcQuotaRead()` — parse composite ID, call DescribeAigcQuotas, set all schema fields
- [x] 2.5 Implement `resourceTencentCloudVodAigcQuotaUpdate()` — call ModifyAigcQuota, then call Read
- [x] 2.6 Implement `resourceTencentCloudVodAigcQuotaDelete()` — call DeleteAigcQuota, poll DescribeAigcQuotas until gone
- [x] 2.7 Add import support using composite ID format

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_vod_aigc_quota` resource in `tencentcloud/provider.go`
- [x] 3.2 Update `tencentcloud/provider.md` with the new resource entry

## 4. Unit Tests

- [x] 4.1 Create `resource_tc_vod_aigc_quota_test.go` with unit test cases using gomonkey to mock cloud API calls
- [x] 4.2 Run `go test -gcflags=all=-l` on the test file and ensure all tests pass

## 5. Documentation

- [x] 5.1 Create `resource_tc_vod_aigc_quota.md` with example usage and import instructions
- [ ] 5.2 Run `make doc` to generate final `website/docs/` documentation (deferred to tfpacer-finalize skill)