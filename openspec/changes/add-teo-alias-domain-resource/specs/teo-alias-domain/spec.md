## ADDED Requirements

### Requirement: 创建别称域名
系统 SHALL 允许用户创建 Teo 别称域名，指定站点 ID、别称域名名称和目标域名。

#### Scenario: 成功创建别称域名
- **WHEN** 用户提供 zone_id、alias_name 和 target_name 参数
- **THEN** 系统调用 CreateAliasDomain 云 API 创建别称域名
- **AND** 系统轮询 DescribeAliasDomains API 直到别称域名创建成功
- **AND** 系统返回资源 ID（格式：zone_id#alias_name）
- **AND** 系统记录资源的 zone_id、alias_name、target_name 和 paused 状态到 Terraform state

#### Scenario: 创建时缺少必需参数
- **WHEN** 用户未提供 zone_id 或 alias_name 或 target_name 参数
- **THEN** 系统返回参数验证错误
- **AND** 资源创建失败

### Requirement: 查询别称域名
系统 SHALL 允许用户查询现有的 Teo 别称域名信息。

#### Scenario: 成功查询别称域名
- **WHEN** 用户通过资源 ID（zone_id#alias_name）查询别称域名
- **THEN** 系统解析 ID 获取 zone_id 和 alias_name
- **AND** 系统调用 DescribeAliasDomains 云 API 查询别称域名信息
- **AND** 系统返回别称域名的 zone_id、alias_name、target_name 和 paused 状态

#### Scenario: 查询不存在的别称域名
- **WHEN** 用户查询的别称域名不存在
- **THEN** 系统返回资源不存在的错误
- **AND** 资源标记为已删除状态

### Requirement: 更新别称域名
系统 SHALL 允许用户更新 Teo 别称域名，包括修改目标域名和暂停/启用状态。

#### Scenario: 成功更新目标域名
- **WHEN** 用户修改 target_name 参数
- **THEN** 系统调用 ModifyAliasDomain 云 API 更新目标域名
- **AND** 系统轮询 DescribeAliasDomains API 直到更新生效
- **AND** 系统更新 Terraform state 中的 target_name

#### Scenario: 成功切换暂停/启用状态
- **WHEN** 用户修改 paused 参数
- **THEN** 系统调用 ModifyAliasDomainStatus 云 API 更新状态
- **AND** 系统轮询 DescribeAliasDomains API 直到状态更新生效
- **AND** 系统更新 Terraform state 中的 paused 状态

#### Scenario: 同时更新多个参数
- **WHEN** 用户同时修改 target_name 和 paused 参数
- **THEN** 系统先调用 ModifyAliasDomain 更新 target_name
- **AND** 系统轮询确认 target_name 更新生效
- **AND** 系统再调用 ModifyAliasDomainStatus 更新 paused 状态
- **AND** 系统轮询确认 paused 状态更新生效
- **AND** 系统更新 Terraform state 中的所有变更

#### Scenario: 未修改任何参数
- **WHEN** 用户未修改任何参数
- **THEN** 系统不调用任何云 API
- **AND** 资源保持不变

### Requirement: 删除别称域名
系统 SHALL 允许用户删除现有的 Teo 别称域名。

#### Scenario: 成功删除别称域名
- **WHEN** 用户请求删除别称域名
- **THEN** 系统调用 DeleteAliasDomain 云 API 删除别称域名
- **AND** 系统轮询 DescribeAliasDomains API 直到别称域名不再存在
- **AND** 系统从 Terraform state 中移除该资源

#### Scenario: 删除已不存在的别称域名
- **WHEN** 用户删除的别称域名已经不存在
- **THEN** 系统返回成功
- **AND** 不执行任何云 API 调用
- **AND** 系统从 Terraform state 中移除该资源

### Requirement: 异步操作轮询
系统 SHALL 对所有异步操作执行轮询，确保操作完全生效后再返回。

#### Scenario: 创建操作轮询成功
- **WHEN** 系统调用 CreateAliasDomain API
- **AND** API 返回成功
- **THEN** 系统开始轮询 DescribeAliasDomains API
- **AND** 系统每 5 秒轮询一次
- **AND** 系统在轮询中检测到别称域名创建成功
- **AND** 系统返回成功

#### Scenario: 更新操作轮询超时
- **WHEN** 系统调用 ModifyAliasDomain API
- **AND** API 返回成功
- **THEN** 系统开始轮询 DescribeAliasDomains API
- **AND** 系统轮询达到默认超时时间（10 分钟）
- **AND** 别称域名仍未更新到预期状态
- **THEN** 系统返回超时错误
- **AND** 提示用户手动检查资源状态

#### Scenario: 删除操作轮询成功
- **WHEN** 系统调用 DeleteAliasDomain API
- **AND** API 返回成功
- **THEN** 系统开始轮询 DescribeAliasDomains API
- **AND** 系统在轮询中检测到别称域名已不存在
- **AND** 系统返回成功

### Requirement: 资源 ID 解析
系统 SHALL 支持从复合资源 ID 中解析出 zone_id 和 alias_name。

#### Scenario: 解析有效的资源 ID
- **WHEN** 系统收到资源 ID "zone-123#example.com"
- **THEN** 系统解析出 zone_id 为 "zone-123"
- **AND** 系统解析出 alias_name 为 "example.com"

#### Scenario: 解析无效的资源 ID
- **WHEN** 系统收到格式错误的资源 ID
- **THEN** 系统返回资源 ID 格式错误的提示

### Requirement: Timeout 配置
系统 SHALL 允许用户自定义异步操作的超时时间。

#### Scenario: 使用默认超时配置
- **WHEN** 用户未配置 Timeouts
- **THEN** 系统使用默认超时时间（Create: 10分钟, Update: 10分钟, Delete: 10分钟）

#### Scenario: 使用自定义超时配置
- **WHEN** 用户配置了 Timeouts.create 为 20 分钟
- **THEN** 系统在创建操作中使用 20 分钟超时时间
- **AND** 其他操作使用默认超时时间

### Requirement: 导入现有资源
系统 SHALL 支持导入已存在的别称域名到 Terraform 管理。

#### Scenario: 成功导入别称域名
- **WHEN** 用户使用 terraform import 命令导入资源
- **AND** 导入 ID 格式为 "zone_id#alias_name"
- **THEN** 系统调用 DescribeAliasDomains API 查询资源信息
- **AND** 系统将资源信息写入 Terraform state
- **AND** 资源成功导入到 Terraform 管理

#### Scenario: 导入不存在的资源
- **WHEN** 用户尝试导入不存在的别称域名
- **THEN** 系统返回资源不存在的错误
- **AND** 导入失败
