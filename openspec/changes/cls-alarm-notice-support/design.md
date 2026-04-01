## Context

`tencentcloud_cls_alarm` 资源当前实现中,`alarm_notice_ids` 被定义为必填字段(Required),用于关联日志服务的告警通知渠道组。但根据 CLS CreateAlarm API 文档和 SDK 定义,`AlarmNoticeIds` 和 `MonitorNotice` 是两个互斥的可选参数:
- `AlarmNoticeIds`: 关联日志服务告警通知渠道组
- `MonitorNotice`: 关联可观测平台通知模板,包含通知规则列表(NoticeId、ContentTmplId、AlarmLevels)

现有实现仅支持 `AlarmNoticeIds`,且错误地将其设为必填,不符合 API 规范,也阻止用户使用可观测平台的统一告警能力。

## Goals / Non-Goals

**Goals:**
- 支持 `monitor_notice` 参数,允许配置可观测平台通知模板
- 修正 `alarm_notice_ids` 为可选字段,与 API 规范保持一致
- 通过 schema 验证确保 `alarm_notice_ids` 和 `monitor_notice` 互斥,避免用户错误配置
- 保持向后兼容:现有使用 `alarm_notice_ids` 的配置继续有效

**Non-Goals:**
- 不涉及 `alarm_notice_ids` 内部结构的变更
- 不涉及其他 CLS 资源的修改
- 不支持同时配置两种通知方式(API 限制)

## Decisions

### 1. Schema 设计

**决策**: 将 `alarm_notice_ids` 从 `Required: true` 改为 `Optional: true`,新增 `monitor_notice` 块,使用 `ExactlyOneOf` 实现互斥约束。

**理由**: 
- Terraform SDK v2 的 `ExactlyOneOf` 能在 plan 阶段捕获配置错误,提供更好的用户体验
- 将 `alarm_notice_ids` 改为可选不影响现有配置(已有配置会继续传值)
- `monitor_notice` 结构设计直接映射 SDK 的 `MonitorNotice` 和 `MonitorNoticeRule`

**备选方案**: 使用 `ConflictsWith` 约束。但 `ExactlyOneOf` 强制至少选择一个,更符合 API 语义(两者必选其一)。

### 2. MonitorNotice 嵌套结构

**决策**: 使用嵌套 block 结构:
```hcl
monitor_notice {
  notices {
    notice_id      = "notice-xxx"
    content_tmpl_id = "tmpl-xxx"
    alarm_levels   = [0, 2]
  }
  notices {
    notice_id    = "notice-yyy"
    alarm_levels = [1]
  }
}
```

**理由**:
- 直接映射 SDK 的 `MonitorNotice.Notices` 数组结构
- `content_tmpl_id` 为可选(API 允许为空使用默认模板)
- `alarm_levels` 使用 `TypeList` + `schema.TypeInt`,支持多个告警级别

**备选方案**: 将 `monitor_notice` 设为 `TypeList` + `MaxItems: 1`。但因其只包含一个 `notices` 数组,使用单个 block 更符合 Terraform 习惯。

### 3. CRUD 实现

**决策**: 
- **Create/Update**: 检查 `d.GetOk("monitor_notice")`,如存在则构建 `MonitorNotice` 对象赋值给 `request.MonitorNotice`
- **Read**: 从 API 响应的 `alarm.MonitorNotice` 读取并通过 `d.Set("monitor_notice", ...)` 回写
- **不需要特殊的互斥处理**: schema 层已通过 `ExactlyOneOf` 保证互斥

**理由**: 
- SDK 已支持 `MonitorNotice` 字段,无需修改 vendor
- 读取逻辑需要处理 `MonitorNotice` 可能为 nil 的情况
- 互斥逻辑交给 Terraform 框架处理,避免重复验证

### 4. 向后兼容性

**决策**: 
- 不修改 `alarm_notice_ids` 的类型(`TypeSet`)和元素类型(`TypeString`)
- 只改变 `Required: true` → `Optional: true`
- 添加 `ExactlyOneOf: []string{"monitor_notice"}`

**理由**: 
- 已有 state 中的 `alarm_notice_ids` 值会被正常读取
- 用户配置中已指定 `alarm_notice_ids` 的继续有效
- schema 变更不触发状态迁移

## Risks / Trade-offs

### Risk 1: 用户未指定任何通知方式
**问题**: 如果用户既不配置 `alarm_notice_ids` 也不配置 `monitor_notice`,API 调用可能失败或创建无效告警。

**缓解措施**: 
- 使用 `ExactlyOneOf` 强制至少配置一个
- 在文档中明确说明两者必须选择其一

### Risk 2: API 响应缺少 MonitorNotice 字段导致读取异常
**问题**: 旧版本告警或仅使用 `AlarmNoticeIds` 的告警,API 返回的 `MonitorNotice` 为 nil。

**缓解措施**: 
- Read 逻辑中检查 `alarm.MonitorNotice != nil` 再进行字段赋值
- 如果为 nil 则跳过 `d.Set("monitor_notice", ...)`

### Trade-off: 两种通知方式不可共存
**说明**: API 限制 `AlarmNoticeIds` 和 `MonitorNotice` 互斥,用户无法同时使用两种通知方式。

**影响**: 
- 已使用 `alarm_notice_ids` 的用户如需切换到 `monitor_notice`,需要修改配置(会触发资源更新)
- 文档需要明确说明互斥关系和切换影响
