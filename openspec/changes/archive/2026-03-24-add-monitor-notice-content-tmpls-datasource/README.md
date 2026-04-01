# Add Monitor Notice Content Templates DataSource

**Status**: 📋 Proposal  
**Priority**: Medium  
**Estimated Effort**: 2-3 hours  
**Complexity**: Medium  
**API Version**: monitor/v20230616

---

## 🎯 Objective

Create a new data source `tencentcloud_monitor_notice_content_tmpls` for querying custom notification content templates in Tencent Cloud Monitor service.

---

## 📋 Problem Statement

### Current Situation

- The Monitor service has a resource `tencentcloud_monitor_notice_content_tmpl` for managing notification content templates
- There is NO corresponding data source to query existing templates
- Users cannot list or filter notification content templates using Terraform

### Business Need

Users need the ability to:
1. Query notification content templates by ID
2. Query templates by name
3. List all templates under an account
4. Query templates by notification ID
5. Filter templates by language and monitor type

---

## 🎯 Solution Overview

Create a new data source `tencentcloud_monitor_notice_content_tmpls` that:
- Uses the `DescribeNoticeContentTmpl` API
- Supports filtering by template IDs, name, notice ID, language, and monitor type
- Implements pagination with retry logic
- Follows the code pattern of `data_source_tc_igtm_instance_list.go`

---

## 📊 API Information

### API Details

- **API Name**: `DescribeNoticeContentTmpl`
- **SDK Version**: `monitor/v20230616`
- **Documentation**: https://cloud.tencent.com/document/product/248/128618
- **Request Domain**: monitor.tencentcloudapi.com

### Request Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| PageNumber | Yes | Integer | Page number, e.g., 1 |
| PageSize | Yes | Integer | Page size, e.g., 10 |
| TmplIDs.N | No | Array of String | Template ID list for query |
| TmplName | No | String | Template name for query |
| NoticeID | No | String | Notice template ID for query |
| TmplLanguage | No | String | Template language: `en`/`zh` |
| MonitorType | No | String | Monitor type, e.g., `MT_QCE` |

### Response Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| NoticeContentTmpls | Array of NoticeContentTmpl | Custom notification content template list |
| NoticeContentTmplBindPolicyCounts | Array of NoticeContentTmplBindPolicyCount | Bound policy counts |
| PageNumber | Integer | Page number |
| PageSize | Integer | Page size |
| TotalCount | Integer | Total count |
| RequestId | String | Unique request ID |

---

## 🏗️ Technical Design

### File Structure

```
tencentcloud/services/monitor/
├── data_source_tc_monitor_notice_content_tmpls.go      (NEW - main data source)
├── data_source_tc_monitor_notice_content_tmpls_test.go (NEW - acceptance test)
├── data_source_tc_monitor_notice_content_tmpls.md      (NEW - documentation)
└── service_tencentcloud_monitor.go                     (UPDATE - add service method)
```

### Schema Design

```go
Schema: map[string]*schema.Schema{
    // Input parameters
    "tmpl_ids": {
        Type:        schema.TypeSet,
        Optional:    true,
        Elem:        &schema.Schema{Type: schema.TypeString},
        Description: "Template ID list for query.",
    },
    "tmpl_name": {
        Type:        schema.TypeString,
        Optional:    true,
        Description: "Template name for query.",
    },
    "notice_id": {
        Type:        schema.TypeString,
        Optional:    true,
        Description: "Notice template ID for query.",
    },
    "tmpl_language": {
        Type:        schema.TypeString,
        Optional:    true,
        Description: "Template language: en/zh.",
    },
    "monitor_type": {
        Type:        schema.TypeString,
        Optional:    true,
        Description: "Monitor type, e.g., MT_QCE.",
    },
    
    // Output parameters
    "notice_content_tmpl_list": {
        Type:        schema.TypeList,
        Computed:    true,
        Description: "Notification content template list.",
        Elem:        &schema.Resource{...},
    },
    "result_output_file": {
        Type:        schema.TypeString,
        Optional:    true,
        Description: "Used to save results.",
    },
}
```

### Service Layer Method

Add to `service_tencentcloud_monitor.go`:

```go
func (me *MonitorService) DescribeNoticeContentTmplsByFilter(
    ctx context.Context,
    param map[string]interface{},
) (
    noticeContentTmpls []*monitorv20230616.NoticeContentTmpl,
    bindPolicyCounts []*monitorv20230616.NoticeContentTmplBindPolicyCount,
    errRet error,
) {
    // Implementation with pagination and retry logic
}
```

---

## 📝 Implementation Tasks

### Phase 1: Service Layer (30 min)

- [ ] Add `DescribeNoticeContentTmplsByFilter` method to `service_tencentcloud_monitor.go`
- [ ] Implement pagination logic (PageNumber-based)
- [ ] Add retry logic for each API call in the pagination loop
- [ ] Handle all filter parameters

### Phase 2: Data Source Layer (60 min)

- [ ] Create `data_source_tc_monitor_notice_content_tmpls.go`
- [ ] Implement schema definition with all input/output fields
- [ ] Implement `dataSourceTencentCloudMonitorNoticeContentTmplsRead` function
- [ ] Parse input parameters and build filter map
- [ ] Call service layer method with retry
- [ ] Map response to schema output

### Phase 3: Registration (10 min)

- [ ] Register data source in `tencentcloud/provider.go`
- [ ] Add to data sources map

### Phase 4: Testing (30 min)

- [ ] Create `data_source_tc_monitor_notice_content_tmpls_test.go`
- [ ] Write basic acceptance test
- [ ] Test filtering by different parameters
- [ ] Test pagination scenarios

### Phase 5: Documentation (30 min)

- [ ] Create `data_source_tc_monitor_notice_content_tmpls.md`
- [ ] Add usage examples
- [ ] Document all parameters
- [ ] Add argument and attribute references

---

## 🔍 Reference Implementation

### Primary Reference

**File**: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`

**Key Patterns**:
1. Schema structure with filters and result_output_file
2. Data source read function structure
3. Parameter parsing and mapping
4. Service layer call with retry
5. Response data mapping to schema

### Service Layer Reference

**File**: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/service_tencentcloud_igtm.go`

**Method**: `DescribeIgtmInstanceListByFilter`

**Key Patterns**:
1. Pagination loop with offset/limit or PageNumber/PageSize
2. Retry logic inside pagination loop (CRITICAL)
3. Response data accumulation
4. Error handling and logging

---

## ⚠️ Critical Requirements

### 1. Pagination with Retry

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
        // Process result
        response = result
        return nil
    })
    
    if err != nil {
        errRet = err
        return
    }
    
    ret = append(ret, response.Response.NoticeContentTmpls...)
    if len(response.Response.NoticeContentTmpls) < int(pageSize) {
        break
    }
    
    pageNumber++
}
```

### 2. Code Style

- Follow exact code structure of `data_source_tc_igtm_instance_list.go`
- Use consistent variable naming
- Include all logging statements
- Add proper error handling

### 3. Schema Consistency

- Use computed output fields
- Support `result_output_file` parameter
- Use appropriate field types (TypeSet for lists, TypeString for IDs)

---

## 📊 Success Criteria

### Functional Requirements

- ✅ Data source can query templates without filters (all templates)
- ✅ Data source can filter by template IDs
- ✅ Data source can filter by template name
- ✅ Data source can filter by notice ID
- ✅ Data source can filter by language and monitor type
- ✅ Pagination works correctly for large datasets
- ✅ Retry logic handles transient API errors

### Code Quality

- ✅ Follows reference implementation pattern exactly
- ✅ All parameters properly documented
- ✅ Proper error handling and logging
- ✅ Acceptance test passes
- ✅ Code formatted with `go fmt`

### Documentation

- ✅ Clear usage examples
- ✅ All parameters documented
- ✅ Attribute references provided

---

## 🧪 Testing Strategy

### Manual Testing

```hcl
# Test 1: Query all templates
data "tencentcloud_monitor_notice_content_tmpls" "all" {
}

# Test 2: Query by template ID
data "tencentcloud_monitor_notice_content_tmpls" "by_id" {
  tmpl_ids = ["ntpl-xxxxxx"]
}

