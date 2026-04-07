# Summary: Monitor Notice Content Templates DataSource

## 📊 Quick Overview

| Item | Details |
|------|---------|
| **Feature** | New data source for querying monitor notification content templates |
| **Resource Name** | `tencentcloud_monitor_notice_content_tmpls` |
| **API** | `DescribeNoticeContentTmpl` (monitor/v20230616) |
| **Complexity** | Medium |
| **Effort** | 2-3 hours |
| **Priority** | Medium |
| **Status** | 📋 Proposal |

---

## 🎯 Objective

Enable users to query and filter notification content templates in Tencent Cloud Monitor service using Terraform data source.

---

## 💡 Value Proposition

### Current Pain Points
- ❌ No way to query existing templates in Terraform
- ❌ Cannot list templates for audit purposes
- ❌ Cannot reference templates in other configurations
- ❌ Lack of visibility into template inventory

### After Implementation
- ✅ Query all templates under account
- ✅ Filter templates by ID, name, notice ID, language, type
- ✅ Reference template data in other resources
- ✅ Export template data to JSON for analysis

---

## 🏗️ Technical Overview

### Files to Create/Modify

| File | Action | Lines (est.) |
|------|--------|--------------|
| `data_source_tc_monitor_notice_content_tmpls.go` | CREATE | ~300 |
| `data_source_tc_monitor_notice_content_tmpls_test.go` | CREATE | ~50 |
| `data_source_tc_monitor_notice_content_tmpls.md` | CREATE | ~150 |
| `service_tencentcloud_monitor.go` | UPDATE | +80 |
| `provider.go` | UPDATE | +1 |

**Total**: 4 new files, 2 updated files, ~580 lines

---

## 🔑 Key Features

### Input Parameters (Filters)

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tmpl_ids` | Set of String | No | Template ID list |
| `tmpl_name` | String | No | Template name |
| `notice_id` | String | No | Notice template ID |
| `tmpl_language` | String | No | Language: en/zh |
| `monitor_type` | String | No | Monitor type |

### Output Attributes

| Attribute | Type | Description |
|-----------|------|-------------|
| `notice_content_tmpl_list` | List | Template list with details |
| ↳ `tmpl_id` | String | Template ID |
| ↳ `tmpl_name` | String | Template name |
| ↳ `monitor_type` | String | Monitor type |
| ↳ `tmpl_language` | String | Template language |
| ↳ `creator` | String | Creator UIN |
| ↳ `create_time` | Int | Create timestamp |
| ↳ `update_time` | Int | Update timestamp |
| ↳ `tmpl_contents_json` | String | Contents (JSON) |
| ↳ `bind_policy_count` | Int | Bound policy count |

---

## 📋 Implementation Phases

### Phase 1: Service Layer (30 min)
- Add `DescribeNoticeContentTmplsByFilter` method
- Implement pagination with PageNumber/PageSize
- **CRITICAL**: Add retry logic INSIDE pagination loop
- Return templates and bind policy counts

### Phase 2: Data Source (60 min)
- Create data source file with schema
- Implement read function
- Parse input parameters
- Call service layer with retry
- Map response to schema output

### Phase 3: Registration (10 min)
- Register data source in provider

### Phase 4: Testing (30 min)
- Create acceptance test
- Write test configurations
- Run and verify tests

### Phase 5: Documentation (30 min)
- Create documentation file
- Add usage examples
- Document all parameters

---

## ⚠️ Critical Requirements

### 1. Follow Reference Pattern

**Primary Reference**: `data_source_tc_igtm_instance_list.go`

**Must Follow**:
- Schema structure
- Parameter parsing pattern
- Service call pattern
- Response mapping pattern
- Error handling pattern

### 2. Pagination with Retry

```go
// ✅ CORRECT: Retry INSIDE loop
for {
    request.PageNumber = &pageNumber
    request.PageSize = &pageSize
    
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        result, e := me.client.UseMonitorV20230616Client().DescribeNoticeContentTmpl(request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        // Accumulate results
        return nil
    })
    
    if err != nil {
        return
    }
    
    if len(results) < int(pageSize) {
        break
    }
    
    pageNumber++
}
```

### 3. API Parameters

| API Parameter | Type | SDK Type |
|---------------|------|----------|
| PageNumber | Integer | `*uint64` |
| PageSize | Integer | `*uint64` |
| TmplIDs.N | Array of String | `[]*string` |
| TmplName | String | `*string` |
| NoticeID | String | `*string` |
| TmplLanguage | String | `*string` |
| MonitorType | String | `*string` |

---

## 🧪 Testing Strategy

### Test Cases

1. **Query all templates** (no filters)
2. **Query by template ID** (single and multiple)
3. **Query by template name**
4. **Query by notice ID**
5. **Filter by language** (zh/en)
6. **Filter by monitor type**
7. **Multiple filters combined**
8. **Empty results** (filter matches nothing)
9. **Export to file**

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
                    tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_notice_content_tmpls.test"),
                    resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_notice_content_tmpls.test", "notice_content_tmpl_list.#"),
                ),
            },
        },
    })
}
```

