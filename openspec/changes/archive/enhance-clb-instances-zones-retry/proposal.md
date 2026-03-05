# Enhance CLB Instances Data Source - Add Zones Field and Retry Logic

**Status**: Completed  
**Author**: Terraform Provider Team  
**Created**: 2026-03-05  
**Completed**: 2026-03-05  
**Change Type**: Enhancement  
**Changelog**: .changelog/3843.txt

## Problem Statement

The `tencentcloud_clb_instances` data source currently has two limitations:

1. **Missing Zones Field**: The CLB API `DescribeLoadBalancers` returns a `Zones` field (introduced for VPC internal load balancers with nearby access mode), but this field is not exposed in the data source's `clb_list` output. Users cannot retrieve zone information where rules are deployed.

2. **Lack of Retry Logic**: The `DescribeLoadBalancerByFilter` function uses a `for` loop to paginate API calls, but these calls are not wrapped in Terraform's retry mechanism. A single API failure during pagination will cause the entire data source read operation to fail, reducing reliability.

## Proposed Solution

### 1. Add Zones Field to Output Schema

Add a new computed field `zones` to the `clb_list` schema:

```go
"zones": {
    Type:        schema.TypeList,
    Computed:    true,
    Elem:        &schema.Schema{Type: schema.TypeString},
    Description: "Zones where rules are deployed for VPC internal load balancers with nearby access mode. Note: This field may return null, indicating no valid values can be obtained.",
},
```

**Data Mapping**: 
- Source: `LoadBalancer.Zones` (`[]*string` from CLB SDK)
- Target: `zones` field in `clb_list` (list of strings)

### 2. Wrap Pagination Loop with Retry Logic

Refactor the `DescribeLoadBalancerByFilter` function to wrap the entire pagination loop within `resource.Retry()`:

**Current Implementation** (no retry):
```go
for {
    request.Offset = &(offset)
    request.Limit = &(pageSize)
    ratelimit.Check(request.GetAction())
    response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
    if err != nil {
        errRet = errors.WithStack(err)
        return
    }
    // ... append results and check pagination
}
```

**Proposed Implementation** (with retry):
```go
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
```

## Benefits

1. **Better Data Completeness**: Users can now retrieve zone information for their CLB instances, which is crucial for understanding resource deployment in multi-AZ scenarios.

2. **Improved Reliability**: Retry logic reduces the impact of transient network issues or API rate limiting, preventing entire operations from failing due to a single API hiccup.

3. **Consistency**: Aligns with Terraform best practices - the data source reader (`dataSourceTencentCloudClbInstancesRead`) already uses retry, but the underlying service function did not.

## Compatibility

- **Backward Compatible**: ✅ Yes
  - Adding a new optional computed field does not break existing configurations
  - Retry logic is internal and transparent to users
  
- **Breaking Changes**: ❌ None

## Testing Requirements

1. **Unit Tests**: Not required (existing service tests cover pagination logic)

2. **Acceptance Tests**: Update existing test case to verify:
   - `zones` field is correctly populated for CLB instances with zone information
   - Existing fields remain unaffected
   
3. **Manual Testing**: 
   - Test with CLB instances that have `Zones` populated
   - Test with CLB instances where `Zones` is null
   - Verify retry behavior with simulated API errors (optional)

## Documentation

Update the data source documentation to include the new `zones` field in the `clb_list` attributes section.

## Implementation Plan

See [tasks.md](./tasks.md) for detailed implementation steps.

## References

- CLB API Documentation: [DescribeLoadBalancers](https://cloud.tencent.com/document/product/214/30685)
- Related Field: `LoadBalancer.Zones` - "私有网络内网负载均衡,就近接入模式下规则所落在的可用区"
