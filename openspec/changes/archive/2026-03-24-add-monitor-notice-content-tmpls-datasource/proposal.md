# Technical Proposal: Monitor Notice Content Templates DataSource

## 1. Executive Summary

### Objective
Add a new data source `tencentcloud_monitor_notice_content_tmpls` to query notification content templates in Tencent Cloud Monitor service.

### Scope
- **In Scope**: Data source implementation, service layer method, pagination with retry, basic filtering
- **Out of Scope**: Template content modification, advanced filtering beyond API support, content decoding

### Key Deliverables
1. New data source file: `data_source_tc_monitor_notice_content_tmpls.go`
2. Service layer method: `DescribeNoticeContentTmplsByFilter`
3. Acceptance test file
4. Documentation file

---

## 2. Background & Context

### Current State

The Monitor service already has:
- ✅ Resource: `tencentcloud_monitor_notice_content_tmpl` (for CRUD operations)
- ✅ Service method: `DescribeNoticeContentTmplByFilter` (basic filtering)
- ❌ No data source for querying templates

### Problem Statement

Users cannot:
- List all notification content templates
- Query templates by multiple criteria
- Use template data in other Terraform configurations
- Reference existing templates in infrastructure code

### User Stories

**Story 1**: As a DevOps engineer, I want to query all notification templates so that I can audit existing configurations.

**Story 2**: As a developer, I want to find templates by name so that I can reference them in alarm policies.

**Story 3**: As a platform admin, I want to filter templates by language and monitor type so that I can manage templates by category.

---

## 3. Technical Requirements

### Functional Requirements

| ID | Requirement | Priority | Verification |
|----|-------------|----------|--------------|
| FR-1 | Query all templates without filters | High | Manual test + acceptance test |
| FR-2 | Filter by template IDs (multiple) | High | Manual test |
| FR-3 | Filter by template name | Medium | Manual test |
| FR-4 | Filter by notice ID | Medium | Manual test |
| FR-5 | Filter by language (en/zh) | Low | Manual test |
| FR-6 | Filter by monitor type | Low | Manual test |
| FR-7 | Support pagination for large datasets | High | Unit test |
| FR-8 | Export results to JSON file | Medium | Manual test |

### Non-Functional Requirements

| ID | Requirement | Target | Verification |
|----|-------------|--------|--------------|
| NFR-1 | API call retry on transient errors | 3 retries | Code review |
| NFR-2 | Response time for 100 templates | < 5s | Performance test |
| NFR-3 | Code coverage | > 80% | Unit tests |
| NFR-4 | Follow existing code patterns | 100% | Code review |
| NFR-5 | Documentation completeness | 100% | Doc review |

---

## 4. Detailed Design

### 4.1 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Terraform Configuration                   │
│  data "tencentcloud_monitor_notice_content_tmpls" "test" {  │
│    tmpl_ids = ["ntpl-xxx"]                                  │
│  }                                                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│           Data Source Layer (data_source_*.go)              │
│  - Parse input parameters                                    │
│  - Build parameter map                                       │
│  - Call service layer                                        │
│  - Map response to schema                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│         Service Layer (service_tencentcloud_monitor.go)     │
│  - DescribeNoticeContentTmplsByFilter()                     │
│  - Handle pagination (PageNumber/PageSize)                   │
│  - Retry logic for each API call                            │
│  - Accumulate results                                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              Tencent Cloud SDK (monitor/v20230616)          │
│  DescribeNoticeContentTmpl API                              │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 Data Flow

