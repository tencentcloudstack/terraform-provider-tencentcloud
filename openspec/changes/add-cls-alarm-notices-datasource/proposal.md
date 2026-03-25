# Proposal: Add tencentcloud_cls_alarm_notices Data Source

**Feature**: New data source for querying CLS alarm notification channel groups  
**Resource Name**: `tencentcloud_cls_alarm_notices`  
**Type**: Data Source  
**Product**: CLS (Cloud Log Service)  
**Priority**: Medium  
**Estimated Effort**: 2-3 hours

---

## 1. Executive Summary

### 1.1 Problem Statement

Currently, the Terraform Provider for Tencent Cloud lacks a data source to query CLS (Cloud Log Service) alarm notification channel groups. Users cannot programmatically retrieve alarm notice configurations for automation, validation, or reference purposes.

**User Pain Points**:
- ❌ Cannot query existing alarm notice configurations in Terraform
- ❌ Cannot reference alarm notice IDs from existing configurations
- ❌ Cannot validate alarm notice settings programmatically
- ❌ Cannot list alarm notices for documentation or auditing

### 1.2 Proposed Solution

Add a new data source `tencentcloud_cls_alarm_notices` that wraps the `DescribeAlarmNotices` API to allow users to query alarm notification channel groups with flexible filtering options.

**Key Benefits**:
- ✅ Query alarm notices by name, ID, or other attributes
- ✅ Reference existing alarm notice IDs in other resources
- ✅ Support filtering by user ID, group ID, and delivery status
- ✅ Enable automation and validation workflows
- ✅ Provide detailed alarm notice configuration data

---

## 2. Background

### 2.1 CLS Alarm Notices Overview

**What are Alarm Notices?**
- Notification channel groups for CLS alarms
- Configure how and when alarm notifications are sent
- Support multiple notification types (email, SMS, WeChat, webhooks)
- Include notification rules, escalation policies, and delivery configs

**Use Cases**:
- Configure alarm notifications for log monitoring
- Set up multi-channel alert delivery
- Implement escalation workflows
- Manage notification recipients and callbacks

---

### 2.2 API Information

**API Name**: `DescribeAlarmNotices`  
**API Version**: `2020-10-16`  
**Request Domain**: `cls.tencentcloudapi.com`  
**Rate Limit**: 20 requests/second

**Supported Filters**:
1. `name` - Alarm notice group name
2. `alarmNoticeId` - Alarm notice ID
3. `uid` - Receiver user ID
4. `groupId` - Receiver user group ID
5. `deliverFlag` - Delivery status (1: Not enabled, 2: Enabled, 3: Delivery exception)

**Pagination**: Yes (Offset + Limit, max 100 per page)

---

## 3. Technical Design

### 3.1 Reference Implementation

This data source will follow the exact pattern of `tencentcloud_igtm_instance_list`:
- File: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`
- Service pattern: Service layer with pagination and retry logic
- Schema design: Filters as input, list as output

---

### 3.2 File Structure

```
tencentcloud/services/cls/
├── data_source_tc_cls_alarm_notices.go      (NEW - data source definition)
└── service_tencentcloud_cls.go              (MODIFY - add service method)
```

---

### 3.3 Schema Design

#### Input Parameters

```hcl
data "tencentcloud_cls_alarm_notices" "example" {
  filters {
    name = "name"
    values = ["my-alarm-notice"]
  }
  
  filters {
    name = "alarmNoticeId"
    values = ["notice-xxxx"]
  }
  
  filters {
    name = "deliverFlag"
    values = ["2"]  # Enabled
  }
  
  has_alarm_shield_count = true
  result_output_file = "alarm_notices.json"
}
```

#### Output Attributes

```hcl
output "alarm_notices" {
  value = data.tencentcloud_cls_alarm_notices.example.alarm_notices
}

