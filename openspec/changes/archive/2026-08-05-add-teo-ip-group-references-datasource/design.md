## Context

当前 Terraform Provider for TencentCloud 中，TEO（EdgeOne）产品缺少查询 IP 分组被引用信息的数据源。用户需要通过 `DescribeIPGroupReferences` 接口分页查询指定站点下某个 IP 组被哪些安全策略、DDoS 防护等实体引用，以便在 Terraform 配置中引用这些数据，用于依赖分析、影响面评估和配置审计。

该接口为同步接口，支持分页查询。入参包括 `ZoneId`（站点 ID，string）、`GroupId`（IP 组 ID，int64）、`Offset`（分页偏移量，默认 0）、`Limit`（单次查询条数，取值范围 [1, 200]，默认 20）。出参包括 `References`（`[]*IPGroupReference` 引用信息列表）和 `TotalCount`（查询结果总数）。

`IPGroupReference` 结构包含以下字段：
- `ZoneId`（string）：站点 ID
- `EntityType`（string）：实体类型，枚举值如 `WebSec.ZonePolicy`、`WebSec.HostPolicy`、`WebSec.Template`、`DDoS.L4Proxy`、`DDoS.L3Transit`
- `EntityId`（string）：实体标识，根据 EntityType 不同代表不同含义
- `EntityName`（string）：实体名称，根据 EntityType 不同代表不同含义
- `SubEntityType`（string）：子实体类型，枚举值如 `WebSec.ExceptionRule`、`WebSec.BasicAccessRule` 等
- `SubEntityId`（string）：子实体标识，根据 SubEntityType 不同代表不同含义
- `SubEntityName`（string）：子实体名称，根据 SubEntityType 不同代表不同含义

现有 TEO 数据源文件位于 `tencentcloud/services/teo/` 目录下，遵循 `data_source_tc_teo_<name>.go` 命名规范。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_ip_group_references` 数据源，封装 `DescribeIPGroupReferences` API
- 支持按 `zone_id` 和 `group_id` 查询指定 IP 组的引用信息列表
- 内部自动分页获取所有引用数据，不暴露 limit/offset 参数给用户
- 在 provider 中注册该数据源
- 提供单元测试（使用 gomonkey mock）和文档

**Non-Goals:**
- 不支持创建、更新或删除 IP 分组引用（这些由对应的安全策略/DDoS 防护资源负责）
- 不暴露分页参数 limit/offset 给 Terraform 用户
- 不修改已有的 IP 分组相关资源/数据源

## Decisions

1. **数据源文件命名**: `data_source_tc_teo_ip_group_references.go`，遵循现有命名规范 `data_source_tc_teo_<name>.go`
2. **分页策略**: 接口 `DescribeIPGroupReferences` 支持分页，Limit 取值范围为 [1, 200]，默认 20。数据源内部自动分页获取所有引用数据，不暴露 limit/offset 给用户。使用 Limit=200（API 注释中标注的最大值）逐页获取，若单页返回条数小于 limit 则结束翻页，否则 offset 递增继续获取。
3. **ID 生成**: 由于该数据源为查询类数据源，使用 `helper.BuildToken()` 生成数据源 ID，与同类数据源（如 `tencentcloud_teo_security_ip_group_content`）保持一致。
4. **GroupId 类型**: 云 API 中 `GroupId` 为 `int64` 类型，在 Terraform Schema 中使用 `TypeInt` 对应。
5. **References schema 结构**: 云 API 返回的 `References` 为 `[]*IPGroupReference` 列表。遵循项目规范"禁止创建列表型数据嵌套层"，将 `references` 定义为 `TypeList`，其元素为包含各引用字段的 `schema.Resource`，使每个字段（entity_type、entity_id、entity_name、sub_entity_type、sub_entity_id、sub_entity_name、zone_id）都可被 terraform 单独 set/read，不额外嵌套 `xxx_set` 一层。
6. **字段命名**: 引用条目内的字段使用蛇形命名：`zone_id`、`entity_type`、`entity_id`、`entity_name`、`sub_entity_type`、`sub_entity_id`、`sub_entity_name`，与云 API 的 `IPGroupReference` 字段一一对应。
7. **Retry 逻辑**: 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，失败时使用 `tccommon.RetryError()` 包装错误。retry 块内仅调用云 API 接口，set 操作放到 retry 块外。
8. **空返回处理**: 遵循项目规范，在 retry 块内检查云 API Read 接口是否返回空（`response == nil` / `response.Response == nil` / `len(response.Response.References) == 0`）。若返回空，不直接 `d.SetId("")`，而是直接返回 `tccommon.RetryError` 包装的 `NonRetryableError`，避免因云 API 短暂波动导致本地 state 中的 id 被清空。
9. **代码风格**: 参考 `tencentcloud_igtm_instance_list` 数据源的业务逻辑模式（列表型数据源将响应字段映射到 list 元素）和 `tencentcloud_teo_security_ip_group_content` 数据源的分页/TEO 客户端调用模式。

## Risks / Trade-offs

- [引用数据量可能超过单页] → 使用分页查询确保能获取所有数据，Limit 设为 API 支持的最大值 200，逐页翻页直到返回条数小于 limit
- [API 短暂波动返回空导致 state 丢失] → retry 块内对空返回返回 NonRetryableError 而非清空 id，让外层 retry 继续尝试，最终以"重试耗尽"形式失败，便于人工介入
- [API 变更风险] → 仅使用已通过 vendor 目录验证存在的 `DescribeIPGroupReferences` 接口及 `IPGroupReference` 结构，参数与字段映射已验证
