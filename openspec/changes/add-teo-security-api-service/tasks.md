## 1. Service Layer

- [x] 1.1 Append `DescribeTeoSecurityAPIServiceById(ctx, zoneId, apiServiceId)` to `tencentcloud/services/teo/service_tencentcloud_teo.go` — paginates `DescribeSecurityAPIService` (Limit=100) until the entry with `Id == apiServiceId` is found; returns `*teo.APIService` or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_security_api_service.go` with full schema following `tencentcloud_igtm_strategy` style:
  - Top-level fields: `zone_id` (Required, ForceNew)
  - `api_services` (Required, TypeList, MaxItems:1) with sub-fields mapping 100% to SDK `APIService` struct:
    - `name` (Required, String)
    - `base_path` (Required, String)
    - `id` (Computed, String)

- [x] 2.2 Implement Create:
  - Build `CreateSecurityAPIServiceRequest` with `ZoneId` and `APIServices` (single-element array, no `Id` field)
  - Call `CreateSecurityAPIServiceWithContext`; extract `APIServiceIds[0]` as `apiServiceId`
  - Set resource ID to `strings.Join([]string{zoneId, apiServiceId}, tccommon.FILED_SP)`
  - Call Read

- [x] 2.3 Implement Read:
  - Split ID into `zoneId` and `apiServiceId`
  - Call `DescribeTeoSecurityAPIServiceById`; if nil → `d.SetId("")`
  - Populate `zone_id` and `api_services` block (id, name, base_path) from response

- [x] 2.4 Implement Update:
  - Build `ModifySecurityAPIServiceRequest` with `ZoneId` and `APIServices=[{Id: apiServiceId, Name, BasePath}]`
  - Call `ModifySecurityAPIServiceWithContext`
  - Call Read

- [x] 2.5 Implement Delete:
  - Build `DeleteSecurityAPIServiceRequest` with `ZoneId` and `APIServiceIds=[apiServiceId]`
  - Call `DeleteSecurityAPIServiceWithContext` with Retry

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_security_api_service` in `tencentcloud/provider.go` ResourcesMap, pointing to `teo.ResourceTencentCloudTeoSecurityAPIService()`

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_security_api_service.md` — document all arguments, attributes, and import syntax with example HCL
- [x] 4.2 Create `tencentcloud/services/teo/resource_tc_teo_security_api_service_test.go` — basic acceptance test covering create/update/import/delete following `resource_tc_igtm_strategy_test.go` style
