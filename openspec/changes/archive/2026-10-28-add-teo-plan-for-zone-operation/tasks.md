## 1. Service Layer

- [x] 1.1 `TeoPlanForZone(zoneId, planType string) (resourceNames, dealNames []*string, errRet error)` exists in `tencentcloud/services/teo/service_tencentcloud_teo.go`, building a `CreatePlanForZoneRequest`, setting `ZoneId` and `PlanType`, calling `me.client.UseTeoV20220901Client().CreatePlanForZone(request)`, with nil-safe access to `response.Response`

## 2. Framework Action Implementation

- [x] 2.1 Create `tencentcloud/services/teo/action_tc_teo_plan_for_zone.go` implementing `action.Action` + `action.ActionWithConfigure`, following the framework Action style
- [x] 2.2 Define schema: required `zone_id` (StringAttribute), required `plan_type` (StringAttribute). No output attributes (framework actions cannot surface output values)
- [x] 2.3 `Metadata` sets TypeName to `tencentcloud_teo_plan_for_zone`
- [x] 2.4 `Configure` receives `*sharedmeta.ProviderMeta` and constructs `NewTeoService(meta.Client)`
- [x] 2.5 `Invoke` reads `zone_id`/`plan_type`, validates non-empty, calls `a.client.TeoPlanForZone(...)`, surfaces API errors and empty responses as error diagnostics; logs `ResourceNames`/`DealNames` on success
- [x] 2.6 Reverted/removed the SDKv2 operation resource files (`resource_tc_teo_plan_for_zone_operation.go`, `resource_tc_teo_plan_for_zone_operation_test.go`)

## 3. Provider Registration

- [x] 3.1 Register `teo.NewTeoPlanForZoneAction` in `tencentcloud/framework/registry.go` `actionFactories` and import the teo service package
- [x] 3.2 Move `tencentcloud_teo_plan_for_zone` from the Resource section to a new `Action` sub-section (TEO product) in `tencentcloud/provider.md`
- [x] 3.3 Remove `tencentcloud_teo_plan_for_zone` from the SDKv2 `ResourcesMap` in `tencentcloud/provider.go`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/action_tc_teo_plan_for_zone.md` (renamed from `resource_tc_teo_plan_for_zone.md`), format following the gendoc framework Action conventions (description with TEO product name + Example Usage; no Import section)
- [ ] 4.2 Generate `website/docs/a/teo_plan_for_zone.html.markdown` via `make doc` and confirm the `tencentcloud.erb` Action sidebar entry (finalization stage)

## 5. Unit Test

- [x] 5.1 Create `tencentcloud/services/teo/action_tc_teo_plan_for_zone_test.go`, using gomonkey to mock the TEO API client; cover Invoke success, missing required fields, API error, and empty response, plus Metadata/Schema/Configure validation

## 6. Verification

- [ ] 6.1 Run `gofmt` on changed Go files (finalization stage)
- [x] 6.2 Verify the framework Action implementation compiles and unit tests cover the Invoke path
- [ ] 6.3 Run `make doc` and verify generated website doc is consistent with the `.md` example (finalization stage)