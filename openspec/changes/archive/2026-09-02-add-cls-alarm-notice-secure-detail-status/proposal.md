## Why

`tencentcloud_cls_alarm_notice` 资源当前不支持配置"告警详情安全认证跳转开关"。CLS 云 API 的 `CreateAlarmNotice`、`ModifyAlarmNotice` 接口已支持 `SecureDetailStatus` 参数（枚举值：1 关闭默认、2 开启），`DescribeAlarmNotices` 响应中也已返回该字段。用户需要通过 Terraform 控制告警详情的安全认证跳转行为时，现有 provider 无法满足，必须手动在控制台操作，造成配置漂移。

## What Changes

- 为 `tencentcloud_cls_alarm_notice` 资源新增可选参数 `secure_detail_status`（类型 `schema.TypeInt`），枚举值 1（关闭，默认）和 2（开启），映射到云 API 的 `SecureDetailStatus` 字段
- 在 `resourceTencentCloudClsAlarmNoticeCreate` 中读取 `secure_detail_status` 并设置到 `CreateAlarmNoticeRequest.SecureDetailStatus`
- 在 `resourceTencentCloudClsAlarmNoticeRead` 中从 `DescribeAlarmNotices` 响应的 `AlarmNotice.SecureDetailStatus` 读取并回写到 Terraform state
- 在 `resourceTencentCloudClsAlarmNoticeUpdate` 中将 `secure_detail_status` 加入 `mutableArgs`，变更时设置到 `ModifyAlarmNoticeRequest.SecureDetailStatus`
- 更新资源文档 `resource_tc_cls_alarm_notice.md` 增加 `secure_detail_status` 参数说明与示例

## Capabilities

### New Capabilities
- `alarm-notice-secure-detail-status`: 支持通过 `secure_detail_status` 参数配置 CLS 告警通知渠道组的告警详情安全认证跳转开关，覆盖 Create、Read、Update 全生命周期

### Modified Capabilities
<!-- No existing spec-level behavior changes -->

## Impact

- **向后兼容性**: 此变更向后兼容。`secure_detail_status` 为新增的可选参数，不影响现有配置和 state。
- **涉及文件**:
  - `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go` - schema 定义和 CRUD 逻辑
  - `tencentcloud/services/cls/resource_tc_cls_alarm_notice_test.go` - 补充单元测试用例
  - `tencentcloud/services/cls/resource_tc_cls_alarm_notice.md` - 文档更新
- **云 API 依赖**: 依赖 `cls/v20201016` 包中 `CreateAlarmNoticeRequest`、`ModifyAlarmNoticeRequest` 的 `SecureDetailStatus` 字段，以及 `AlarmNotice`（`DescribeAlarmNotices` 响应）的 `SecureDetailStatus` 字段，均已存在于 vendor 中。
