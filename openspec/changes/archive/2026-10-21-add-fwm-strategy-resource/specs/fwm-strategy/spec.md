## ADDED Requirements

### Requirement: FWM 策略资源生命周期管理

系统 SHALL 提供通用资源 `tencentcloud_fwm_strategy`，通过 FWM 云 API（v20250611）管理单条下发规则策略的完整生命周期：创建、读取、更新、删除、导入。

#### Scenario: 创建单条前区策略
- **WHEN** 用户声明 `tencentcloud_fwm_strategy` 资源且配置 `product`、`sub_account_uin`、`exec_area=pre`、`group_id`、`sequence`
- **THEN** provider 调用 `CreateStrategy`，向 `PreStrategy` 传入 1 个 `{GroupId,Sequence}` 元素，取响应 `StrategyIds[0]` 作为资源 ID 写入 state

#### Scenario: 创建单条后区策略
- **WHEN** 用户配置 `exec_area=post`、`group_id`、`sequence`
- **THEN** provider 调用 `CreateStrategy`，向 `PostStrategy` 传入 1 个元素，取 `StrategyIds[0]` 作为资源 ID

#### Scenario: 拒绝多元素批量创建
- **WHEN** 用户在 `pre_strategy` 或 `post_strategy` 配置多个元素
- **THEN** provider 返回 error，提示用户拆分为多个 `tencentcloud_fwm_strategy` 资源实例，不静默只取首个

#### Scenario: 创建返回空 ID 时报错
- **WHEN** `CreateStrategy` 响应的 `StrategyIds` 为 nil、长度为 0 或首个元素为空字符串
- **THEN** provider 返回 `NonRetryableError`，不写入空 ID

### Requirement: 单实例精确查询读取

系统 SHALL 优先使用 `DescribeStrategy`（按 `StrategyId` 精确单查）作为 Read 接口。

#### Scenario: 读取存在的策略
- **WHEN** provider 执行 Read 且策略存在
- **THEN** provider 调用 `DescribeStrategy` 获取详情并回填 state 中可读字段

#### Scenario: 策略不存在时移除 state
- **WHEN** `DescribeStrategy` 返回 `ResourceNotFound` 或 `response.Response.Strategy` 为 nil
- **THEN** provider 打印 `log.Printf("[CRUD] fwm_strategy id=%s", d.Id())` 后调用 `d.SetId("")` 移除 state，不返回 error

### Requirement: 不可回读字段用 state 保留

由于 `DescribeStrategy` 不回传 `Product` 与 `ExecArea`，系统 SHALL 将 `product`、`sub_account_uin`、`exec_area` 标为 `Required` + `ForceNew`，且 Read 时不对它们调用 `d.Set()`，保留 state 中用户配置值。

#### Scenario: product 不可回读
- **WHEN** provider 执行 Read
- **THEN** 不对 `product` 调用 `d.Set()`，state 保留用户配置值；修改 `product` 触发重建

#### Scenario: exec_area 不可回读
- **WHEN** provider 执行 Read
- **THEN** 不对 `exec_area` 调用 `d.Set()`，state 保留用户配置值；修改 `exec_area` 触发重建

#### Scenario: sub_account_uin 不可回读
- **WHEN** provider 执行 Read
- **THEN** 不对 `sub_account_uin` 调用 `d.Set()`，state 保留用户配置值；修改 `sub_account_uin` 触发重建

### Requirement: 增量语义更新

系统 SHALL 通过 `ModifyStrategy` 执行更新，采用增量语义：只把发生变更的字段写入请求，未变更字段保持 API 指针 nil（不改）。

#### Scenario: 仅更新 group_id
- **WHEN** 用户修改 `group_id` 但不改 `sequence`、`receive_account`、`receive_group`
- **THEN** provider 仅在 `ModifyStrategy` 请求中设置 `GroupId`，`Sequence`/`ReceiveAccount`/`ReceiveGroup` 保持 nil

#### Scenario: 更新多个字段
- **WHEN** 用户同时修改多个可变字段
- **THEN** provider 只把变更字段写入 `ModifyStrategy` 请求

### Requirement: 接收对象不对称回读

入参 `receive_account` 是 UIN 字符串数组、`receive_group` 是账号组 ID 数组；响应 `ReceiveAccount` 是对象数组（含 `ReceiverType`）。系统 SHALL 按 `ReceiverType` 拆分回填。

#### Scenario: 回读账号类型接收对象
- **WHEN** 响应 `ReceiveAccount` 元素 `ReceiverType == 0`
- **THEN** provider 将该元素 `Uin` 收集到 `receive_account` 集合并回填 state

#### Scenario: 回读账号组类型接收对象
- **WHEN** 响应 `ReceiveAccount` 元素 `ReceiverType == 1`
- **THEN** provider 将该元素 `Uin` 收集到 `receive_group` 集合并回填 state

#### Scenario: 不比对 Nickname
- **WHEN** 回读接收对象集合
- **THEN** 只比对 `Uin` 集合，不比对 `Nickname`

### Requirement: 幂等删除

系统 SHALL 通过 `DeleteStrategy` 删除策略，传入 `StrategyIds` 数组（单个删除传 1 个元素）。不存在的 ID 静默成功。

#### Scenario: 删除存在的策略
- **WHEN** provider 执行 Delete 且策略存在
- **THEN** provider 调用 `DeleteStrategy` 传入单元素 `StrategyIds` 数组

#### Scenario: 删除不存在的策略幂等
- **WHEN** provider 执行 Delete 但策略已不存在
- **THEN** provider 不报错（服务端静默成功），destroy 天然幂等

### Requirement: 同步接口无轮询

四个接口（Create/Describe/Modify/Delete）均为同步接口，系统 SHALL NOT 生成 waitReady/async 轮询逻辑。

#### Scenario: 创建后立即可读
- **WHEN** Create 成功返回
- **THEN** provider 直接执行 Read 回填 state，不轮询等待

### Requirement: 导入支持

系统 SHALL 支持 import，导入时直接使用 `StrategyId` 作为资源 ID。

#### Scenario: 导入现有策略
- **WHEN** 用户执行 `terraform import tencentcloud_fwm_strategy.example fwms-xxxx`
- **THEN** provider 将 `fwms-xxxx` 作为 `StrategyId`，调用 `DescribeStrategy` 读取并回填 state

### Requirement: 只读字段

系统 SHALL 将服务端返回的只读字段标为 `Computed`：`group_name`、`rule_count`、`rule_status`、`err_msg`、`error_type`、`create_time`、`update_time`、`create_by`、`create_name`、`update_by`、`update_name`、`strategy_id`。

#### Scenario: 只读字段不可由用户设置
- **WHEN** 用户在 HCL 中未设置这些字段
- **THEN** Read 时从 `DescribeStrategy` 响应回填这些字段到 state