```
User Input Parameters
  ├─ tmpl_ids ([]string)
  ├─ tmpl_name (string)
  ├─ notice_id (string)
  ├─ tmpl_language (string)
  └─ monitor_type (string)
          │
          ▼
    Build paramMap
          │
          ▼
Service Layer: DescribeNoticeContentTmplsByFilter
    │
    ├─ Page 1: API Call with Retry
    │     └─ Collect results
    ├─ Page 2: API Call with Retry
    │     └─ Collect results
    └─ ...until no more pages
          │
          ▼
    Return Results
    ├─ NoticeContentTmpls []
    └─ BindPolicyCounts []
          │
          ▼
    Map to Schema
    ├─ notice_content_tmpl_list
    │     ├─ tmpl_id
    │     ├─ tmpl_name
    │     ├─ tmpl_contents (complex nested)
    │     ├─ creator
    │     ├─ create_time
    │     └─ ...
    └─ Export to file (optional)
```

### 4.3 Schema Definition

#### Input Schema

```go
"tmpl_ids": {
    Type: schema.TypeSet,
    Optional: true,
    Elem: &schema.Schema{Type: schema.TypeString},
    Description: "Template ID list for query. Example: [\"ntpl-3r1spzjn\"]",
},
"tmpl_name": {
    Type: schema.TypeString,
    Optional: true,
    Description: "Template name for query. Example: \"custom-template\"",
},
"notice_id": {
    Type: schema.TypeString,
    Optional: true,
    Description: "Notice template ID for query. Example: \"notice-xxx\"",
},
"tmpl_language": {
    Type: schema.TypeString,
    Optional: true,
    Description: "Template language: en/zh. Example: \"zh\"",
},
"monitor_type": {
    Type: schema.TypeString,
    Optional: true,
    Description: "Monitor type. Example: \"MT_QCE\"",
},
```

#### Output Schema

```go
"notice_content_tmpl_list": {
    Type: schema.TypeList,
    Computed: true,
    Description: "Notification content template list.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "tmpl_id": {
                Type: schema.TypeString,
                Computed: true,
                Description: "Template ID.",
            },
            "tmpl_name": {
                Type: schema.TypeString,
                Computed: true,
                Description: "Template name.",
            },
            "monitor_type": {
                Type: schema.TypeString,
                Computed: true,
                Description: "Monitor type.",
            },
            "tmpl_language": {
                Type: schema.TypeString,
                Computed: true,
                Description: "Template language.",
            },
            "creator": {
                Type: schema.TypeString,
                Computed: true,
                Description: "Creator UIN.",
            },
            "last_modifier": {
                Type: schema.TypeString,
                Computed: true,
                Description: "Last modifier UIN.",
            },
            "create_time": {
                Type: schema.TypeInt,
                Computed: true,
                Description: "Create time (Unix timestamp).",
            },
            "update_time": {
                Type: schema.TypeInt,
                Computed: true,
                Description: "Update time (Unix timestamp).",
            },
            "tmpl_contents_json": {
                Type: schema.TypeString,
                Computed: true,
                Description: "Template contents in JSON format.",
            },
            "bind_policy_count": {
                Type: schema.TypeInt,
                Computed: true,
                Description: "Number of bound alarm policies.",
            },
        },
    },
},
```

**Note**: The `tmpl_contents` field is extremely complex with nested structures. We will serialize it to JSON string for simplicity.

### 4.4 Service Layer Implementation

#### Method Signature

```go
func (me *MonitorService) DescribeNoticeContentTmplsByFilter(
    ctx context.Context,
    param map[string]interface{},
) (
    noticeContentTmpls []*monitorv20230616.NoticeContentTmpl,
    bindPolicyCounts map[string]*monitorv20230616.NoticeContentTmplBindPolicyCount,
    errRet error,
)
```

#### Implementation Pattern

