# Add tencentcloud_cls_alarm_notices Data Source

**Status**: ✅ Completed (Code Phase)  
**Priority**: Medium  
**Implementation Date**: 2026-03-24  
**Actual Effort**: 20 minutes (vs. estimated 3 hours)  
**Complexity**: Medium

---

## 📋 Quick Summary

Add a new data source `tencentcloud_cls_alarm_notices` to query CLS (Cloud Log Service) alarm notification channel groups using the `DescribeAlarmNotices` API.

**What**: New Terraform data source for CLS alarm notices  
**Why**: Enable programmatic querying and referencing of alarm notice configurations  
**How**: Follow exact pattern of `tencentcloud_igtm_instance_list` with pagination and retry logic

---

## 🎯 Problem Statement

Users cannot query existing CLS alarm notice configurations in Terraform, which limits:
- Automation workflows requiring alarm notice IDs
- Validation of alarm notice settings
- Documentation and auditing of notification configurations
- Reference of existing notices in other resources

---

## ✨ Proposed Solution

Create a data source that wraps the `DescribeAlarmNotices` API with:

### Key Features
- ✅ Query all alarm notices or filter by specific criteria
- ✅ Support 5 filter types (name, ID, user, group, delivery status)
- ✅ Return complete alarm notice configuration data
- ✅ Handle pagination for large result sets
- ✅ Optional alarm shield count statistics
- ✅ Export results to JSON file

### User Benefits
- 🚀 Reference existing alarm notice IDs in Terraform
- 🔍 Discover and validate alarm notice configurations
- 📊 Export alarm notice data for auditing
- 🔧 Enable automation with existing infrastructure

---

## 📝 Usage Example

```hcl
# Query all alarm notices
data "tencentcloud_cls_alarm_notices" "all" {}

# Filter by name
data "tencentcloud_cls_alarm_notices" "by_name" {
  filters {
    name   = "name"
    values = ["prod-alarm"]
  }
}

# Filter by delivery status
data "tencentcloud_cls_alarm_notices" "enabled" {
  filters {
    name   = "deliverFlag"
    values = ["2"]  # Enabled
  }
  has_alarm_shield_count = true
}

# Reference alarm notice ID
data "tencentcloud_cls_alarm_notices" "existing" {
  filters {
    name   = "alarmNoticeId"
    values = ["notice-xxxx"]
  }
}

resource "tencentcloud_cls_alarm" "example" {
  notice_id = data.tencentcloud_cls_alarm_notices.existing.alarm_notices[0].alarm_notice_id
  # ... other configuration
}

# Export to file
data "tencentcloud_cls_alarm_notices" "export" {
  result_output_file = "alarm_notices.json"
}

# Output for reference
output "notice_ids" {
  value = [
    for notice in data.tencentcloud_cls_alarm_notices.all.alarm_notices :
    notice.alarm_notice_id
  ]
}
```

---

## 🔍 Supported Filters

| Filter Name | Type | Description | Example Values |
|-------------|------|-------------|----------------|
| `name` | String | Alarm notice group name | `"prod-alarm"`, `"test-alarm"` |
| `alarmNoticeId` | String | Alarm notice ID | `"notice-xxxx-yyyy"` |
| `uid` | String | Receiver user ID | `"100001234567"` |
| `groupId` | String | Receiver group ID | `"group-xxxx"` |
| `deliverFlag` | String | Delivery status | `"1"` (Not enabled), `"2"` (Enabled), `"3"` (Exception) |

**Filter Constraints**:
- Maximum 10 filters per request
- Maximum 5 values per filter
- Multiple values within a filter = OR logic
- Multiple filters = AND logic

---

## 📊 Output Attributes

### Main Attributes
- `alarm_notices` - List of alarm notice objects
- `total_count` - Total number of matching notices

### Alarm Notice Object
```hcl
{
  name                       = "test-alarm"
  alarm_notice_id            = "notice-xxxx"
  create_time                = "2025-08-06 15:47:00"
  update_time                = "2025-08-06 15:47:00"
  notice_receivers           = [...]  # Notification receivers
  web_callbacks              = [...]  # Webhook configurations
  tags                       = [...]  # Resource tags
  jump_domain                = "https://console.cloud.tencent.com"
  notice_rules               = [...]  # Notification rules
  deliver_status             = 1
  deliver_flag               = 2      # 1: Not enabled, 2: Enabled, 3: Exception
  alarm_shield_status        = 2
  alarm_shield_count         = {...}  # Shield statistics (if requested)
  callback_prioritize        = true
}
```