# Output structure:
# alarm_notices = [
#   {
#     name                       = "test-alarm"
#     alarm_notice_id            = "notice-xxxx"
#     create_time                = "2025-08-06 15:47:00"
#     update_time                = "2025-08-06 15:47:00"
#     notice_receivers           = [...]
#     web_callbacks              = [...]
#     tags                       = [...]
#     jump_domain                = "https://console.cloud.tencent.com"
#     notice_rules               = [...]
#     deliver_status             = 1
#     deliver_flag               = 1
#     alarm_notice_deliver_config = {...}
#     alarm_shield_status        = 2
#     alarm_shield_count         = {...}
#     callback_prioritize        = true
#   }
# ]
```

---

### 3.4 Detailed Schema Definition

#### Data Source Schema

```go
"filters": {
    Type:        schema.TypeList,
    Optional:    true,
    Description: "Filter conditions. Maximum 10 filters, each filter can have up to 5 values.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "name": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "Filter field name. Supported values: name, alarmNoticeId, uid, groupId, deliverFlag.",
            },
            "values": {
                Type:        schema.TypeSet,
                Required:    true,
                Description: "Filter values. Maximum 5 values per filter.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    },
},

"has_alarm_shield_count": {
    Type:        schema.TypeBool,
    Optional:    true,
    Description: "Whether to return alarm shield count statistics. Default: false.",
},

"alarm_notices": {
    Type:        schema.TypeList,
    Computed:    true,
    Description: "Alarm notice list.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "name": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Alarm notice group name.",
            },
            "alarm_notice_id": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Alarm notice group ID.",
            },
            "create_time": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Creation time.",
            },
            "update_time": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Update time.",
            },
            "notice_receivers": {
                Type:        schema.TypeList,
                Computed:    true,
                Description: "Notification receivers.",
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "receiver_type": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Receiver type.",
                        },
                        "receiver_ids": {
                            Type:        schema.TypeSet,
                            Computed:    true,
                            Description: "Receiver IDs.",
                            Elem: &schema.Schema{Type: schema.TypeString},
                        },
                        "receiver_channels": {
                            Type:        schema.TypeSet,
                            Computed:    true,
                            Description: "Notification channels.",
                            Elem: &schema.Schema{Type: schema.TypeString},
                        },
                    },
                },
            },
            "web_callbacks": {
                Type:        schema.TypeList,
                Computed:    true,
                Description: "Web callback configurations.",
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "url": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Callback URL.",
                        },
                        "callback_type": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Callback type (WeCom, DingTalk, etc).",
                        },
                        "method": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "HTTP method.",
                        },
                        "headers": {
                            Type:        schema.TypeSet,
                            Computed:    true,
                            Description: "HTTP headers.",
                            Elem: &schema.Schema{Type: schema.TypeString},
                        },
                        "body": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Request body.",
                        },
                        "index": {
                            Type:        schema.TypeInt,
                            Computed:    true,
                            Description: "Callback index.",
                        },
                    },
                },
            },
            "tags": {
                Type:        schema.TypeList,
                Computed:    true,
                Description: "Tags.",
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "key": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Tag key.",
                        },
                        "value": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "Tag value.",
                        },
                    },
                },
            },
            "jump_domain": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Jump domain for console link.",
            },
            "deliver_status": {
                Type:        schema.TypeInt,
                Computed:    true,
                Description: "Delivery status.",
            },
            "deliver_flag": {
                Type:        schema.TypeInt,
                Computed:    true,
                Description: "Delivery flag. 1: Not enabled, 2: Enabled, 3: Delivery exception.",
            },
            "alarm_shield_status": {
                Type:        schema.TypeInt,
                Computed:    true,
                Description: "Alarm shield status.",
            },
            "alarm_shield_count": {
                Type:        schema.TypeList,
                Computed:    true,
                Description: "Alarm shield count statistics.",
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "total_count": {
                            Type:        schema.TypeInt,
                            Computed:    true,
                            Description: "Total count.",
                        },
                        "invalid_count": {
                            Type:        schema.TypeInt,
                            Computed:    true,
                            Description: "Invalid count.",
                        },
                        "valid_count": {
                            Type:        schema.TypeInt,
                            Computed:    true,
                            Description: "Valid count.",
                        },
                        "expire_count": {
                            Type:        schema.TypeInt,
                            Computed:    true,
                            Description: "Expired count.",
                        },
                    },
                },
            },
            "callback_prioritize": {
                Type:        schema.TypeBool,
                Computed:    true,
                Description: "Whether callback is prioritized.",
            },
        },
    },
},

