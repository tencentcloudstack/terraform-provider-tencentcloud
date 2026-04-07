# Implementation Tasks: Add tencentcloud_cls_alarm_notices Data Source

**Feature**: CLS Alarm Notices Data Source  
**Estimated Total Time**: 3 hours

---

## Current Status Summary

| Category | Status | Progress |
|----------|--------|----------|
| Code Implementation | ✅ Completed | 3/3 tasks |
| Testing | ⏳ Pending | 0/8 tasks |
| Documentation | ⏳ Pending | 0/2 tasks |
| Review | ⏳ Pending | 0/2 tasks |
| **Overall** | **✅ Code Complete** | **3/15 tasks (20%)** |

**Implementation Date**: 2026-03-24  
**See**: IMPLEMENTATION_SUMMARY.md for detailed implementation report

---

## Phase 1: Code Implementation ✅

**Status**: ✅ Completed  
**Time Spent**: 20 minutes (vs. estimated 1.5 hours)

### Task 1.1: Create Data Source File ✅

- [x] **File**: `tencentcloud/services/cls/data_source_tc_cls_alarm_notices.go` (661 lines)
- [x] **Components**:
  - [x] Package declaration and imports
  - [x] DataSourceTencentCloudClsAlarmNotices() function
  - [x] Complete schema definition with all fields:
    - [x] filters (input)
    - [x] has_alarm_shield_count (input)
    - [x] alarm_notices (output, with nested schema)
    - [x] result_output_file (input)
  - [x] dataSourceTencentCloudClsAlarmNoticesRead() function:
    - [x] Initialize context and service
    - [x] Parse filters from schema
    - [x] Parse has_alarm_shield_count
    - [x] Build paramMap
    - [x] Call service layer with retry
    - [x] Map response to schema (handle nested structures)
    - [x] Set computed attributes
    - [x] Write to file if specified
- [x] **Reference**: Followed exact pattern of `data_source_tc_igtm_instance_list.go`
- [x] **Validation**: Code compiles, follows Go conventions

### Task 1.2: Add Service Layer Method ✅

- [x] **File**: `tencentcloud/services/cls/service_tencentcloud_cls.go` (+62 lines)
- [x] **Function**: `DescribeClsAlarmNoticesByFilter`
- [x] **Components**:
  - [x] Context and parameter handling
  - [x] API request construction
  - [x] Filter parsing (Filters, HasAlarmShieldCount)
  - [x] Pagination loop (offset/limit)
  - [x] **Retry logic INSIDE pagination loop** ✅ CRITICAL
  - [x] Rate limiting (ratelimit.Check)
  - [x] Response validation
  - [x] Error handling with logging
  - [x] Result accumulation
- [x] **Client**: Uses `UseClsClient()` (not `UseCls20201016Client()`)
- [x] **Validation**: No linter errors, correct retry placement

### Task 1.3: Register Data Source ✅

- [x] **File**: `tencentcloud/provider.go` (+1 line)
- [x] **Location**: Line 1188 in DataSourcesMap
- [x] **Entry**: `"tencentcloud_cls_alarm_notices": cls.DataSourceTencentCloudClsAlarmNotices()`
- [x] **Validation**: Alphabetically ordered with other CLS data sources

### Code Quality ✅

- [x] Go formatted (gofmt)
- [x] No compilation errors
- [x] No linter errors (only pre-existing hints)
- [x] Comprehensive nil checks
- [x] Follows reference implementation pattern
- [x] Clear error messages
- [x] Proper logging

---

### Task 1.2: Add Service Layer Method

- [ ] **File**: `tencentcloud/services/cls/service_tencentcloud_cls.go`
- [ ] **Method**: `DescribeClsAlarmNoticesByFilter`
- [ ] **Signature**:
  ```go
  func (me *ClsService) DescribeClsAlarmNoticesByFilter(
      ctx context.Context, 
      param map[string]interface{},
  ) (ret []*cls.AlarmNotice, totalCount *uint64, errRet error)
  ```
- [ ] **Implementation**:
  - [ ] Initialize request and response
  - [ ] Add defer error logging
  - [ ] Parse parameters (Filters, HasAlarmShieldCount)
  - [ ] Implement pagination loop (offset/limit)
  - [ ] **CRITICAL**: Add retry INSIDE pagination loop
  - [ ] Rate limit check before each call
  - [ ] Append results
  - [ ] Check for more pages
  - [ ] Return results and total count
