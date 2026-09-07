## Context

`tencentcloud_cls_alarm_notice` 资源管理 CLS（日志服务）的告警通知渠道组。当前资源已支持 `name`、`type`、`notice_receivers`、`web_callbacks`、`jump_domain`、`deliver_status`、`deliver_config`、`alarm_shield_status`、`callback_prioritize`、`notice_rules`、`tags` 等参数，但尚未支持云 API 已提供的 `SecureDetailStatus` 字段。

CLS 云 API（`cls/v20201016` 包）中，`SecureDetailStatus` 字段定义如下：
- `CreateAlarmNoticeRequest.SecureDetailStatus *uint64`：枚举值 1（关闭，默认）、2（开启），描述为"告警详情安全认证跳转开关，未传时默认关闭"。
- `ModifyAlarmNoticeRequest.SecureDetailStatus *uint64`：同上。
- `AlarmNotice.SecureDetailStatus *uint64`（`DescribeAlarmNotices` 响应结构体）：描述为"告警详情需要安全认证登录开关"。

三个接口均已包含该字段，vendor 中无需任何改动。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_cls_alarm_notice` 资源 schema 中新增 `secure_detail_status` 可选参数
- 在 Create、Read、Update 方法中正确映射 `SecureDetailStatus` 字段
- 保持向后兼容，不影响现有参数和 state
- 补充单元测试用例和文档

**Non-Goals:**
- 不修改 `tags` 参数及其相关逻辑（`tags` 已存在并已通过 tag service 实现）
- 不涉及其他 CLS 资源的修改
- 不修改 vendor 目录下的云 API 定义

## Decisions

### 1. Schema 定义

**决策**: 新增 `secure_detail_status` 字段，类型 `schema.TypeInt`，`Optional: true`，`Computed: true`，映射到云 API 的 `SecureDetailStatus`。

**理由**:
- 云 API 中 `SecureDetailStatus` 为 `*uint64`，未传时服务端默认为 1（关闭），因此 Terraform 侧设为 Optional + Computed，与现有 `deliver_status`、`alarm_shield_status` 等同类型字段保持一致。
- 使用 `schema.TypeInt` 而非 `schema.TypeString`，因为 API 字段为数值枚举，TypeInt 更自然且与同资源中 `deliver_status`、`alarm_shield_status` 风格一致。

**备选方案**: 使用 `schema.TypeBool`。但云 API 使用 uint64 枚举（1/2），保留 TypeInt 更贴合 API 语义，避免布尔与数值之间的转换歧义。

### 2. Create 实现

**决策**: 使用 `d.GetOkExists("secure_detail_status")` 读取参数值，通过 `helper.IntUint64()` 转换后赋值给 `request.SecureDetailStatus`。

**理由**:
- 使用 `GetOkExists` 而非 `GetOk`，因为 `0` 是 `TypeInt` 的零值，而 `SecureDetailStatus` 的有效值从 1 开始；使用 `GetOkExists` 可区分"用户未配置"与"用户配置为 0"的场景。这与现有 `deliver_status`、`alarm_shield_status` 的处理方式一致。
- `helper.IntUint64` 是项目中 int → *uint64 的标准转换函数。

### 3. Read 实现

**决策**: 在 `resourceTencentCloudClsAlarmNoticeRead` 中，先判断 `alarmNotice.SecureDetailStatus != nil`，非 nil 时调用 `d.Set("secure_detail_status", alarmNotice.SecureDetailStatus)` 回写 state。

**理由**:
- 遵循项目规范：set 前先判断 response 字段是否为 nil，避免 nil 指针问题。
- `AlarmNotice.SecureDetailStatus` 类型为 `*uint64`，`d.Set` 可直接接受。

### 4. Update 实现

**决策**: 将 `secure_detail_status` 加入 `mutableArgs` 数组，使 `d.HasChange("secure_detail_status")` 时触发 `needChange`；在 `needChange` 块内使用 `d.GetOkExists("secure_detail_status")` 读取并设置到 `request.SecureDetailStatus`。

**理由**:
- 该字段在云 API 的 `ModifyAlarmNotice` 接口中支持修改，属于可变参数。
- 与现有 `deliver_status`、`alarm_shield_status`、`callback_prioritize` 的处理模式一致。

## Risks / Trade-offs

### Risk 1: 旧版本告警通知可能不返回 SecureDetailStatus 字段
**问题**: 早期创建的告警通知渠道组，API 响应中 `SecureDetailStatus` 可能为 nil。

**缓解措施**: Read 逻辑中检查 `alarmNotice.SecureDetailStatus != nil` 后再 set，nil 时跳过，不会引发错误。

### Trade-off: 字段为 Computed 导致首次 apply 后 state 出现值
**说明**: 由于设为 `Optional + Computed`，用户未配置时，apply 后 state 会写入服务端返回的默认值（通常为 1）。

**影响**: 这是预期行为，与 `deliver_status`、`alarm_shield_status` 一致，不影响向后兼容性。