---

## 🏗️ Implementation Overview

### Files to Create/Modify

| File | Type | Lines | Description |
|------|------|-------|-------------|
| `data_source_tc_cls_alarm_notices.go` | NEW | ~400 | Data source definition and read logic |
| `service_tencentcloud_cls.go` | MODIFY | ~80 | Add `DescribeClsAlarmNoticesByFilter` method |
| `provider.go` | MODIFY | 1 | Register data source |

**Total**: ~481 lines across 3 files

---

### Key Implementation Details

#### 1. Data Source Structure
```go
func DataSourceTencentCloudClsAlarmNotices() *schema.Resource {
    return &schema.Resource{
        Read: dataSourceTencentCloudClsAlarmNoticesRead,
        Schema: map[string]*schema.Schema{
            "filters": {...},
            "has_alarm_shield_count": {...},
            "alarm_notices": {...},
            "total_count": {...},
            "result_output_file": {...},
        },
    }
}
```

#### 2. Service Layer Method
```go
func (me *ClsService) DescribeClsAlarmNoticesByFilter(
    ctx context.Context, 
    param map[string]interface{},
) (ret []*cls.AlarmNotice, totalCount *uint64, errRet error) {
    // Pagination loop
    for {
        // CRITICAL: Retry INSIDE pagination loop
        err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
            // API call with rate limiting
            result, e := me.client.UseClsClient().DescribeAlarmNotices(request)
            // ...
        })
        // Append results and check for more pages
    }
}
```

---

## 🎯 Reference Implementation

This implementation follows the **EXACT** pattern of:

**File**: `tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`

**Key Patterns to Follow**:
1. ✅ Schema structure with filters as input, list as output
2. ✅ Parse filters into API-compatible format
3. ✅ Service layer with pagination
4. ✅ **Retry logic INSIDE pagination loop** (CRITICAL)
5. ✅ Rate limiting before each API call
6. ✅ Comprehensive nil checks
7. ✅ Map response to Terraform schema
8. ✅ Optional result output file

---

## 📋 API Information

**API**: `DescribeAlarmNotices`  
**Version**: `2020-10-16`  
**Domain**: `cls.tencentcloudapi.com`  
**Rate Limit**: 20 requests/second  
**Documentation**: https://cloud.tencent.com/document/api/614/56462

**Request Parameters**:
- `Filters` - Array of filters (optional)
- `Offset` - Pagination offset (default: 0)
- `Limit` - Page size (default: 20, max: 100)
- `HasAlarmShieldCount` - Return shield count (optional, default: false)

**Response**:
- `AlarmNotices` - Array of alarm notice objects
- `TotalCount` - Total number of matching notices
- `RequestId` - Unique request identifier

---

## ✅ Testing Plan

### Manual Test Scenarios
1. ✅ Query all alarm notices (no filters)
2. ✅ Filter by name
3. ✅ Filter by alarm notice ID
4. ✅ Filter by delivery status
5. ✅ Multiple filters combined
6. ✅ Export to file
7. ✅ Reference in other resources
8. ✅ Pagination with large datasets (> 100 items)

### Validation Points
- Data source executes without errors
- Filters work correctly (AND/OR logic)
- All fields populated properly
- Nested structures parsed correctly
- Pagination handles large result sets
- Output file contains valid JSON
- No state drift on repeated runs

---

## 📚 Documentation

### Data Source Documentation
```markdown
# tencentcloud_cls_alarm_notices

Use this data source to query CLS alarm notification channel groups.

## Example Usage
[See examples above]

## Argument Reference
- `filters` - (Optional) Filter conditions
- `has_alarm_shield_count` - (Optional) Return shield count
- `result_output_file` - (Optional) Save results to file

## Attributes Reference
- `alarm_notices` - Alarm notice list
- `total_count` - Total number of notices
```

