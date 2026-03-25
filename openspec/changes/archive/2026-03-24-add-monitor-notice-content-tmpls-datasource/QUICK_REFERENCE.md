# Quick Reference: Monitor Notice Content Templates DataSource

## 🚀 TL;DR

**What**: New data source to query monitor notification content templates  
**Why**: Enable Terraform to discover and reference existing templates  
**How**: Follow `data_source_tc_igtm_instance_list.go` pattern  
**Time**: 2-3 hours  
**Complexity**: Medium

---

## 📋 Checklist

### Before You Start
- [ ] Read `README.md` for overview
- [ ] Review `proposal.md` for technical details
- [ ] Check reference implementation: `data_source_tc_igtm_instance_list.go`
- [ ] Verify API access: [DescribeNoticeContentTmpl](https://cloud.tencent.com/document/product/248/128618)

### Implementation Steps
1. [ ] **Service Layer** - Add method to `service_tencentcloud_monitor.go`
2. [ ] **Data Source** - Create `data_source_tc_monitor_notice_content_tmpls.go`
3. [ ] **Registration** - Add to `provider.go`
4. [ ] **Testing** - Create test file and run
5. [ ] **Documentation** - Create `.md` file

### After Implementation
- [ ] Run `go fmt` on all modified files
- [ ] Run acceptance tests
- [ ] Verify manual testing
- [ ] Review code quality
- [ ] Update documentation

---

## 🔑 Key Points

### 1. API Information

```
API Name: DescribeNoticeContentTmpl
SDK: monitor/v20230616
Method: me.client.UseMonitorV20230616Client().DescribeNoticeContentTmpl(request)
```

### 2. Pagination Pattern

**Type**: PageNumber/PageSize (NOT Offset/Limit)

```go
pageNumber := uint64(1)
pageSize := uint64(100)

for {
    request.PageNumber = &pageNumber
    request.PageSize = &pageSize
    
    // ⚠️ CRITICAL: Retry INSIDE loop
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        result, e := me.client.UseMonitorV20230616Client().DescribeNoticeContentTmpl(request)
        // ...
        return nil
    })
    
    // Accumulate results
    // Break if no more pages
    
    pageNumber++
}
```

### 3. Input Parameters

| Parameter | Schema Type | SDK Type | Required |
|-----------|-------------|----------|----------|
| tmpl_ids | TypeSet | []*string | No |
| tmpl_name | TypeString | *string | No |
| notice_id | TypeString | *string | No |
| tmpl_language | TypeString | *string | No |
| monitor_type | TypeString | *string | No |

### 4. Output Structure

```go
"notice_content_tmpl_list": {
    Type: schema.TypeList,
    Computed: true,
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "tmpl_id":             {Type: schema.TypeString, Computed: true},
            "tmpl_name":           {Type: schema.TypeString, Computed: true},
            "monitor_type":        {Type: schema.TypeString, Computed: true},
            "tmpl_language":       {Type: schema.TypeString, Computed: true},
            "creator":             {Type: schema.TypeString, Computed: true},
            "last_modifier":       {Type: schema.TypeString, Computed: true},
            "create_time":         {Type: schema.TypeInt, Computed: true},
            "update_time":         {Type: schema.TypeInt, Computed: true},
            "tmpl_contents_json":  {Type: schema.TypeString, Computed: true},
            "bind_policy_count":   {Type: schema.TypeInt, Computed: true},
        },
    },
}
```

---

## 📝 Code Snippets

### Service Layer Method Signature

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

### Data Source Function

```go
func DataSourceTencentCloudMonitorNoticeContentTmpls() *schema.Resource {
    return &schema.Resource{
        Read:   dataSourceTencentCloudMonitorNoticeContentTmplsRead,
        Schema: map[string]*schema.Schema{
            // Input parameters
            "tmpl_ids": {
                Type:     schema.TypeSet,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
                Description: "Template ID list for query.",
            },
            // ... other input parameters
            
            // Output
            "notice_content_tmpl_list": {
                Type:     schema.TypeList,
                Computed: true,
                Description: "Notification content template list.",
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        // ... output fields
                    },
                },
            },
            "result_output_file": {
                Type:     schema.TypeString,
                Optional: true,
                Description: "Used to save results.",
            },
        },
    }
}
```

### Parameter Parsing Pattern

```go
paramMap := make(map[string]interface{})

// Parse tmpl_ids (Set -> []*string)
if v, ok := d.GetOk("tmpl_ids"); ok {
    tmplIDsSet := v.(*schema.Set).List()
    tmplIDs := make([]*string, 0, len(tmplIDsSet))
    for _, item := range tmplIDsSet {
        tmplIDs = append(tmplIDs, helper.String(item.(string)))
    }
    paramMap["TmplIDs"] = tmplIDs
}

// Parse simple string parameters
if v, ok := d.GetOk("tmpl_name"); ok {
    paramMap["TmplName"] = helper.String(v.(string))
}
```

### Response Mapping Pattern

```go
tmplList := make([]map[string]interface{}, 0, len(respData))
for _, tmpl := range respData {
    tmplMap := map[string]interface{}{}
    
    if tmpl.TmplID != nil {
        tmplMap["tmpl_id"] = tmpl.TmplID
    }
    if tmpl.TmplName != nil {
        tmplMap["tmpl_name"] = tmpl.TmplName
    }
    // ... map other fields
    
    // Serialize complex nested structure to JSON
    if tmpl.TmplContents != nil {
        if jsonBytes, err := json.Marshal(tmpl.TmplContents); err == nil {
            tmplMap["tmpl_contents_json"] = string(jsonBytes)
        }
    }
    
    // Look up bind policy count
    if tmpl.TmplID != nil && bindPolicyCounts[*tmpl.TmplID] != nil {
        if bindPolicyCounts[*tmpl.TmplID].Count != nil {
            tmplMap["bind_policy_count"] = *bindPolicyCounts[*tmpl.TmplID].Count
        }
    }
    
    tmplList = append(tmplList, tmplMap)
}

_ = d.Set("notice_content_tmpl_list", tmplList)
```

---

## 🧪 Testing

### Basic Test Configuration

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "test" {
}

output "template_count" {
  value = length(data.tencentcloud_monitor_notice_content_tmpls.test.notice_content_tmpl_list)
}
```

### Test with Filters

```hcl
data "tencentcloud_monitor_notice_content_tmpls" "filtered" {
  tmpl_ids      = ["ntpl-3r1spzjn"]
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
                    tcacctest.AccCheckTencentCloudDataSourceID(
                        "data.tencentcloud_monitor_notice_content_tmpls.test",
                    ),
                    resource.TestCheckResourceAttrSet(
                        "data.tencentcloud_monitor_notice_content_tmpls.test",
                        "notice_content_tmpl_list.#",
                    ),
                ),
            },
        },
    })
}
```

### Run Tests

```bash
# Format code
go fmt ./tencentcloud/services/monitor/data_source_tc_monitor_notice_content_tmpls*.go

