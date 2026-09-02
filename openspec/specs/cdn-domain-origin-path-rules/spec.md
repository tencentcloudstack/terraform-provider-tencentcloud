# cdn-domain-origin-path-rules Specification

## Purpose
TBD - created by archiving change add-cdn-domain-origin-path-rules. Update Purpose after archive.
## Requirements
### Requirement: CDN 域名 origin 支持 path_rules 回源路径重写规则
`tencentcloud_cdn_domain` 资源的 `origin` 配置块 SHALL 支持 `path_rules` 子列表参数，用于配置回源路径重写规则。`path_rules` 列表中的每个元素 SHALL 包含 `regex`（通配符匹配开关，bool）、`path`（匹配的 URL 路径，string）、`server_name`（路径匹配时回源的 Host 头部，string）、`forward_uri`（路径匹配时回源的 URI 路径，string）四个字段，所有字段均为可选。

#### Scenario: Create 时填充 path_rules
- **WHEN** 用户在 `origin` 配置块中配置了 `path_rules` 并执行 `terraform apply` 创建 CDN 域名
- **THEN** Create 方法（`AddCdnDomain`）的 `request.Origin.PathRules` 中 SHALL 包含用户配置的 `regex`、`path`、`server_name`、`forward_uri` 字段值

#### Scenario: Read 时从 API 响应写入 path_rules
- **WHEN** 执行 `terraform plan` 或 `terraform refresh` 读取 CDN 域名配置
- **THEN** Read 方法（`DescribeDomainsConfig`）SHALL 将 `response.Domains.Origin.PathRules` 中的 `Regex`、`Path`、`ServerName`、`ForwardUri` 字段写入 state 的 `origin.path_rules` 中（若 API 返回非 nil）

#### Scenario: Update 时填充 path_rules
- **WHEN** 用户修改 `origin` 配置块中的 `path_rules` 并执行 `terraform apply` 更新 CDN 域名
- **THEN** Update 方法（`UpdateDomainConfig`）的 `request.Origin.PathRules` 中 SHALL 包含用户配置的 `regex`、`path`、`server_name`、`forward_uri` 字段值

#### Scenario: 未配置 path_rules 时不影响现有行为
- **WHEN** 用户未配置 `origin.path_rules`
- **THEN** Create/Update 方法 SHALL 不设置 `request.Origin.PathRules`（或设置为空），现有行为不受影响

#### Scenario: path_rules 字段均为可选
- **WHEN** 用户配置 `path_rules` 时仅填写部分字段（如仅 `path` 和 `regex`）
- **THEN** 未填写的字段 SHALL 不被填充到 API 请求中，且 SHALL 能从 API 响应中读取实际值写入 state

