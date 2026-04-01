## ADDED Requirements

### Requirement: TEO 源站组支持配置 HostHeader
tencentcloud_teo_origin_group 资源 SHALL 允许用户配置 `host_header` 字段，用于在源站类型为 HTTP 时指定回源请求的 Host Header。

#### Scenario: 创建源站组时指定 HostHeader
- **WHEN** 用户在 tencentcloud_teo_origin_group 资源中配置了 `host_header` 字段，且源站类型为 HTTP
- **THEN** Provider SHALL 在创建源站组时将 HostHeader 参数传递给 CreateOriginGroup API
- **THEN** Provider SHALL 在读取源站组状态时从 DescribeOriginGroup API 响应中读取并返回 HostHeader 值

#### Scenario: 创建源站组时不指定 HostHeader
- **WHEN** 用户在 tencentcloud_teo_origin_group 资源中未配置 `host_header` 字段
- **THEN** Provider SHALL 在创建源站组时不传递 HostHeader 参数给 CreateOriginGroup API
- **THEN** Provider SHALL 在读取源站组状态时不返回 HostHeader 值或返回空值

#### Scenario: 更新源站组的 HostHeader
- **WHEN** 用户修改 tencentcloud_teo_origin_group 资源的 `host_header` 字段值
- **THEN** Provider SHALL 调用 ModifyOriginGroup API 更新 HostHeader 参数
- **THEN** Provider SHALL 在读取源站组状态时返回更新后的 HostHeader 值

#### Scenario: 移除源站组的 HostHeader
- **WHEN** 用户在 tencentcloud_teo_origin_group 资源中将 `host_header` 字段设置为 null 或从配置中移除
- **THEN** Provider SHALL 在更新时清除 HostHeader 配置
- **THEN** Provider SHALL 在读取源站组状态时不返回 HostHeader 值

#### Scenario: HostHeader 字段属性定义
- **WHEN** 用户查询 tencentcloud_teo_origin_group 资源的 schema 定义
- **THEN** `host_header` 字段 SHALL 定义为 Optional（可选）
- **THEN** `host_header` 字段 SHALL 类型为 String
- **THEN** `host_header` 字段 SHALL 不设置为 Computed

### Requirement: HostHeader 参数仅在 HTTP 类型源站时生效
tencentcloud_teo_origin_group 资源的 HostHeader 配置 SHALL 仅在源站类型为 HTTP 时生效。

#### Scenario: 源站类型为 HTTPS 时配置 HostHeader
- **WHEN** 用户在 tencentcloud_teo_origin_group 资源中配置了 `host_header` 字段，但源站类型为 HTTPS
- **THEN** Provider SHALL 仍然将 HostHeader 参数传递给 API（由服务端处理是否生效）
- **THEN** API 服务端 SHALL 忽略或拒绝该配置（取决于 API 实现）

#### Scenario: 源站类型为 HTTP 时配置 HostHeader
- **WHEN** 用户在 tencentcloud_teo_origin_group 资源中配置了 `host_header` 字段，且源站类型为 HTTP
- **THEN** Provider SHALL 将 HostHeader 参数传递给 API
- **THEN** API 服务端 SHALL 应用该配置到回源请求

### Requirement: Provider 保持向后兼容性
tencentcloud_teo_origin_group 资源的 HostHeader 功能 SHALL 不影响现有用户配置和 state。

#### Scenario: 现有配置不包含 HostHeader 字段
- **WHEN** 用户已有的 tencentcloud_teo_origin_group 资源配置不包含 `host_header` 字段
- **THEN** Provider SHALL 在 refresh 操作中不产生 diff
- **THEN** Provider SHALL 不需要用户手动修改配置

#### Scenario: 现有 state 升级
- **WHEN** 用户升级到包含 HostHeader 功能的 Provider 版本
- **THEN** 现有 state SHALL 自动兼容新的 schema
- **THEN** 不需要手动迁移或修改 state 文件