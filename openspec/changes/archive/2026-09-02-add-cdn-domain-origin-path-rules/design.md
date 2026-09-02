## Context

`tencentcloud_cdn_domain` 是 CDN 域名管理资源，其 `origin` 配置块目前包含 `origin_type`、`origin_list`、`server_name`、`cos_private_access`、`origin_pull_protocol`、`backup_origin_type`、`backup_origin_list`、`backup_server_name`、`origin_company` 等字段，但不包含回源路径重写规则（`PathRules`）。

云 API 的 `Origin` 结构体中已有 `PathRules []*PathRule` 字段，`PathRule` 结构体包含 `Regex`（bool，通配符匹配开关）、`Path`（string，匹配的 URL 路径）、`ServerName`（string，回源 Host 头部）、`ForwardUri`（string，回源 URI 路径）等字段。`AddCdnDomain`、`DescribeDomainsConfig`、`UpdateDomainConfig` 三个接口均通过 `Origin.PathRules` 支持该功能。

约束：必须保持向后兼容，新增参数均为 Optional + Computed。

## Goals / Non-Goals

**Goals:**
- 在 `origin` 配置块中新增 `path_rules` 子列表参数，支持配置回源路径重写规则
- `path_rules` 包含 `regex`、`path`、`server_name`、`forward_uri` 四个可选字段
- Create/Read/Update 方法均支持该参数的完整生命周期

**Non-Goals:**
- 不新增 `PathRule` 中的其他字段（`Origin`、`OriginArea`、`RequestHeaders`、`FullMatch`），本次仅覆盖需求中指定的 4 个字段
- 不修改 `origin` 配置块中已有的字段
- 不新增独立的顶层参数

## Decisions

### 决策 1：`path_rules` 作为 `origin` 配置块的嵌套子列表

**选择**：将 `path_rules` 定义为 `origin` 配置块内部的 `TypeList` 子列表，每个元素为一个 `schema.Resource`，包含 `regex`、`path`、`server_name`、`forward_uri` 四个字段。

**理由**：云 API 中 `PathRules` 是 `Origin` 结构体的嵌套字段，语义上属于源站配置的一部分。将其作为 `origin` 的子列表符合 API 结构，也与现有 `origin` 块中其他字段保持一致的组织方式。

### 决策 2：四个字段均为 Optional + Computed

**选择**：`regex`（TypeBool）、`path`（TypeString）、`server_name`（TypeString）、`forward_uri`（TypeString）均设为 Optional + Computed。

**理由**：云 API 注释明确标注这些字段"可能返回 null"，且均为可选入参。Optional + Computed 确保用户未配置时能从 API 响应中读取实际值，避免 import 后产生 plan diff。

### 决策 3：仅在 `d.HasChange("origin")` 时填充 Update 请求

**选择**：在 Update 方法中，`path_rules` 的填充逻辑放在已有的 `if d.HasChange("origin")` 块内，与 `origin` 块中其他字段一起处理。

**理由**：`path_rules` 是 `origin` 的子字段，当 `origin` 发生变更时整体重新提交。这与现有 `origin` 块中 `origin_type`、`origin_list` 等字段的处理方式一致，避免单独对子列表做 `HasChange` 判断引入复杂性。

## Risks / Trade-offs

- **[风险] API 返回的 `PathRule` 中包含未暴露的字段** → 缓解：本次仅映射需求中指定的 4 个字段，其余字段（`Origin`、`OriginArea`、`RequestHeaders`、`FullMatch`）不纳入 schema，API 返回的这些字段会被忽略，不影响已映射字段的功能
- **[风险] 用户配置 `path_rules` 但 API 未生效** → 缓解：`PathRule` 的 `Path` 和 `Regex` 是核心匹配字段，`ServerName` 和 `ForwardUri` 是回源行为字段，均为 API 支持的独立字段，不依赖未暴露的字段即可工作