---

## 📖 Usage Examples

### Example 1: Query All Templates

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "all" {
}

output "template_count" {
  value = length(data.tencentcloud_monitor_notice_content_tmpls.all.notice_content_tmpl_list)
}
```

### Example 2: Query by Template ID

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "specific" {
  tmpl_ids = ["ntpl-3r1spzjn"]
}

output "template_name" {
  value = data.tencentcloud_monitor_notice_content_tmpls.specific.notice_content_tmpl_list[0].tmpl_name
}
```

### Example 3: Filter by Language

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "zh_templates" {
  tmpl_language = "zh"
}
```

### Example 4: Multiple Filters

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "filtered" {
  tmpl_language = "zh"
  monitor_type  = "MT_QCE"
  tmpl_name     = "production"
}
```

### Example 5: Export to File

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "export" {
  result_output_file = "templates.json"
}
```

---

## 📊 Success Metrics

### Functional Criteria

- ✅ Can query templates without filters
- ✅ All filter parameters work correctly
- ✅ Pagination handles 100+ templates
- ✅ Retry logic handles transient errors
- ✅ Export to file works
- ✅ All acceptance tests pass

### Code Quality Criteria

- ✅ Follows reference implementation pattern
- ✅ Proper error handling and logging
- ✅ All fields have nil checks
- ✅ Code formatted with `go fmt`
- ✅ No linter warnings
- ✅ Documentation complete

### Performance Criteria

- ✅ Query 100 templates in < 5 seconds
- ✅ Memory usage reasonable
- ✅ No unnecessary API calls

---

## 🚀 Implementation Timeline

| Day | Activities | Deliverables |
|-----|------------|--------------|
| Day 1 | Service layer + Data source | Code files |
| Day 1 | Testing + Registration | Tests pass |
| Day 1 | Documentation | Docs complete |

**Total**: 1 day (2-3 hours of development)

---

## 🔗 Related Resources

### Existing Resources
- `resource_tc_monitor_notice_content_tmpl` - Template management resource
- `data_source_tc_monitor_alarm_notices` - Similar data source pattern
- `data_source_tc_igtm_instance_list` - Reference implementation

### API Documentation
- [DescribeNoticeContentTmpl API](https://cloud.tencent.com/document/product/248/128618)
- [Monitor API Overview](https://cloud.tencent.com/document/product/248)

### Code References
- Service: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/service_tencentcloud_igtm.go`
- DataSource: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`

---

## ⚡ Quick Start Commands

```bash
# Navigate to project root
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud

# Create service layer method
# Edit: tencentcloud/services/monitor/service_tencentcloud_monitor.go

# Create data source
# Create: tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls.go

# Register in provider
# Edit: tencentcloud/provider.go

# Format code
go fmt ./tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls*.go
go fmt ./tencentcloud/services/monitor/service_tencentcloud_monitor.go

# Run tests
cd tencentcloud
TF_ACC=1 go test -v -run TestAccTencentCloudMonitorNoticeContentTmplsDataSource ./services/monitor/
```

---

## 🎯 Key Takeaways

1. **New Capability**: Users can query notification templates via Terraform
2. **Pattern Following**: Strictly follows `data_source_tc_igtm_instance_list.go` pattern
3. **Critical Detail**: Retry logic MUST be inside pagination loop
4. **Complexity**: Medium - requires pagination + complex nested structures
5. **Effort**: 2-3 hours for complete implementation
6. **Impact**: Improves visibility and automation of template management

---

**Proposal Created**: 2026-03-24  
**Status**: 📋 Ready for Implementation  
**Next Steps**: Review proposal → Approve → Begin implementation
