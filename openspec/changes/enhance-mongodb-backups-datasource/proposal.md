# Proposal: Enhance MongoDB Backups Data Source

## Change ID
`enhance-mongodb-backups-datasource`

## Overview
Enhance the existing `tencentcloud_mongodb_instance_backups` data source to include missing fields from the `DescribeDBBackups` API and add pagination support. The current implementation is missing several important fields (BackId, DeleteTime, BackupRegion, RestoreTime) and doesn't support pagination parameters (Limit, Offset).

## Problem Statement
The current `tencentcloud_mongodb_instance_backups` data source only exposes a subset of fields available in the MongoDB `DescribeDBBackups` API response. This limits users' ability to:

1. **Identify specific backups** - Missing `back_id` field prevents programmatic backup identification
2. **Plan backup lifecycle** - Missing `delete_time` and `restore_time` fields prevent backup retention planning
3. **Cross-region backup management** - Missing `backup_region` field prevents identifying backup locations
4. **Control query size** - Missing pagination parameters may cause performance issues with large backup lists

### Current Gaps

**Missing Response Fields:**
- `back_id` (Integer) - Backup record ID for unique identification
- `delete_time` (String) - Backup deletion time for lifecycle management
- `backup_region` (String) - Backup location for cross-region scenarios
- `restore_time` (String) - Supported restore time for point-in-time recovery

**Missing Request Parameters:**
- `limit` (Integer) - Page size control (max 100, default: all)
- `offset` (Integer) - Pagination offset (min 0, default: 0)

## Proposed Solution
Enhance the existing data source by:

1. **Add missing fields to schema** - Include all fields from the API response
2. **Add pagination parameters** - Support `limit` and `offset` for query control
3. **Update service layer** - Pass pagination parameters to API calls
4. **Update documentation** - Document new fields and parameters with examples

### Design Decisions

1. **Backward Compatibility**: All new fields are `Computed`, and pagination parameters are `Optional`, ensuring no breaking changes
2. **Field Naming**: Follow existing Terraform conventions (snake_case) and existing patterns in the data source
3. **Pagination Default**: Service layer already implements automatic pagination; new parameters allow user control when needed

## Implementation Scope

### Changes Required

1. **Schema Changes** (`data_source_tc_mongodb_instance_backups.go`)
   - Add input parameters: `limit`, `offset`
   - Add output fields: `back_id`, `delete_time`, `backup_region`, `restore_time`

2. **Service Layer** (`service_tencentcloud_mongodb.go`)
   - Update `DescribeMongodbInstanceBackupsByFilter` to accept pagination parameters
   - Update field mapping to include new response fields

3. **Documentation** (`data_source_tc_mongodb_instance_backups.md`)
   - Document new input parameters with usage examples
   - Document new output fields with descriptions
   - Add pagination example

4. **Tests** (`data_source_tc_mongodb_instance_backups_test.go`)
   - Add test case for pagination parameters
   - Verify new fields are properly populated

### Out of Scope
- Changing existing field types or behaviors
- Adding filtering capabilities beyond what API supports
- Creating new data sources or resources

## Dependencies
- **API Version**: MongoDB API v20190725 (already in use)
- **SDK**: tencentcloud-sdk-go mongodb package (no version change needed)
- **Related Changes**: None - standalone enhancement

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Breaking existing configurations | High | All new fields are Computed/Optional |
| Pagination behavior change | Medium | Keep default behavior (query all) unchanged |
| SDK field type mismatch | Low | Verify SDK types match API documentation |

## Success Criteria
- [x] All missing API fields are exposed in the data source
- [x] Pagination parameters work as documented
- [x] Existing configurations continue to work without changes
- [x] Documentation is complete with examples
- [x] Tests pass for new functionality
- [x] No breaking changes to existing behavior

## Timeline Estimate
- Schema updates: 1 hour
- Service layer updates: 1 hour
- Documentation: 30 minutes
- Testing: 1 hour
- **Total: ~3.5 hours**

## References
- API Documentation: https://cloud.tencent.com/document/api/240/38574
- Current Implementation: `tencentcloud/services/mongodb/data_source_tc_mongodb_instance_backups.go`
- Service Layer: `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go`
