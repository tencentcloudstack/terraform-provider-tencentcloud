## Context

腾讯云配置审计（Config）产品提供账号组（Aggregator）能力。当前 Terraform Provider 缺少查询账号组列表的数据源。本变更新增 `tencentcloud_config_list_aggregators` 数据源，封装云 API `ListAggregators`（config v20220802）。

云 API `ListAggregators` 仅包含分页入参 `Limit`/`Offset`，无业务过滤参数，因此数据源不暴露额外过滤字段，内部自动分页拉取全量数据（与 `DescribeConfigCompliancePacksByFilter` 一致的分页模式）。

参考实现风格：
- 数据源业务代码：`tencentcloud/services/config/data_source_tc_config_compliance_packs.go`（同产品数据源，分页 + 顶层 schema 列表展开模式）
- 服务层分页方法：`service_tencentcloud_config.go` 中 `DescribeConfigCompliancePacksByFilter`（for 循环 + `resource.Retry` + `total`/`Offset` 累加）
- 数据源注册：`provider.go` 中 config 数据源区块

## Goals / Non-Goals

**Goals:**
- 新增数据源 `tencentcloud_config_list_aggregators`，调用 `ListAggregators` 查询并返回账号组列表。
- 内部自动分页获取全部账号组（不向用户暴露 limit/offset）。
- 字段映射严格遵循 vendor 中 `Aggregator` struct 与 `ListAggregatorsResponseParams` 的字段定义。
- 提供单元测试（gomonkey mock 云 API）。

**Non-Goals:**
- 不实现账号组的增删改（无 CRUD 写接口资源）。
- 不暴露云 API 的 `Limit`/`Offset` 参数给用户。

## Decisions

### 1. 字段映射（基于 vendor 校验）

`ListAggregatorsRequest` 入参（vendor models.go:3468-3474）：
- `Limit *uint64`（每页数量）
- `Offset *uint64`（起始偏移）

`ListAggregatorsResponseParams` 出参（vendor models.go:3507-3516）：
- `Total *uint64` → schema `total`（TypeInt，Computed）
- `Items []*Aggregator` → schema `items`（TypeList，Computed，Elem 为 Resource）

`Aggregator` struct 字段（vendor models.go:654-682）→ schema `items` 子字段映射：
| 云API字段(JsonPath) | Go 字段 | 类型 | SchemaName |
|---|---|---|---|
| `response.Response.Total` | `Total *uint64` | uint64 | `total` |
| `response.Response.Items.Name` | `Name *string` | string | `name` |
| `response.Response.Items.Description` | `Description *string` | string | `description` |
| `response.Response.Items.OwnerUin` | `OwnerUin *uint64` | uint64 | `owner_uin` |
| `response.Response.Items.CreateTime` | `CreateTime *string` | string | `create_time` |
| `response.Response.Items.AccountCount` | `AccountCount *uint64` | uint64 | `account_count` |
| `response.Response.Items.Type` | `Type *string` | string | `type` |
| `response.Response.Items.AccountGroupId` | `AccountGroupId *string` | string | `account_group_id` |
| `response.Response.Items.AggregatorStatus` | `AggregatorStatus *uint64` | uint64 | `aggregator_status` |
| `response.Response.Items.MemberName` | `MemberName *string`（可为 null） | string | `member_name` |

说明：本次需求 JsonPath 列出了 `response.Response.Items` → `items`，按本项目数据源规则（规则13），列表字段展开为顶层 `items` 子结构，每个 `Aggregator` 元素的字段平铺到 `items` 的 Elem Schema 中。`total` 作为顶层 Computed 字段。

### 2. 服务层分页方法

在 `service_tencentcloud_config.go` 新增 `DescribeConfigListAggregatorsByFilter(ctx, paramMap) (items []*configv20220802.Aggregator, total uint64, errRet error)`，复用 `DescribeConfigCompliancePacksByFilter` 的分页骨架：
- `request.Limit = helper.Uint64(100)`、`request.Offset = helper.Uint64(0)`
- for 循环：retry 调用 `ListAggregatorsWithContext`，累加 `pageItems`，当 `len(items) >= total` 时 break，否则 `*request.Offset += *request.Limit`
- 返回 items 与 total（数据源需设置 `total` schema）

### 3. 数据源 Read 方法

参照 `data_source_tc_config_compliance_packs.go` 与项目数据源规则（规则14）：
- retry 块内调用 `DescribeConfigListAggregatorsByFilter`，**必须**检查返回是否为空（`items == nil` 且 `total == 0`），若为空**直接返回 `NonRetryableError`**，不 `d.SetId("")`。
- 外层 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")`。
- flatten 后 `_ = d.Set("items", itemsList)`、`_ = d.Set("total", int(total))`（设置前判断非 nil）。
- `d.SetId(helper.BuildToken())`，可选 `result_output_file` 写文件。

### 4. provider.go 注册

在 config 数据源区块（约 line 1422-1428）新增：
```go
"tencentcloud_config_list_aggregators": config.DataSourceTencentCloudConfigListAggregators(),
```

### 5. 单元测试（gomonkey mock）

按项目要求，新增数据源使用 gomonkey mock 云 API（不使用 TF 测试套件）。参考 `tencentcloud/services/config` 下其他数据源的 mock 测试风格，mock `UseConfigV20220802Client().ListAggregatorsWithContext`，验证 flatten 与 set 逻辑。

## Risks / Trade-offs

- **[云 API 返回空导致 state 清空]** → 按规则14，数据源 Read retry 块内返回 `NonRetryableError`，不在 retry 内 `d.SetId("")`，避免 API 波动清空 state。
- **[vendor 字段可为 null]** → `MemberName` 等字段可能为 null，flatten 时逐一判 nil 后再 set，避免 panic。
- **[分页总量大]** → 内部自动分页，单页 Limit 取 100（云 API 注释未标注明确最大值，参照同产品数据源取 100）。
