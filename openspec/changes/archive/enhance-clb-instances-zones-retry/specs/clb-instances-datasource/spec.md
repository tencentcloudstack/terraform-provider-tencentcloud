# Technical Specification: CLB Instances Data Source Enhancement

## Overview

This specification details the technical implementation for enhancing the `tencentcloud_clb_instances` data source with zones field support and improved retry logic.

## Background

### Current Implementation

The `tencentcloud_clb_instances` data source queries CLB instances using the `DescribeLoadBalancers` API. The implementation consists of:

1. **Data Source Layer** (`data_source_tc_clb_instances.go`): 
   - Defines schema and handles Terraform state
   - Calls service layer with retry wrapper

2. **Service Layer** (`service_tencentcloud_clb.go`):
   - Implements `DescribeLoadBalancerByFilter` function
   - Handles API pagination without retry logic

### API Changes

The CLB API's `DescribeLoadBalancers` endpoint has added a `Zones` field to the `LoadBalancer` response structure:

```go
type LoadBalancer struct {
    // ... existing fields ...
    
    // 私有网络内网负载均衡,就近接入模式下规则所落在的可用区
    // 注意:此字段可能返回 null,表示取不到有效值。
    Zones []*string `json:"Zones,omitnil,omitempty" name:"Zones"`
}
```

## Requirements

### Functional Requirements

1. **FR-1**: Expose `Zones` field in data source output
   - Field name: `zones`
   - Type: List of strings
   - Computed: Yes (no user input required)
   - Nullable: Yes (may return null from API)

2. **FR-2**: Implement retry logic for API pagination
   - Timeout: Use `tccommon.ReadRetryTimeout` constant
   - Scope: Entire pagination loop
   - Error handling: Use `tccommon.RetryError()` helper

### Non-Functional Requirements

1. **NFR-1**: Backward Compatibility
   - Existing configurations must continue to work without modification
   - No changes to input parameters
   - No changes to existing output fields

2. **NFR-2**: Performance
   - Retry logic should not significantly impact execution time for successful calls
   - Pagination efficiency should remain unchanged

3. **NFR-3**: Reliability
   - Transient API errors should be automatically retried
   - Permanent errors should fail gracefully with clear error messages

## Detailed Design

### 1. Schema Definition

**Location**: `tencentcloud/services/clb/data_source_tc_clb_instances.go`

**Implementation**:

```go
// Add to clb_list Elem Schema (after "numerical_vpc_id", before closing brace at line ~185)
"zones": {
    Type:        schema.TypeList,
    Computed:    true,
    Elem:        &schema.Schema{Type: schema.TypeString},
    Description: "Zones where rules are deployed for VPC internal load balancers with nearby access mode. Note: This field may return null, indicating no valid values can be obtained.",
},
```

**Rationale**:
- Use `TypeList` with `TypeString` elements to match API's `[]*string` type
- Mark as `Computed: true` (no user input needed)
- Description clearly indicates nullable nature

### 2. Data Mapping

**Location**: `tencentcloud/services/clb/data_source_tc_clb_instances.go` (line ~287)

**Implementation**:

```go
// Add after numerical_vpc_id mapping, before clbList append
if clbInstance.Zones != nil {
    mapping["zones"] = helper.StringsInterfaces(clbInstance.Zones)
}
```

**Rationale**:
- Check for `nil` before mapping to avoid nil pointer issues
- Use `helper.StringsInterfaces()` to convert `[]*string` to `[]interface{}`
- Consistent with other list field mappings (e.g., `clb_vips`, `security_groups`)

### 3. Service Layer Retry Logic

**Location**: `tencentcloud/services/clb/service_tencentcloud_clb.go` (line 77-126)

**Current Implementation**:

```go
func (me *ClbService) DescribeLoadBalancerByFilter(ctx context.Context, params map[string]interface{}) (clbs []*clb.LoadBalancer, errRet error) {
    logId := tccommon.GetLogId(ctx)
    request := clb.NewDescribeLoadBalancersRequest()
    
    // ... parameter handling ...
    
    offset := int64(0)
    pageSize := int64(CLB_PAGE_LIMIT)
    clbs = make([]*clb.LoadBalancer, 0)
    for {
        request.Offset = &(offset)
        request.Limit = &(pageSize)
        ratelimit.Check(request.GetAction())
        response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
        if err != nil {
            errRet = errors.WithStack(err)
            return
        }
        log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
            logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

        if response == nil || len(response.Response.LoadBalancerSet) < 1 {
            break
        }

        clbs = append(clbs, response.Response.LoadBalancerSet...)

        if int64(len(response.Response.LoadBalancerSet)) < pageSize {
            break
        }
        offset += pageSize
    }
    return
}
```

**Proposed Implementation**:

```go
func (me *ClbService) DescribeLoadBalancerByFilter(ctx context.Context, params map[string]interface{}) (clbs []*clb.LoadBalancer, errRet error) {
    logId := tccommon.GetLogId(ctx)
    request := clb.NewDescribeLoadBalancersRequest()
    
    // ... parameter handling remains the same ...
    
    offset := int64(0)
    pageSize := int64(CLB_PAGE_LIMIT)
    clbs = make([]*clb.LoadBalancer, 0)
    
    for {
        request.Offset = &(offset)
        request.Limit = &(pageSize)
        
        var response *clb.DescribeLoadBalancersResponse
        err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
            ratelimit.Check(request.GetAction())
            result, e := me.client.UseClbClient().DescribeLoadBalancers(request)
            if e != nil {
                return tccommon.RetryError(e)
            }
            log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
                logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

            if result == nil || result.Response == nil || result.Response.LoadBalancerSet == nil {
                return resource.NonRetryableError(fmt.Errorf("DescribeLoadBalancers response is nil"))
            }

            response = result
            return nil
        })

        if err != nil {
            log.Printf("[CRITAL]%s DescribeLoadBalancerByFilter failed, reason:%+v", logId, err)
            errRet = err
            return
        }

        if len(response.Response.LoadBalancerSet) < 1 {
            break
        }

        clbs = append(clbs, response.Response.LoadBalancerSet...)

        if int64(len(response.Response.LoadBalancerSet)) < pageSize {
            break
        }
        offset += pageSize
    }
    
    return
}
```