"total_count": {
    Type:        schema.TypeInt,
    Computed:    true,
    Description: "Total number of alarm notices matching the filter.",
},

"result_output_file": {
    Type:        schema.TypeString,
    Optional:    true,
    Description: "Used to save results.",
},
```

---

## 4. Implementation Details

### 4.1 Data Source File

**File**: `tencentcloud/services/cls/data_source_tc_cls_alarm_notices.go`

**Key Components**:

1. **Schema Definition**
   - Filters (optional): name, values
   - has_alarm_shield_count (optional): boolean
   - alarm_notices (computed): list of alarm notice objects
   - total_count (computed): integer
   - result_output_file (optional): string

2. **Read Function**
   ```go
   func dataSourceTencentCloudClsAlarmNoticesRead(d *schema.ResourceData, meta interface{}) error {
       // Initialize context and service
       // Parse filters
       // Call service layer with retry
       // Map response to schema
       // Set output and file (if specified)
   }
   ```

3. **Filter Parsing**
   - Parse filters from schema
   - Convert to []*cls.Filter
   - Support multiple filter combinations

4. **Response Mapping**
   - Map AlarmNotice objects to schema
   - Handle nested structures (receivers, callbacks, tags)
   - Handle nil values gracefully

---

### 4.2 Service Layer Method

**File**: `tencentcloud/services/cls/service_tencentcloud_cls.go`

**Method**: `DescribeClsAlarmNoticesByFilter`

**Implementation Pattern**:
```go
func (me *ClsService) DescribeClsAlarmNoticesByFilter(
    ctx context.Context, 
    param map[string]interface{},
) (ret []*cls.AlarmNotice, totalCount *uint64, errRet error) {
    var (
        logId    = tccommon.GetLogId(ctx)
        request  = cls.NewDescribeAlarmNoticesRequest()
        response = cls.NewDescribeAlarmNoticesResponse()
    )

    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", 
                logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()

    // Parse parameters
    for k, v := range param {
        if k == "Filters" {
            request.Filters = v.([]*cls.Filter)
        }
        if k == "HasAlarmShieldCount" {
            request.HasAlarmShieldCount = v.(*bool)
        }
    }

    // Pagination loop with retry inside
    var (
        offset uint64 = 0
        limit  uint64 = 100
    )
    for {
        request.Offset = &offset
        request.Limit = &limit
        
        // IMPORTANT: Retry inside pagination loop
        err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
            ratelimit.Check(request.GetAction())
            result, e := me.client.UseClsClient().DescribeAlarmNotices(request)
            if e != nil {
                return tccommon.RetryError(e)
            }
            
            log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", 
                logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

            if result == nil || result.Response == nil {
                return resource.NonRetryableError(
                    fmt.Errorf("DescribeAlarmNotices failed, Response is nil."))
            }

            response = result
            return nil
        })

        if err != nil {
            errRet = err
            return
        }

        ret = append(ret, response.Response.AlarmNotices...)
        totalCount = response.Response.TotalCount
        
        if len(response.Response.AlarmNotices) < int(limit) {
            break
        }

        offset += limit
    }

    return
}
```

**Key Design Points**:
1. **Pagination with Retry**: Retry logic is inside the for loop (as required)
2. **Error Handling**: Proper defer logging and error returns
3. **Rate Limiting**: Check rate limit before each API call
4. **Response Validation**: Check for nil responses
5. **TotalCount**: Return total count for reference

---

## 5. Code Implementation

### 5.1 Files to Create/Modify

| File | Type | Lines | Description |
|------|------|-------|-------------|
| `data_source_tc_cls_alarm_notices.go` | NEW | ~400 | Data source definition |
| `service_tencentcloud_cls.go` | MODIFY | ~80 | Add service method |
| **Total** | - | **~480** | **2 files** |

---

### 5.2 Dependencies

**SDK Package**:
```go
import (
    cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
)
```

**Existing CLS Service**: `ClsService` struct already exists

---

## 6. Testing Strategy

### 6.1 Manual Testing Scenarios

#### Test Case 1: Query All Alarm Notices
```hcl
data "tencentcloud_cls_alarm_notices" "all" {}

