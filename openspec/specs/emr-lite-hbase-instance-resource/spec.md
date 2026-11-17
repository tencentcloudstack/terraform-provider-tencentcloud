# emr-lite-hbase-instance-resource Specification

## Purpose
TBD - created by archiving change support-lite-hbase-disk-size-update. Update Purpose after archive.
## Requirements
### Requirement: lite_hbase_instance 支持 disk_size 字段扩容更新
`tencentcloud_lite_hbase_instance` 资源 SHALL 支持通过修改 `disk_size` 参数对轻量 HBase 实例进行在线磁盘扩容。调用 `ModifySLInstance` 接口执行变配，并通过 `DescribeSLInstance` 接口轮询实例状态等待扩容完成。

#### Scenario: 用户扩容 lite_hbase_instance 磁盘容量
- **GIVEN** 一个已创建的 lite_hbase_instance 实例，当前 disk_size 为 100 GB
- **WHEN** 用户修改 Terraform 配置中的 `disk_size` 参数为 200 GB
- **THEN** Provider 应调用 ModifySLInstance API 执行扩容
- **AND** 等待实例状态从 22（变配中）变为 2（运行中/变配完成）
- **AND** 将新容量同步到 Terraform 状态

#### Scenario: disk_size 未变更时不触发扩容
- **GIVEN** 一个已创建的 lite_hbase_instance 实例
- **WHEN** 用户未修改 `disk_size` 参数
- **THEN** Provider 不应调用 ModifySLInstance API
- **AND** 不应触发异步状态等待

### Requirement: disk_size 从不可变参数列表中移除
`disk_size` 参数 SHALL 从 Update 函数的 `immutableArgs` 列表中移除，使修改该字段不再返回 `argument cannot be changed` 错误。

#### Scenario: disk_size 参数可被修改
- **GIVEN** 资源 Update 函数的 immutableArgs 列表
- **WHEN** 检查 `disk_size` 是否在列表中
- **THEN** `disk_size` 不应在 immutableArgs 列表中
- **AND** Update 函数应能检测到 `d.HasChange("disk_size")`

### Requirement: disk_size 扩容异步状态等待
Provider SHALL 在调用 ModifySLInstance API 后通过现有的 `SLInstanceStateRefreshFunc` 轮询实例状态，等待状态值变为 `2`（运行中/变配完成）。

#### Scenario: 异步扩容操作进行中
- **GIVEN** ModifySLInstance API 调用成功
- **WHEN** 通过 DescribeSLInstance 查询实例状态
- **AND** Status 字段值为 22（变配中/ADJUSTING）
- **THEN** Provider 应继续等待并重新轮询

#### Scenario: 异步扩容操作完成
- **GIVEN** ModifySLInstance API 调用成功
- **WHEN** 通过 DescribeSLInstance 查询实例状态
- **AND** Status 字段值为 2（运行中/RUNNING）
- **THEN** Provider 应停止等待并返回成功
- **AND** Update 操作完成

### Requirement: disk_size 扩容 API 调用参数构造
Provider SHALL 在检测到 `disk_size` 变更时构造 `ModifySLInstanceRequest`，设置 `InstanceId`（取自 `d.Id()`）、`Zone`（取自 `zone_settings` 第一个元素的 `zone`，接口必传）与目标 `DiskSize`（取自 schema 新值），并通过 `ModifySLInstanceWithContext` 调用接口，调用过程使用 `resource.Retry` + `tccommon.WriteRetryTimeout` 进行重试。若取不到 `zone_settings.0.zone`，SHALL 在调用接口前返回错误。

#### Scenario: 构造 ModifySLInstanceRequest 并调用接口
- **GIVEN** 检测到 `d.HasChange("disk_size")` 为 true
- **WHEN** Provider 执行 disk_size 变更处理
- **THEN** 应构造 ModifySLInstanceRequest
- **AND** 设置 InstanceId 为 d.Id() 的值
- **AND** 设置 Zone 为 `zone_settings` 第一个元素的 `zone`
- **AND** 设置 DiskSize 为 schema 中的新值
- **AND** 使用 resource.Retry 包装调用，超时使用 tccommon.WriteRetryTimeout

#### Scenario: 取不到 zone 时不调用接口
- **GIVEN** 检测到 `d.HasChange("disk_size")` 为 true
- **WHEN** `zone_settings` 为空或其第一个元素没有 `zone`
- **THEN** SHALL 在调用 ModifySLInstance 之前返回错误

#### Scenario: API 调用失败处理
- **GIVEN** ModifySLInstance API 返回错误
- **WHEN** 错误为可重试类型（如网络错误、限流）
- **THEN** 使用 tccommon.RetryError 包装错误
- **AND** resource.Retry 继续重试

### Requirement: disk_size 与 zone_settings.node_num 变更共存处理
当 `disk_size` 与 `zone_settings` 中的 `node_num` 同时变更时，Provider SHALL 先处理 `node_num` 变更（现有逻辑），再处理 `disk_size` 变更，两次变配均需独立等待状态恢复为 2。

#### Scenario: disk_size 与 node_num 同时变更
- **GIVEN** 用户同时修改 disk_size 和 zone_settings.node_num
- **WHEN** Provider 执行 Update
- **THEN** 应先处理 node_num 变更并等待状态恢复为 2
- **AND** 再处理 disk_size 变更并等待状态恢复为 2

### Requirement: 文档更新
资源 `.md` 示例文档 SHALL 同步更新，补充 `disk_size` 支持扩容（仅扩容）的使用说明。

#### Scenario: 示例文档补充扩容说明
- **GIVEN** `tencentcloud/services/emr/resource_tc_lite_hbase_instance.md`
- **WHEN** 更新文档
- **THEN** 应体现 disk_size 支持修改（扩容）的能力

