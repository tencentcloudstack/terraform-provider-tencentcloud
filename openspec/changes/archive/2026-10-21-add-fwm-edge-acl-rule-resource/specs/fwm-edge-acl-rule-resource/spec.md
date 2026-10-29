## ADDED Requirements

### Requirement: 资源管理一条 Edge ACL 规则的完整生命周期
Provider SHALL 提供资源 `tencentcloud_fwm_edge_acl_rule`，通过 `fwm/v20250611` SDK 的 `CreateEdgeAclRule`、`DescribeEdgeAclRules`、`ModifyEdgeAclRule`、`DeleteEdgeAclRule` 四个同步接口管理一条腾讯云云防火墙互联网边界（Edge）ACL 规则。

#### Scenario: Create 同步创建单条规则并绑定 RuleId
- **WHEN** 用户定义资源并提供 `group_id`、`order_index`、`direction`、`source_content`、`source_type`、`target_content`、`target_type`、`protocol`、`port`、`rule_action`、`scope` 等字段
- **THEN** Provider SHALL 构造 `CreateEdgeAclRuleRequest`，`Rules` 数组固定放入 1 个元素，调用接口后取返回 `RuleIds` 数组的首元素作为 `RuleId`，并以 `<GroupId>#<RuleId>` 作为资源 ID 执行 `d.SetId`，随后调用 Read 回填 state；不生成任何异步轮询逻辑

#### Scenario: Read 通过列表查询精准定位并回填
- **WHEN** Provider 执行 Read
- **THEN** 它 SHALL 拆分资源 ID 为 `GroupId` 与 `RuleId` 两段，以 `Direction`（0=出向，1=入向）与 `Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]` 调用 `DescribeEdgeAclRules`（`Limit` 取 1000），取 `Rules` 数组首元素回填所有非空字段

#### Scenario: Read 返回 ResourceNotFound 时清除 state
- **WHEN** `DescribeEdgeAclRules` 返回的 `Rules` 为空（规则已被外部删除）
- **THEN** Provider SHALL 打印 `log.Printf("[CRUD] fwm_edge_acl_rule id=%s", d.Id())` 保留现场，随后 `d.SetId("")` 并返回 nil（导入/读场景下不报错）

#### Scenario: Update 全量覆盖修改可变字段
- **WHEN** 除 ForceNew 字段外的任意可变规则属性变更
- **THEN** Provider SHALL 构造 `ModifyEdgeAclRuleRequest`，`Rule` 字段携带 `RuleId` 与所有可写字段（全量覆盖语义），调用接口后调用 Read；不生成异步轮询

#### Scenario: Delete 单条删除并幂等处理已删除
- **WHEN** 资源被销毁
- **THEN** Provider SHALL 构造 `DeleteEdgeAclRuleRequest`，`RuleIds` 只传 1 个元素调用接口；若返回 `ResourceNotFound` SHALL 视为删除成功（幂等），其他错误 SHALL 经 `tccommon.RetryError` 包装重试

### Requirement: Schema 字段映射 CRUD 接口参数
资源 schema SHALL 包含 `CreateEdgeAclRule`/`ModifyEdgeAclRule` 的 `Rule`/`Rules` 中所有可写参数，以及 `DescribeEdgeAclRules` 返回的全部只读字段。

#### Scenario: 可写字段一一映射
- **WHEN** 用户配置资源
- **THEN** schema SHALL 包含顶层可写字段：`group_id`（Required）、`order_index`（Required，可变）、`direction`（Required+ForceNew）、`source_content`、`source_type`（ForceNew）、`target_content`、`target_type`（ForceNew）、`port`、`protocol`、`rule_action`、`description`、`scope`、`belong_member_id`、`param_template_id`

#### Scenario: 只读字段仅 Describe 回填
- **WHEN** Describe 返回结果
- **THEN** schema SHALL 包含 Computed 只读字段：`rule_id`、`sequence`、`create_time`、`update_time`、`source_name`、`target_name`、`param_template_name`、`country_name`、`city_name`、`instance_name`，且用户不可写

#### Scenario: ForceNew 字段触发重建
- **WHEN** 用户变更 `direction`、`source_type` 或 `target_type`
- **THEN** Terraform SHALL 规划 destroy-and-recreate，因这些字段无法原地更新

