# Summary: Add VOD Sub Applications Data Source

## Change Overview
**Change ID**: `add-vod-sub-applications-datasource`  
**Type**: New Feature - Data Source  
**Status**: Proposal Stage  
**Complexity**: Low-Medium  
**Estimated Effort**: 1-2 days

## What's Being Added
A new Terraform data source `tencentcloud_vod_sub_applications` that enables users to query and discover existing Tencent Cloud VOD sub-applications programmatically.

## Problem Statement
The provider currently has a `tencentcloud_vod_sub_application` resource for managing VOD sub-applications, but lacks the corresponding data source for querying existing applications. This creates gaps in:
- Infrastructure discovery
- Referencing existing sub-applications in configurations
- Tag-based filtering and organization
- Programmatic application enumeration

## Solution
Implement a comprehensive data source that wraps the `DescribeSubAppIds` API with:

### Query Capabilities
- **Name filtering**: Query by exact application name
- **Tag filtering**: Filter applications by resource tags
- **Pagination**: Handle large result sets with offset/limit parameters
- **Full data access**: Return all available fields from the API

### Data Fields Returned
| Field | Type | Description |
|-------|------|-------------|
| `sub_app_id` | Integer | Application unique ID |
| `sub_app_id_name` | String | Current application name |
| `name` | String | Legacy name (backward compatibility) |
| `description` | String | Application description |
| `create_time` | String | ISO 8601 creation timestamp |
| `status` | String | On/Off/Destroying/Destroyed |
| `mode` | String | fileid or fileid+path mode |
| `storage_regions` | List | Enabled storage regions |
| `tags` | Map | Resource tags |

## Key Use Cases

### 1. Query All Sub-Applications
```hcl
data "tencentcloud_vod_sub_applications" "all" {}

output "app_count" {
  value = length(data.tencentcloud_vod_sub_applications.all.sub_application_info_set)
}
```

### 2. Find Application by Name
```hcl
data "tencentcloud_vod_sub_applications" "prod_app" {
  name = "ProductionVideoApp"
}

resource "tencentcloud_vod_super_player_config" "config" {
  sub_app_id = data.tencentcloud_vod_sub_applications.prod_app.sub_application_info_set[0].sub_app_id
  name       = "player-config"
}
```

### 3. Filter by Tags
```hcl
data "tencentcloud_vod_sub_applications" "by_env" {
  tags = {
    Environment = "Production"
    Team        = "VideoTeam"
  }
}
```

### 4. Pagination for Large Accounts
```hcl
data "tencentcloud_vod_sub_applications" "batch" {
  offset = 0
  limit  = 50
}
```

## Technical Implementation

### Files to Create
1. **`data_source_tc_vod_sub_applications.go`** (Data source implementation)
2. **`data_source_tc_vod_sub_applications_test.go`** (Acceptance tests)
3. **`data_source_tc_vod_sub_applications.md`** (User documentation)

### Files to Modify
1. **`service_tencentcloud_vod.go`** (Add `DescribeSubApplicationsByFilter` method)
2. **`extension_vod.go`** (Register new data source)

### Implementation Phases
1. **Service Layer** (5 tasks): API integration with pagination and retry logic
2. **Data Source** (8 tasks): Schema definition and state mapping
3. **Registration** (2 tasks): Provider integration
4. **Testing** (6 tasks): Unit and acceptance tests
5. **Documentation** (4 tasks): User guide and examples
6. **Quality** (5 tasks): Linting, formatting, validation
7. **Integration** (2 tasks): Final validation

**Total: 32 tasks**

## Requirements Summary
The spec defines 9 key requirements:

| ID | Requirement | Priority |
|----|-------------|----------|
| DS-VOD-SUBAPPS-001 | Query by name filter | High |
| DS-VOD-SUBAPPS-002 | Query by tags | High |
| DS-VOD-SUBAPPS-003 | Support pagination | Medium |
| DS-VOD-SUBAPPS-004 | Return complete information | High |
| DS-VOD-SUBAPPS-005 | Export to JSON file | Low |
| DS-VOD-SUBAPPS-006 | Handle empty results | Medium |
| DS-VOD-SUBAPPS-007 | API retry logic | High |
| DS-VOD-SUBAPPS-008 | Rate limiting | Medium |
| DS-VOD-SUBAPPS-009 | Resource integration | High |

## API Details
- **API Name**: `DescribeSubAppIds`
- **API Version**: `2018-07-17`
- **API Documentation**: https://cloud.tencent.com/document/product/266/36304
- **SDK Package**: Already available in vendor directory
- **Rate Limit**: 100 requests/second
- **Pagination**: Max 200 results per request

## Benefits
1. ✅ **Terraform Best Practice**: Completes resource/data source pair
2. ✅ **Infrastructure Discovery**: Query existing applications programmatically
3. ✅ **Dynamic Configuration**: Reference applications in other resources
4. ✅ **Tag Organization**: Leverage cloud-native tagging
5. ✅ **No Breaking Changes**: Pure addition, no impact on existing code

## Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|-----------|
| API Rate Limiting | 🟡 Medium | Implement retry with exponential backoff |
| Field Compatibility | 🟢 Low | All fields from stable SDK |
| Empty Results | 🟢 Low | Follow existing data source patterns |
| Pagination Bugs | 🟢 Low | Use proven patterns from other VOD sources |

## Dependencies
- ✅ TencentCloud VOD SDK (already in vendor)
- ✅ No new external dependencies
- ✅ No breaking changes
- ✅ Compatible with existing resources

## Testing Strategy
1. **Basic Query**: Query all applications without filters
2. **Name Filter**: Query by exact name match
3. **Tag Filter**: Query by tag key-value pairs
4. **Pagination**: Test offset/limit parameters
5. **Empty Results**: Verify graceful handling
6. **Integration**: Use data source with other resources

## Success Metrics
- [ ] All 32 implementation tasks completed
- [ ] All acceptance tests pass
- [ ] Documentation includes 4+ usage examples
- [ ] Zero linting errors
- [ ] API pagination handles 1000+ applications
- [ ] Query latency < 5 seconds for typical queries

## Timeline
- **Proposal Review**: 1 day
- **Implementation**: 1-2 days
- **Testing & Documentation**: 0.5 day
- **Code Review**: 0.5 day
- **Total**: ~3-4 days

## Related Work
- **Reference**: `data_source_tc_vod_super_player_configs.go` (similar pattern)
- **Context**: `align-vod-sub-application-params` (SubAppIdInfo structure knowledge)
- **Resource**: `resource_tc_vod_sub_application.go` (resource counterpart)

## Next Steps
1. ✅ Proposal created and validated
2. ⏳ Await proposal approval
3. ⏳ Begin Phase 1: Service layer implementation
4. ⏳ Continue through 7 implementation phases
5. ⏳ Final validation and merge

## Questions for Review
None - the proposal is comprehensive and ready for approval.

---

**Created**: 2025-02-02  
**Author**: AI Assistant  
**Reviewer**: TBD
