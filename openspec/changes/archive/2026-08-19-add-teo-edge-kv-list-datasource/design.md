## Context

TEO EdgeKV 已在 Provider 中落地了 `tencentcloud_teo_edge_kv`（单键 CRUD，基于 `EdgeKVPut`/`EdgeKVGet`/`EdgeKVDelete`）与 `tencentcloud_teo_edge_kv_namespace`（命名空间管理，基于 `CreateEdgeKVNamespace`/`DescribeEdgeKVNamespaces`/`ModifyEdgeKVNamespace`/`DeleteEdgeKVNamespace`）两个资源。但当前缺少枚举命名空间下全部键名的能力。

腾讯云 TEO SDK（`teo/v20220901`）已提供 `EdgeKVList` 接口：

- 请求 `EdgeKVListRequest`：`ZoneId`、`Namespace`、`Prefix`、`Cursor`、`Limit`（默认 20，最大 1000）。
- 响应 `EdgeKVListResponseParams`：`Keys []*string`、`Cursor *string`（空字符串表示已遍历完）、`RequestId`。

接口为同步查询接口（非异步），无需轮询 Read 接口确认生效。

当前状态：
- `vendor/.../teo/v20220901` 中 `EdgeKVListRequest`/`EdgeKVListResponse` 类型已就绪。
- `TeoService` 已存在，通过 `me.client.UseTeoClient()` 获取客户端。
- 现有数据源参考实现：`data_source_tc_igtm_instance_list.go`（service 层 `DescribeXxxByFilter` + 轻量 Read 模式）。

