## Context

腾讯云 Billing 产品的账单明细（L3-明细账单）数据目前无法通过 Terraform Provider 查询。用户在自动化成本分析、账单核对场景下，需要在 Terraform 配置中引用账单明细数据。

云 API `DescribeBillDetail`（billing/v20180709）支持按月份/时间区间、产品、资源、交易类型等条件查询账单明细列表，单次最多返回 300 条，通过 Offset 分页或 Context 上下文翻页。

当前 Provider 的 `tencentcloud/services/billing/` 目录下尚无 data source 文件，本次为该目录首个数据源。

参考资源: `tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`（数据源业务逻辑模板）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_billing_bill_detail` 数据源，支持通过 `DescribeBillDetail` 查询账单明细列表
- 内部自动分页获取所有明细数据，不暴露 limit/offset 给用户
- 返回完整的明细字段（含组件明细、关联订单、标签等嵌套结构）
- 在 Provider 中注册该数据源

**Non-Goals:**
- 不实现账单明细的创建/更新/删除（接口本身为只读查询）
- 不处理账单下载链接（L0/L1/L2 账单由其他接口/资源覆盖）
- 不暴露分页参数给用户（内部封装自动翻页）
- 不修改其他已有 billing 资源的行为

## Decisions

### 决策 1: 数据源命名与文件组织
- 数据源名称: `tencentcloud_billing_bill_detail`
- 文件: `tencentcloud/services/billing/data_source_tc_billing_bill_detail.go`
- 文档: `tencentcloud/services/billing/data_source_tc_billing_bill_detail.md`
- 测试: `tencentcloud/services/billing/data_source_tc_billing_bill_detail_test.go`
- **理由**: 遵循 Provider 现有命名规范（`data_source_tc_<product>_<name>`），billing 作为一级产品目录。

### 决策 2: 内部自动分页
- 云 API 单次最大 Limit=300，使用 Offset 翻页
- 数据源不向用户暴露 limit/offset 参数
- 在 service 层封装 `DescribeBillingBillDetailByFilter` 方法，循环调用直到取完所有数据
- 同时支持返回的 Context 上下文字段（用于 Month>=2023-05 加速翻页），但作为只读结果返回，不要求用户传入
- **理由**: 遵循项目硬约束「数据源分页不暴露 limit/offset，内部自动分页获取所有数据」。

### 决策 3: Schema 结构（展开列表，不平铺）
- 顶层输出字段 `detail_set`（TypeList）展开 `response.Response.DetailSet` 数组
- `detail_set` 内每个元素包含 BillDetail 的所有字段（business_code_name、product_code_name、resource_id 等）
- `detail_set` 内的 `component_set`（TypeList）展开 `BillDetailComponent` 数组
- `component_set` 内的 `component_config`（TypeList）展开 `BillDetailComponentConfig` 数组
- `detail_set` 内的 `tags`（TypeList）展开 `BillTagInfo` 数组
- `detail_set` 内的 `associated_order`（单对象）展开 `BillDetailAssociatedOrder`
- `detail_set` 内的 `price_info`（TypeList of String）
- 顶层返回 `total`（总记录数）和 `context`（上下文）
- **理由**: 遵循项目规则「资源列表型数据展开，把列表中每个元素的所有参数平铺到资源参数 schema 顶层」，同时 DetailSet 本身就是列表，必须保留列表层级，但列表元素内的嵌套对象/数组按各自类型建模。

### 决策 4: 入参映射
- 所有入参均为 Optional（云 API 也均非必填，但 month 与 begin_time/end_time 二选一）
- `need_record_num` 为 int64 类型
- `project_id` 为 int64 类型
- **理由**: 完全遵循云 API 入参定义；`period_type`、`product_code` 已从官方文档移除且 SDK 标注 Deprecated/未开放，实测传参不生效，故从 schema 中移除。

### 决策 5: 错误处理与重试
- Read 方法使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 API 调用
- 失败时使用 `tccommon.RetryError(e)` 包装返回
- `Response` 为 nil（真正的异常）时返回 `NonRetryableError`；`DetailSet` 为空（查询无结果）时**正常返回空列表、不报错**（按 resource_id 或组合条件过滤查不到数据是常见合法场景，报错会让 plan 直接失败）
- retry 块外设置 schema 字段、SetId
- **理由**: 遵循项目数据源 Read 重试与空响应处理规则；空结果属于正常业务场景而非错误。

### 决策 6: 测试方式
- 新增资源使用 mock（gomonkey）对云 API 进行 mock，只测试业务代码逻辑
- 不使用 terraform 测试套件
- **理由**: 遵循项目规则「新增 terraform 资源的测试使用 mock，不使用 terraform 测试套件」。

## Risks / Trade-offs

- [账单明细数据量大] → 通过内部自动分页（每次 300 条循环）获取全部数据，大数据量场景下 Read 耗时较长，但这是只读数据源的固有特性；用户可通过 month/resource_id 等条件缩小范围
- [Context 翻页加速字段] → 本次作为只读结果返回顶层 `context` 字段，不要求用户循环传入，简化使用模型；如用户需要跨 plan 加速可自行引用该输出
- [部分入参已废弃（period_type/product_code）] → 保留为可选参数并标注 deprecated，避免破坏存量使用习惯，同时不主动推荐
- [嵌套结构层级较深（detail_set → component_set → component_config）] → 按云 API 实际结构建模为多层 TypeList，schema 定义较复杂但语义清晰
