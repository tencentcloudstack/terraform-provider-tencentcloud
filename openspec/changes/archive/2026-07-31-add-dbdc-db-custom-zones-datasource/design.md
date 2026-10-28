## Context

需要为 DBDC（数据库专用集群）服务新增一个只读数据源 `tencentcloud_dbdc_db_custom_zones`，用于查询当前地域下 DBDC 支持的可用区列表及其售卖状态。该数据源调用腾讯云 DBDC `DescribeDBCustomZones` API。

当前 Provider 的 DBDC 服务（`tencentcloud/services/dbdc/`）已有 4 个数据源：`dbdc_db_custom_clusters`、`dbdc_db_custom_nodes`、`dbdc_db_custom_cluster_nodes`、`dbdc_db_custom_images`，以及 `resource_tc_dbdc_db_custom_cluster`、`resource_tc_dbdc_db_custom_node`、`resource_tc_dbdc_node_to_db_custom_cluster_attachment` 等资源。用户在创建 DBDC 集群和节点时需要知道哪些可用区处于 SELL（正常售卖）状态，但目前缺少此查询能力，可能导致选用 SOLD_OUT（售罄）可用区而创建失败。

云 API 现状（已核实 vendor 中 SDK）：
- `DescribeDBCustomZones`（`dbdc/v20201029`）请求无任何参数
- 响应 `response.Response.ZoneSet` 为 `[]*ZoneInfo`
- `ZoneInfo` 含两个字段：`Zone`（string，可用区）、`ZoneState`（string，枚举 SELL/SOLD_OUT）
- 接口无分页参数，一次返回全部可用区

本次新增为独立的只读数据源，不修改任何现有资源/数据源的 schema，向后完全兼容，风险较低。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_dbdc_db_custom_zones` 数据源，调用 `DescribeDBCustomZones` API 查询可用区列表
- 在 `zone_set` 列表中平铺每个可用区的 `zone` 和 `zone_state` 字段
- 在 service 层封装 `DescribeDBCustomZonesByFilter` 方法，复用 DBDC service 客户端，并遵循现有 DBDC 数据源（如 `DescribeDBCustomImagesByFilter`）的实现模式
- 在 retry 块内对返回空做严格检查（返回 `NonRetryableError`，不清空 id）
- 在 `provider.go` 与 `provider.md` 中注册数据源
- 提供文档模板与单元测试（gomonkey mock 云 API）

**Non-Goals:**
- 不支持任何入参过滤（API 本身无入参，仅返回当前地域全部可用区）
- 不实现分页（API 无分页参数，一次返回全部）
- 不修改现有 DBDC 资源/数据源的 schema 或行为

## Decisions

### 1. Schema 设计

数据源仅含可选的 `result_output_file` 与 Computed 的 `zone_set` 列表：

- `result_output_file` (Optional, String)：用于保存结果，沿用 Provider 数据源通用约定
- `zone_set` (Computed, List)：
  - `zone` (Computed, String)：可用区，如 `ap-guangzhou-3`
  - `zone_state` (Computed, String)：可用区状态，枚举 `SELL`（正常售卖）/ `SOLD_OUT`（售罄）

**理由：** 与现有 `data_source_tc_dbdc_db_custom_images.go`（`image_set` 列表展开字段）的 schema 风格保持一致；按代码生成要求，列表型返回结果展开为字段平铺，不额外再包一层"列表"schema。无入参对应 API 无入参的事实。

### 2. Service 层封装

在 `service_tencentcloud_dbdc.go` 新增 `DescribeDBCustomZonesByFilter` 方法，签名为：
```go
func (me *DbdcService) DescribeDBCustomZonesByFilter(ctx context.Context, param map[string]interface{}) (ret []*dbdcv20201029.ZoneInfo, totalCount int64, errRet error)
```

实现要点（参考现有 `DescribeDBCustomImagesByFilter`）：
- 使用 `dbdcv20201029.NewDescribeDBCustomZonesRequest()` 构造请求
- 无入参映射（API 无参数），`param` 仅用于保持签名一致
- 单次调用即可返回全部数据，无需分页循环
- 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 `me.client.UseDbdcV20201029Client().DescribeDBCustomZones(request)`
- 对返回空严格检查：`result == nil || result.Response == nil || result.Response.ZoneSet == nil` 时，打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 并返回 `resource.NonRetryableError`，**不**调用 `d.SetId("")`
- 错误统一用 `tccommon.RetryError(e)` 包装

**理由：** 与现有 DBDC 数据源 service 方法实现模式一致；无分页是因为 API 无分页参数。保留 `param map[string]interface{}` 与 `totalCount int64` 返回值以与同类方法签名对齐，便于未来扩展。

### 3. 数据源 Read 方法

`dataSourceTencentCloudDbdcDbCustomZonesRead` 实现要点：
- `defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_zones.read")()`
- `defer tccommon.InconsistentCheck(d, meta)()`
- 通过 `service.DescribeDBCustomZonesByFilter(ctx, paramMap)` 获取数据
- 外层用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 service 调用（service 内已有 retry 时外层不再二次 retry；按现有 `data_source_tc_dbdc_db_custom_images.go` 模式，外层保留 retry 块并将错误交由 `tccommon.RetryError(e)` 处理）
- 遍历 `respData`，对每个 `ZoneInfo` 的 `Zone`/`ZoneState` 做 nil 判断后再写入 map（遵循"先判 nil 再 set"要求）
- `_ = d.Set("zone_set", zoneSetList)`
- `d.SetId(helper.BuildToken())`
- 支持 `result_output_file` 写文件