# Run acceptance test
cd tencentcloud
TF_ACC=1 go test -v -run TestAccTencentCloudMonitorNoticeContentTmplsDataSource ./services/monitor/

# Run all monitor data source tests
TF_ACC=1 go test -v ./services/monitor/ -run DataSource
```

---

## ⚠️ Common Pitfalls

### ❌ WRONG: Retry Outside Loop

```go
// ❌ BAD: No retry for subsequent pages
err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    for {
        result, e := client.DescribeNoticeContentTmpl(request)
        // Only first call is retried!
        pageNumber++
    }
    return nil
})
```

### ✅ CORRECT: Retry Inside Loop

```go
// ✅ GOOD: Each page call is retried
for {
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        result, e := client.DescribeNoticeContentTmpl(request)
        // Each page gets retry protection
        return nil
    })
    if err != nil {
        return
    }
    pageNumber++
}
```

### ❌ WRONG: Missing Nil Checks

```go
// ❌ BAD: Panic if pointer is nil
tmplMap["tmpl_name"] = *tmpl.TmplName
```

### ✅ CORRECT: Always Check Nil

```go
// ✅ GOOD: Safe nil check
if tmpl.TmplName != nil {
    tmplMap["tmpl_name"] = tmpl.TmplName
}
```

---

## 📁 File Locations

```
/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/

