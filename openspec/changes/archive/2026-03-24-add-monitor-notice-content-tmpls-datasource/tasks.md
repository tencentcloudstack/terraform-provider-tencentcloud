# Implementation Tasks: Monitor Notice Content Templates DataSource

## Task Breakdown

### Phase 1: Service Layer Implementation ⏱️ 30 min

#### Task 1.1: Add Service Method
**File**: `tencentcloud/services/monitor/service_tencentcloud_monitor.go`

- [ ] Add new method `DescribeNoticeContentTmplsByFilter`
- [ ] Method signature with context and param map
- [ ] Return templates and bind policy counts
- [ ] Add proper logging (defer with error check)

**Acceptance Criteria**:
- Method compiles without errors
- Follows naming conventions
- Includes error logging

---

#### Task 1.2: Implement Parameter Parsing
**File**: `tencentcloud/services/monitor/service_tencentcloud_monitor.go`

- [ ] Parse `TmplIDs` from paramMap
- [ ] Parse `TmplName` from paramMap
- [ ] Parse `NoticeID` from paramMap
- [ ] Parse `TmplLanguage` from paramMap
- [ ] Parse `MonitorType` from paramMap

**Acceptance Criteria**:
- All parameters properly type-casted
- Nil checks for optional parameters
- Parameters set on request object

---

#### Task 1.3: Implement Pagination Loop
**File**: `tencentcloud/services/monitor/service_tencentcloud_monitor.go`

- [ ] Initialize pageNumber = 1, pageSize = 100
- [ ] Create for loop for pagination
- [ ] Set PageNumber and PageSize on request
- [ ] Break condition when no more data
- [ ] Increment pageNumber after each iteration

**Acceptance Criteria**:
- Loop structure correct
- Page size configurable
- Proper break condition

---

#### Task 1.4: Implement Retry Logic
**File**: `tencentcloud/services/monitor/service_tencentcloud_monitor.go`

**CRITICAL**: Retry MUST be inside the pagination loop!

- [ ] Wrap API call in `resource.Retry(tccommon.ReadRetryTimeout, ...)`
- [ ] Add `ratelimit.Check(request.GetAction())`
- [ ] Call `me.client.UseMonitorV20230616Client().DescribeNoticeContentTmpl(request)`
- [ ] Return `tccommon.RetryError(e)` on error
- [ ] Log success with request/response bodies
- [ ] Check for nil response
- [ ] Accumulate results in slice
- [ ] Store bind policy counts in map

**Acceptance Criteria**:
- Retry logic inside for loop ✅ CRITICAL
- Rate limiting check present
- Proper error handling
- Debug logging included
- Nil checks for response

**Code Pattern**:
```go
for {
    request.PageNumber = &pageNumber
    request.PageSize = &pageSize
    
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        result, e := me.client.UseMonitorV20230616Client().DescribeNoticeContentTmpl(request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        // ... accumulate results
        return nil
    })
    
    if err != nil {
        errRet = err
        return
    }
    
    // Check if more pages
    if len(currentPageResults) < int(pageSize) {
        break
    }
    
    pageNumber++
}
```

---

#### Task 1.5: Result Accumulation
**File**: `tencentcloud/services/monitor/service_tencentcloud_monitor.go`

- [ ] Append templates to result slice
- [ ] Store bind policy counts in map (TmplID as key)
- [ ] Handle nil values gracefully

**Acceptance Criteria**:
- Results accumulated correctly
- Map stores bind counts by template ID
- No data loss during pagination

---

### Phase 2: Data Source Implementation ⏱️ 60 min

#### Task 2.1: Create Data Source File
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

- [ ] Create new file
- [ ] Add package declaration
- [ ] Add imports (context, schema, helper, resource, monitorv20230616, etc.)
- [ ] Add file header comment

**Acceptance Criteria**:
- File created in correct directory
- Proper package and imports
- No syntax errors

---

#### Task 2.2: Define Schema
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

**Input Fields**:
- [ ] `tmpl_ids` - TypeSet of String, Optional
- [ ] `tmpl_name` - TypeString, Optional
- [ ] `notice_id` - TypeString, Optional
- [ ] `tmpl_language` - TypeString, Optional
- [ ] `monitor_type` - TypeString, Optional
- [ ] `result_output_file` - TypeString, Optional

**Output Fields**:
- [ ] `notice_content_tmpl_list` - TypeList, Computed
  - [ ] `tmpl_id` - TypeString, Computed
  - [ ] `tmpl_name` - TypeString, Computed
  - [ ] `monitor_type` - TypeString, Computed
  - [ ] `tmpl_language` - TypeString, Computed
  - [ ] `creator` - TypeString, Computed
  - [ ] `last_modifier` - TypeString, Computed
  - [ ] `create_time` - TypeInt, Computed
  - [ ] `update_time` - TypeInt, Computed
  - [ ] `tmpl_contents_json` - TypeString, Computed
  - [ ] `bind_policy_count` - TypeInt, Computed

