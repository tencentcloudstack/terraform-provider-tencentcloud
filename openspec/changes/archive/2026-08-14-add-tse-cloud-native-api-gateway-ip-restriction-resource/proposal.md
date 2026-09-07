## Why

腾讯云 TSE（微服务引擎）云原生 API 网关已提供 IP 访问控制能力（白名单/黑名单），相关云 API（`CreateOrModifyCloudNativeAPIGatewayIPRestriction`、`DescribeCloudNativeAPIGatewayIPRestriction`、`DeleteCloudNativeAPIGatewayIPRestriction`）已在 `tencentcloud-sdk-go/tencentcloud/tse/v20201207` 发布，但 Terraform Provider 尚未暴露该能力，导致用户只能在控制台手动配置网关 IP 限制策略，无法将其纳入基础设施即代码（IaC）管理。

本 change 新增通用型（RESOURCE_KIND_GENERAL）资源 `tencentcloud_tse_cloud_native_api_gateway_ip_restriction`，让用户在 HCL 中声明网关的 IP 访问控制策略，由 Provider 调用上述三个接口完成资源的创建（或编辑）、查询、删除。

## What Changes

- 新增通用型资源 `tencentcloud_tse_cloud_native_api_gateway_ip_restriction`，完整覆盖 CRUD 生命周期。
- 新增 Schema 字段，严格按云 API 入参映射：
  - `gateway_id`（TypeString, Required, ForceNew）：网关 ID。
  - `source_type`（TypeString, Required, ForceNew）：访问控制绑定的资源类型，route|service。
  - `source_id`（TypeString, Required, ForceNew）：路由或服务 ID。
  - `enabled`（TypeBool, Optional, Computed）：是否启用插件。
  - `restriction_type`（TypeString, Optional, Computed）：访问控制类型，whiteList|blackList。
  - `address_list`（TypeSet of TypeString, Optional, Computed）：IP/CIDR 列表（顺序不敏感）。
- 资源 ID 为复合 ID：`gateway_id#source_type#source_id`，使用 `tccommon.FILED_SP`（"#"）作为分隔符。这三个字段均为 ForceNew，改变任意一个即重建资源。
- Create 与 Update 共用 `CreateOrModifyCloudNativeAPIGatewayIPRestriction` 接口（API 为 upsert 语义）。
- Read 调用 `DescribeCloudNativeAPIGatewayIPRestriction` 查询当前配置，刷新 state。
- Delete 调用 `DeleteCloudNativeAPIGatewayIPRestriction` 删除。
- 资源支持 import（import 时使用复合 ID）。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册新资源。
- 新增 `.md` 资源文档与 `_test.go` 单元测试（使用 gomonkey mock 云 API）。

## Capabilities

### New Capabilities

- `tse-cloud-native-api-gateway-ip-restriction-resource`: 新增 `tencentcloud_tse_cloud_native_api_gateway_ip_restriction` 资源的 schema、CRUD（Create=CreateOrModify / Read=Describe / Update=CreateOrModify / Delete=Delete）行为、复合 ID 约定、字段约束、文档与测试规范，作为 TSE 云原生网关 IP 访问控制能力的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go / provider.md 增加注册条目。 -->

## Impact

- 新文件：
  - `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction.go`
  - `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction.md`
  - `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction_test.go`
  - `website/docs/r/tse_cloud_native_api_gateway_ip_restriction.html.markdown`（由 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/provider.go`：在 tse 资源注册段追加 `"tencentcloud_tse_cloud_native_api_gateway_ip_restriction": tse.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction()` 一行。
  - `tencentcloud/provider.md`：在 TSE Resource 段追加 `tencentcloud_tse_cloud_native_api_gateway_ip_restriction`，使 gendoc 索引可识别。
- 依赖：现有 SDK `tencentcloud-sdk-go/tencentcloud/tse/v20201207` 已包含三个接口及其 model，无需升级 SDK。
- 不修改任何既有资源/数据源 schema，对现有用户零影响。