约束：
- 数据源不向用户暴露 `limit` 参数；内部固定取云 API 注释标注的最大值 1000，并通过 `Cursor` 自动循环直到游标为空字符串。
- RESOURCE_KIND_DATASOURCE 资源 Read 方法的 retry 块内必须检查 API 是否返回空，返回空时不直接 `d.SetId("")`，而是返回 `NonRetryableError`，避免短暂波动清空 state。
- 向后兼容：纯新增数据源，不触碰任何现有资源。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_teo_edge_kv_list` 数据源，让用户能在 Terraform 中声明式查询某命名空间下的全部键名。
- 支持按 `prefix` 前缀过滤、按 `cursor` 续续遍历。
- 对用户屏蔽 `Limit`/分页细节，内部自动循环直至游标为空字符串。
- 通过 gomonkey mock 的单元测试覆盖 Read 主流程与空响应路径。

**Non-Goals:**
- 不提供单键的值读取（已有 `tencentcloud_teo_edge_kv` 资源 Read 覆盖）。
- 不暴露 `limit` 参数给用户。
- 不修改 `tencentcloud_teo_edge_kv` / `tencentcloud_teo_edge_kv_namespace` 资源的任何行为。
- 不实现异步轮询（`EdgeKVList` 为同步接口）。

## Decisions

### Decision 1: 在 service 层新增 `DescribeTeoEdgeKvListByFilter` 方法，Read 函数保持轻量

**选择**：参照 `data_source_tc_igtm_instance_list.go` / `data_source_tc_teo_content_quota.go` 的模式，在 `service_tencentcloud_teo.go` 中新增 `DescribeTeoEdgeKvListByFilter(ctx, param map[string]interface{}) (keys []*string, cursor *string, err error)` service 层方法，将 `EdgeKVList` 调用、retry、分页循环全部封装在 service 层。数据源 `dataSourceTencentCloudTeoEdgeKvListRead` 仅负责组装 `paramMap`、调用 service 方法、`d.Set` 回填结果、设置 ID 与输出文件，不再内联 retry/分页逻辑。

**备选**：在 Read 函数内直接调用云 API，分页与 retry 内联（参照 `data_source_tc_teo_security_ip_group_content.go`）。

**理由**：
- 项目规范要求 RESOURCE_KIND_DATASOURCE 资源的代码生成格式严格参考 `tencentcloud_igtm_instance_list` 资源的业务逻辑代码，即 service 层 `DescribeXxxByFilter` + 轻量 Read 模式。
- 同服务的 `data_source_tc_teo_content_quota.go` 也采用此模式，保持一致性。
- service 层封装后，retry 与分页细节集中在一处，Read 函数更简洁；由于 service 方法内部已有 retry，Read 函数无需再包一层 retry（避免 retry 嵌套）。

### Decision 2: service 层内部固定 Limit=1000，通过 Cursor 自动循环遍历

**选择**：在 `DescribeTeoEdgeKvListByFilter` 的 for 循环中固定 `request.Limit = helper.IntInt64(1000)`（云 API 注释标注的最大值），每次循环将上一次响应的 `Cursor` 填入 `request.Cursor`，直到响应 `Cursor` 为空字符串（或 nil）时跳出循环。将累计的 `Keys` 与最后一次的 `Cursor` 返回给 Read 函数，由 Read 一次性 `d.Set("keys", ...)` 与 `d.Set("cursor", ...)`。

**备选**：将 `cursor` 仅作为入参透传，单次调用即返回，不做内部循环。

**理由**：
- 项目硬约束：数据源不暴露 limit/offset 给用户，内部实现自动分页获取所有数据。
- `EdgeKVList` 使用游标分页（非 offset/limit 翻页），`Cursor` 为空字符串表示遍历完成，循环终止条件明确。
- 把全部键名聚合后一次性 set，符合数据源"一次查询返回完整结果"的语义。

### Decision 3: `cursor` 既是入参也是出参

**选择**：`cursor` 字段定义为 `Optional + Computed`（TypeString）。作为入参时用于续续遍历；作为出参时在 Read 末尾用最后一次响应的 `Cursor` 回填，便于用户在 HCL 中读取并用于下一次查询。

**备选**：仅作为入参，不回填。

**理由**：
- 云 API 响应的 `Cursor` 是游标遍历的核心产出，用户需要它判断是否还有更多数据并继续遍历。
- `Optional + Computed` 是 Provider 处理"可选入参 + 由云端回填"字段的标准模式，与现有数据源一致。

### Decision 4: service 层 retry 块内对空响应返回 NonRetryableError

**选择**：在 `DescribeTeoEdgeKvListByFilter` 内的 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中，若 `result == nil || result.Response == nil`，返回 `resource.NonRetryableError(...)`；在 retry 失败路径上保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。遍历循环中 `Keys` 为空但 `Cursor` 非空字符串时继续翻页，不视为错误。Read 函数收到 service 返回的 error 后直接返回，不调用 `d.SetId("")`，从而保留 state id。

**备选**：返回 nil 让外层继续，或直接 `d.SetId("")`。

**理由**：
- 遵循项目规范：RESOURCE_KIND_DATASOURCE 资源 Read 方法的 retry 块内必须检查 API 是否返回空，返回空时不直接 `d.SetId("")`，而是返回 `NonRetryableError`，避免短暂波动清空 state。
- 空的 `Keys` 列表是合法的业务结果（命名空间下暂无键），不应视为 API 错误；仅 `response`/`Response` 结构为 nil 才视为异常。

### Decision 5: 数据源 ID 使用 `helper.BuildToken()`

**选择**：Read 成功后 `d.SetId(helper.BuildToken())`，与 `data_source_tc_igtm_instance_list`、`data_source_tc_teo_security_ip_group_content` 等现有数据源一致。

**理由**：数据源无真实云端 ID，使用随机 token 作为 Terraform state 占位 ID 是项目通用模式。

### Decision 6: `keys` 使用 TypeList of String 平铺到顶层

**选择**：`keys` 定义为 `schema.TypeList`，`Elem: &schema.Schema{Type: schema.TypeString}`，直接平铺在资源 schema 顶层，不再嵌套一层"列表项对象"。

**理由**：`Keys` 本身就是 `[]*string`，元素为标量字符串，无需嵌套对象 schema；遵循"列表展开平铺到顶层"的规范。

## Risks / Trade-offs

- **Risk**：命名空间下键名数量极大时，内部自动循环可能产生多次 API 调用 → **Mitigation**：每次取 Limit=1000，单命名空间键名量级通常远小于此；循环以游标空字符串终止，无死循环风险。
- **Risk**：`cursor` 同时为入参和出参，用户传入过期的 `cursor` 可能导致从中间位置开始遍历，遗漏前段键名 → **Mitigation**：文档中明确说明 `cursor` 用于续续遍历，首次查询不填写；这是云 API 游标的标准语义。
- **Trade-off**：新增 service 层方法 `DescribeTeoEdgeKvListByFilter`，service 文件变长 → 可接受，与 `data_source_tc_igtm_instance_list` / `data_source_tc_teo_content_quota` 模式一致，便于维护。

## Migration Plan

- 纯新增数据源，无 state 迁移需求。
- 存量配置不受影响：未引用该数据源的现有 TF 配置无任何 plan diff。
- 回滚：删除新增文件并移除 provider.go/provider.md 中的注册条目即可，无残留 state。

## Open Questions

- 无
