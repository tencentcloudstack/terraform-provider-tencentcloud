## 1. Service Layer

- [x] 1.1 Add `TeoPlanForZone(zoneId, planType string) (resourceNames, dealNames []*string, errRet error)` method to `tencentcloud/services/teo/service_tencentcloud_teo.go`, following the `TeoIdentifyZone` pattern: build `CreatePlanForZoneRequest`, set `ZoneId` and `PlanType`, ratelimit check, call `me.client.UseTeoV20220901Client().CreatePlanForZone(request)`, nil-safe access to `response.Response`

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_plan_for_zone_operation.go` with constructor `ResourceTencentCloudTeoPlanForZone()`, operation-type code style following `tencentcloud_teo_identify_zone_operation`
- [x] 2.2 Define schema: required ForceNew `zone_id` (TypeString), required ForceNew `plan_type` (TypeString); computed `resource_names` (TypeList, elem TypeString), computed `deal_names` (TypeList, elem TypeString)
- [x] 2.3 Implement Create: build `CreatePlanForZoneRequest` from `zone_id` and `plan_type`, call `TeoPlanForZone` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`, after retry success set `d.SetId(helper.BuildToken())` and nil-safe set computed `resource_names` and `deal_names`, then call Read
- [x] 2.4 Implement Read: no-op (return nil)
- [x] 2.5 Implement Delete: no-op (return nil)
- [x] 2.6 Ensure nil-safe response access: if response or response.Response is nil, return NonRetryableError; all response-value access is nil-safe

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_plan_for_zone` in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add the resource entry to `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_plan_for_zone_operation.md` example file, format following `resource_tc_teo_identify_zone_operation.md` (description with TEO product name + Example Usage; no Import section for OPERATION type)
- [ ] 4.2 Generate `website/docs/r/teo_plan_for_zone.html.markdown` via `make doc` (do not hand-write) and confirm `website/tencentcloud.erb` link entry

## 5. Unit Test

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_plan_for_zone_operation_test.go`, using gomonkey to mock the TEO API client, naming/format following `resource_tc_teo_identify_zone_operation_test.go`; cover Create with successful response

## 6. Verification

- [ ] 6.1 Run `gofmt` on changed Go files
- [x] 6.2 Run `go test -gcflags=all=-l ./tencentcloud/services/teo/ -run TestAccTencentCloudTeoPlanForZone` to confirm unit tests pass
- [ ] 6.3 Run `make doc` and verify generated website doc is consistent with the `.md` example