├── tencentcloud/
│   ├── provider.go                                           (UPDATE: +1 line)
│   └── services/
│       └── monitor/
│           ├── service_tencentcloud_monitor.go               (UPDATE: +80 lines)
│           ├── data_source_tc_monitor_notice_content_tmpls.go     (NEW: ~300 lines)
│           ├── data_source_tc_monitor_notice_content_tmpls_test.go (NEW: ~50 lines)
│           └── data_source_tc_monitor_notice_content_tmpls.md     (NEW: ~150 lines)
└── openspec/
    └── changes/
        └── add-monitor-notice-content-tmpls-datasource/
            ├── README.md
            ├── proposal.md
            ├── tasks.md
            ├── SUMMARY.md
            └── QUICK_REFERENCE.md (this file)
```

---

## 🔍 Reference Files

### Primary Reference
```
/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go
```

**What to copy**:
- Overall file structure
- Schema definition pattern
- Read function structure
- Parameter parsing logic
- Service call with retry
- Response mapping pattern

### Service Layer Reference
```
/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/service_tencentcloud_igtm.go
```

**Method**: `DescribeIgtmInstanceListByFilter` (line 398-460)

**What to copy**:
- Pagination loop structure
- Retry logic placement (INSIDE loop)
- Result accumulation
- Break condition

---

## 🎯 Success Criteria

### Must Have ✅
- [ ] Follows reference implementation pattern exactly
- [ ] Retry logic inside pagination loop
- [ ] All parameters properly parsed
- [ ] All fields have nil checks
- [ ] Acceptance test passes
- [ ] Code formatted with `go fmt`
- [ ] Documentation complete

### Nice to Have 🌟
- [ ] Multiple acceptance tests for different filters
- [ ] Performance test with 100+ templates
- [ ] Integration test with actual API

---

## 📞 Help & Resources

### Documentation
- **API Docs**: https://cloud.tencent.com/document/product/248/128618
- **Provider Docs**: Check existing monitor data sources for patterns

### Code Examples
- **Similar DataSource**: `data_source_tc_monitor_alarm_notices.go`
- **Reference Pattern**: `data_source_tc_igtm_instance_list.go`
- **Service Pattern**: `service_tencentcloud_igtm.go`

### Debugging
```bash
# Enable debug logging
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log

# Run Terraform
terraform plan

# Check logs
cat terraform.log | grep -i "monitor.*notice"
```

---

## 🚦 Status Indicators

### Current Status: 📋 Proposal

- ✅ **README.md** - Complete
- ✅ **proposal.md** - Complete
- ✅ **tasks.md** - Complete
- ✅ **SUMMARY.md** - Complete
- ✅ **QUICK_REFERENCE.md** - Complete

### Next Status: 🚧 Implementation

After implementation starts, update status in README.md

### Final Status: ✅ Completed

After all tasks complete and tests pass

---

## ⏱️ Time Estimates

| Phase | Minimum | Maximum | Average |
|-------|---------|---------|---------|
| Service Layer | 20 min | 40 min | 30 min |
| Data Source | 45 min | 90 min | 60 min |
| Registration | 5 min | 15 min | 10 min |
| Testing | 20 min | 45 min | 30 min |
| Documentation | 20 min | 45 min | 30 min |
| **Total** | **110 min** | **235 min** | **160 min** |

**Realistic Estimate**: 2-3 hours (including breaks and debugging)

---

**Quick Reference Created**: 2026-03-24  
**Last Updated**: 2026-03-24  
**Status**: Ready for Use

---

## 💡 Pro Tips

1. **Start with service layer** - Get the API call working first
2. **Test incrementally** - Don't wait until everything is done
3. **Use debug logs** - `log.Printf` is your friend
4. **Check nil everywhere** - Pointers can be nil in Go
5. **Follow the pattern** - When in doubt, copy reference implementation
6. **Format early, format often** - Run `go fmt` after every change

---

🎉 **Ready to implement? Start with Phase 1 in tasks.md!**
