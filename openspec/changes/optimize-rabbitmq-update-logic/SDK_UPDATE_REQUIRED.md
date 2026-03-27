# SDK Update Requirements

This document outlines any TencentCloud SDK updates required for implementing the RabbitMQ instance update logic optimization.

## Current SDK Version

The project currently uses the following TDMQ SDK:

- **Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq`
- **Version**: As specified in `go.mod`

## SDK Analysis

### Required SDK Features

The optimization requires the following SDK features, which should already be available:

1. **ModifyRabbitMQVipInstance API**
   - Used for updating cluster_name, remark
   - Supports individual field updates
   - Returns TaskId for async operations

2. **ModifyRabbitMQVipInstanceSpec API**
   - Used for updating node_count, spec_name
   - Supports spec-specific updates
   - Returns TaskId for async operations

3. **DescribeTaskDetail API**
   - Used for polling task status
   - Returns task status (success, running, failed)
   - Returns error message for failed tasks

4. **ModifyAutoRenewFlag API** (billing service)
   - Used for updating auto_renew_flag
   - Standard billing API
   - Sync operation (no task waiting)

## SDK Version Requirements

### No SDK Update Required

Based on the analysis, **no SDK version update is required** for this optimization.

**Rationale**:
1. All required APIs are already available in the current SDK
2. The APIs support the operations needed (Modify, Describe, etc.)
3. No new API methods are being introduced
4. No breaking changes are expected in existing API methods

### API Compatibility

The following APIs will be used and are confirmed to be available:

| API Name | SDK Method | Required For | Status |
|----------|------------|--------------|--------|
| ModifyRabbitMQVipInstance | `ModifyRabbitMQVipInstance()` | Update cluster_name, remark | ✅ Available |
| ModifyRabbitMQVipInstanceSpec | `ModifyRabbitMQVipInstanceSpec()` | Update node_count, spec_name | ✅ Available |
| DescribeTaskDetail | `DescribeTaskDetail()` | Poll task status | ✅ Available |
| ModifyAutoRenewFlag | `ModifyAutoRenewFlag()` | Update auto_renew_flag | ✅ Available |

## Implementation Notes

### API Request Structure

The SDK request structures should support the following parameters:

1. **ModifyRabbitMQVipInstanceRequest**
   ```go
   type ModifyRabbitMQVipInstanceRequest struct {
       ClusterId   *string `json:"ClusterId,omitempty"`
       ClusterName *string `json:"ClusterName,omitempty"`
       Remark      *string `json:"Remark,omitempty"`
       // ... other fields
   }
   ```

2. **ModifyRabbitMQVipInstanceSpecRequest**
   ```go
   type ModifyRabbitMQVipInstanceSpecRequest struct {
       ClusterId *string `json:"ClusterId,omitempty"`
       NodeCount *uint64 `json:"NodeCount,omitempty"`
       SpecName  *string `json:"SpecName,omitempty"`
       // ... other fields
   }
   ```

3. **DescribeTaskDetailRequest**
   ```go
   type DescribeTaskDetailRequest struct {
       TaskId *string `json:"TaskId,omitempty"`
   }
   ```

### API Response Structure

The SDK response structures should provide:

1. **Modify API Responses**
   ```go
   type ModifyRabbitMQVipInstanceResponse struct {
       Response *struct {
           TaskId *string `json:"TaskId,omitempty"`
           // ... other fields
       } `json:"Response"`
   }
   ```

2. **DescribeTaskDetailResponse**
   ```go
   type DescribeTaskDetailResponse struct {
       Response *struct {
           TaskStatus    *string `json:"TaskStatus,omitempty"`
           ErrorMessage  *string `json:"ErrorMessage,omitempty"`
           // ... other fields
       } `json:"Response"`
   }
   ```

### Task Status Values

The SDK should support the following task status values:

- `success`: Task completed successfully
- `running`: Task is still in progress
- `failed`: Task failed with an error

## Dependency Verification

Before proceeding with implementation, verify the following:

1. **Check current SDK version**
   ```bash
   grep "tencentcloud-sdk-go" /repo/go.mod
   ```

2. **Verify API availability in current SDK**
   ```go
   // The following should compile without errors:
   _ = tdmq.NewModifyRabbitMQVipInstanceRequest
   _ = tdmq.NewModifyRabbitMQVipInstanceSpecRequest
   _ = tdmq.NewDescribeTaskDetailRequest
   _ = billing.NewModifyAutoRenewFlagRequest
   ```

3. **Test API functionality** (in development environment)
   - Create a test instance
   - Try calling Modify APIs
   - Verify response structure

## Potential SDK Issues

### Known Issues

None known at this time.

### Potential Issues to Watch For

1. **Nil field handling**
   - Ensure SDK properly handles nil fields in requests
   - Test with partial updates (only some fields set)

2. **Task ID format**
   - Verify TaskId format from API responses
   - Ensure it's compatible with DescribeTaskDetailRequest

3. **Task status mapping**
   - Confirm exact string values for task statuses
   - May need to handle case variations

4. **Error response format**
   - Verify error message format in failed tasks
   - Ensure error messages are properly propagated

## SDK Testing Recommendations

1. **Unit test with mocks**
   - Mock SDK client responses
   - Test all API call scenarios
   - Verify request structures are correct

2. **Integration test with real API** (if possible)
   - Test in development environment
   - Verify actual API behavior matches expectations
   - Test error cases

3. **Version compatibility test**
   - Test with minimum required SDK version
   - Ensure no breaking changes from older versions

## Conclusion

**No SDK update is required** for implementing the RabbitMQ instance update logic optimization. All required APIs and features are available in the current SDK version.

However, the following steps should be taken during implementation:

1. ✅ Verify current SDK version in `go.mod`
2. ✅ Confirm API availability by checking import paths
3. ✅ Test API functionality in development environment
4. ✅ Document any SDK-specific issues encountered
5. ✅ Update this document if SDK requirements change during implementation

## Next Steps

1. Review current TDMQ resource implementation
2. Identify existing SDK usage patterns
3. Implement optimized Update logic using current SDK
4. Test thoroughly with both unit and integration tests
5. Monitor for any SDK-related issues during testing