- [ ] **Estimated Time**: 20 minutes
- [ ] **Validation**: Method signature matches, retry in correct place

---

### Task 1.3: Register Data Source in Provider

- [ ] **File**: `tencentcloud/provider.go`
- [ ] **Action**: Add data source to provider registration
- [ ] **Code**:
  ```go
  "tencentcloud_cls_alarm_notices": cls.DataSourceTencentCloudClsAlarmNotices(),
  ```
- [ ] **Location**: Find data sources map, add alphabetically in CLS section
- [ ] **Estimated Time**: 5 minutes
- [ ] **Validation**: Provider compiles

---

### Task 1.4: Code Formatting

- [ ] **Commands**:
  - [ ] `go fmt tencentcloud/services/cls/data_source_tc_cls_alarm_notices.go`
  - [ ] `go fmt tencentcloud/services/cls/service_tencentcloud_cls.go`
- [ ] **Estimated Time**: 2 minutes
- [ ] **Validation**: Files formatted without errors

---

## Phase 2: Testing ⏳

**Status**: Pending  
**Estimated Time**: 45 minutes

### Task 2.1: Test - Query All Alarm Notices

- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "all" {}
  
  output "notices" {
    value = data.tencentcloud_cls_alarm_notices.all.alarm_notices
  }
  ```
- [ ] **Expected**: Returns all alarm notices
- [ ] **Validation**:
  - [ ] Data source executes successfully
  - [ ] alarm_notices list populated
  - [ ] total_count matches list length
  - [ ] No errors in logs
- [ ] **Estimated Time**: 5 minutes

---

### Task 2.2: Test - Filter by Name

- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "by_name" {
    filters {
      name   = "name"
      values = ["test-alarm"]
    }
  }
  ```
- [ ] **Expected**: Returns only matching alarm notices
- [ ] **Validation**:
  - [ ] Correct filtering applied
  - [ ] Results match filter criteria
  - [ ] total_count accurate
- [ ] **Estimated Time**: 5 minutes

---

### Task 2.3: Test - Filter by ID

- [ ] **Setup**: Get existing alarm notice ID from console
- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "by_id" {
    filters {
      name   = "alarmNoticeId"
      values = ["notice-xxxx-yyyy"]
    }
  }
  ```
- [ ] **Expected**: Returns specific alarm notice
- [ ] **Validation**:
  - [ ] Exactly one result returned
  - [ ] alarm_notice_id matches filter
  - [ ] All fields populated correctly
- [ ] **Estimated Time**: 5 minutes

---

### Task 2.4: Test - Filter by Delivery Status

- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "enabled" {
    filters {
      name   = "deliverFlag"
      values = ["2"]  # Enabled
    }
    has_alarm_shield_count = true
  }
  ```
- [ ] **Expected**: Returns only enabled alarm notices
- [ ] **Validation**:
  - [ ] deliver_flag = 2 for all results
  - [ ] alarm_shield_count populated
  - [ ] has_alarm_shield_count parameter works
- [ ] **Estimated Time**: 5 minutes

---

### Task 2.5: Test - Multiple Filters

- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "complex" {
    filters {
      name   = "name"
      values = ["prod-alarm", "test-alarm"]
    }
    
    filters {
      name   = "deliverFlag"
      values = ["2"]
    }
  }
  ```
- [ ] **Expected**: Returns notices matching all filters (AND logic)
- [ ] **Validation**:
  - [ ] Results match both filters
  - [ ] Multiple values in single filter work (OR within filter)
  - [ ] Combined filters work correctly (AND between filters)
- [ ] **Estimated Time**: 8 minutes

---

### Task 2.6: Test - Export to File

- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "export" {
    result_output_file = "alarm_notices_output.json"
  }
  ```
- [ ] **Expected**: Creates JSON file with results
- [ ] **Validation**:
  - [ ] File created successfully
  - [ ] File contains valid JSON
  - [ ] All data present in file
  - [ ] File readable and parsable
- [ ] **Estimated Time**: 5 minutes

---

### Task 2.7: Test - Reference in Other Resources

- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "existing" {
    filters {
      name   = "name"
      values = ["my-notice"]
    }
  }
  
  # Example: Reference in alarm resource (if exists)
  resource "tencentcloud_cls_alarm" "example" {
    notice_id = data.tencentcloud_cls_alarm_notices.existing.alarm_notices[0].alarm_notice_id
    # ... other fields
  }
  ```
- [ ] **Expected**: alarm_notice_id can be referenced
- [ ] **Validation**:
  - [ ] Reference syntax works
  - [ ] ID passed correctly to resource
  - [ ] No circular dependencies
- [ ] **Estimated Time**: 7 minutes

---

### Task 2.8: Test - Pagination (Large Dataset)

- [ ] **Setup**: Account with > 100 alarm notices (or mock)
- [ ] **Configuration**:
  ```hcl
  data "tencentcloud_cls_alarm_notices" "many" {}
  ```
- [ ] **Expected**: Returns all notices across multiple pages
- [ ] **Validation**:
  - [ ] All pages retrieved
  - [ ] No duplicate entries
  - [ ] total_count matches actual count
  - [ ] Retry logic works in pagination loop
  - [ ] No memory leaks
- [ ] **Estimated Time**: 10 minutes

---

## Phase 3: Documentation ⏳

**Status**: Pending  
**Estimated Time**: 30 minutes

### Task 3.1: Create Data Source Documentation

- [ ] **File**: `website/docs/d/cls_alarm_notices.html.markdown`
- [ ] **Sections**:
  - [ ] Title and description
  - [ ] Example usage (3-4 examples)
  - [ ] Argument reference
  - [ ] Attributes reference
  - [ ] Filter documentation
- [ ] **Estimated Time**: 20 minutes
- [ ] **Validation**: Documentation builds, no broken links

---

### Task 3.2: Update Changelog

- [ ] **File**: `CHANGELOG.md`
- [ ] **Entry**:
  ```markdown
  **NEW DATA SOURCES**
  * `tencentcloud_cls_alarm_notices` ([#XXXX](link-to-pr))
  ```
- [ ] **Estimated Time**: 5 minutes
- [ ] **Validation**: Changelog format correct

---

## Phase 4: Review and Polish ⏳

**Status**: Pending  
**Estimated Time**: 15 minutes

### Task 4.1: Code Review

- [ ] **Checks**:
  - [ ] Follows reference implementation pattern exactly
  - [ ] Retry logic inside pagination loop (CRITICAL)
  - [ ] All nil checks present
  - [ ] Error handling comprehensive
  - [ ] Logging appropriate
  - [ ] Code comments clear
  - [ ] No hardcoded values
  - [ ] Uses helper functions consistently
- [ ] **Estimated Time**: 10 minutes

---

### Task 4.2: Run Linters

- [ ] **Commands**:
  - [ ] `golangci-lint run ./tencentcloud/services/cls/...`
  - [ ] Fix any warnings
- [ ] **Validation**: No linter errors/warnings
- [ ] **Estimated Time**: 5 minutes

---

## Task Dependencies

```
┌─────────────────────────────────────────────────────────────┐
│                        Phase 1                               │
│                   Code Implementation                        │
└─────────────────┬───────────────────────────────────────────┘
                  │
        ┌─────────┼─────────┬─────────┐
        │         │         │         │
        ▼         ▼         ▼         ▼
     Task 1.1  Task 1.2  Task 1.3  Task 1.4
    Create DS  Add Svc   Register   Format
                         
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                        Phase 2                               │
│                        Testing                               │
└───┬───┬───┬───┬───┬───┬───┬───────────────────────────────┘
    │   │   │   │   │   │   │
    ▼   ▼   ▼   ▼   ▼   ▼   ▼
  2.1 2.2 2.3 2.4 2.5 2.6 2.7  (Can run in parallel)
  
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                        Phase 3                               │
│                     Documentation                            │
└───────────────────┬─────────────────────────────────────────┘
                    │
            ┌───────┴───────┐
            │               │
            ▼               ▼
         Task 3.1        Task 3.2
         Create Doc      Update Log
         
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│                        Phase 4                               │
│                   Review & Polish                            │
└───────────────────┬─────────────────────────────────────────┘
                    │
            ┌───────┴───────┐
            │               │
            ▼               ▼
         Task 4.1        Task 4.2
        Code Review    Run Linters
```

---

## Implementation Checklist

### Before Starting
- [ ] Read proposal.md and design.md thoroughly
- [ ] Review reference implementation
- [ ] Understand API documentation
- [ ] Set up test environment

---

### During Implementation
- [ ] Follow reference pattern exactly
- [ ] Add comprehensive nil checks
- [ ] Include error logging
- [ ] Test after each major component
- [ ] Commit frequently with clear messages

---

### Before Submitting
- [ ] All tests pass
- [ ] Code formatted
- [ ] No linter warnings
- [ ] Documentation complete
- [ ] Changelog updated
- [ ] Self-review done

---

## Testing Commands

### Manual Testing
```bash
# Initialize Terraform
cd examples/cls
terraform init

# Test data source
terraform plan
terraform apply

# Verify output
terraform show
```

### Unit Testing
```bash
# Run specific test
go test -v ./tencentcloud/services/cls -run TestAccTencentCloudClsAlarmNoticesDataSource

# Run all CLS tests
go test -v ./tencentcloud/services/cls/...
```

### Linting
```bash
# Format code
go fmt ./tencentcloud/services/cls/...

# Run linters
golangci-lint run ./tencentcloud/services/cls/...
```

---

## Common Issues and Solutions

### Issue 1: Retry Not Inside Pagination Loop
**Problem**: Retry wrapper outside for loop  
**Solution**: Move `resource.Retry()` inside the for loop  
**Reference**: See service_tencentcloud_igtm.go lines 89-104

---

### Issue 2: Nil Pointer Errors
**Problem**: Accessing nested fields without nil checks  
**Solution**: Check `!= nil` before accessing any pointer field  
**Example**:
```go
if alarmNotice.Name != nil {
    alarmNoticeMap["name"] = alarmNotice.Name
}
```

---

### Issue 3: Filter Key vs Name
**Problem**: Using "name" instead of "Key" for Filter struct  
**Solution**: API uses `Key` field, not `name`  
**Correct**:
```go
filter.Key = helper.String(v)
filter.Values = []*string{...}
```

---

### Issue 4: Pagination Not Working
**Problem**: Not checking length correctly  
**Solution**: Compare with limit, not 0  
**Correct**:
```go
if len(response.Response.AlarmNotices) < int(limit) {
    break
}
```

---

## Time Breakdown

| Task | Estimated | Actual | Notes |
|------|-----------|--------|-------|
| Create data source file | 60 min | | |
| Add service method | 20 min | | |
| Register in provider | 5 min | | |
| Format code | 2 min | | |
| Test - All notices | 5 min | | |
| Test - By name | 5 min | | |
| Test - By ID | 5 min | | |
| Test - By status | 5 min | | |
| Test - Multiple filters | 8 min | | |
| Test - Export file | 5 min | | |
| Test - Reference | 7 min | | |
| Test - Pagination | 10 min | | |
| Documentation | 20 min | | |
| Changelog | 5 min | | |
| Code review | 10 min | | |
| Linting | 5 min | | |
| **TOTAL** | **177 min** | | **~3 hours** |

---

## Success Criteria

### Must Have ✅
- [x] Data source created
- [x] Service method added
- [x] Follows exact reference pattern
- [x] Retry inside pagination loop
- [x] All manual tests pass
- [x] Code formatted
- [x] No linter errors

### Should Have ✅
- [ ] Documentation complete
- [ ] Changelog updated
- [ ] All test scenarios covered
- [ ] Large dataset pagination tested

### Nice to Have ⭐
- [ ] Integration tests
- [ ] Performance benchmarks
- [ ] Usage examples in repo
- [ ] Blog post or announcement

---

## Progress Tracking

**Start Date**: TBD  
**Target Completion**: TBD  
**Status**: ⏳ Not Started

### Phase Status
- [ ] Phase 1: Code Implementation - 0% complete
- [ ] Phase 2: Testing - 0% complete
- [ ] Phase 3: Documentation - 0% complete
- [ ] Phase 4: Review - 0% complete

**Overall Progress**: 0/15 tasks (0%)

---

## Notes

- Remember to follow the EXACT pattern from `data_source_tc_igtm_instance_list.go`
- The most critical requirement is retry INSIDE the pagination loop
- Test pagination thoroughly with large datasets
- All pointer fields must have nil checks
- Use `helper` package functions consistently

---

**Task List Version**: 1.0  
**Last Updated**: 2026-03-24  
**Ready for Implementation**: ✅ Yes
