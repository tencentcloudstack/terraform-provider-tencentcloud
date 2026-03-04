# Design: Monitor Notice Content Template Resource

## Overview

This resource manages Tencent Cloud Monitor notification content templates through Terraform. It allows users to define custom notification templates for different channels (WeWork Robot, DingDing Robot, Feishu Robot, etc.) with full CRUD support.

## SDK Information

**Package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616`

**Import Statement**:
```go
import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

    tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)
```

## Technical Decisions

### 1. Resource Identifier Strategy

**Decision**: Use composite ID format `tmplID#tmplName`

**Rationale**:
- API returns `TmplID` from Create operation
- Query API requires both `TmplID` and `TmplName` for precise lookups
- Delete API only needs `TmplID`
- Composite ID provides complete context for all operations
- Follows pattern from `resource_tc_igtm_strategy.go`

**Implementation**:
```go
// Create: Set ID after API returns TmplID
tmplId := *response.Response.TmplID
tmplName := d.Get("tmpl_name").(string)
d.SetId(strings.Join([]string{tmplId, tmplName}, tccommon.FILED_SP))

// Read/Update/Delete: Split ID
idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
if len(idSplit) != 2 {
    return fmt.Errorf("id is broken,%s", d.Id())
}
tmplId := idSplit[0]
tmplName := idSplit[1]
```

### 2. Schema Design

**Complex Nested Structure**:
The `tmpl_contents` field contains deeply nested objects representing different notification channels and their configurations.

**Schema Structure**:
```
tmpl_contents (TypeList, MaxItems: 1)
  └─ matching_status (TypeList of TypeString) - e.g., ["Trigger"]
  └─ template (TypeList, MaxItems: 1)
      ├─ q_cloud_yehe (TypeList, MaxItems: 1)
      │   ├─ sms (TypeList, MaxItems: 1)
      │   │   ├─ title_tmpl (TypeString)
      │   │   └─ content_tmpl (TypeString)
      │   ├─ email (TypeList, MaxItems: 1)
      │   ├─ voice (TypeList, MaxItems: 1)
      │   └─ site (TypeList, MaxItems: 1)
      ├─ we_work_robot (TypeList, MaxItems: 1)
      │   ├─ title_tmpl (TypeString)
      │   └─ content_tmpl (TypeString)
      ├─ ding_ding_robot (TypeList, MaxItems: 1)
      └─ fei_shu_robot (TypeList, MaxItems: 1)
```

**Key Schema Fields**:
- `tmpl_name` (Required, ForceNew) - Template name
- `monitor_type` (Required, ForceNew) - Monitor type, e.g., "MT_QCE"
- `tmpl_language` (Required, ForceNew) - Template language: "en" or "zh"
- `tmpl_contents` (Required) - Complex nested template content structure
- `tmpl_id` (Computed) - Template ID returned by API

### 3. API Integration

**Create Flow**:
1. Call `CreateNoticeContentTmpl` with `TmplName`, `MonitorType`, `TmplLanguage`, `TmplContents`
2. API returns `TmplID`
3. Set resource ID as `tmplID#tmplName`
4. Call Read to populate state

**Read Flow**:
1. Parse composite ID to extract `tmplID` and `tmplName`
2. Call `DescribeNoticeContentTmpl` with:
   - `TmplIDs`: `[tmplID]` (string array)
   - `TmplName`: `tmplName` (optional filter)
   - `PageNumber`: 1
   - `PageSize`: 10
3. Check if template exists in response
4. Populate state with returned data

**Update Flow**:
1. Parse composite ID to extract `tmplID` and `tmplName`
2. Detect changes using `d.HasChange()`
3. Call `ModifyNoticeContentTmpl` with:
   - `TmplID`: `tmplID`
   - `TmplName`: updated name if changed
   - `TmplContents`: updated content
4. Update resource ID if `tmpl_name` changed
5. Call Read to refresh state

**Delete Flow**:
1. Parse composite ID to extract `tmplID`
2. Call `DeleteNoticeContentTmpls` with `TmplIDs: [tmplID]`

### 4. Error Handling

**Retry Strategy**:
- Use `resource.Retry` with `tccommon.WriteRetryTimeout` for Create/Update/Delete
- Use `resource.Retry` with `tccommon.ReadRetryTimeout` for Read
- Apply `tccommon.RetryError()` for transient failures

**Not Found Handling**:
```go
if respData == nil {
    log.Printf("[WARN]%s resource `tencentcloud_monitor_notice_content_tmpl` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
    d.SetId("")
    return nil
}
```

### 5. Type Conversions

**Following helper patterns**:
- `helper.String()` for string pointers
- `helper.IntInt64()` / `helper.IntUint64()` for integer conversions
- Manual iteration for complex nested structures
- Safe type assertions with `ok` checks

### 6. Validation

**ForceNew Fields**:
- `tmpl_name`: Changing requires recreation (API limitation)
- `monitor_type`: Fixed at creation
- `tmpl_language`: Fixed at creation

**Rationale**: These fields are fundamental identifiers that cannot be modified without affecting template binding relationships.

## Code Organization

### File Structure
```
tencentcloud/services/monitor/
├── resource_tc_monitor_notice_content_tmpl.go  (new)
├── resource_tc_monitor_notice_content_tmpl.md  (new, documentation)
└── service_tencentcloud_monitor.go             (add service method)
```

### Naming Conventions
- Resource function: `ResourceTencentCloudMonitorNoticeContentTmpl()`
- CRUD functions: `resourceTencentCloudMonitorNoticeContentTmpl{Create|Read|Update|Delete}`
- Service function: `DescribeMonitorNoticeContentTmplById(ctx, tmplId, tmplName)`

## Testing Strategy

**Acceptance Test Structure**:
```go
func TestAccTencentCloudMonitorNoticeContentTmpl_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { tcacctest.AccPreCheck(t) },
        Providers:    tcacctest.AccProviders,
        CheckDestroy: testAccCheckMonitorNoticeContentTmplDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccMonitorNoticeContentTmpl_basic,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckMonitorNoticeContentTmplExists("tencentcloud_monitor_notice_content_tmpl.test"),
                    resource.TestCheckResourceAttrSet("tencentcloud_monitor_notice_content_tmpl.test", "tmpl_id"),
                    resource.TestCheckResourceAttr("tencentcloud_monitor_notice_content_tmpl.test", "tmpl_name", "tf-test-template"),
                ),
            },
            {
                Config: testAccMonitorNoticeContentTmpl_update,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("tencentcloud_monitor_notice_content_tmpl.test", "tmpl_contents.0.template.0.we_work_robot.0.title_tmpl", "Updated Title"),
                ),
            },
            {
                ResourceName:      "tencentcloud_monitor_notice_content_tmpl.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}
```

## Reference Implementation

Primary reference: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_strategy.go`

**Key patterns to follow**:
1. Composite ID handling with `strings.Split()` and `strings.Join()`
2. Complex nested object mapping between Schema and SDK types
3. Context lifecycle management with `tccommon.NewResourceLifeCycleHandleFuncContext()`
4. Retry logic with proper error handling
5. Nil safety checks for all pointer dereferences
6. Consistent logging patterns

## Security Considerations

- Template content may contain sensitive information
- Consider marking sensitive fields with `Sensitive: true` if needed
- Ensure proper validation of template formats to prevent injection
- Log API responses without exposing sensitive template content
