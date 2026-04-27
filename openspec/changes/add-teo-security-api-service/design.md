# Design: tencentcloud_teo_security_api_service Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style:

```
provider.go
    └─ tencentcloud/services/teo/resource_tc_teo_security_api_service.go  (CRUD handlers)
           └─ tencentcloud/services/teo/service_tencentcloud_teo.go (DescribeTeoSecurityAPIServiceById)
                  └─ teo SDK v20220901
```

## Resource ID

Composite: `<zone_id>#<api_service_id>` (using `tccommon.FILED_SP`), e.g. `zone-123sfakjf#apisrv-1232382313`.

## Key Constraint

`APIServices` array in Create/Modify is limited to **1 element**. The `api_services` field uses `TypeList` with `MaxItems: 1` to align with the SDK `APIService` struct — all `APIService` fields are nested inside this block.

## Schema

### Required

| Field | Type | ForceNew | Description |
|---|---|---|---|
| `zone_id` | String | Yes | Site ID |
| `api_services` | List (MaxItems:1) | No | API service configuration block |

### api_services sub-fields (100% mapping to SDK APIService struct)

| Sub-field | Type | Required | SDK Field | Description |
|---|---|---|---|---|
| `name` | String | Yes | `Name` | API service name |
| `base_path` | String | Yes | `BasePath` | API service base path, e.g. `/tt` |
| `id` | String | Computed | `Id` | API service ID returned by API |

## Read Logic

Call `DescribeSecurityAPIService` with `ZoneId`, paginate (Limit=100) until the entry with `APIService.Id == apiServiceId` is found.
If not found → resource deleted → `d.SetId("")`.

## Update Logic

Call `ModifySecurityAPIService` with `ZoneId` and `APIServices=[{Id: apiServiceId, Name: ..., BasePath: ...}]`.

## Delete Logic

Call `DeleteSecurityAPIService` with `ZoneId` and `APIServiceIds=[apiServiceId]`.

## Key SDK Types

```go
// teo v20220901
type APIService struct {
    Id       *string   // Computed
    Name     *string   // Required
    BasePath *string   // Required
}

CreateSecurityAPIServiceRequest {
    ZoneId      *string
    APIServices []*APIService  // length must be 1
}
// Response: APIServiceIds []*string

ModifySecurityAPIServiceRequest {
    ZoneId      *string
    APIServices []*APIService  // Id field required
}

DeleteSecurityAPIServiceRequest {
    ZoneId        *string
    APIServiceIds []*string
}

DescribeSecurityAPIServiceRequest {
    ZoneId  *string
    Limit   *int64
    Offset  *int64
}
// Response: APIServices []*APIService, TotalCount *int64
```
