## Context

腾讯云监控告警策略 API（v20180724）的 `CreateAlarmPolicy` 和 `ModifyAlarmPolicyNotice` 接口已支持 `HierarchicalNotices`（告警分级通知规则配置）和 `NoticeContentTmplBindInfos`（通知内容模板绑定信息）参数。当前 Terraform 资源 `tencentcloud_monitor_alarm_policy` 尚未暴露这两个参数。

`AlarmHierarchicalNotice` 结构：
- `NoticeId`（*string）：通知模板 ID
- `Classification`（[]*string）：通知等级列表，如 `["Remind","Serious"]`

`NoticeContentTmplBindInfo` 结构：
- `ContentTmplID`（*string）：通知内容模板 ID
- `NoticeID`（*string）：通知模板 ID

两个参数在 Create 和 Update（ModifyAlarmPolicyNotice）中都受支持，在 DescribeAlarmPolicy 的响应 `AlarmPolicy` 中也返回了这两个字段。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_monitor_alarm_policy` 的 schema 中新增 `hierarchical_notices`（TypeList）参数
- 在 `tencentcloud_monitor_alarm_policy` 的 schema 中新增 `notice_content_tmpl_bind_infos`（TypeList）参数
- 在 Create 方法中将这两个参数传递给 `CreateAlarmPolicy` API
- 在 Read 方法中从 `DescribeAlarmPolicy` 响应中读取这两个参数
- 在 Update 方法中通过 `ModifyAlarmPolicyNotice` API 更新这两个参数

**Non-Goals:**
- 不修改现有 schema 字段的行为
- 不改变 resource ID 格式
- 不添加新的 API 调用

## Decisions

1. **新增参数均为 Optional + Computed**：与现有 `notice_ids` 参数风格一致，保持向后兼容。

2. **Update 逻辑**：当 `hierarchical_notices` 或 `notice_content_tmpl_bind_infos` 发生变化时，通过 `ModifyAlarmPolicyNotice` API 更新。由于 `ModifyAlarmPolicyNotice` 同时支持 `HierarchicalNotices` 和 `NoticeContentTmplBindInfos` 参数，可以直接使用。

3. **Schema 结构**：
   - `hierarchical_notices`：TypeList，元素为 Resource，包含 `notice_id`（TypeString, Required）和 `classification`（TypeSet of TypeString, Optional）
   - `notice_content_tmpl_bind_infos`：TypeList，元素为 Resource，包含 `content_tmpl_id`（TypeString, Required）和 `notice_id`（TypeString, Required）

4. **Read 时的 nil 检查**：在读取响应中 `HierarchicalNotices` 和 `NoticeContentTmplBindInfos` 字段时，需要先检查是否为 nil，遵循项目规范。

## Risks / Trade-offs

- 新增参数均为 Optional，不影响现有用户配置
- 如果 API 返回的字段为 nil，Read 方法中不设置对应字段，保持 state 一致性