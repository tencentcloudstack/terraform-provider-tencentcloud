# Proposal: Add VOD Sub Applications Data Source

## Change ID
`add-vod-sub-applications-datasource`

## Overview
Add a new data source `tencentcloud_vod_sub_applications` to query Tencent Cloud VOD (Video on Demand) sub-application lists using the `DescribeSubAppIds` API.

## Motivation
Currently, the provider has a `tencentcloud_vod_sub_application` resource for managing individual VOD sub-applications, but lacks a corresponding data source to query existing sub-applications. This limitation prevents users from:

1. **Discovering existing sub-applications**: Users cannot list and filter sub-applications programmatically in their Terraform configurations
2. **Referencing sub-applications**: Without a data source, users cannot reference existing sub-applications for use in other resources
3. **Tag-based filtering**: The API supports querying sub-applications by tags, but there's no way to leverage this capability
4. **Pagination support**: Users cannot efficiently query large numbers of sub-applications with proper pagination

## Proposed Solution
Implement a new data source `tencentcloud_vod_sub_applications` that:

1. **Query capabilities**:
   - Filter by application name
   - Filter by resource tags
   - Support pagination with offset/limit parameters

2. **Response fields** (based on `SubAppIdInfo` structure):
   - `sub_app_id`: Sub-application ID
   - `sub_app_id_name`: Sub-application name
   - `name`: Legacy name field (for backward compatibility)
   - `description`: Sub-application description
   - `create_time`: Creation time in ISO 8601 format
   - `status`: Application status (On/Off/Destroying/Destroyed)
   - `mode`: Application mode (fileid/fileid+path)
   - `storage_regions`: Enabled storage regions list
   - `tags`: Resource tags bound to the sub-application

3. **Implementation approach**:
   - Follow existing VOD data source patterns (e.g., `data_source_tc_vod_super_player_configs.go`)
   - Use the VOD service layer for API calls
   - Implement pagination with retry logic (similar to recent VOD resource updates)
   - Support `result_output_file` for exporting query results

## API Reference
- **API**: `DescribeSubAppIds`
- **Documentation**: https://cloud.tencent.com/document/product/266/36304
- **SDK Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717`

### Request Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| Name | String | No | Application name filter |
| Tags | Array of ResourceTag | No | Tag filter for application list |
| Offset | Integer | No | Pagination offset (default: 0) |
| Limit | Integer | No | Max results per page (default: 200, max: 200) |

### Response Structure
| Field | Type | Description |
|-------|------|-------------|
| SubAppIdInfoSet | Array of SubAppIdInfo | Application information collection |
| TotalCount | Integer | Total number of applications |

## Benefits
1. **Consistency**: Aligns with Terraform best practices of having both resource and data source for managed entities
2. **Flexibility**: Enables users to query and reference existing VOD sub-applications
3. **Automation**: Supports infrastructure discovery and dynamic configuration
4. **Tag support**: Leverages cloud-native tag-based organization

## Implementation Scope
This change affects:
- **New files**:
  - `tencentcloud/services/vod/data_source_tc_vod_sub_applications.go`
  - `tencentcloud/services/vod/data_source_tc_vod_sub_applications_test.go`
  - `tencentcloud/services/vod/data_source_tc_vod_sub_applications.md`
  
- **Modified files**:
  - `tencentcloud/services/vod/service_tencentcloud_vod.go` (add service method)
  - `tencentcloud/services/vod/extension_vod.go` (register data source)

## Alternatives Considered
1. **Extend existing resource**: Adding query functionality to the resource itself would violate Terraform's resource vs data source separation principle
2. **Use external scripts**: Users could query via CLI/API, but this defeats the purpose of infrastructure-as-code

## Dependencies
- Requires TencentCloud VOD SDK (already available in vendor)
- No breaking changes to existing resources
- No new external dependencies

## Success Criteria
1. Users can query VOD sub-applications with optional filters
2. Pagination works correctly for large result sets
3. Tag filtering returns accurate results
4. All fields from `SubAppIdInfo` are properly mapped
5. Documentation includes practical usage examples
6. Tests cover common query scenarios

## Risks and Mitigations
| Risk | Impact | Mitigation |
|------|--------|-----------|
| API rate limiting | Medium | Implement proper retry logic with backoff |
| Field compatibility | Low | All fields come from stable SDK structures |
| Empty result handling | Low | Follow existing pattern from other VOD data sources |

## Timeline
Estimated implementation: 1-2 days
- Day 1: Implement data source, service method, and tests
- Day 2: Documentation, validation, and final testing

## References
- Existing VOD data source: `data_source_tc_vod_super_player_configs.go`
- VOD resource implementation: `resource_tc_vod_sub_application.go`
- API documentation: https://cloud.tencent.com/document/product/266/36304
- Related change: `align-vod-sub-application-params` (context on SubAppIdInfo structure)
