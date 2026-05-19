## 1. Service Layer

- [ ] 1.1 Add `DescribeLiveOriginStreamInfo(ctx, domainName)` to `service_tencentcloud_live.go` — calls `DescribeOriginStreamInfo`, returns `*live.DescribeOriginStreamInfoResponseParams`, handles `ResourceNotFound` by returning nil

## 2. Resource Implementation

- [ ] 2.1 Create `resource_tc_live_origin_stream_info_config.go` with `ResourceTencentCloudLiveOriginStreamInfoConfig()` — schema definition with all `ModifyOriginStreamInfo` fields (required + optional + computed `status`)
- [ ] 2.2 Implement `Create`: set ID to `domain_name`, delegate to `Update`
- [ ] 2.3 Implement `Read`: call `DescribeLiveOriginStreamInfo`, map all fields to state; if nil, set ID to ""
- [ ] 2.4 Implement `Update`: build `ModifyOriginStreamInfoRequest` from state, call with retry, then poll `DescribeOriginStreamInfo` until `Status == 1` or `3`
- [ ] 2.5 Implement `Delete`: no-op, return nil
- [ ] 2.6 Implement `customization_rules` nested schema (`TypeList`) with all `OriginStreamCustomizationRule` fields and build/flatten helpers

## 3. Provider Registration

- [ ] 3.1 Register `tencentcloud_live_origin_stream_info_config` in `provider.go`

## 4. Documentation

- [ ] 4.1 Create `resource_tc_live_origin_stream_info_config.md` with example HCL and import instruction

## 5. Tests

- [ ] 5.1 Create `resource_tc_live_origin_stream_info_config_test.go` with `TestAccTencentCloudLiveOriginStreamInfoConfigResource_basic` covering create, update, and import steps

## 6. Verification

- [ ] 6.1 Run `go build ./tencentcloud/services/live/` — confirm no compile errors
- [ ] 6.2 Run `go build ./tencentcloud/` — confirm provider compiles
