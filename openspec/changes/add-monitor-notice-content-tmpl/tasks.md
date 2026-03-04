# Implementation Tasks

## ⚠️ IMPORTANT: SDK Dependency Issue

**Status**: Implementation is BLOCKED by SDK version availability.

The required APIs are documented for Monitor API v2023-06-16, but the Tencent Cloud Go SDK currently only provides `monitor/v20180724`. The v20230616 package does not exist in the vendor directory.

**Required APIs**:
- `CreateNoticeContentTmpl`
- `DescribeNoticeContentTmpl`
- `ModifyNoticeContentTmpl`
- `DeleteNoticeContentTmpls`

**Action Required**: Update `github.com/tencentcloud/tencentcloud-sdk-go` to include monitor/v20230616 APIs.

## 1. Implementation

### Resource Development
- [x] Create `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl.go`
  - [x] Define schema with all required and optional fields
  - [x] Implement placeholder CRUD functions with error messages
  - [x] Add import support with `schema.ImportStatePassthrough`
  - ⚠️ **BLOCKED**: Actual API calls cannot be implemented without SDK

### Service Layer
- [x] Add service methods to `tencentcloud/services/monitor/service_tencentcloud_monitor.go`
  - [x] Implement `DescribeMonitorNoticeContentTmplById` placeholder method
  - [x] Add proper error handling and logging
  - ⚠️ **BLOCKED**: Cannot call API without SDK

### Provider Registration
- [x] Register resource in `tencentcloud/provider.go`
  - [x] Add `"tencentcloud_monitor_notice_content_tmpl": monitor.ResourceTencentCloudMonitorNoticeContentTmpl()`

### Documentation
- [x] Create resource documentation markdown file
  - [x] Add usage examples
  - [x] Document all schema fields with descriptions
  - [x] Add import examples

## 2. Testing

- [ ] Create basic acceptance test
- [ ] Test create operation
- [ ] Test read operation
- [ ] Test update operation
- [ ] Test delete operation
- [ ] Test import functionality
- [ ] Test composite ID parsing

⚠️ **BLOCKED**: Cannot test without SDK implementation

## 3. Validation

- [x] Run `go fmt` on new files
- [x] Run linter checks (1 deprecation hint only, acceptable)
- [x] Verify no breaking changes to existing code
- [ ] Test with real Tencent Cloud account

⚠️ **BLOCKED**: Cannot test against real account without SDK

## Implementation Notes

### SDK Information
- **Required Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616`
- **API Version**: 2023-06-16
- **Import alias**: `monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"`
- **Current Status**: ❌ Package does not exist in vendor

### Composite ID Format
- Use format: `tmplID#tmplName`
- Split using `tccommon.FILED_SP` constant
- Example: `ntpl-3r1spzjn#MyTemplate`

### API Mapping
- **Create**: `CreateNoticeContentTmpl` → returns `TmplID`
- **Read**: `DescribeNoticeContentTmpl` with `TmplIDs` (array) and `TmplName` parameters
- **Update**: `ModifyNoticeContentTmpl` with `TmplID`, `TmplName`, and `TmplContents`
- **Delete**: `DeleteNoticeContentTmpls` with `TmplIDs` (array)

### Reference Implementation
Follow the code patterns from:
- `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_strategy.go`
- Composite ID handling
- Nested complex object handling
- Error handling and retry logic
- Context usage pattern

## Next Steps

1. **Contact SDK Team**: Request addition of monitor/v20230616 APIs
2. **Update SDK**: Run `go get github.com/tencentcloud/tencentcloud-sdk-go@latest`
3. **Implement Full CRUD**: Replace placeholder functions with real API calls
4. **Add Type Definitions**: Use actual SDK types instead of placeholders
5. **Test**: Run acceptance tests with real Tencent Cloud account
6. **Deploy**: After successful testing

## Files Modified

✅ Completed:
- `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl.go` (placeholder)
- `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl.md`
- `tencentcloud/services/monitor/service_tencentcloud_monitor.go` (placeholder)
- `tencentcloud/provider.go`

📋 Additional Documentation:
- `openspec/changes/add-monitor-notice-content-tmpl/IMPLEMENTATION_NOTES.md`
