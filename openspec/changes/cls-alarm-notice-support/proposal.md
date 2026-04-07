## Why

`tencentcloud_cls_alarm` 资源当前仅支持通过 `alarm_notice_ids` 配置日志服务的告警通知渠道组,但腾讯云 CLS CreateAlarm 接口已支持 `MonitorNotice` 参数用于对接可观测平台的通知模板。用户需要使用可观测平台的统一告警通知能力时,无法通过 Terraform 配置。此外,现有实现将 `alarm_notice_ids` 设置为必填字段,但实际 API 接口中 `AlarmNoticeIds` 和 `MonitorNotice` 是互斥的可选参数,现有设计与 API 行为不一致。

## What Changes

- 将 `alarm_notice_ids` 字段从 `Required` 修改为 `Optional`
- 新增 `monitor_notice` schema 支持配置可观测平台通知模板
  - 支持 `notices` 列表,每个元素包含 `notice_id`、`content_tmpl_id`、`alarm_levels` 字段
- 在 `resourceTencentCloudClsAlarmCreate` 中添加 `monitor_notice` 参数处理逻辑,将其映射到 `CreateAlarmRequest.MonitorNotice`
- 在 `resourceTencentCloudClsAlarmRead` 中添加 `monitor_notice` 读取逻辑
- 在 `resourceTencentCloudClsAlarmUpdate` 中添加 `monitor_notice` 更新逻辑
- 添加 schema 互斥约束:`alarm_notice_ids` 和 `monitor_notice` 不能同时配置(使用 `ExactlyOneOf`)
- 更新资源文档说明参数互斥关系和使用方式

## Capabilities

### New Capabilities
- `monitor-notice-support`: 支持配置可观测平台的通知模板 MonitorNotice 参数,包括通知规则列表(notice_id、content_tmpl_id、alarm_levels)

### Modified Capabilities
<!-- No existing spec-level behavior changes -->

## Impact

- **向后兼容性**: 此变更向后兼容。现有使用 `alarm_notice_ids` 的配置不受影响,只是将其从必填改为可选。
- **互斥验证**: 通过 Terraform schema 的 `ExactlyOneOf` 约束确保用户不会同时配置两个参数,避免 API 调用错误。
- **涉及文件**:
  - `tencentcloud/services/cls/resource_tc_cls_alarm.go` - schema 和 CRUD 逻辑
  - `tencentcloud/services/cls/resource_tc_cls_alarm_test.go` - 添加新参数的测试用例
  - `website/docs/r/cls_alarm.html.markdown` - 文档更新
