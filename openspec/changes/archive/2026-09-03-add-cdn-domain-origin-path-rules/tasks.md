## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/cdn/resource_tc_cdn_domain.go` 的 `origin` 配置块 schema 中新增 `path_rules` 子列表参数（TypeList, Optional, Computed），每个元素为 `schema.Resource`，包含 `regex`（TypeBool, Optional, Computed）、`path`（TypeString, Optional, Computed）、`server_name`（TypeString, Optional, Computed）、`forward_uri`（TypeString, Optional, Computed）四个字段

## 2. CRUD 函数实现

- [x] 2.1 在 Create 方法（`resourceTencentCloudCdnDomainCreate`）的 `origin` 处理块中，读取 `origin["path_rules"]` 并填充 `request.Origin.PathRules`（`[]*cdn.PathRule`），仅填充用户配置的非空字段
- [x] 2.2 在 Read 方法（`resourceTencentCloudCdnDomainRead`）的 `origin` 处理块中，将 `domainConfig.Origin.PathRules` 中的 `Regex`、`Path`、`ServerName`、`ForwardUri` 写入 `origin` map 的 `path_rules` 字段（若 API 返回非 nil）
- [x] 2.3 在 Update 方法（`resourceTencentCloudCdnDomainUpdate`）的 `if d.HasChange("origin")` 块中，读取 `origin["path_rules"]` 并填充 `request.Origin.PathRules`，与 Create 方法逻辑一致
- [x] 2.4 在 Update 方法调用 `UpdateDomainConfig` 之前，当配置了 `path_rules` 时先通过 `DescribeDomainsConfig`（`CdnService.DescribeDomainsConfigByDomain`）查询云上当前的 `Origin.PathRules`；对 tfstate 中每条规则，若其 `{regex, path, server_name, forward_uri}` 四元组（nil 与空字符串视为等同）已存在于查询结果中，则从结果中获取 `FullMatch`、`OriginArea`、`Origin`、`RequestHeaders` 四个参数并填充到对应的 `cdn.PathRule` 中；若四元组不存在，则仅提交用户配置的字段

## 3. 测试

- [x] 3.1 在 `tencentcloud/services/cdn/resource_tc_cdn_domain_test.go` 中补充 `path_rules` 参数的验收测试用例（`TestAccTencentCloudCdnDomainResource_OriginPathRules`），覆盖 Create/Read/Update 流程中 `path_rules` 的配置、读取与变更

## 4. 文档

- [x] 4.1 更新 `tencentcloud/services/cdn/resource_tc_cdn_domain.md` 示例文件，在 `origin` 配置块的示例中补充 `path_rules` 的用法示例
