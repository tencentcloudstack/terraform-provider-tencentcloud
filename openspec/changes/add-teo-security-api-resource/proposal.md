# Add tencentcloud_teo_security_api_resource Resource

## What

Add a new Terraform resource `tencentcloud_teo_security_api_resource` for managing Tencent Cloud EdgeOne (TEO) API security resources. This resource manages a single API resource entry under a site zone, supporting full CRUD lifecycle.

## Why

TEO API protection allows users to define API resources (path, methods, service bindings, request constraints) for fine-grained security control. Currently no Terraform resource exists to manage these API resources, requiring manual portal operations. This resource enables infrastructure-as-code management of TEO API resources.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateSecurityAPIResource` | `APIResources` array length is limited to **1**; returns `APIResourceIds` array |
| Read | `DescribeSecurityAPIResource` | Filter by zone and paginate; match by `Id` field in response |
| Update | `ModifySecurityAPIResource` | Pass `Id` + updated fields in `APIResources[0]` |
| Delete | `DeleteSecurityAPIResource` | Pass `ZoneId` + `APIResourceIds[<id>]` |

## Resource ID

Composite: `<zone_id>#<api_resource_id>` (using `tccommon.FILED_SP`), e.g. `zone-123sfakjf#apires-1232382313`.
`api_resource_id` is taken from `APIResourceIds[0]` returned by `CreateSecurityAPIResource`.
