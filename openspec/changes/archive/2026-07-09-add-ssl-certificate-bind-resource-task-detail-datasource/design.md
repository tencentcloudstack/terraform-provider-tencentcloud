## Context

SSL 服务提供证书关联云资源的异步查询能力。用户通过 `CreateCertificateBindResourceSyncTask` 创建查询任务后，需要通过 `DescribeCertificateBindResourceTaskDetail` 接口查询任务结果详情。当前 Terraform Provider 缺少对应数据源，用户无法在 IaC 中获取证书关联的各云资源（CLB、CDN、WAF、DDoS、Live、VOD、TKE、APIGateway、TCB、TEO、TSE、COS、TDMQ、MQTT、GAAP、SCF）的详细绑定信息。

该接口为异步查询接口：返回 `Status` 字段（0=查询中，1=查询成功，2=查询异常），需要轮询直到任务完成。参考实现遵循 `tencentcloud_igtm_instance_list` 数据源模式（`RESOURCE_KIND_DATASOURCE`）。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_ssl_certificate_bind_resource_task_detail` 数据源，支持通过任务 ID 查询证书关联云资源任务结果详情
- 正确处理异步任务状态，轮询直到 Status != 0
- 完整映射 16 类云资源关联详情的嵌套结构
- 遵循现有 SSL 数据源代码风格（参考 `data_source_tc_ssl_describe_host_clb_instance_list.go`）

**Non-Goals:**
- 不创建/修改/删除任务（本数据源仅读取）
- 不暴露 limit/offset 分页参数给用户（内部自动处理分页，Limit 取云 API 最大值 100）
- 不支持 `DescribeCertificateBindResourceTaskResult` 接口（该接口返回任务列表概览，非本数据源范围）

## Decisions

### Decision 1: 异步任务轮询策略

`DescribeCertificateBindResourceTaskDetail` 返回 `Status` 字段（0=查询中，1=查询成功，2=查询异常）。

**选择**: 在 service 层实现轮询，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹。在 retry 块内检查 `Status`，若为 0（查询中）则返回 `tccommon.RetryError` 以触发重试；若为 1 或 2 则视为完成并返回结果。

**理由**: 符合项目"最终一致性重试"模式，且与 `tccommon.ReadRetryTimeout` 超时控制一致。Status=2（异常）时不重试，直接返回结果让用户查看 `Error` 字段。

### Decision 2: 分页处理

该接口入参含 `Limit`（最大 100）和 `Offset`，但出参中各云资源列表（如 `CLB`、`CDN` 等）为 `[]*XxxInstanceList`，每个 `XxxInstanceList` 内含 `InstanceList`（详情数组）和 `TotalCount`，分页针对的是每个云资源地域下的实例总数。

**选择**: service 层内部循环分页，`Limit` 设为 100（云 API 最大值），循环累加 `Offset` 直到所有云资源详情获取完毕。不向用户暴露 limit/offset。

**理由**: 遵循项目"数据源分页:不暴露 limit/offset 参数给用户,内部实现自动分页获取所有数据"约束。

### Decision 3: Schema 结构设计（列表展开）

根据项目规则"资源列表型数据禁止再嵌套一层 schema，应将列表中每个元素的参数平铺"，但本数据源出参本身是多个独立的云资源类型列表（CLB、CDN 等互不相同的结构），每类云资源是 `[]*XxxInstanceList`，每个 `XxxInstanceList` 含 `Region`、`TotalCount`、`Error`、`InstanceList`（详情数组）。

**选择**: 每个 `XxxInstanceList` 作为 schema 顶层字段（`clb`、`cdn` 等），类型为 `TypeList`，其 `Elem` 为 `Resource`，展开 `Region`、`TotalCount`、`Error` 字段及 `InstanceList`（嵌套详情）。详情结构按各云资源 SDK 类型完整映射。

**理由**: 出参为多种异构云资源列表的集合，无法"展开"为同一层级；每种云资源作为独立的顶层列表字段，符合 Terraform 数据源扁平化原则。

### Decision 4: 空 ID 处理

数据源无唯一业务 ID，采用 `helper.BuildToken()` 生成数据源 ID（参考 `tencentcloud_igtm_instance_list`）。

### Decision 5: 单元测试方式

根据项目规则，新增 terraform 资源（含数据源）的测试使用 mock（gomonkey）方法对云 API 进行 mock 处理，仅进行业务代码逻辑的单元测试，使用 `go test -gcflags=all=-l` 跑通。

## Risks / Trade-offs

- **[Risk] 异步任务长时间未完成** → 通过 `tccommon.ReadRetryTimeout` 超时控制，超时后返回错误，便于人工介入
- **[Risk] 各云资源详情结构差异大、嵌套深** → 严格按照 vendor SDK 中各 `XxxInstanceDetail` 结构体字段映射，避免遗漏
- **[Risk] 接口返回空（response/Response 为空）** → 根据 RESOURCE_KIND_DATASOURCE 规则，retry 块内检查返回空时直接返回 `NonRetryableError`，避免清空 state 中的 id
- **[Trade-off] 分页合并多次请求结果** → 增加了请求次数但保证数据完整性，符合"获取所有数据"约束
