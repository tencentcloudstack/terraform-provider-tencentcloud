# Change Proposal: Add Monitor Notice Content Template Resource

## 📋 Overview

This proposal adds a new Terraform resource `tencentcloud_monitor_notice_content_tmpl` to manage Tencent Cloud Monitor notification content templates.

## 📁 Files Created

```
openspec/changes/add-monitor-notice-content-tmpl/
├── README.md                                    # This file
├── proposal.md                                  # High-level proposal
├── tasks.md                                     # Implementation checklist
├── design.md                                    # Technical design decisions
└── specs/
    └── monitor-notice-content-tmpl/
        └── spec.md                              # Detailed requirements
```

## 🎯 Key Features

### Resource Identifier
- **Format**: `tmplID#tmplName` (composite ID)
- **Example**: `ntpl-3r1spzjn#MyTemplate`

### Supported Operations
- ✅ **Create**: `CreateNoticeContentTmpl` API
- ✅ **Read**: `DescribeNoticeContentTmpl` API
- ✅ **Update**: `ModifyNoticeContentTmpl` API
- ✅ **Delete**: `DeleteNoticeContentTmpls` API
- ✅ **Import**: Support standard Terraform import

### Schema Highlights

```hcl
resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "my-custom-template"
  monitor_type  = "MT_QCE"
  tmpl_language = "zh"

  tmpl_contents {
    matching_status = ["Trigger"]
    
    template {
      we_work_robot {
        title_tmpl   = "告警通知"
        content_tmpl = "告警详情：{{.Content}}"
      }
      
      ding_ding_robot {
        title_tmpl   = "告警通知"
        content_tmpl = "告警详情：{{.Content}}"
      }
    }
  }
}
```

## 🔧 Technical Decisions

### 1. Composite ID Strategy
- Combines `tmplID` (from API) and `tmplName` (from user config)
- Enables efficient querying and deletion
- Follows pattern from `resource_tc_igtm_strategy.go`

### 2. ForceNew Fields
- `tmpl_name`: Template name change requires recreation
- `monitor_type`: Monitor type is immutable
- `tmpl_language`: Language setting is immutable

**Rationale**: These fields are fundamental identifiers in the Monitor system.

### 3. Complex Nested Schema
- Supports multiple notification channels (WeWork, DingDing, Feishu, etc.)
- Uses `MaxItems: 1` for single-object nesting
- Proper type conversion with helper functions

## 📚 API Documentation

| Operation | API Endpoint | Documentation |
|-----------|--------------|---------------|
| Create | CreateNoticeContentTmpl | https://cloud.tencent.com/document/api/248/128272 |
| Read | DescribeNoticeContentTmpl | https://cloud.tencent.com/document/api/248/128618 |
| Update | ModifyNoticeContentTmpl | https://cloud.tencent.com/document/api/248/128617 |
| Delete | DeleteNoticeContentTmpls | https://cloud.tencent.com/document/api/248/128619 |

## 🎓 Implementation Reference

**Primary Reference**: 
```
/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_strategy.go
```

**Key Patterns**:
- Composite ID handling
- Complex nested object mapping
- Context lifecycle management
- Retry logic with error handling
- Nil safety checks

## ✅ Implementation Checklist

See [`tasks.md`](./tasks.md) for detailed task breakdown.

**High-level phases**:
1. ✅ Create proposal documentation (DONE)
2. ⏳ Implement resource code
3. ⏳ Add service layer methods
4. ⏳ Register in provider
5. ⏳ Write tests
6. ⏳ Create documentation

## 🧪 Testing Strategy

### Acceptance Tests
- Basic CRUD operations
- Complex nested structure handling
- Import functionality
- Error scenarios (not found, malformed ID)

### Test Configuration Example
```hcl
resource "tencentcloud_monitor_notice_content_tmpl" "test" {
  tmpl_name     = "tf-test-template"
  monitor_type  = "MT_QCE"
  tmpl_language = "zh"

  tmpl_contents {
    matching_status = ["Trigger"]
    template {
      we_work_robot {
        title_tmpl   = "Test Template"
        content_tmpl = "Test Content"
      }
    }
  }
}
```

## 🚀 Next Steps

1. **Review and Approve**: Team review of this proposal
2. **Validate Spec**: Run `openspec validate add-monitor-notice-content-tmpl --strict`
3. **Implementation**: Follow tasks in `tasks.md`
4. **Testing**: Write and run acceptance tests
5. **Documentation**: Create user-facing documentation
6. **Archive**: After deployment, move to archive

## 📝 Notes

- Monitor API version: 2023-06-16
- Rate limit: 20 requests/second
- Composite ID separator: `tccommon.FILED_SP` (typically `#`)
- Template content may include sensitive information - consider marking fields as sensitive if needed

## 🔗 Related Resources

- Monitor Alarm Policy: `tencentcloud_monitor_alarm_policy`
- Monitor Alarm Notice: `tencentcloud_monitor_alarm_notice`
- Monitor Grafana Instance: `tencentcloud_monitor_grafana_instance`

---

**Proposal Status**: ⏳ Pending Review

**Created**: 2026-02-04

**Author**: AI Assistant following OpenSpec standards
