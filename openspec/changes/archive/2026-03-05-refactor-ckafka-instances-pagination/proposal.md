# Refactor CKafka Instances Data Source - Deprecate Pagination Parameters and Encapsulate API Logic

**Status**: Completed  
**Author**: Terraform Provider Team  
**Created**: 2026-03-05  
**Completed**: 2026-03-05  
**Change Type**: Enhancement / Refactoring  
**Changelog**: .changelog/3844.txt

## Problem Statement

The `tencentcloud_ckafka_instances` data source currently has design issues that expose internal implementation details to users:

1. **Exposed Pagination Parameters**: The data source exposes `offset` and `limit` parameters (lines 60-71 in `data_source_tc_ckafka_instances.go`), allowing users to manually control pagination. This creates several problems:
   - Users must manually implement pagination logic to retrieve all instances
   - It's inconsistent with Terraform best practices where data sources should return complete datasets
   - The default limit of 10 instances may cause users to miss data without realizing it

2. **Lack of Service Layer Abstraction**: The data source directly calls `DescribeInstancesDetail` API (line 309) without going through a service layer function. This violates the established pattern in the provider where:
   - Data sources should call service layer functions
   - Service functions should handle pagination automatically
   - This pattern is already used in other data sources (e.g., CLB, IGTM)

## Proposed Solution

### 1. Deprecate Offset and Limit Parameters

Mark both `offset` and `limit` fields as deprecated in the schema, and **remove their `Default` attributes** since these parameters no longer have any effect:

```go
"offset": {
    Type:        schema.TypeInt,
    Optional:    true,
    Deprecated:  "This parameter is deprecated and will be removed in a future version. The data source now automatically retrieves all instances.",
    Description: "The page start offset, default is `0`.",
},
"limit": {
    Type:        schema.TypeInt,
    Optional:    true,
    Deprecated:  "This parameter is deprecated and will be removed in a future version. The data source now automatically retrieves all instances.",
    Description: "The number of pages, default is `10`.",
},
```

**Note**: We use deprecation rather than direct removal to maintain backward compatibility. The `Default` attributes are removed to avoid misleading users that these parameters still have an effect.

### 2. Create Service Layer Function

Add a new `DescribeInstancesByFilter` function in `service_tencentcloud_ckafka.go` following the established pattern:

```go
func (me *CkafkaService) DescribeInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret []*ckafka.InstanceDetail, errRet error) {
    var (
        logId    = tccommon.GetLogId(ctx)
        request  = ckafka.NewDescribeInstancesDetailRequest()
        response = ckafka.NewDescribeInstancesDetailResponse()
    )

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    // Build request from param map
    for k, v := range param {
        if k == "instance_ids" {
            request.InstanceIdList = v.([]*string)
        }
        if k == "search_word" {
            request.SearchWord = v.(*string)
        }
        if k == "tag_key" {
            request.TagKey = v.(*string)
        }
        if k == "status" {
            request.Status = v.([]*int64)
        }
        if k == "filters" {
            request.Filters = v.([]*ckafka.Filter)
        }
    }

    // Automatic pagination with retry logic
    var (
        offset int64 = 0
        limit  int64 = 100
    )
    for {
        request.Offset = &offset
        request.Limit = &limit
        err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
            ratelimit.Check(request.GetAction())
            result, e := me.client.UseCkafkaClient().DescribeInstancesDetail(request)
            if e != nil {
                return tccommon.RetryError(e)
            } else {
                log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
            }

            if result == nil || result.Response == nil || result.Response.Result == nil || result.Response.Result.InstanceList == nil {
                return resource.NonRetryableError(fmt.Errorf("Describe instances failed, Response is nil."))
            }

            response = result
            return nil
        })

        if err != nil {
            errRet = err
            return
        }

        ret = append(ret, response.Response.Result.InstanceList...)
        if len(response.Response.Result.InstanceList) < int(limit) {
            break
        }

        offset += limit
    }

    return
}
```

### 3. Refactor Data Source Read Function

Update `dataSourceTencentCloudCkafkaInstancesRead` to use the new service function:

**Before** (direct API call):
```go
response, err := ckafkaService.client.UseCkafkaClient().DescribeInstancesDetail(request)
if err != nil {
    return err
}
var kafkaInstanceDetails []*ckafka.InstanceDetail
if response.Response.Result != nil {
    kafkaInstanceDetails = response.Response.Result.InstanceList
}
```

**After** (service layer call):
```go
// Build param map
param := make(map[string]interface{})
if v, ok := d.GetOk("instance_ids"); ok {
    param["instance_ids"] = helper.InterfacesStringsPoint(v.([]interface{}))
}
// ... other parameters ...

// Call service function with retry and pagination handled automatically
kafkaInstanceDetails, err := ckafkaService.DescribeInstancesByFilter(ctx, param)
if err != nil {
    return err
}
```

## Benefits

1. **Improved User Experience**: Users no longer need to worry about pagination - all instances are automatically retrieved.

2. **Consistency**: Aligns with Terraform best practices and other data sources in the provider (e.g., `tencentcloud_clb_instances`, IGTM data sources).

3. **Better Reliability**: Centralized retry logic in the service layer improves resilience against transient API failures.

4. **Maintainability**: Service layer abstraction makes it easier to modify API call logic without touching multiple data sources.

5. **Backward Compatibility**: Deprecation warnings guide users to remove obsolete parameters while existing configurations continue to work.

## Compatibility

- **Backward Compatible**: ✅ Yes
  - Existing configurations with `offset` and `limit` will continue to work
  - Deprecation warnings will guide users to update their configurations
  - The new automatic pagination will retrieve all results regardless of these parameters
  
- **Breaking Changes**: ❌ None (parameters are deprecated, not removed)

## Migration Path

Users with configurations like:
```hcl
data "tencentcloud_ckafka_instances" "example" {
  offset = 0
  limit  = 50
}
```

Will see deprecation warnings and should update to:
```hcl
data "tencentcloud_ckafka_instances" "example" {
  # offset and limit are no longer needed
  # All instances are retrieved automatically
}
```

## Testing Requirements

1. **Unit Tests**: Not required (existing service tests cover pagination logic patterns)

2. **Manual Testing**: 
   - Test with accounts having > 100 instances to verify automatic pagination
   - Test with various filter combinations
   - Verify backward compatibility with existing `offset` and `limit` parameters
   - Confirm deprecation warnings are displayed

3. **Documentation Updates**: 
   - Add deprecation notice to `offset` and `limit` parameter descriptions
   - Update examples to remove these parameters
   - Add migration notes

## Implementation Plan

See [tasks.md](./tasks.md) for detailed implementation steps.

## References

- Pattern reference: `DescribeIgtmMonitorsByFilter` in `tencentcloud/services/igtm/service_tencentcloud_igtm.go` (lines ~100-160)
- Similar refactoring: `tencentcloud_clb_instances` data source uses service layer with automatic pagination
- CKafka API Documentation: [DescribeInstancesDetail](https://cloud.tencent.com/document/product/597/40835)
