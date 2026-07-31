## Context

`dbdc` (DB Dedicated Cloud / DB Custom) 服务已有集群、节点、镜像等数据源和资源，但缺少查询集群资源汇总信息的数据源。腾讯云 API `DescribeDBCustomClusterResources` 已存在于 vendored SDK `tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`，`UseDbdcV20201029Client()` 已可用。该接口接收 `ClusterId` 入参，返回单个集群的资源汇总对象（非列表），包含节点总数以及 CPU/内存/Pods 的容量、可分配量、申请量、上限和可用余量。

现有 `dbdc` 数据源均采用 `service.<Xxx>ByFilter` 模式处理列表型响应。本数据源的响应是单对象，不是列表，因此 schema 设计将出参字段平铺到顶层，嵌套的 `MetaResource` 对象展开为 TypeList（MaxItems: 1）嵌套结构。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_dbdc_db_custom_cluster_resources` 数据源，通过 `cluster_id` 查询集群资源汇总信息
- 在 service 层新增 `DescribeDBCustomClusterResources` 方法封装 SDK 调用
- 正确处理 `MetaResource` 嵌套对象（Cpu/Memory/Pods），展开为 Terraform 嵌套 schema
- 遵循项目规范：retry 包装、nil 检查、NonRetryableError 处理空响应

**Non-Goals:**
- 不支持分页（该 API 返回单对象，无分页需求）
- 不修改现有 dbdc 资源或数据源
- 不暴露 `limit`/`offset` 参数

## Decisions

### Decision 1: 响应字段平铺到顶层 schema

**Rationale**: `DescribeDBCustomClusterResources` 返回的是单个资源汇总对象，不是列表。根据项目规范，不创建额外的 `xxx_set`/`xxx_list` 嵌套层，将 `node_count`、`capacity`、`allocatable`、`requests`、`limits`、`available` 直接作为顶层 Computed 字段。

### Decision 2: MetaResource 嵌套对象使用 TypeList (MaxItems: 1)

**Rationale**: `MetaResource` 包含 `Cpu *float64`、`Memory *float64`、`Pods *uint64` 三个字段。使用 `TypeList` + `MaxItems: 1` 的嵌套 `schema.Resource` 来映射，与项目中对单对象嵌套的惯例一致。读取时取列表第一项（`[0]`），并对 nil/空列表做安全检查。

**Alternatives considered**:
- 使用 `TypeMap` — 不合适，因为字段类型混合 float64 和 int64，且语义明确为固定结构
- 将 cpu/memory/pods 也平铺到顶层（如 `capacity_cpu`） — 字段名冗长且失去结构化语义

### Decision 3: Service 层方法签名

`DescribeDBCustomClusterResources(ctx context.Context, clusterId string) (ret *dbdcv20201029.DescribeDBCustomClusterResourcesResponseParams, errRet error)`

**Rationale**: 与现有 `DescribeDBCustomClusterById` 方法风格一致，入参为单个 ID 字符串，返回 `ResponseParams` 指针。内部使用 `resource.Retry(ReadRetryTimeout)` 包装，retry 块内检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`。

### Decision 4: Read 函数空响应处理

**Rationale**: 根据项目规范，数据源 Read 方法 retry 块内若云 API 返回空（`response == nil` / `response.Response == nil`），不直接 `d.SetId("")`，而是返回 `NonRetryableError`，让外层 retry 继续尝试，最终以"重试耗尽"形式失败，避免因云 API 短暂波动导致 state 丢失。

### Decision 5: 数据源 ID

使用 `helper.BuildToken()` 生成 ID，与 `data_source_tc_dbdc_db_custom_clusters.go` 等现有数据源一致。

## Risks / Trade-offs

- [Risk] `MetaResource` 字段可能为 nil → Read 时对每个嵌套字段做 nil 检查，nil 则跳过 set，避免 panic
- [Risk] API 返回空响应 → retry 块内返回 `NonRetryableError`，外层打印 `[DATASOURCE] read empty` 日志
- [Trade-off] 使用 TypeList (MaxItems:1) 而非单对象 — 与项目其他嵌套对象惯例对齐，牺牲少许简洁性换取一致性
