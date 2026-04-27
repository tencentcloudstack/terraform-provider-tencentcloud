## 1. Service Layer

- [x] 1.1 Append `DescribeTeoSecurityAPIResourceById(ctx, zoneId, apiResourceId)` to `tencentcloud/services/teo/service_tencentcloud_teo.go` — paginates `DescribeSecurityAPIResource` (Limit=20) until the entry with `Id == apiResourceId` is found; returns `*teo.APIResource` or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_security_api_resource.go` with full schema following `tencentcloud_igtm_strategy` style:
  - Schema fields: `zone_id` (Required, ForceNew), `name` (Required), `path` (Required), `api_service_ids` (Optional, List of String), `methods` (Optional, List of String), `request_constraint` (Optional, String), `api_resource_id` (Computed)

- [x] 2.2 Implement Create:
  - Build `CreateSecurityAPIResourceRequest` with `ZoneId` and `APIResources` (single-element array)
  - Call `CreateSecurityAPIResourceWithContext`; extract `APIResourceIds[0]` as `apiResourceId`
  - Set resource ID to `strings.Join([]string{zoneId, apiResourceId}, tccommon.FILED_SP)`
  - Call Read

- [x] 2.3 Implement Read:
  - Split ID into `zoneId` and `apiResourceId`
  - Call `DescribeTeoSecurityAPIResourceById`; if nil → `d.SetId("")`
  - Populate all schema fields: `zone_id`, `name`, `path`, `api_service_ids`, `methods`, `request_constraint`, `api_resource_id`

- [x] 2.4 Implement Update:
  - Build `ModifySecurityAPIResourceRequest` with `ZoneId` and `APIResources=[{Id, Name, Path, APIServiceIds, Methods, RequestConstraint}]`
  - Call `ModifySecurityAPIResourceWithContext`
  - Call Read

- [x] 2.5 Implement Delete:
  - Build `DeleteSecurityAPIResourceRequest` with `ZoneId` and `APIResourceIds=[apiResourceId]`
  - Call `DeleteSecurityAPIResourceWithContext` with Retry

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_security_api_resource` in `tencentcloud/provider.go` ResourcesMap, pointing to `teo.ResourceTencentCloudTeoSecurityAPIResource()`

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_security_api_resource.md` — document all arguments, attributes, and import syntax with example HCL
- [x] 4.2 Create `tencentcloud/services/teo/resource_tc_teo_security_api_resource_test.go` — basic acceptance test covering create/update/import/delete following `resource_tc_igtm_strategy_test.go` style
