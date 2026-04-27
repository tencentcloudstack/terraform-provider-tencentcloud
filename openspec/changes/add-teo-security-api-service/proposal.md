# Add tencentcloud_teo_security_api_service Resource

## What

Add a new Terraform resource `tencentcloud_teo_security_api_service` for managing Tencent Cloud EdgeOne (TEO) API security services. This resource manages a single API service entry under a site zone, supporting full CRUD lifecycle.

## Why

TEO API protection allows users to define API services (name, base path) as logical groupings for API security policy enforcement. Currently no Terraform resource exists to manage these API services, requiring manual portal operations. This resource enables infrastructure-as-code management of TEO API services.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateSecurityAPIService` | `APIServices` array length limited to **1**; returns `APIServiceIds` array |
| Read | `DescribeSecurityAPIService` | Paginate by zone; match by `Id` field in response |
| Update | `ModifySecurityAPIService` | Pass `Id` + updated fields in `APIServices[0]` |
| Delete | `DeleteSecurityAPIService` | Pass `ZoneId` + `APIServiceIds=[<id>]` |

## Resource ID

Composite: `<zone_id>#<api_service_id>` (using `tccommon.FILED_SP`), e.g. `zone-123sfakjf#apisrv-1232382313`.
`api_service_id` is taken from `APIServiceIds[0]` returned by `CreateSecurityAPIService`.