output "all_notices" {
  value = data.tencentcloud_cls_alarm_notices.all.alarm_notices
}
```
**Expected**: Returns all alarm notices

---

#### Test Case 2: Filter by Name
```hcl
data "tencentcloud_cls_alarm_notices" "by_name" {
  filters {
    name   = "name"
    values = ["test-alarm"]
  }
}
```
**Expected**: Returns alarm notices matching name

---

#### Test Case 3: Filter by ID
```hcl
data "tencentcloud_cls_alarm_notices" "by_id" {
  filters {
    name   = "alarmNoticeId"
    values = ["notice-xxxx-yyyy"]
  }
}
```
**Expected**: Returns specific alarm notice

---

#### Test Case 4: Filter by Delivery Status
```hcl
data "tencentcloud_cls_alarm_notices" "enabled" {
  filters {
    name   = "deliverFlag"
    values = ["2"]  # Enabled
  }
  has_alarm_shield_count = true
}
```
**Expected**: Returns only enabled alarm notices with shield count

---

#### Test Case 5: Multiple Filters
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
**Expected**: Returns alarm notices matching all filters

---

#### Test Case 6: Export to File
```hcl
data "tencentcloud_cls_alarm_notices" "export" {
  result_output_file = "alarm_notices.json"
}
```
**Expected**: Creates JSON file with results

---

#### Test Case 7: Reference in Other Resources
```hcl
data "tencentcloud_cls_alarm_notices" "existing" {
  filters {
    name   = "name"
    values = ["my-notice"]
  }
}

resource "tencentcloud_cls_alarm" "example" {
  # Reference the alarm notice ID
  notice_id = data.tencentcloud_cls_alarm_notices.existing.alarm_notices[0].alarm_notice_id
}
```
**Expected**: Successfully references alarm notice ID

---

#### Test Case 8: Pagination Handling
```hcl
# Test with account that has > 100 alarm notices
data "tencentcloud_cls_alarm_notices" "many" {}
```
**Expected**: Returns all alarm notices across multiple pages

---

### 6.2 Test Validation Points

For each test:
- ✅ Data source executes without errors
- ✅ Returns expected number of results
- ✅ All fields are correctly populated
- ✅ Nested structures (callbacks, receivers) are properly parsed
- ✅ No state drift on subsequent runs
- ✅ result_output_file creates valid JSON
- ✅ Pagination works correctly (> 100 results)

---

## 7. Documentation

### 7.1 Data Source Documentation

**File**: `website/docs/d/cls_alarm_notices.html.markdown`

```markdown
---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_alarm_notices"
sidebar_current: "docs-tencentcloud-datasource-cls-alarm-notices"
description: |-
  Use this data source to query CLS alarm notification channel groups.
---

# tencentcloud_cls_alarm_notices

Use this data source to query CLS alarm notification channel groups.

## Example Usage

### Query all alarm notices

```hcl
data "tencentcloud_cls_alarm_notices" "all" {}
```

### Filter by name

```hcl
data "tencentcloud_cls_alarm_notices" "by_name" {
  filters {
    name   = "name"
    values = ["my-alarm-notice"]
  }
}
```

### Filter by delivery status

```hcl
data "tencentcloud_cls_alarm_notices" "enabled" {
  filters {
    name   = "deliverFlag"
    values = ["2"]  # Enabled
  }
  has_alarm_shield_count = true
}
```

### Reference alarm notice ID

```hcl
data "tencentcloud_cls_alarm_notices" "existing" {
  filters {
    name   = "alarmNoticeId"
    values = ["notice-xxxx"]
  }
}

resource "tencentcloud_cls_alarm" "example" {
  notice_id = data.tencentcloud_cls_alarm_notices.existing.alarm_notices[0].alarm_notice_id
}
```

## Argument Reference

* `filters` - (Optional) Filter conditions. Maximum 10 filters, each filter can have up to 5 values.
  * `name` - (Required) Filter field name. Supported values: `name`, `alarmNoticeId`, `uid`, `groupId`, `deliverFlag`.
  * `values` - (Required) Filter values. Maximum 5 values per filter.
