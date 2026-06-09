# 任务清单：fix-cls-alarm-notice-missing-params

## 1. 新增顶层简单字段到 schema

**文件**: `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go`

- [x] 在 `tags` 字段之前，新增以下 schema 字段：
  - `jump_domain`（TypeString, Optional）：查询数据链接
  - `deliver_status`（TypeInt, Optional）：投递日志开关，1=关闭，2=开启
  - `alarm_shield_status`（TypeInt, Optional）：免登录操作告警开关，1=关闭，2=开启
  - `callback_prioritize`（TypeBool, Optional）：统一自定义回调参数优先级
- [x] 新增 `deliver_config` 嵌套块（TypeList, MaxItems:1, Optional），包含子字段：
  - `region`（TypeString, Required）
  - `topic_id`（TypeString, Required）
  - `scope`（TypeInt, Optional）
- [x] 新增 `notice_rules` 嵌套块（TypeList, Optional），包含子字段：
  - `rule`（TypeString, Optional）：匹配规则 JSON 串
  - `notice_receivers`（TypeList, Optional）：复用 NoticeReceiver 字段结构
  - `web_callbacks`（TypeList, Optional）：复用 WebCallback 字段结构
  - `escalate`（TypeBool, Optional）：告警升级开关
  - `type`（TypeInt, Optional）：告警升级条件，1=无人认领且未恢复，2=未恢复
  - `interval`（TypeInt, Optional）：告警升级间隔（分钟）

---

## 2. 补齐 Create 模块入参

**文件**: `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go`

- [x] 在 `web_callbacks` 处理代码之后，`resource.Retry` 之前，新增以下参数填充：
  - `jump_domain` → `request.JumpDomain`
  - `deliver_status` → `request.DeliverStatus`（`d.GetOkExists`）
  - `deliver_config` → `request.DeliverConfig`（解析嵌套块，赋值 `Region`、`TopicId`、`Scope`）
  - `alarm_shield_status` → `request.AlarmShieldStatus`（`d.GetOkExists`）
  - `callback_prioritize` → `request.CallbackPrioritize`（`d.GetOkExists`）
  - `notice_rules` → `request.NoticeRules`（解析嵌套块，赋值 `Rule`、`NoticeReceivers`、`WebCallbacks`、`Escalate`、`Type`、`Interval`）
- [x] 执行 `go fmt ./tencentcloud/services/cls/`

---

## 3. 补齐 Update 模块入参

**文件**: `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go`

- [x] 将新增字段加入 `mutableArgs` 列表：`"jump_domain"`, `"deliver_status"`, `"deliver_config"`, `"alarm_shield_status"`, `"callback_prioritize"`, `"notice_rules"`
- [x] 在 update 的 `web_callbacks` 处理代码之后，`resource.Retry` 之前，新增与 create 相同的参数填充逻辑（复用相同代码结构）
- [x] 执行 `go fmt ./tencentcloud/services/cls/`

---

## 4. 补齐 Read 模块返参

**文件**: `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go`

- [x] 在 `web_callbacks` 读取代码之后，`tags` 读取之前，新增以下返参读取：
  - `alarmNotice.JumpDomain` → `d.Set("jump_domain", ...)`
  - `alarmNotice.DeliverStatus` → `d.Set("deliver_status", ...)`
  - `alarmNotice.AlarmShieldStatus` → `d.Set("alarm_shield_status", ...)`
  - `alarmNotice.CallbackPrioritize` → `d.Set("callback_prioritize", ...)`
  - `alarmNotice.AlarmNoticeDeliverConfig.DeliverConfig` → `d.Set("deliver_config", ...)`（构建 map 后 set）
  - `alarmNotice.NoticeRules` → `d.Set("notice_rules", ...)`（遍历构建 list）
- [x] 执行 `go fmt ./tencentcloud/services/cls/`

---

## 5. 编译验证

- [x] `go build ./tencentcloud/services/cls/` 确认编译通过

---

## 6. 补齐 notice_rules 中缺失的 EscalateNotice 字段（post-apply 补充）

**背景**: 经过字段对齐分析，发现 `NoticeRule.EscalateNotice`（`*EscalateNoticeInfo`，最多 5 层递归）未接入。采用方案 D（扁平有序列表 → SDK 链式结构）实现。

**文件**: `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go`

- [x] schema: `notice_rules` 中新增 `escalate_notices`（TypeList, MaxItems:5）嵌套块，子字段包含 `notice_receivers`、`web_callbacks`、`escalate`、`type`、`interval`
- [x] Create 模块: 调用 `buildEscalateNoticeChain(escalateList)` 将有序列表组装为链式 `*EscalateNoticeInfo`
- [x] Update 模块: 同 Create 逻辑
- [x] Read 模块: 调用 `flattenEscalateNoticeChain(info)` 将链式结构展开为有序列表写入 state
- [x] 新增辅助函数 `buildEscalateNoticeChain` 和 `flattenEscalateNoticeChain`
- [x] 所有注释改为英文
- [x] `go fmt` + `go build` 编译通过

---

## 总结

- **预计工作量**：中等（约 2 小时）
- **风险等级**：低（纯增量，不修改已有字段逻辑）
- **破坏性变更**：无
- **状态**: 已完成
