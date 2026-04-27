# Design: tencentcloud_teo_security_api_resource Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style:

```
provider.go
    └─ tencentcloud/services/teo/resource_tc_teo_security_api_resource.go  (CRUD handlers)
           └─ tencentcloud/services/teo/service_tencentcloud_teo.go (DescribeTeoSecurityAPIResourceById)
                  └─ teo SDK v20220901
```

## Resource ID

Composite: `<zone_id>#<api_resource_id>` (using `tccommon.FILED_SP`), e.g. `zone-123sfakjf#apires-1232382313`.

## Key Constraint

`APIResources` array in Create/Modify is limited to **1 element**. The schema reflects this as a flat resource (not a list) — all `APIResource` fields are top-level schema fields.

## Schema

### Required

| Field | Type | ForceNew | Description |
|---|---|---|---|
| `zone_id` | String | Yes | Site ID |
| `name` | String | No | API resource name |
| `path` | String | No | API path, e.g. `/ava` |

### Optional

| Field | Type | Description |
|---|---|---|
| `api_service_ids` | List of String | Associated API service ID list |
| `methods` | List of String | Allowed HTTP methods: `GET`, `POST`, `PUT`, `HEAD`, `PATCH`, `OPTIONS`, `DELETE` |
| `request_constraint` | String | Request matching rule expression |

### Computed

| Field | Type | Description |
|---|---|---|
| `api_resource_id` | String | API resource ID returned by the API (same as resource ID suffix) |

## Read Logic

Call `DescribeSecurityAPIResource` with `ZoneId`, paginate and find the entry where `APIResource.Id == apiResourceId`.
If not found → resource deleted → `d.SetId("")`.

## Update Logic

Call `ModifySecurityAPIResource` with `ZoneId` and `APIResources=[{Id: apiResourceId, Name: ..., Path: ..., ...}]`.

## Delete Logic

Call `DeleteSecurityAPIResource` with `ZoneId` and `APIResourceIds=[apiResourceId]`.

## Key SDK Types

```go
// teo v20220901
type APIResource struct {
    Id                *string
    Name              *string
    APIServiceIds     []*string
    Path              *string
    Methods           []*string
    RequestConstraint *string
}

CreateSecurityAPIResourceRequest {
    ZoneId       *string
    APIResources []*APIResource  // length must be 1
}
// Response: APIResourceIds []*string

ModifySecurityAPIResourceRequest {
    ZoneId       *string
    APIResources []*APIResource  // Id field required
}

DeleteSecurityAPIResourceRequest {
    ZoneId         *string
    APIResourceIds []*string
}

DescribeSecurityAPIResourceRequest {
    ZoneId  *string
    Limit   *int64
    Offset  *int64
}
// Response: APIResources []*APIResource, TotalCount *int64
```