**理由：** 严格遵循代码生成要求中 DATASOURCE 资源 Read 方法的空返回检查与字段 nil 判断规则。

### 4. 错误处理

- service 层：API 调用失败用 `tccommon.RetryError(e)` 包装；返回空用 `NonRetryableError`
- 数据源层：外层 retry 失败路径直接返回 `reqErr`，不调用 `d.SetId("")`，避免短暂波动清空 state id

### 5. 测试策略

按代码生成要求，新增的 terraform 资源使用 mock（gomonkey）对云 API 进行 mock 处理，只做业务代码逻辑的单元测试：
- mock `DbdcService.DescribeDBCustomZonesByFilter` 返回构造的 `[]*dbdcv20201029.ZoneInfo`
- 验证 `zone_set` 的字段映射正确（zone / zone_state）
- 验证空返回场景下 Read 方法行为
- **禁止**通过 `go test` 实际执行测试用例，仅保证代码可正确构建

### 6. 文档

`data_source_tc_dbdc_db_custom_zones.md` 内容：
- 一句话描述并带云产品名称：`Use this data source to query available zones and their sale status from the TencentCloud DBDC product.`
- Example Usage（`data "tencentcloud_dbdc_db_custom_zones" "example" {}`）
- 不手写 Argument Reference / Attribute Reference（由工具自动生成）
- 无 Import 部分（DATASOURCE 资源没有 Import）

## Risks / Trade-offs

**[风险] API 返回的 `ZoneState` 字段可能新增枚举值** → schema 中 `zone_state` 定义为 String，透传 API 返回值，天然兼容新枚举；文档中列出当前已知枚举（SELL/SOLD_OUT）

**[风险] 云 API 短暂波动导致查询失败** → 通过 retry 机制重试；返回空时返回 `NonRetryableError` 让外层 retry 继续尝试，避免误清 state id

**[取舍] 不暴露任何过滤参数给用户** → API 本身无入参，无法过滤；保持数据源简洁，与 API 行为一致

**[取舍] 保留 service 方法 `param map[string]interface{}` 与 `totalCount int64` 返回值签名** → 与现有同类 DBDC 数据源 service 方法签名对齐，便于维护与未来扩展，即使当前不使用