```go
func (me *MonitorService) DescribeNoticeContentTmplsByFilter(ctx context.Context, param map[string]interface{}) (
    noticeContentTmpls []*monitorv20230616.NoticeContentTmpl,
    bindPolicyCounts map[string]*monitorv20230616.NoticeContentTmplBindPolicyCount,
    errRet error,
) {
    logId := tccommon.GetLogId(ctx)
    request := monitorv20230616.NewDescribeNoticeContentTmplRequest()

    // Parse parameters
    if v, ok := param["TmplIDs"]; ok {
        request.TmplIDs = v.([]*string)
    }
    if v, ok := param["TmplName"]; ok {
        request.TmplName = v.(*string)
    }
    if v, ok := param["NoticeID"]; ok {
        request.NoticeID = v.(*string)
    }
    if v, ok := param["TmplLanguage"]; ok {
        request.TmplLanguage = v.(*string)
    }
    if v, ok := param["MonitorType"]; ok {
        request.MonitorType = v.(*string)
    }

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
                logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    // Pagination loop
    var (
        pageNumber uint64 = 1
        pageSize   uint64 = 100
    )
    
    bindPolicyCounts = make(map[string]*monitorv20230616.NoticeContentTmplBindPolicyCount)

    for {
        request.PageNumber = &pageNumber
        request.PageSize = &pageSize

        // Retry logic for THIS page
        err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
            ratelimit.Check(request.GetAction())
            result, e := me.client.UseMonitorV20230616Client().DescribeNoticeContentTmpl(request)
            if e != nil {
                return tccommon.RetryError(e)
            }

            log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
                logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

            if result == nil || result.Response == nil {
                return resource.NonRetryableError(fmt.Errorf("Response is nil"))
            }

            // Accumulate results
            if result.Response.NoticeContentTmpls != nil {
                noticeContentTmpls = append(noticeContentTmpls, result.Response.NoticeContentTmpls...)
            }

            // Store bind policy counts
            if result.Response.NoticeContentTmplBindPolicyCounts != nil {
                for _, bindCount := range result.Response.NoticeContentTmplBindPolicyCounts {
                    if bindCount.TmplId != nil {
                        bindPolicyCounts[*bindCount.TmplId] = bindCount
                    }
                }
            }

            // Check if more pages
            if result.Response.NoticeContentTmpls == nil || 
               len(result.Response.NoticeContentTmpls) < int(pageSize) {
                return nil
            }

            return nil
        })

        if err != nil {
            errRet = err
            return
        }

        // Break if no more data
        if len(noticeContentTmpls) == 0 || 
           len(noticeContentTmpls)%int(pageSize) != 0 {
            break
        }

        pageNumber++
    }

    return
}
```

### 4.5 Data Source Implementation

#### Main Read Function

