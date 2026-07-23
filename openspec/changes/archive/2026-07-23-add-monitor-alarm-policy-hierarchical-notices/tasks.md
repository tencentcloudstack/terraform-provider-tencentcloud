## 1. Schema 定义

- [x] 1.1 在 resource_tc_monitor_alarm_policy.go 的 Resource schema 中新增 `hierarchical_notices` 参数（TypeList, Optional, Computed），元素包含 `notice_id`（Required, TypeString）和 `classification`（Optional, TypeSet of TypeString）
- [x] 1.2 在 resource_tc_monitor_alarm_policy.go 的 Resource schema 中新增 `notice_content_tmpl_bind_infos` 参数（TypeList, Optional, Computed），元素包含 `content_tmpl_id`（Required, TypeString）和 `notice_id`（Required, TypeString）

## 2. Create 函数修改

- [x] 2.1 在 resourceTencentMonitorAlarmPolicyCreate 函数中，从 `hierarchical_notices` 参数读取数据并构建 `[]*monitor.AlarmHierarchicalNotice` 传递给 `CreateAlarmPolicy` API
- [x] 2.2 在 resourceTencentMonitorAlarmPolicyCreate 函数中，从 `notice_content_tmpl_bind_infos` 参数读取数据并构建 `[]*monitor.NoticeContentTmplBindInfo` 传递给 `CreateAlarmPolicy` API

## 3. Read 函数修改

- [x] 3.1 在 resourceTencentMonitorAlarmPolicyRead 函数中，从 `DescribeAlarmPolicy` 响应读取 `HierarchicalNotices` 并设置到 `d.Set("hierarchical_notices", ...)` ，注意 nil 检查
- [x] 3.2 在 resourceTencentMonitorAlarmPolicyRead 函数中，从 `DescribeAlarmPolicy` 响应读取 `NoticeContentTmplBindInfos` 并设置到 `d.Set("notice_content_tmpl_bind_infos", ...)` ，注意 nil 检查

## 4. Update 函数修改

- [x] 4.1 在 resourceTencentMonitorAlarmPolicyUpdate 函数中，当 `hierarchical_notices` 或 `notice_content_tmpl_bind_infos` 发生变化时，通过 `ModifyAlarmPolicyNotice` API 更新这两个参数

## 5. 测试

- [x] 5.1 在 resource_tc_monitor_alarm_policy_test.go 中补充单元测试用例，覆盖新增参数的 CRUD 场景

## 6. 文档

- [x] 6.1 更新 resource_tc_monitor_alarm_policy.md 文件，添加新增参数的使用示例