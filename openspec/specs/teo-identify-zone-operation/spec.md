# teo-identify-zone-operation Specification

## Purpose
TBD - created by archiving change add-teo-identify-zone-operation. Update Purpose after archive.
## Requirements
### Requirement: 获取站点认证配置信息
当用户提供站点名称时，系统 SHALL 调用 IdentifyZone 云 API 获取该站点的认证配置信息，包括 DNS 校验信息和文件校验信息。

#### Scenario: 成功获取站点认证配置
- **WHEN** 用户提供有效的站点名称（zone_name）
- **THEN** 系统返回 DNS 校验信息（ascription），包含 Subdomain、RecordType、RecordValue
- **THEN** 系统返回文件校验信息（file_ascription），包含 IdentifyPath、IdentifyContent

### Requirement: 获取子域名认证配置信息
当用户提供站点名称和子域名时，系统 SHALL 调用 IdentifyZone 云 API 获取该子域名的认证配置信息。

#### Scenario: 成功获取子域名认证配置
- **WHEN** 用户提供有效的站点名称（zone_name）和子域名（domain）
- **THEN** 系统返回该子域名的 DNS 校验信息和文件校验信息

### Requirement: 处理缺少必需参数的错误
当用户缺少必需参数时，系统 SHALL 返回明确的错误信息。

#### Scenario: 缺少 zone_name 参数
- **WHEN** 用户未提供 zone_name 参数
- **THEN** 系统返回错误信息，提示 zone_name 是必需参数

### Requirement: 处理云 API 调用错误
当云 API 调用失败时，系统 SHALL 返回详细的错误信息。

#### Scenario: 站点不存在
- **WHEN** 提供的站点名称不存在
- **THEN** 系统返回错误信息，说明站点不存在或无效

#### Scenario: 权限不足
- **WHEN** 当前账号对站点没有操作权限
- **THEN** 系统返回权限不足的错误信息

### Requirement: 验证参数格式
系统 SHALL 验证传入参数的格式符合云 API 的要求。

#### Scenario: zone_name 格式验证
- **WHEN** 用户提供格式不正确的 zone_name
- **THEN** 系统返回参数格式错误的提示

#### Scenario: domain 格式验证
- **WHEN** 用户提供格式不正确的 domain
- **THEN** 系统返回参数格式错误的提示

### Requirement: 返回认证信息的结构
系统 SHALL 以结构化的方式返回认证信息，便于用户理解和使用。

#### Scenario: DNS 校验信息结构
- **WHEN** 系统返回 ascription 信息
- **THEN** ascription 包含三个字段：Subdomain（主机记录）、RecordType（记录类型）、RecordValue（记录值）

#### Scenario: 文件校验信息结构
- **WHEN** 系统返回 file_ascription 信息
- **THEN** file_ascription 包含两个字段：IdentifyPath（文件校验目录）、IdentifyContent（文件校验内容）

### Requirement: 一次性操作特性
由于站点认证是一次性操作，系统 SHALL 不提供读取、更新、删除接口。

#### Scenario: 仅支持创建操作
- **WHEN** 用户尝试使用 Terraform 的读取、更新、删除操作
- **THEN** 系统提示该资源仅支持创建操作