```go
func dataSourceTencentCloudMonitorNoticeContentTmplsRead(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("data_source.tencentcloud_monitor_notice_content_tmpls.read")()
    defer tccommon.InconsistentCheck(d, meta)()

    var (
        logId   = tccommon.GetLogId(nil)
        ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        service = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
    )

    paramMap := make(map[string]interface{})

    // Parse tmpl_ids
    if v, ok := d.GetOk("tmpl_ids"); ok {
        tmplIDsSet := v.(*schema.Set).List()
        tmplIDs := make([]*string, 0, len(tmplIDsSet))
        for _, item := range tmplIDsSet {
            tmplIDs = append(tmplIDs, helper.String(item.(string)))
        }
        paramMap["TmplIDs"] = tmplIDs
    }

    // Parse tmpl_name
    if v, ok := d.GetOk("tmpl_name"); ok {
        paramMap["TmplName"] = helper.String(v.(string))
    }

    // Parse notice_id
    if v, ok := d.GetOk("notice_id"); ok {
        paramMap["NoticeID"] = helper.String(v.(string))
    }

    // Parse tmpl_language
    if v, ok := d.GetOk("tmpl_language"); ok {
        paramMap["TmplLanguage"] = helper.String(v.(string))
    }

    // Parse monitor_type
    if v, ok := d.GetOk("monitor_type"); ok {
        paramMap["MonitorType"] = helper.String(v.(string))
    }

    // Call service layer with retry
    var (
        respData         []*monitorv20230616.NoticeContentTmpl
        bindPolicyCounts map[string]*monitorv20230616.NoticeContentTmplBindPolicyCount
    )

    reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        result, bindCounts, e := service.DescribeNoticeContentTmplsByFilter(ctx, paramMap)
        if e != nil {
            return tccommon.RetryError(e)
        }
        respData = result
        bindPolicyCounts = bindCounts
        return nil
    })

    if reqErr != nil {
        return reqErr
    }

    // Map response to schema
    tmplList := make([]map[string]interface{}, 0, len(respData))
    for _, tmpl := range respData {
        tmplMap := map[string]interface{}{}

        if tmpl.TmplID != nil {
            tmplMap["tmpl_id"] = tmpl.TmplID
        }
        if tmpl.TmplName != nil {
            tmplMap["tmpl_name"] = tmpl.TmplName
        }
        if tmpl.MonitorType != nil {
            tmplMap["monitor_type"] = tmpl.MonitorType
        }
        if tmpl.TmplLanguage != nil {
            tmplMap["tmpl_language"] = tmpl.TmplLanguage
        }
        if tmpl.Creator != nil {
            tmplMap["creator"] = tmpl.Creator
        }
        if tmpl.LastModifier != nil {
            tmplMap["last_modifier"] = tmpl.LastModifier
        }
        if tmpl.CreateTime != nil {
            tmplMap["create_time"] = tmpl.CreateTime
        }
        if tmpl.UpdateTime != nil {
            tmplMap["update_time"] = tmpl.UpdateTime
        }

        // Serialize complex TmplContents to JSON
        if tmpl.TmplContents != nil {
            if jsonBytes, err := json.Marshal(tmpl.TmplContents); err == nil {
                tmplMap["tmpl_contents_json"] = string(jsonBytes)
            }
        }

        // Get bind policy count
        if tmpl.TmplID != nil && bindPolicyCounts[*tmpl.TmplID] != nil {
            if bindPolicyCounts[*tmpl.TmplID].Count != nil {
                tmplMap["bind_policy_count"] = *bindPolicyCounts[*tmpl.TmplID].Count
            }
        }

        tmplList = append(tmplList, tmplMap)
    }

    _ = d.Set("notice_content_tmpl_list", tmplList)
    d.SetId(helper.BuildToken())

    // Export to file
    if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
        if e := tccommon.WriteToFile(output.(string), d); e != nil {
            return e
        }
    }

    return nil
}
```

---

## 5. Error Handling

### API Errors

| Error Type | Handling Strategy |
|------------|-------------------|
| Network timeout | Retry with exponential backoff |
| Rate limiting | Retry with ratelimit.Check |
| Invalid parameters | Return immediately with error |
| Permission denied | Return immediately with error |
| Resource not found | Return empty list (not error) |

### Edge Cases

| Case | Handling |
|------|----------|
| No templates exist | Return empty list, no error |
| Filter matches no templates | Return empty list, no error |
| API returns nil | Log warning, return empty list |
| Pagination returns 0 items | Break pagination loop |

---

## 6. Testing Strategy

### Unit Tests

```go
// Test parameter parsing
func TestParseParameters(t *testing.T) {
    // Test tmpl_ids parsing
    // Test other parameters
}

// Test pagination logic
func TestPagination(t *testing.T) {
    // Mock API with multiple pages
    // Verify all pages collected
}

// Test retry logic
func TestRetryOnError(t *testing.T) {
    // Mock API failure
    // Verify retry attempts
}
```

### Acceptance Tests