**Key Changes**:
1. Keep pagination loop structure (`for { ... }`) unchanged
2. Wrap **only the API call** with `resource.Retry()` inside each iteration
3. Declare `response` variable before retry block to capture result
4. Return `tccommon.RetryError(e)` for retryable API errors
5. Return `resource.NonRetryableError()` for nil response validation
6. Set `response` variable on successful API call
7. Check error after each retry completes and return early if failed
8. Add error logging with proper context

**Rationale**:
- Retry only the API call (not pagination logic) - aligns with best practices
- `tccommon.RetryError()` automatically retries on common retryable errors
- `resource.NonRetryableError()` prevents retry on malformed responses
- Timeout value (`tccommon.ReadRetryTimeout`) consistent with other read operations
- Pagination state (`offset`, `pageSize`) maintained outside retry for clarity
- Each page fetch is independently retried, improving reliability

### 4. Import Requirements

Ensure `resource` package is imported in `service_tencentcloud_clb.go`:

```go
import (
    // ... existing imports ...
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)
```

**Verification**: Check line 14 - this import should already exist in the file.

## Testing Strategy

### Unit Testing

Not required - service layer functions are typically tested via acceptance tests in this codebase.

### Acceptance Testing

Update or create test configuration to verify:

1. **Zones Field Population**:
   - Query CLB instances with zones information
   - Verify `zones` field is populated in output
   - Verify field is absent/null when API returns null

2. **Existing Fields Verification**:
   - Ensure all existing fields remain unaffected
   - Check backward compatibility

3. **Retry Logic** (optional manual testing):
   - Simulate API errors (e.g., via network issues)
   - Verify retry behavior and eventual success/failure

### Example Configuration

**File**: `examples/tencentcloud-clb-instances-zones/main.tf`

```hcl
data "tencentcloud_clb_instances" "example" {
  network_type = "INTERNAL"
}

output "clb_zones" {
  value = [
    for clb in data.tencentcloud_clb_instances.example.clb_list : {
      clb_id = clb.clb_id
      name   = clb.clb_name
      zones  = clb.zones
    }
  ]
  description = "CLB instances with their zones information"
}
```

## Error Handling

### API Errors

1. **Retryable Errors**: 
   - Network timeouts
   - Rate limiting (429)
   - Temporary server errors (5xx)
   - Handled by: `tccommon.RetryError()` wrapper

2. **Non-Retryable Errors**:
   - Authentication errors (401/403)
   - Invalid parameters (400)
   - Handled by: Return error immediately, fail operation

### Null/Empty Handling

- **Zones Field**: Check for `nil` before mapping to avoid panic
- **Empty Response**: Existing logic handles empty response sets

## Backward Compatibility Analysis

### Breaking Changes
❌ None

### Compatible Changes
✅ New optional computed field (`zones`)
✅ Internal retry logic (transparent to users)

### Migration Required
❌ No migration needed

### State Compatibility
✅ New field added to state without affecting existing fields
✅ Existing state files remain valid

## Performance Considerations

### Expected Impact
- **Negligible**: Retry logic only activates on API errors
- **Successful Calls**: No additional overhead
- **Failed Calls**: Automatic retries improve success rate without user intervention

### Resource Usage
- **Memory**: Minimal increase (one additional list field per CLB instance)
- **Network**: No additional API calls for successful operations
- **Execution Time**: Slightly longer on retry scenarios (acceptable trade-off for reliability)

## Security Considerations

- No authentication/authorization changes
- No sensitive data exposed in new field
- Retry logic does not log sensitive information

## Documentation Requirements

### Data Source Documentation

Update `website/docs/d/clb_instances.html.markdown`:

Add to `clb_list` attributes section:

```markdown
* `zones` - (Optional) Zones where rules are deployed for VPC internal load balancers with nearby access mode. Note: This field may return null, indicating no valid values can be obtained.
```

### Changelog

Create `.changelog/<next-number>.txt`:

```
release-note:enhancement
data-source/tencentcloud_clb_instances: add `zones` field to `clb_list` output and improve reliability with retry logic in underlying API calls
```

## References

### API Documentation
- [DescribeLoadBalancers API](https://cloud.tencent.com/document/product/214/30685)
- LoadBalancer.Zones field: "私有网络内网负载均衡,就近接入模式下规则所落在的可用区"

### SDK References
- `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317/models.go`
  - Line 6409: `Zones []*string` field definition

### Codebase Patterns
- Similar list field: `security_groups` (line 256)
- Retry pattern: Data source reader (line 223-231)
- Helper function: `helper.StringsInterfaces()` for type conversion

## Appendix: Code Locations

| Component | File | Line Range |
|-----------|------|------------|
| Data Source Schema | `data_source_tc_clb_instances.go` | 52-189 |
| Data Source Reader | `data_source_tc_clb_instances.go` | 192-307 |
| Service Function | `service_tencentcloud_clb.go` | 77-126 |
| LoadBalancer Struct | `vendor/.../models.go` | 6236-6429 |
| Zones Field Definition | `vendor/.../models.go` | 6409 |