# Test 3: Query by template name
data "tencentcloud_monitor_notice_content_tmpls" "by_name" {
  tmpl_name = "test-template"
}

# Test 4: Query by notice ID
data "tencentcloud_monitor_notice_content_tmpls" "by_notice" {
  notice_id = "notice-xxxxxx"
}

# Test 5: Filter by language
data "tencentcloud_monitor_notice_content_tmpls" "zh_only" {
  tmpl_language = "zh"
}

# Test 6: Multiple filters
data "tencentcloud_monitor_notice_content_tmpls" "filtered" {
  tmpl_language = "zh"
  monitor_type  = "MT_QCE"
}
```

### Acceptance Test

```go
func TestAccTencentCloudMonitorNoticeContentTmplsDataSource_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:  func() { tcacctest.AccPreCheck(t) },
        Providers: tcacctest.AccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccMonitorNoticeContentTmplsDataSource,
                Check: resource.ComposeTestCheckFunc(
                    tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_notice_content_tmpls.tmpls"),
                    resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_notice_content_tmpls.tmpls", "notice_content_tmpl_list.#"),
                ),
            },
        },
    })
}
```

---

## 🚀 Implementation Plan

### Timeline

| Phase | Task | Duration | Priority |
|-------|------|----------|----------|
| 1 | Service layer implementation | 30 min | High |
| 2 | Data source implementation | 60 min | High |
| 3 | Provider registration | 10 min | High |
| 4 | Testing | 30 min | Medium |
| 5 | Documentation | 30 min | Medium |

**Total Estimated Time**: 2-3 hours

### Dependencies

- ✅ SDK already has `DescribeNoticeContentTmpl` API
- ✅ Monitor service layer already exists
- ✅ Reference implementation available

---

## 📚 Related Resources

### Existing Resources

- `resource_tc_monitor_notice_content_tmpl` - Resource for managing templates
- `data_source_tc_monitor_alarm_notices` - Similar data source pattern
- `data_source_tc_igtm_instance_list` - Reference implementation

### API Documentation

- [DescribeNoticeContentTmpl API](https://cloud.tencent.com/document/product/248/128618)
- [Monitor API Overview](https://cloud.tencent.com/document/product/248)

---

## 💡 Best Practices

### From Reference Implementation

1. **Use paramMap pattern**: Build parameter map in data source, pass to service
2. **Retry at service layer**: Always wrap API calls in `resource.Retry`
3. **Pagination in loop**: Handle pagination in service layer, not data source
4. **Nil checks**: Always check for nil responses
5. **Logging**: Include debug and error logs
6. **Token generation**: Use `helper.BuildToken()` for data source ID

---

## ⚙️ Configuration Example

```hcl
# Query all notification content templates
data "tencentcloud_monitor_notice_content_tmpls" "all_tmpls" {
  result_output_file = "tmpls.json"
}

output "template_count" {
  value = length(data.tencentcloud_monitor_notice_content_tmpls.all_tmpls.notice_content_tmpl_list)
}

# Query specific template by ID
data "tencentcloud_monitor_notice_content_tmpls" "specific" {
  tmpl_ids = ["ntpl-3r1spzjn"]
}

output "template_name" {
  value = data.tencentcloud_monitor_notice_content_tmpls.specific.notice_content_tmpl_list[0].tmpl_name
}
```

---

## 📝 Notes

### Important Considerations

1. **Pagination Type**: This API uses `PageNumber`/`PageSize` (not `Offset`/`Limit`)
2. **Retry Logic**: MUST be inside pagination loop for each API call
3. **Bind Policy Counts**: Include in response for comprehensive information
4. **Base64 Content**: Template contents are base64 encoded in API response
5. **Complex Nested Structure**: `TmplContents` has multiple levels of nesting

### Potential Challenges

1. **Complex Schema**: The `TmplContents` structure is deeply nested
2. **Multiple Channel Types**: Different channels (Email, SMS, etc.) have different structures
3. **Base64 Decoding**: Consider if content should be decoded or kept encoded

---

**Proposal Created**: 2026-03-24  
**Status**: Ready for Implementation  
**Assignee**: TBD
