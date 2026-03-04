# Implementation Notes

## SDK Version Issue

### Problem
The required API endpoints are documented for Monitor API v2023-06-16:
- `CreateNoticeContentTmpl`
- `DescribeNoticeContentTmpl`
- `ModifyNoticeContentTmpl`
- `DeleteNoticeContentTmpls`

However, the Tencent Cloud Go SDK in the vendor directory only contains:
- `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724`

The v20230616 version does not exist in the current SDK.

### Solutions

#### Option 1: Wait for SDK Update (Recommended)
- Contact Tencent Cloud SDK team to add v20230616 APIs to the Go SDK
- Update `go.mod` to use newer SDK version when available
- This is the cleanest and most maintainable solution

#### Option 2: Implement Custom API Calls
- Use the generic HTTP client to make direct API calls
- Handle request signing manually
- Less maintainable and more error-prone

#### Option 3: Check if APIs exist in v20180724
- Verify if these APIs were back-ported to v20180724
- May have different field names or structures

### Current Implementation Status

✅ **Completed**:
1. Resource file created with correct structure
2. Service method implemented
3. Client connection methods added (UseMonitorV20230616Client)
4. Provider registration complete
5. Documentation created

❌ **Blocked**:
- Cannot compile due to missing SDK package
- Need SDK update from `github.com/tencentcloud/tencentcloud-sdk-go`

### Next Steps

1. **Verify SDK availability**:
   ```bash
   go get github.com/tencentcloud/tencentcloud-sdk-go@latest
   go mod tidy
   ```

2. **Check if v20230616 is available in latest SDK**

3. **If not available**:
   - File issue with Tencent Cloud SDK team
   - Request addition of v20230616 Monitor APIs
   - Provide API documentation links

4. **Once SDK is updated**:
   - Update vendor dependencies
   - Fix import paths if needed
   - Run linter checks
   - Test with real Tencent Cloud account

### Files Modified

- ✅ `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl.go`
- ✅ `tencentcloud/services/monitor/resource_tc_monitor_notice_content_tmpl.md`
- ✅ `tencentcloud/services/monitor/service_tencentcloud_monitor.go`
- ✅ `tencentcloud/connectivity/client.go`
- ✅ `tencentcloud/provider.go`

### Code Structure

The implementation follows best practices from `resource_tc_igtm_strategy.go`:
- Composite ID format: `tmplID#tmplName`
- Proper context lifecycle management
- Retry logic with error handling
- JSON marshaling for complex nested structures
- Nil safety checks

### Alternative Approach (If SDK not available)

If the SDK is not updated in time, we could implement a simplified version:

```go
// Use generic API call through common client
request := tchttp.NewCommonRequest("monitor", "2023-06-16", "CreateNoticeContentTmpl")
request.SetOctetStreamParameters(requestBody)
response, err := client.SendCommonRequest(request)
```

However, this approach loses type safety and IDE support.

## Recommendation

**Action Required**: Update SDK before this feature can be used.

Contact: Tencent Cloud SDK Team
Repository: https://github.com/TencentCloud/tencentcloud-sdk-go
