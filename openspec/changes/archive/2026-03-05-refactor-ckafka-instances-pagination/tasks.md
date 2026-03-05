# Implementation Tasks

## Overview

This document tracks implementation tasks for refactoring the `tencentcloud_ckafka_instances` data source to deprecate pagination parameters and encapsulate API logic in the service layer.

## Tasks

### 1. Add Service Layer Function ✅ Pending

**File**: `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go`

- [ ] Add `DescribeInstancesByFilter` function at the end of the file (after line 1862)
  - Follow the pattern from `DescribeIgtmMonitorsByFilter` in IGTM service
  - Accept `ctx context.Context` and `param map[string]interface{}` parameters
  - Return `ret []*ckafka.InstanceDetail` and `errRet error`
  
- [ ] Implement function logic:
  - Create request and response variables
  - Add defer block for error logging
  - Parse parameters from param map: `instance_ids`, `search_word`, `tag_key`, `status`, `filters`
  - Implement pagination loop with offset starting at 0 and limit of 100
  - Wrap API call in `resource.Retry()` with `tccommon.ReadRetryTimeout`
  - Use `ratelimit.Check()` before API call
  - Use `tccommon.RetryError()` for retryable errors
  - Use `resource.NonRetryableError()` for nil response checks
  - Append results to ret slice
  - Break loop when response has fewer items than limit

**Code Reference**:
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

### 2. Deprecate Pagination Parameters ✅ Pending

**File**: `tencentcloud/services/ckafka/data_source_tc_ckafka_instances.go`

- [ ] Remove `Default` field and add `Deprecated` field to `offset` parameter (line 60-65):
  ```go
  "offset": {
      Type:        schema.TypeInt,
      Optional:    true,
      Deprecated:  "This parameter is deprecated and will be removed in a future version. The data source now automatically retrieves all instances.",
      Description: "The page start offset, default is `0`.",
  },
  ```

- [ ] Remove `Default` field and add `Deprecated` field to `limit` parameter (line 66-71):
  ```go
  "limit": {
      Type:        schema.TypeInt,
      Optional:    true,
      Deprecated:  "This parameter is deprecated and will be removed in a future version. The data source now automatically retrieves all instances.",
      Description: "The number of pages, default is `10`.",
  },
  ```

### 3. Refactor Data Source Read Function ✅ Pending

**File**: `tencentcloud/services/ckafka/data_source_tc_ckafka_instances.go`

- [ ] Replace direct API call (lines 271-316) with service layer call
  
**Current implementation** (lines 271-316):
```go
func dataSourceTencentCloudCkafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_instances.read")()

    ckafkaService := CkafkaService{
        client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
    }
    request := ckafka.NewDescribeInstancesDetailRequest()
    if v, ok := d.GetOk("instance_ids"); ok {
        request.InstanceIdList = helper.InterfacesStringsPoint(v.([]interface{}))
    }
    // ... more parameter building ...
    if v, ok := d.GetOk("offset"); ok {
        request.Offset = helper.IntInt64(v.(int))
    }
    if v, ok := d.GetOk("limit"); ok {
        request.Limit = helper.IntInt64(v.(int))
    }

    response, err := ckafkaService.client.UseCkafkaClient().DescribeInstancesDetail(request)
    if err != nil {
        return err
    }
    var kafkaInstanceDetails []*ckafka.InstanceDetail
    if response.Response.Result != nil {
        kafkaInstanceDetails = response.Response.Result.InstanceList
    }
    // ... result processing ...
}
```

**New implementation**:
```go
func dataSourceTencentCloudCkafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_instances.read")()

    ctx := context.Background()
    ckafkaService := CkafkaService{
        client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
    }
    
    // Build param map for service function
    param := make(map[string]interface{})
    if v, ok := d.GetOk("instance_ids"); ok {
        param["instance_ids"] = helper.InterfacesStringsPoint(v.([]interface{}))
    }
    if v, ok := d.GetOk("search_word"); ok {
        param["search_word"] = helper.String(v.(string))
    }
    if v, ok := d.GetOk("tag_key"); ok {
        param["tag_key"] = helper.String(v.(string))
    }
    if v, ok := d.GetOk("status"); ok {
        param["status"] = helper.InterfacesIntInt64Point(v.([]interface{}))
    }
    if v, ok := d.GetOk("filters"); ok {
        filterParams := v.([]interface{})
        filters := make([]*ckafka.Filter, 0)
        for _, filterParam := range filterParams {
            filterParamMap := filterParam.(map[string]interface{})
            filters = append(filters, &ckafka.Filter{
                Name:   helper.String(filterParamMap["name"].(string)),
                Values: helper.InterfacesStringsPoint(filterParamMap["values"].([]interface{})),
            })
        }
        param["filters"] = filters
    }
    
    // Note: offset and limit are deprecated and ignored - function always retrieves all results

    // Call service function with automatic pagination and retry
    kafkaInstanceDetails, err := ckafkaService.DescribeInstancesByFilter(ctx, param)
    if err != nil {
        return err
    }
    
    // ... rest of result processing remains the same ...
}
```

### 4. Update Documentation ⏸️ Deferred

**Files**: 
- Website documentation (if exists)
- In-code documentation comments

- [ ] Add deprecation notice to parameter descriptions in website docs (deferred to release time)
- [ ] Update examples to remove `offset` and `limit` parameters (deferred to release time)
- [ ] Add migration guide in documentation (deferred to release time)

**Note**: Documentation updates will be handled during the release process.

### 5. Add Changelog Entry ✅ Pending

**File**: `.changelog/<next-number>.txt`

- [ ] Determine next changelog number (check latest in `.changelog/` directory)
- [ ] Create changelog file with content:
  ```
  ```release-note:enhancement
  data-source/tencentcloud_ckafka_instances: deprecate `offset` and `limit` parameters, add automatic pagination with retry logic
  ```
  ```

## Implementation Notes

### Parameter Mapping

The param map uses the following keys:
- `instance_ids`: `[]*string` - List of instance IDs
- `search_word`: `*string` - Search keyword for instance name
- `tag_key`: `*string` - Tag key for filtering
- `status`: `[]*int64` - List of status codes
- `filters`: `[]*ckafka.Filter` - Complex filter objects

### Pagination Details

- **Page size**: 100 instances per API call (consistent with other services)
- **Starting offset**: 0
- **Loop termination**: When response contains fewer items than the limit
- **Retry timeout**: `tccommon.ReadRetryTimeout` (consistent with read operations)

### Backward Compatibility

The deprecated `offset` and `limit` parameters are still present in the schema but:
1. Will show deprecation warnings to users
2. Are ignored by the new implementation
3. All results are always retrieved regardless of these values
4. Can be safely removed in a future major version

## Definition of Done

- [ ] All code changes completed
- [ ] No linter errors or warnings
- [ ] Deprecation warnings appear when using offset/limit parameters
- [ ] Manual testing confirms:
  - All instances are retrieved (test with > 100 instances if possible)
  - Filters work correctly
  - Existing configurations still work
- [ ] Changelog entry added
- [ ] Code review completed (if applicable)