#### Scenario: group_id 不可在 Update 中变更
- **WHEN** Update 执行
- **THEN** `group_id` SHALL 仅作为父资源引用传入 `ModifyEdgeAclRuleRequest` 的 `GroupId`，服务端不支持修改 `GroupId`，schema 中 `group_id` 不设 ForceNew 但变更不会下发（属父引用）

### Requirement: 关键字段错位显式映射
Provider SHALL 对入参与出参命名/结构错位的字段做显式映射，避免 import 后 plan 永久 diff。

#### Scenario: OrderIndex ↔ Sequence 映射
- **WHEN** 写入规则时
- **THEN** Provider SHALL 使用入参 `Rule.OrderIndex`；回读时 SHALL 用出参 `Rules[].Sequence` 填充 schema 字段 `order_index`

#### Scenario: BelongMemberId ↔ BelongMember.MemberId 映射
- **WHEN** 回读规则时
- **THEN** Provider SHALL 从出参 `Rules[].BelongMember.MemberId` 取值，仅当其非空才写入 schema 字段 `belong_member_id`；写入时使用入参 `Rule.BelongMemberId`

#### Scenario: Scope 与 RuleAction 直传不二次映射
- **WHEN** 读写规则
- **THEN** `scope`（all/serial/side）与 `rule_action`（accept/drop/log）入参与出参均为同名字符串，Provider SHALL 直接透传，不做整数/"0" 之类二次映射

### Requirement: 复合资源 ID 与导入
资源 ID SHALL 为字符串 `<GroupId>#<RuleId>`，使用 `tccommon.FILED_SP` 作为分隔符，并支持两段式导入。

#### Scenario: 创建成功设置复合 ID
- **WHEN** `CreateEdgeAclRule` 成功且 `RuleIds[0]` 非空
- **THEN** Provider SHALL 调用 `d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))`

#### Scenario: 导入支持两段式 ID
- **WHEN** 用户执行 `terraform import` 传入 `GroupId:RuleId` 或 `GroupId#RuleId`
- **THEN** Provider SHALL 支持两段式 ID 解析并通过 `schema.ImportStatePassthrough` 复用 Read

#### Scenario: 非法 ID 格式报错
- **WHEN** Read/Update/Delete 拆分 `d.Id()` 后段数不等于 2
- **THEN** Provider SHALL 返回 `fmt.Errorf("id is broken,%s", d.Id())`

### Requirement: 所有 SDK 调用包裹重试逻辑与防御性空值处理
所有 SDK API 调用 SHALL 在 `resource.Retry(...)` 内执行，且每个返回值在解引用前做防御性空检查。

#### Scenario: 读取使用 ReadRetryTimeout
- **WHEN** `DescribeEdgeAclRules` 在 Read 或 service helper 中被调用
- **THEN** 该调用 SHALL 包裹在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中并使用 `tccommon.RetryError` 处理可重试错误

#### Scenario: 写操作使用 WriteRetryTimeout
- **WHEN** `CreateEdgeAclRule`、`ModifyEdgeAclRule`、`DeleteEdgeAclRule` 被调用
- **THEN** 该调用 SHALL 包裹在 `resource.Retry(tccommon.WriteRetryTimeout, ...)`

#### Scenario: Create 返回空值不可重试报错
- **WHEN** `CreateEdgeAclRule` 返回 `result == nil` 或 `result.Response == nil` 或 `RuleIds` 为空或 `RuleIds[0]` 为空字符串
- **THEN** Provider SHALL 返回 `resource.NonRetryableError`，避免写入空 id

#### Scenario: 空字段跳过 set
- **WHEN** Read 处理某个为 nil 的响应字段
- **THEN** 它 SHALL 跳过对应的 `d.Set(...)` 调用而非解引用空指针

### Requirement: bool/int 可选字段使用 GetOkExists
对于 bool/int 类型的可选 schema 字段，Provider SHALL 使用 `d.GetOkExists(...)` 取值，避免 false/0 被当作未设置导致状态漂移。

#### Scenario: direction 取值
- **WHEN** 读取 `direction`（int）字段构造请求
- **THEN** Provider SHALL 使用 `d.GetOkExists("direction")` 而非 `d.GetOk`

#### Scenario: order_index 取值
- **WHEN** 读取 `order_index`（int）字段构造请求
- **THEN** Provider SHALL 使用 `d.GetOkExists("order_index")`