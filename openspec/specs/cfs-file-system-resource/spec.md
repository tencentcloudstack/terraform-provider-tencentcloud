# CFS 文件系统资源规格

## Requirement: CFS 文件系统加密参数支持
`tencentcloud_cfs_file_system` 资源 SHALL 新增 `encrypted` 参数，用于在创建文件系统时设置是否加密，并在读取时同步加密状态到 state。

### Scenario: 用户创建加密文件系统
- **GIVEN** 用户在 Terraform 配置中设置 `encrypted = true`
- **WHEN** 执行 `terraform apply` 创建 CFS 文件系统
- **THEN** Provider 应将 `Encrypted` 参数设置为 `true` 并传入 `CreateCfsFileSystem` 接口
- **AND** 创建完成后 Read 操作应从 `DescribeCfsFileSystems` 返回的 `FileSystemInfo.Encrypted` 同步加密状态到 state

### Scenario: 用户创建不加密文件系统
- **GIVEN** 用户在 Terraform 配置中不设置或设置 `encrypted = false`
- **WHEN** 执行 `terraform apply` 创建 CFS 文件系统
- **THEN** Provider 应将 `Encrypted` 参数设置为 `false` 传入 `CreateCfsFileSystem` 接口
- **AND** 文件系统创建为不加密状态

### Scenario: 修改 encrypted 触发资源重建
- **GIVEN** 一个已创建的 CFS 文件系统
- **WHEN** 用户修改 Terraform 配置中的 `encrypted` 参数
- **THEN** Provider 应触发资源重建（ForceNew）
- **AND** 因为云 API 不支持更新加密状态

### Scenario: Read 操作同步加密状态
- **GIVEN** 一个已创建的 CFS 文件系统
- **WHEN** 执行 `terraform refresh` 或 Read 操作
- **THEN** Provider 应调用 `DescribeCfsFileSystems` 获取文件系统详情
- **AND** 在 `FileSystemInfo.Encrypted` 非 nil 时将其同步到 state 的 `encrypted` 字段

## Requirement: encrypted 参数 Schema 定义
`encrypted` 参数 SHALL 定义为 `TypeBool`、`Optional`、`ForceNew: true`。

### Scenario: Schema 字段属性
- **GIVEN** 资源 Schema 定义
- **WHEN** 检查 `encrypted` 字段属性
- **THEN** 该字段类型应为 `schema.TypeBool`
- **AND** 应为 `Optional`
- **AND** 应包含 `ForceNew: true`