* `has_alarm_shield_count` - (Optional) Whether to return alarm shield count statistics. Default: `false`.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

* `alarm_notices` - Alarm notice list.
  * `name` - Alarm notice group name.
  * `alarm_notice_id` - Alarm notice group ID.
  * `create_time` - Creation time.
  * `update_time` - Update time.
  * `notice_receivers` - Notification receivers.
  * `web_callbacks` - Web callback configurations.
  * `tags` - Tags.
  * `jump_domain` - Jump domain for console link.
  * `deliver_status` - Delivery status.
  * `deliver_flag` - Delivery flag. 1: Not enabled, 2: Enabled, 3: Delivery exception.
  * `alarm_shield_status` - Alarm shield status.
  * `alarm_shield_count` - Alarm shield count statistics (when `has_alarm_shield_count` is true).
  * `callback_prioritize` - Whether callback is prioritized.
* `total_count` - Total number of alarm notices matching the filter.
```

---

## 8. Success Criteria

### 8.1 Must Have
- [x] Data source definition created
- [x] Service layer method implemented
- [x] Follows exact pattern of reference implementation
- [x] Pagination with retry in for loop
- [x] All filters supported
- [x] Manual tests pass
- [x] Code formatted (go fmt)
- [x] Documentation created

### 8.2 Should Have
- [ ] Integration tests (optional but recommended)
- [ ] Multiple filter combinations tested
- [ ] Error handling validated
- [ ] Large dataset pagination tested (> 100 records)

### 8.3 Nice to Have
- [ ] Examples in documentation
- [ ] Changelog entry
- [ ] Code review by team

---

## 9. Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| API response structure mismatch | Low | Medium | Validate against API documentation |
| Pagination issues | Low | Medium | Test with large datasets |
| Filter parameter errors | Low | Low | Validate filter names and values |
| Nil pointer errors | Low | Medium | Add comprehensive nil checks |
| SDK version incompatibility | Low | Low | Use correct SDK package version |

---

## 10. Timeline

| Phase | Duration | Tasks |
|-------|----------|-------|
| **Development** | 1.5 hours | Code implementation |
| **Testing** | 0.5 hours | Manual testing |
| **Documentation** | 0.5 hours | Write documentation |
| **Review** | 0.5 hours | Code review and polish |
| **Total** | **3 hours** | **Complete implementation** |

---

## 11. References

### 11.1 API Documentation
- **DescribeAlarmNotices**: https://cloud.tencent.com/document/api/614/56462
- **API Version**: 2020-10-16

### 11.2 Reference Implementation
- **File**: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`
- **Pattern**: Data source with filters, pagination, and retry

### 11.3 Related Resources
- `tencentcloud_cls_alarm_notice` - Resource for managing alarm notices
- `tencentcloud_cls_alarm` - Resource for managing alarms

---

## 12. Appendix

### 12.1 Filter Examples

**Supported Filter Keys**:
- `name`: Alarm notice group name (string)
- `alarmNoticeId`: Alarm notice ID (string)
- `uid`: Receiver user ID (string)
- `groupId`: Receiver user group ID (string)
- `deliverFlag`: Delivery status (string: "1", "2", "3")

**Filter Limitations**:
- Maximum 10 filters per request
- Maximum 5 values per filter
- Values are exact match (no fuzzy search for this API)

---

### 12.2 API Response Sample

```json
{
  "Response": {
    "TotalCount": 1,
    "AlarmNotices": [
      {
        "Name": "test-alarm",
        "AlarmNoticeId": "notice-xxxx",
        "CreateTime": "2025-08-06 15:47:00",
        "UpdateTime": "2025-08-06 15:47:00",
        "DeliverFlag": 2,
        "DeliverStatus": 1,
        "AlarmShieldStatus": 2,
        "CallbackPrioritize": true
      }
    ],
    "RequestId": "xxxx-xxxx-xxxx"
  }
}
```

---

**Proposal Status**: ✅ Ready for Implementation  
**Next Steps**: Review proposal → Implement code → Test → Document → Merge