**Acceptance Criteria**:
- All fields properly typed
- Descriptions clear and accurate
- Optional vs Required correctly set
- Computed fields marked

---

#### Task 2.3: Implement DataSource Function
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

- [ ] Create `DataSourceTencentCloudMonitorNoticeContentTmpls()` function
- [ ] Return `&schema.Resource` with Read function and Schema
- [ ] Set Read function to `dataSourceTencentCloudMonitorNoticeContentTmplsRead`

**Acceptance Criteria**:
- Function signature correct
- Returns valid schema.Resource
- Read function wired up

---

#### Task 2.4: Implement Read Function
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

**Function**: `dataSourceTencentCloudMonitorNoticeContentTmplsRead`

- [ ] Add defer `tccommon.LogElapsed`
- [ ] Add defer `tccommon.InconsistentCheck`
- [ ] Initialize logId, ctx, service
- [ ] Create paramMap

**Acceptance Criteria**:
- Function signature: `func(d *schema.ResourceData, meta interface{}) error`
- Proper initialization
- Service instance created

---

#### Task 2.5: Parse Input Parameters
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

- [ ] Parse `tmpl_ids` from schema (TypeSet to []*string)
- [ ] Parse `tmpl_name` from schema
- [ ] Parse `notice_id` from schema
- [ ] Parse `tmpl_language` from schema
- [ ] Parse `monitor_type` from schema
- [ ] Build paramMap with parsed values

**Acceptance Criteria**:
- All parameters parsed correctly
- Type conversions done properly
- helper.String() used for SDK pointers

**Code Pattern**:
```go
if v, ok := d.GetOk("tmpl_ids"); ok {
    tmplIDsSet := v.(*schema.Set).List()
    tmplIDs := make([]*string, 0, len(tmplIDsSet))
    for _, item := range tmplIDsSet {
        tmplIDs = append(tmplIDs, helper.String(item.(string)))
    }
    paramMap["TmplIDs"] = tmplIDs
}
```

---

#### Task 2.6: Call Service Layer with Retry
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

- [ ] Declare result variables
- [ ] Wrap service call in `resource.Retry`
- [ ] Call `service.DescribeNoticeContentTmplsByFilter(ctx, paramMap)`
- [ ] Store results in variables
- [ ] Return `tccommon.RetryError(e)` on error
- [ ] Check for error after retry

**Acceptance Criteria**:
- Retry wrapper present
- Service call correct
- Error handling proper

---

#### Task 2.7: Map Response to Schema
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

- [ ] Create result list `make([]map[string]interface{}, 0, len(respData))`
- [ ] Loop through templates
- [ ] Map each field with nil checks
- [ ] Serialize `TmplContents` to JSON string
- [ ] Look up bind policy count from map
- [ ] Append to result list
- [ ] Set `notice_content_tmpl_list` on schema

**Acceptance Criteria**:
- All fields mapped
- Nil checks for all pointers
- JSON serialization works
- Bind count lookup correct

---

#### Task 2.8: Set ID and Export
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go`

- [ ] Set data source ID using `helper.BuildToken()`
- [ ] Check for `result_output_file` parameter
- [ ] Call `tccommon.WriteToFile` if present
- [ ] Return nil on success

**Acceptance Criteria**:
- ID set correctly
- File export works (if parameter present)
- No errors on success path

---

### Phase 3: Provider Registration ⏱️ 10 min

#### Task 3.1: Register DataSource
**File**: `tencentcloud/provider.go`

- [ ] Find data sources map
- [ ] Add entry: `"tencentcloud_monitor_notice_content_tmpls": monitor.DataSourceTencentCloudMonitorNoticeContentTmpls()`
- [ ] Ensure alphabetical order
- [ ] Run `go fmt`

**Acceptance Criteria**:
- Data source registered
- Alphabetically sorted
- Code formatted

---

### Phase 4: Testing ⏱️ 30 min

#### Task 4.1: Create Test File
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls_test.go`

- [ ] Create new file
- [ ] Add package and imports
- [ ] Add test configuration constants

**Acceptance Criteria**:
- File created
- Proper imports
- Test configs defined

---

#### Task 4.2: Write Basic Acceptance Test
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls_test.go`

- [ ] Create `TestAccTencentCloudMonitorNoticeContentTmplsDataSource_basic`
- [ ] Setup test case with PreCheck and Providers
- [ ] Add test step with config
- [ ] Add checks for data source ID and list length

**Acceptance Criteria**:
- Test compiles
- Test case structure correct
- Basic checks present

**Test Config**:
```go
const testAccMonitorNoticeContentTmplsDataSource = `
data "tencentcloud_monitor_notice_content_tmpls" "test" {
}

output "tmpl_count" {
  value = length(data.tencentcloud_monitor_notice_content_tmpls.test.notice_content_tmpl_list)
}
`
```

---

#### Task 4.3: Run Tests
**Commands**:
```bash
# Format code
go fmt ./tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls*.go

