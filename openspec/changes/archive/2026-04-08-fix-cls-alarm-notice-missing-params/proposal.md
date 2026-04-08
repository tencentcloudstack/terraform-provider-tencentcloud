# 变更提案：补齐 tencentcloud_cls_alarm_notice 资源缺失的接口参数

## 变更类型

**功能增强** — 针对 `tencentcloud_cls_alarm_notice` 资源，补齐 Create/Update 接口的缺失入参，并在 Read 模块同步新增对应返参的读取。

## Why

当前 `resource_tc_cls_alarm_notice.go` 的 Create（`CreateAlarmNotice`）和 Update（`ModifyAlarmNotice`）模块只接入了 `Name`、`Type`、`NoticeReceivers`、`WebCallbacks` 四个参数，以下接口参数**均未接入**：

| 参数 | SDK 类型 | 说明 |
|------|---------|------|
| `NoticeRules` | `[]*NoticeRule` | 高级模式通知规则（与简易模式互斥） |
| `JumpDomain` | `*string` | 查询数据链接 |
| `DeliverStatus` | `*uint64` | 投递日志开关（1关/2开） |
| `DeliverConfig` | `*DeliverConfig` | 投递日志配置（DeliverStatus=2 时必填） |
| `AlarmShieldStatus` | `*uint64` | 免登录操作告警开关（1关/2开） |
| `CallbackPrioritize` | `*bool` | 统一自定义回调参数开关 |

缺失这些参数导致用户无法通过 Terraform 配置高级通知规则、投递日志、告警屏蔽等功能。

## What Changes

### 新增字段（schema + create/update/read 全量接入）

#### 顶层简单字段

| Terraform 字段名 | SDK 字段 | 类型 | 说明 |
|-----------------|---------|------|------|
| `jump_domain` | `JumpDomain` | `string` | 查询数据链接，http/https 开头，不能以 `/` 结尾 |
| `deliver_status` | `DeliverStatus` | `int`（uint64） | 投递日志开关：1=关闭，2=开启 |
| `alarm_shield_status` | `AlarmShieldStatus` | `int`（uint64） | 免登录操作告警开关：1=关闭，2=开启 |
| `callback_prioritize` | `CallbackPrioritize` | `bool` | 自定义回调参数优先级 |

#### `deliver_config` 嵌套块（对应 `DeliverConfig` 结构体）

| 子字段 | SDK 字段 | 类型 | 说明 |
|--------|---------|------|------|
| `region` | `Region` | `string` | 投递目标地域 |
| `topic_id` | `TopicId` | `string` | 投递目标日志主题 ID |
| `scope` | `Scope` | `int`（uint64） | 投递数据范围：0=全部，1=仅告警触发及恢复 |

#### `notice_rules` 嵌套块（对应 `[]*NoticeRule`，高级模式）

| 子字段 | SDK 字段 | 类型 | 说明 |
|--------|---------|------|------|
| `rule` | `Rule` | `string` | 匹配规则 JSON 串 |
| `notice_receivers` | `NoticeReceivers` | 复用现有 `NoticeReceiver` schema | 规则关联的通知接收者 |
| `web_callbacks` | `WebCallbacks` | 复用现有 `WebCallback` schema | 规则关联的回调 |
| `escalate` | `Escalate` | `bool` | 告警升级开关 |
| `type` | `Type` | `int`（uint64） | 告警升级条件：1=无人认领且未恢复，2=未恢复 |
| `interval` | `Interval` | `int`（uint64） | 告警升级间隔（分钟） |

> `NoticeRule.EscalateNotice`（`*EscalateNoticeInfo`）结构较复杂，本次暂不接入，后续可单独扩展。

### 修改文件

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go` | schema 新增字段；create/update 新增参数填充；read 新增返参读取 |

### 向后兼容性

✅ 完全向后兼容：
- 新增字段均为 `Optional`，不影响已有配置
- `NoticeRules` 与 `Type`/`NoticeReceivers`/`WebCallbacks` 为互斥模式，由用户自行选择（接口侧行为）