---

## ⏱️ Timeline

| Phase | Duration | Tasks | Status |
|-------|----------|-------|--------|
| **Code Implementation** | 20 minutes | Create files, implement logic | ✅ Completed |
| **Testing** | 45 minutes | Manual testing scenarios | ⏳ Pending |
| **Documentation** | 30 minutes | Write docs, update changelog | ⏳ Pending |
| **Review** | 15 minutes | Code review, linting | ⏳ Pending |
| **Total** | **3 hours** | **15 tasks** | **60% complete** |

---

## 🎯 Success Criteria

### ✅ Must Have (MVP) - COMPLETED
- [x] Data source created and registered
- [x] Service method implemented with pagination
- [x] Retry logic INSIDE pagination loop
- [x] All filters supported
- [x] Code formatted (go fmt)
- [x] No linter errors

### ⏳ Should Have - PENDING
- [ ] Manual tests pass
- [ ] Documentation complete
- [ ] Changelog updated
- [ ] All test scenarios validated
- [ ] Large dataset pagination tested

### Nice to Have
- [ ] Integration tests
- [ ] Performance benchmarks
- [ ] Usage examples in repository

---

## ⚠️ Critical Requirements

### 1. Retry Logic Placement
**CRITICAL**: Retry logic MUST be inside the pagination loop, not outside.

❌ **Wrong**:
```go
err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    for { /* pagination */ }
})
```

✅ **Correct**:
```go
for {
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        // API call here
    })
}
```

### 2. Nil Checks
All pointer fields must be checked before access:
```go
if alarmNotice.Name != nil {
    alarmNoticeMap["name"] = alarmNotice.Name
}
```

### 3. Filter Structure
API uses `Key` and `Values` (not `name` and `value`):
```go
filter.Key = helper.String(filterName)
filter.Values = []*string{...}
```

---

## 🔗 Related Resources

### Existing Resources
- `tencentcloud_cls_alarm_notice` - Resource for managing alarm notices
- `tencentcloud_cls_alarm` - Resource for managing alarms
- `tencentcloud_cls_logset` - Resource for log sets
- `tencentcloud_cls_topic` - Resource for log topics

### Related Data Sources
- `tencentcloud_cls_logsets` - Query log sets
- `tencentcloud_cls_topics` - Query log topics
- `tencentcloud_cls_machines` - Query machine groups

---

## 📊 Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| API structure mismatch | Low | Medium | Validate against API docs |
| Pagination issues | Low | Medium | Test with large datasets |
| Filter errors | Low | Low | Validate filter names |
| Nil pointer errors | Medium | Medium | Comprehensive nil checks |

**Overall Risk**: 🟢 Low

---

## 📂 Document Structure

This proposal consists of 4 documents:

1. **README.md** (this file) - Quick overview and summary
2. **proposal.md** - Detailed proposal with requirements
3. **design.md** - Technical design and implementation details
4. **tasks.md** - Step-by-step implementation tasks

---

## 🚀 Next Steps

1. **Review**: Review all proposal documents
2. **Approve**: Get approval to proceed
3. **Implement**: Follow tasks.md step-by-step
4. **Test**: Execute all test scenarios
5. **Document**: Complete documentation
6. **Submit**: Create pull request

---

## 💡 Key Takeaways

1. **Pattern**: Follow `tencentcloud_igtm_instance_list` exactly
2. **Retry**: Place retry logic INSIDE pagination loop (CRITICAL)
3. **Nil Checks**: Check all pointer fields before access
4. **Testing**: Test pagination with large datasets
5. **Filters**: Support 5 filter types with AND/OR logic
6. **Documentation**: Provide clear examples and descriptions

---

## 📞 Questions?

- **API Documentation**: https://cloud.tencent.com/document/api/614/56462
- **Reference Implementation**: `tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`
- **Service Pattern**: `tencentcloud/services/igtm/service_tencentcloud_igtm.go` (lines 398-459)

---

**Proposal Created**: 2026-03-24  
**Implementation Completed**: 2026-03-24  
**Status**: ✅ Code Complete - Ready for Testing  
**See**: IMPLEMENTATION_SUMMARY.md for details

Ready for testing phase! 🎉