# Run acceptance test
cd tencentcloud
TF_ACC=1 go test -v -run TestAccTencentCloudMonitorNoticeContentTmplsDataSource_basic ./services/monitor/

# Or run all monitor data source tests
TF_ACC=1 go test -v ./services/monitor/ -run DataSource
```

- [ ] Run `go fmt`
- [ ] Run acceptance test
- [ ] Verify test passes
- [ ] Check test output

**Acceptance Criteria**:
- Code formatted
- All tests pass
- No errors in output

---

### Phase 5: Documentation ⏱️ 30 min

#### Task 5.1: Create Documentation File
**File**: `tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.md`

- [ ] Create new file
- [ ] Add title and description
- [ ] Add usage examples
- [ ] Add argument reference
- [ ] Add attribute reference

**Acceptance Criteria**:
- Documentation complete
- Examples work
- All parameters documented

---

#### Task 5.2: Documentation Content

**Required Sections**:
- [ ] Title: `# tencentcloud_monitor_notice_content_tmpls`
- [ ] Description paragraph
- [ ] Example Usage section
  - [ ] Query all templates
  - [ ] Query by ID
  - [ ] Filter by name
  - [ ] Filter by language
  - [ ] Multiple filters
- [ ] Argument Reference section
  - [ ] All input parameters
  - [ ] Types and descriptions
- [ ] Attribute Reference section
  - [ ] Output structure
  - [ ] Nested attributes

**Acceptance Criteria**:
- All sections present
- Examples tested
- Clear descriptions

---

## Task Summary

### Progress Tracking

| Phase | Tasks | Status | Time |
|-------|-------|--------|------|
| 1. Service Layer | 5 tasks | ⬜ Pending | 30 min |
| 2. Data Source | 8 tasks | ⬜ Pending | 60 min |
| 3. Registration | 1 task | ⬜ Pending | 10 min |
| 4. Testing | 3 tasks | ⬜ Pending | 30 min |
| 5. Documentation | 2 tasks | ⬜ Pending | 30 min |
| **Total** | **19 tasks** | **0% Complete** | **2.5 hrs** |

---

## Verification Checklist

### Code Quality ✅

- [ ] Code follows reference implementation pattern
- [ ] All functions have proper error handling
- [ ] Logging statements present (debug and error)
- [ ] Nil checks for all pointers
- [ ] Code formatted with `go fmt`
- [ ] No linter warnings

### Functionality ✅

- [ ] Can query all templates (no filters)
- [ ] Can filter by template IDs
- [ ] Can filter by template name
- [ ] Can filter by notice ID
- [ ] Can filter by language
- [ ] Can filter by monitor type
- [ ] Pagination works (tested with 100+ templates)
- [ ] Retry logic works on transient errors
- [ ] Result export works

### Testing ✅

- [ ] Acceptance test passes
- [ ] Manual testing successful
- [ ] Edge cases handled (empty results, nil values)

### Documentation ✅

- [ ] README.md complete
- [ ] proposal.md complete
- [ ] tasks.md complete (this file)
- [ ] API documentation file complete
- [ ] Examples work

---

## Critical Notes

### ⚠️ CRITICAL: Retry Logic Placement

The retry logic **MUST** be inside the pagination loop:

```go
// ✅ CORRECT
for {
    request.PageNumber = &pageNumber
    request.PageSize = &pageSize
    
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        // API call here
        return nil
    })
    
    if err != nil {
        return
    }
    
    pageNumber++
}

// ❌ WRONG
err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    for {
        // API call here - NO RETRY ON SUBSEQUENT PAGES!
    }
    return nil
})
```

### 📋 Code Pattern Reference

Always refer to:
- `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`
- `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/service_tencentcloud_igtm.go`

---

## Implementation Order

**Recommended Order**:
1. Service layer (Phase 1) - Foundation
2. Data source (Phase 2) - Core logic
3. Registration (Phase 3) - Make it discoverable
4. Testing (Phase 4) - Verify it works
5. Documentation (Phase 5) - Make it usable

**Dependencies**:
- Phase 2 depends on Phase 1
- Phase 3 depends on Phase 2
- Phase 4 depends on Phase 3
- Phase 5 can be done in parallel with Phase 4

---

**Tasks Created**: 2026-03-24  
**Status**: Ready for Implementation  
**Total Estimated Time**: 2.5 - 3 hours
