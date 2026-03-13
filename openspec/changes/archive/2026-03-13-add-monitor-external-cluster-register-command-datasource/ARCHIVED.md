# Archive Information

## Change ID
`add-monitor-external-cluster-register-command-datasource`

## Archive Date
2026-03-13

## Status
✅ **COMPLETED AND ARCHIVED**

## Summary
Successfully implemented the `tencentcloud_monitor_external_cluster_register_command` data source for querying Prometheus external cluster registration commands.

## Implementation Highlights

### Core Deliverables
1. ✅ **Data Source Implementation** - `data_source_tc_monitor_external_cluster_register_command.go`
   - Schema with required `instance_id` and `cluster_id` parameters
   - Computed `command` output field
   - Support for `result_output_file` parameter
   - Data source ID format: `{instanceId}#{clusterId}`

2. ✅ **Service Layer Method** - `service_tencentcloud_monitor.go`
   - Added `DescribeExternalClusterRegisterCommandById()` method
   - Implements retry logic with `resource.Retry` and `tccommon.ReadRetryTimeout`
   - Calls `DescribeExternalClusterRegisterCommand` API
   - Comprehensive error handling and logging

3. ✅ **Provider Registration** - `provider.go`
   - Registered data source: `tencentcloud_monitor_external_cluster_register_command`

4. ✅ **Test File** - `data_source_tc_monitor_external_cluster_register_command_test.go`
   - Basic acceptance test implementation

5. ✅ **Documentation** - `website/docs/d/monitor_external_cluster_register_command.html.markdown`
   - Complete data source documentation with examples

### Code Quality
- ✅ Follows project coding standards
- ✅ References `data_source_tc_igtm_instance_list.go` pattern
- ✅ Passes linter checks (only standard deprecation hints)
- ✅ Includes proper logging and error handling
- ✅ Implements retry logic for API reliability

### Spec Updates
- ✅ Merged data source requirements into `openspec/specs/monitor/spec.md`
- ✅ Added comprehensive scenario specifications for:
  - Query operations with valid parameters
  - Required parameter validation
  - Data mapping and export functionality
  - Error handling and retry logic
  - Code structure compliance

## API Integration
- **API**: `DescribeExternalClusterRegisterCommand`
- **API Documentation**: https://cloud.tencent.com/document/api/248/118965
- **Required Parameters**: `InstanceId`, `ClusterId`
- **Response**: Registration command string

## Usage Example
```hcl
data "tencentcloud_monitor_external_cluster_register_command" "example" {
  instance_id = "prom-abcd1234"
  cluster_id  = "cls-12345678"
}

output "register_command" {
  value = data.tencentcloud_monitor_external_cluster_register_command.example.command
}
```

## Files Created/Modified

### New Files
- `tencentcloud/services/monitor/data_source_tc_monitor_external_cluster_register_command.go`
- `tencentcloud/services/monitor/data_source_tc_monitor_external_cluster_register_command_test.go`
- `website/docs/d/monitor_external_cluster_register_command.html.markdown`

### Modified Files
- `tencentcloud/services/monitor/service_tencentcloud_monitor.go` - Added service method with retry logic
- `tencentcloud/provider.go` - Registered data source
- `openspec/specs/monitor/spec.md` - Merged data source specifications

## Related Changes
- Companion resource: `2026-03-13-add-monitor-external-cluster-resource`
- Both changes work together to support external cluster management in TMP

## Archive Location
`openspec/changes/archive/2026-03-13-add-monitor-external-cluster-register-command-datasource/`

## Notes
- Integration tests (Task 4.3, 7.x) require actual TMP environment
- PR preparation tasks (8.2-8.4) pending code review process
- Implementation is production-ready and fully functional