```go
const testAccMonitorNoticeContentTmplsDataSource = `
data "tencentcloud_monitor_notice_content_tmpls" "test" {
}
`

const testAccMonitorNoticeContentTmplsDataSourceWithFilter = `
data "tencentcloud_monitor_notice_content_tmpls" "test" {
  tmpl_language = "zh"
  monitor_type  = "MT_QCE"
}
`

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

## 7. Performance Considerations

### Optimization Strategies

1. **Batch Size**: Use 100 items per page (balance between API calls and memory)
2. **Parallel Requests**: Not implemented (API doesn't support parallel queries with same filters)
3. **Caching**: Not implemented (Terraform data sources are read-only, no caching needed)
4. **Memory**: Accumulate results in slices (reasonable for expected data volumes)

### Expected Performance

| Scenario | Templates | API Calls | Time (est.) |
|----------|-----------|-----------|-------------|
| Small | 1-100 | 1 | < 1s |
| Medium | 100-500 | 5 | 2-3s |
| Large | 500-1000 | 10 | 4-5s |

---

## 8. Security Considerations

### Data Protection

- Template contents may contain sensitive information (API keys, webhooks)
- No additional masking required (user must secure Terraform state)

### Access Control

- Uses existing Tencent Cloud CAM permissions
- Requires `monitor:DescribeNoticeContentTmpl` permission

### Audit

- All API calls logged with request/response bodies
- Terraform state tracks data source reads

---

## 9. Migration & Rollback

### Deployment Steps

1. Add service layer method
2. Create data source file
3. Register in provider
4. Add tests
5. Update documentation
6. Release version

### Rollback Plan

- New data source, no migration needed
- To rollback: Remove data source registration from provider
- No impact on existing resources

---

## 10. Documentation Requirements

### Provider Documentation

```markdown
# tencentcloud_monitor_notice_content_tmpls

Use this data source to query monitor notification content templates.

## Example Usage

### Query all templates

\`\`\`hcl
data "tencentcloud_monitor_notice_content_tmpls" "all" {
}
\`\`\`

### Query by template ID

\`\`\`hcl
data "tencentcloud_monitor_notice_content_tmpls" "by_id" {
  tmpl_ids = ["ntpl-3r1spzjn"]
}
\`\`\`

### Filter by language

\`\`\`hcl
data "tencentcloud_monitor_notice_content_tmpls" "zh_only" {
  tmpl_language = "zh"
}
\`\`\`

## Argument Reference

* `tmpl_ids` - (Optional, Set of String) Template ID list for query.
* `tmpl_name` - (Optional, String) Template name for query.
* `notice_id` - (Optional, String) Notice template ID for query.
* `tmpl_language` - (Optional, String) Template language: `en`/`zh`.
* `monitor_type` - (Optional, String) Monitor type, e.g., `MT_QCE`.
* `result_output_file` - (Optional, String) Used to save results.

## Attribute Reference

* `notice_content_tmpl_list` - Notification content template list.
  * `tmpl_id` - Template ID.
  * `tmpl_name` - Template name.
  * `monitor_type` - Monitor type.
  * `tmpl_language` - Template language.
  * `creator` - Creator UIN.
  * `last_modifier` - Last modifier UIN.
  * `create_time` - Create time (Unix timestamp).
  * `update_time` - Update time (Unix timestamp).
  * `tmpl_contents_json` - Template contents in JSON format.
  * `bind_policy_count` - Number of bound alarm policies.
```

---

## 11. Open Questions & Risks

### Open Questions

1. **Q**: Should we decode base64 content in template?  
   **A**: No, keep as-is. Users can decode if needed.

2. **Q**: Should we flatten the nested `TmplContents` structure?  
   **A**: No, serialize to JSON string for simplicity.

3. **Q**: Maximum number of templates expected?  
   **A**: Typically < 100 per account. Pagination handles larger cases.

### Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| API rate limiting | Medium | Medium | Retry logic + ratelimit.Check |
| Large response payload | Low | Low | Pagination with reasonable page size |
| Schema changes in API | Low | High | Monitor API updates, update schema |
| Complex nested structure | Medium | Low | Serialize to JSON string |

---

## 12. Success Metrics

### Acceptance Criteria

- [ ] All acceptance tests pass
- [ ] Code coverage > 80%
- [ ] Documentation complete
- [ ] Manual testing successful
- [ ] Code review approved
- [ ] Performance acceptable (< 5s for 100 templates)

### Post-Launch Metrics

- User adoption rate
- Error rate in production
- Performance in real-world scenarios

---

**Proposal Status**: Ready for Review  
**Last Updated**: 2026-03-24  
**Next Review**: TBD
