## ADDED Requirements

### Requirement: CFS 文件系统容量扩容支持
`tencentcloud_cfs_file_system` 资源 SHALL 支持通过修改 `capacity` 参数对 Turbo 系列文件系统进行在线扩容。

#### Scenario: 用户扩容 Turbo 标准型文件系统
- **GIVEN** 一个已创建的 Turbo 标准型文件系统（storage_type = "TB"），当前容量为 40,960 GiB
- **WHEN** 用户修改 Terraform 配置中的 `capacity` 参数为 61,440 GiB（增加 20,480 GiB）
- **THEN** Provider 应调用 ScaleUpFileSystem API 执行扩容
- **AND** 等待文件系统状态变为 `available`
- **AND** 将新容量同步到 Terraform 状态

#### Scenario: 用户扩容 Turbo 性能型文件系统
- **GIVEN** 一个已创建的 Turbo 性能型文件系统（storage_type = "TP"），当前容量为 20,480 GiB
- **WHEN** 用户修改 `capacity` 参数为 30,720 GiB（增加 10,240 GiB）
- **THEN** Provider 应调用 ScaleUpFileSystem API 执行扩容
- **AND** 等待扩容操作完成
- **AND** Read 操作应返回更新后的容量值

#### Scenario: 拒绝缩容操作
- **GIVEN** 一个容量为 40,960 GiB 的 Turbo 文件系统
- **WHEN** 用户尝试将 `capacity` 修改为 20,480 GiB（缩容）
- **THEN** Provider 应在 Update 函数中返回错误
- **AND** 错误信息应明确说明不支持缩容操作

#### Scenario: 拒绝非 Turbo 文件系统扩容
- **GIVEN** 一个标准型文件系统（storage_type = "SD" 或 "HP"）
- **WHEN** 用户尝试修改 `capacity` 参数
- **THEN** Provider 应返回错误
- **AND** 错误信息应说明仅 Turbo 系列支持扩容

#### Scenario: 验证扩容步长 - Turbo 标准型
- **GIVEN** 一个 Turbo 标准型文件系统（TB），当前容量为 40,960 GiB
- **WHEN** 用户尝试将容量修改为 50,000 GiB（非 20,480 的倍数）
- **THEN** Provider 应返回错误
- **AND** 错误信息应说明扩容步长必须是 20,480 GiB

#### Scenario: 验证扩容步长 - Turbo 性能型
- **GIVEN** 一个 Turbo 性能型文件系统（TP），当前容量为 20,480 GiB
- **WHEN** 用户尝试将容量修改为 25,000 GiB（非 10,240 的倍数）
- **THEN** Provider 应返回错误
- **AND** 错误信息应说明扩容步长必须是 10,240 GiB

### Requirement: 扩容状态等待
Provider SHALL 在调用 ScaleUpFileSystem API 后通过 DescribeCfsFileSystems 接口轮询 LifeCycleState 字段，等待状态变为 `available`。

#### Scenario: 异步扩容操作进行中
- **GIVEN** ScaleUpFileSystem API 调用成功
- **WHEN** 调用 DescribeCfsFileSystems 查询文件系统状态
- **AND** LifeCycleState 字段值为 `expanding`
- **THEN** Provider 应继续等待并重新轮询
- **AND** 使用 Update timeout 作为最大等待时间（默认 30 分钟）

#### Scenario: 异步扩容操作完成
- **GIVEN** ScaleUpFileSystem API 调用成功
- **WHEN** 调用 DescribeCfsFileSystems 查询文件系统状态
- **AND** LifeCycleState 字段值为 `available`
- **THEN** Provider 应停止等待并返回成功
- **AND** Update 操作完成

#### Scenario: 扩容超时
- **GIVEN** ScaleUpFileSystem API 调用成功
- **WHEN** 在 Update timeout 时间内 LifeCycleState 未变为 `available`
- **THEN** Provider 应返回 timeout 错误
- **AND** 错误信息应包含当前状态和已等待时间
- **AND** 建议用户检查云端资源状态

#### Scenario: 扩容过程中遇到异常状态
- **GIVEN** ScaleUpFileSystem API 调用成功
- **WHEN** LifeCycleState 字段值为非预期状态（既非 `expanding` 也非 `available`）
- **THEN** Provider 应返回错误
- **AND** 错误信息应包含异常状态值

### Requirement: Schema 更新
`capacity` 参数 SHALL 从 immutableArgs 列表中移除，以支持更新操作。

#### Scenario: capacity 参数可被修改
- **GIVEN** 资源 Schema 定义
- **WHEN** 检查 `capacity` 字段属性
- **THEN** 该字段应为 `Optional` 和 `Computed`
- **AND** 不应包含 `ForceNew: true`
- **AND** Update 函数应能检测到 `d.HasChange("capacity")`

### Requirement: Update Timeout 支持
资源 SHALL 支持 `update` timeout 配置，用于控制扩容操作的最大等待时间。

#### Scenario: 使用默认 update timeout
- **GIVEN** 用户未在 timeouts 块中指定 update 值
- **WHEN** 执行扩容操作
- **THEN** Provider 应使用默认值 30 分钟作为超时时间

#### Scenario: 使用自定义 update timeout
- **GIVEN** 用户在 timeouts 块中配置 `update = "60m"`
- **WHEN** 执行扩容操作
- **THEN** Provider 应使用 60 分钟作为超时时间

### Requirement: 服务层扩容方法
CfsService SHALL 提供 ScaleUpFileSystem 方法用于执行扩容操作。

#### Scenario: 调用 ScaleUpFileSystem API
- **GIVEN** 文件系统 ID 和目标容量
- **WHEN** 调用 `cfsService.ScaleUpFileSystem(ctx, fsId, targetCapacity)`
- **THEN** 应构造 ScaleUpFileSystemRequest
- **AND** 设置 FileSystemId 和 TargetCapacity 参数
- **AND** 调用 API 并返回结果

#### Scenario: API 调用失败处理
- **GIVEN** ScaleUpFileSystem API 返回错误
- **WHEN** 错误为可重试类型（如网络错误）
- **THEN** 使用 tccommon.RetryError 包装错误
- **AND** 使用 resource.Retry 进行重试

### Requirement: 文档更新
文档 SHALL 明确说明 `capacity` 参数支持修改，但仅限于扩容操作。

#### Scenario: 参数说明更新
- **GIVEN** `cfs_file_system.html.markdown` 文档
- **WHEN** 查看 `capacity` 参数描述
- **THEN** 应包含"支持修改（仅扩容，不支持缩容）"说明
- **AND** 应说明扩容步长要求
- **AND** 应说明仅 Turbo 系列支持扩容

#### Scenario: 添加扩容示例
- **GIVEN** 文档中的示例代码
- **WHEN** 添加扩容使用场景
- **THEN** 应包含 Turbo 文件系统扩容示例
- **AND** 应展示 timeouts 块的 update 配置
- **AND** 应说明扩容的容量计算方式
