## Context

EdgeOne (TEO) 推理服务（Inference Service）通过云 API `DescribeInferenceServiceDeploymentLogs`（teo v20220901）提供部署日志查询能力。当前 Terraform Provider 中没有任何数据源覆盖该接口，用户无法在 Terraform 中拉取推理服务部署日志。

该接口为同步只读查询接口（非异步），入参包含站点 ID、推理服务 ID、部署记录 ID 及可选的时间范围与排序字段，出参为部署日志列表（每条含 `LogMessage` 与 `Timestamp`）。接口支持 `Offset`/`Limit` 分页，`Limit` 最大值为 1000。

本变更遵循项目既有的 RESOURCE_KIND_DATASOURCE 实现模式，参考 `tencentcloud_teo_config_group_versions` 数据源（同为 teo 服务、列表型返回、内部自动分页）。

## Goals / Non-Goals

**Goals:**
- 新增数据源 `tencentcloud_teo_inference_service_deployment_logs`，支持通过 `zone_id`、`service_id`、`record_id` 查询推理服务部署日志。
- 支持可选的时间范围检索（`start_time`、`end_time`）与排序（`sort_by`、`sort_order`）参数。
- 在服务层实现内部自动分页（`Limit` 取云 API 注释最大值 1000），不向用户暴露 `limit`/`offset` 参数。
- 在 Read 方法中按规则使用 `tccommon.ReadRetryTimeout` 进行重试包装，并在 retry 块内对云 API 返回空的情况返回 `NonRetryableError`，避免误清 state。
- 在 provider 中注册该数据源并生成对应文档。

**Non-Goals:**
- 不创建/修改/删除推理服务或部署记录（本变更仅为只读数据源，不涉及 CRD 操作）。
- 不暴露分页参数 `offset`/`limit` 给用户。
- 不引入新的外部依赖。
- 不修改任何现有资源或数据源的 schema。

## Decisions

### 决策 1：列表字段采用嵌套结构 `deployment_log_info_set`，内部元素字段平铺

**选择**：在 schema 中使用 `deployment_log_info_set`（TypeList）作为日志列表容器，其中每个元素为 TypeResource，元素内的 `log_message`、`timestamp` 直接作为该 Resource 的字段，不再额外嵌套一层"列表型数据"结构。

**理由**：
- 云 API 出参 `DeploymentLogInfoSet` 本身就是数组，每个 `InferenceServiceDeploymentLogInfo` 含 `LogMessage` 与 `Timestamp` 两个字段。
- 参考既有数据源（如 `tencentcloud_teo_config_group_versions` 的 `config_group_version_infos`、`tencentcloud_igtm_instance_list` 的 `instance_set`），列表型返回统一使用一层 TypeList + 内部 TypeResource 元素的模式，元素字段直接平铺，不再套一层容器。
- 该模式既符合云 API 响应结构，也便于 Terraform 对每个字段单独 set/read。

**备选方案**：将 `log_message`、`timestamp` 直接平铺到资源 schema 顶层（即去掉 `deployment_log_info_set` 这一层）。未采用，因为云 API 返回的是日志数组（多条日志），顶层平铺会丢失"列表"语义，无法表达多条日志，且与同类型数据源的既有实现不一致。

### 决策 2：分页在服务层内部循环实现，Limit 取 1000

**选择**：在 `DescribeTeoInferenceServiceDeploymentLogs` 服务方法中循环调用云 API，`Offset` 从 0 递增、`Limit` 固定为 1000（云 API 注释标注的最大值），直到某次返回数量小于 `Limit` 时退出循环。

**理由**：
- 项目硬约束要求数据源分页不暴露 `limit`/`offset` 给用户，内部自动分页获取所有数据。
- 参考既有实现 `DescribeTeoConfigGroupVersionsByFilter`（limit=100）与 `DescribeTeoZonesByFilter`（pageSize=100）的分页循环模式。
- 取 1000 作为单页大小可在日志条数较多时减少请求次数，符合"分页字段给定云 API 注释最大值"的要求。

### 决策 3：Read 方法中 retry 块内对空返回返回 NonRetryableError

**选择**：在数据源 Read 的 `resource.Retry` 块内，调用服务方法后若结果为空（`len(respData) == 0`），直接返回 `NonRetryableError`，并在外层 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。

**理由**：
- 需求规则明确要求 RESOURCE_KIND_DATASOURCE 的 Read retry 块内遇到空返回时不直接 `d.SetId("")`，而是返回 `NonRetryableError`，避免云 API 短暂波动导致 state 被清空。
- 让上层任务以"重试耗尽"形式失败，便于人工介入排障。

### 决策 4：数据源 ID 使用 helper.DataResourceIdsHash 或 helper.BuildToken

**选择**：使用 `helper.DataResourceIdsHash(ids)` 生成数据源 ID（其中 `ids` 由各日志条目的 `timestamp` 组成），若日志列表为空则使用空列表生成。

**理由**：
- 参考 `tencentcloud_teo_config_group_versions` 使用 `helper.DataResourceIdsHash(ids)` 的模式。
- 数据源无真实主键，使用基于返回内容计算的 hash 作为 ID 可保证 plan/refresh 的稳定性。

## Risks / Trade-offs

- **[风险] 日志条数极大时单次 Read 拉取耗时长** → 通过内部自动分页（每页 1000 条）控制单次请求规模；日志按时间范围检索可由用户通过 `start_time`/`end_time` 收敛范围。
- **[风险] 云 API 短暂波动导致 Read 失败** → 使用 `tccommon.ReadRetryTimeout` 重试；对空返回返回 `NonRetryableError` 以避免误清 state。
- **[权衡] 不暴露分页参数** → 牺牲了用户手动控制分页的灵活性，换取了配置简洁性，符合项目数据源分页硬约束。