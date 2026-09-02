## Why

`tencentcloud_cdn_domain` 资源的 `origin` 配置块缺少回源路径重写规则（`PathRules`）参数，用户无法通过 Terraform 配置 CDN 域名的路径匹配回源规则（包括通配符匹配、路径匹配、回源 Host 头部、回源 URI 重写）。云 API 的 `AddCdnDomain`、`DescribeDomainsConfig`、`UpdateDomainConfig` 接口均已支持 `Origin.PathRules` 下的 `Regex`、`Path`、`ServerName`、`ForwardUri` 字段，但 Terraform 资源尚未暴露这些参数。

## What Changes

- 在 `tencentcloud_cdn_domain` 资源的 `origin` 配置块中新增 `path_rules` 子列表参数（Optional + Computed），用于配置回源路径重写规则
- `path_rules` 列表中的每个元素包含 4 个可选字段：`regex`（是否开启通配符匹配，bool）、`path`（匹配的 URL 路径，string）、`server_name`（路径匹配时回源的 Host 头部，string）、`forward_uri`（路径匹配时回源的 URI 路径，string）
- 在资源的 Create 方法（`AddCdnDomain`）中填充 `request.Origin.PathRules`
- 在资源的 Read 方法（`DescribeDomainsConfig`）中将 `response.Domains.Origin.PathRules` 写入 state
- 在资源的 Update 方法（`UpdateDomainConfig`）中填充 `request.Origin.PathRules`
- 更新单元测试文件补充测试用例
- 更新资源文档（.md 文件）

## Capabilities

### New Capabilities
- `cdn-domain-origin-path-rules`: `tencentcloud_cdn_domain` 资源的 `origin` 配置块支持回源路径重写规则配置，包含通配符匹配开关、匹配路径、回源 Host 头部、回源 URI 路径四个字段

### Modified Capabilities

（无现有 spec 需要修改）

## Impact

- 修改文件：`tencentcloud/services/cdn/resource_tc_cdn_domain.go`
- 修改测试：`tencentcloud/services/cdn/resource_tc_cdn_domain_test.go`
- 修改文档：`tencentcloud/services/cdn/resource_tc_cdn_domain.md`
- 向后兼容：新增参数均为 Optional + Computed，不影响现有配置和 state
- 无新增依赖，云 API SDK 已支持相关字段
