## Context

当前 Terraform Provider for TencentCloud 中，TEO（EdgeOne）产品缺少查询安全 IP 组中 IP 列表的数据源。用户需要通过 `DescribeSecurityIPGroupContent` 接口分页查询 IP 组中的 IP 或网段列表及总数，以便在 Terraform 配置中引用这些数据。

该接口为同步接口，支持分页查询，入参包括 `ZoneId`（站点 ID）和 `GroupId`（IP 组 ID），出参包括 `IPTotalCount`（IP 总数）和 `IPList`（IP 列表）。

现有 TEO 数据源文件位于 `tencentcloud/services/teo/` 目录下，遵循 `data_source_tc_teo_<name>.go` 命名规范。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_security_ip_group_content` 数据源，封装 `DescribeSecurityIPGroupContent` API
- 支持按 `zone_id` 和 `group_id` 查询指定 IP 组中的 IP 列表
- 内部自动分页获取所有 IP 数据，不暴露 limit/offset 参数给用户
- 在 provider 中注册该数据源
- 提供单元测试（使用 gomonkey mock）和文档

**Non-Goals:**
- 不支持创建、更新或删除安全 IP 组（这些由其他资源/接口负责）
- 不暴露分页参数 limit/offset 给 Terraform 用户
- 不修改已有的安全 IP 组相关资源

## Decisions

1. **数据源文件命名**: `data_source_tc_teo_security_ip_group_content.go`，遵循现有命名规范
2. **分页策略**: 接口 `DescribeSecurityIPGroupContent` 支持分页，Limit 最大值为 100000。数据源内部自动分页获取所有 IP 数据，不暴露 limit/offset 给用户。使用 Limit=100000（API 注释中的最大值）一次性获取所有数据，若数据量超过单次限制则自动翻页。
3. **ID 生成**: 由于该数据源为查询类数据源，使用 `helper.BuildToken()` 生成数据源 ID，与同类数据源保持一致。
4. **GroupId 类型**: 云 API 中 `GroupId` 为 `int64` 类型，在 Terraform Schema 中使用 `TypeInt` 对应。
5. **IPList 类型**: 云 API 返回的 `IPList` 为 `[]*string` 类型，在 Terraform Schema 中使用 `TypeList` + `TypeString` 元素。
6. **IPTotalCount 类型**: 云 API 返回的 `IPTotalCount` 为 `int64` 类型，在 Terraform Schema 中使用 `TypeInt`。
7. **Retry 逻辑**: 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，失败时使用 `tccommon.RetryError()` 包装错误。
8. **代码风格**: 参考 `tencentcloud_teo_default_certificate` 数据源的分页模式和 `tencentcloud_igtm_instance_list` 数据源的业务逻辑模式。

## Risks / Trade-offs

- [IP 组数据量可能很大] → 使用分页查询确保能获取所有数据，Limit 设为 API 支持的最大值 100000
- [API 变更风险] → 仅使用已验证存在的 `DescribeSecurityIPGroupContent` 接口，参数映射已通过 vendor 目录验证
