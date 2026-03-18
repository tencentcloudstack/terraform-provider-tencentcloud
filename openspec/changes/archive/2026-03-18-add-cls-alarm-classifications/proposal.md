# Proposal: 为 CLS Alarm 资源添加 Classifications 参数支持

## Why

腾讯云 CLS 告警 API 已支持 `Classifications` 参数用于告警分类管理,但当前 `tencentcloud_cls_alarm` 资源尚未暴露此字段。用户无法通过 Terraform 配置告警的分类信息,需要补充此能力以支持完整的告警管理场景。

## What Changes

- 在 `tencentcloud_cls_alarm` 资源的 schema 中新增 `classifications` 字段(可选,支持修改)
- 在 Create 逻辑中支持设置 `Classifications` 参数(调用 `CreateAlarm` API)
- 在 Update 逻辑中支持修改 `Classifications` 参数(调用 `ModifyAlarm` API)
- 在 Read 逻辑中读取 `Classifications` 字段(从 `DescribeAlarms` API 响应中获取)
- 在 Delete 逻辑中无需特殊处理(调用 `DeleteAlarm` API)

## Capabilities

### New Capabilities
- `cls-alarm-classifications`: 为 CLS 告警资源添加分类参数的创建、读取和更新支持

### Modified Capabilities
<!-- 此变更不涉及现有 spec 的 requirement 变更,仅是新增 Optional 字段 -->

## Impact

**受影响的文件**:
- `tencentcloud/services/cls/resource_tc_cls_alarm.go` - 资源定义和 CRUD 逻辑
- `tencentcloud/services/cls/resource_tc_cls_alarm_test.go` - 验收测试
- `website/docs/r/cls_alarm.html.markdown` - 用户文档

**API 调用**:
- `CreateAlarm` - 新增 `Classifications` 参数
- `ModifyAlarm` - 新增 `Classifications` 参数
- `DescribeAlarms` - 读取响应中的 `Classifications` 字段
- `DeleteAlarm` - 无变更

**兼容性**:
- ✅ 向后兼容:新字段为 Optional,不影响现有配置
- ✅ 不破坏 state:字段可计算,升级后自动同步
- ✅ 无需数据迁移
