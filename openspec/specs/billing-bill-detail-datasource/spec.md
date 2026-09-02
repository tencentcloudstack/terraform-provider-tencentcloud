# billing-bill-detail-datasource Specification

## Purpose
定义 `tencentcloud_billing_bill_detail` 数据源，用于通过 Terraform 查询腾讯云计费账单明细（DescribeBillDetail），支持按月份或时间区间查询、多条件过滤、自动分页及结果输出到文件。

## Requirements
### Requirement: Data Source Schema Definition

数据源 `tencentcloud_billing_bill_detail` SHALL 支持以下可选输入参数，并 MUST 返回以下输出属性。

**Input Parameters (All Optional):**
- `period_type` (String): 周期类型，byUsedTime 按计费周期 / byPayTime 按扣费周期（已废弃，保留兼容）
- `month` (String): 月份，格式 yyyy-mm；与 begin_time/end_time 二选一
- `begin_time` (String): 周期开始时间，格式 yyyy-mm-dd hh:ii:ss；与 end_time 成对出现，且必须同月
- `end_time` (String): 周期结束时间，格式 yyyy-mm-dd hh:ii:ss
- `need_record_num` (Int): 是否返回总记录数，1=需要，0=不需要
- `product_code` (String): 子产品编码（已废弃，保留兼容）
- `pay_mode` (String): 付费模式，prePay 包年包月 / postPay 按量计费
- `resource_id` (String): 资源 ID
- `action_type` (String): 交易类型名称
- `project_id` (Int): 项目 ID
- `business_code` (String): 产品名称代码
- `context` (String): 上下文信息，用于翻页加速
- `payer_uin` (String): 支付者账号 ID
- `result_output_file` (String): 结果输出文件路径

**Output Attributes:**
- `detail_set` (List): 账单明细列表，每个元素包含:
  - `business_code_name` (String): 产品名称
  - `product_code_name` (String): 子产品名称
  - `pay_mode_name` (String): 计费模式
  - `project_name` (String): 项目名称
  - `region_name` (String): 地域
  - `zone_name` (String): 可用区
  - `resource_id` (String): 资源 ID
  - `resource_name` (String): 资源别名
  - `action_type_name` (String): 交易类型
  - `order_id` (String): 订单 ID
  - `bill_id` (String): 交易 ID
  - `pay_time` (String): 扣费时间
  - `fee_begin_time` (String): 开始使用时间
  - `fee_end_time` (String): 结束使用时间
  - `component_set` (List): 组件列表，每个元素包含 component_code_name、item_code_name、single_price、specified_price、price_unit、used_amount、used_amount_unit、real_total_measure、deducted_measure、time_span、time_unit_name、cost、discount、reduce_type、real_cost、voucher_pay_amount、cash_pay_amount、incentive_pay_amount、transfer_pay_amount、item_code、component_code、contract_price、instance_type、ri_time_span、original_cost_with_ri、sp_deduction_rate、sp_deduction、original_cost_with_sp、blended_discount 字段，以及 `component_config`（List，元素含 name、value）
  - `payer_uin` (String): 支付者 UIN
  - `owner_uin` (String): 使用者 UIN
  - `operate_uin` (String): 操作者 UIN
  - `tags` (List): 标签列表，元素含 tag_key、tag_value
  - `business_code` (String): 产品编码
  - `product_code` (String): 子产品编码
  - `action_type` (String): 交易类型编码
  - `region_id` (String): 地域 ID
  - `project_id` (Int): 项目 ID
  - `price_info` (List of String): 价格属性
  - `associated_order` (List): 关联交易单据，元素含 prepay_purchase、prepay_renew、prepay_modify_up、reverse_order、new_order、original
  - `formula` (String): 计算说明
  - `formula_url` (String): 计费规则链接
  - `bill_day` (String): 账单归属日
  - `bill_month` (String): 账单归属月
  - `id` (String): 账单记录 ID
  - `region_type` (String): 国内国际编码
  - `region_type_name` (String): 国内国际名称
  - `reserve_detail` (String): 备注属性
  - `discount_object` (String): 优惠对象
  - `discount_type` (String): 优惠类型
  - `discount_content` (String): 优惠内容
- `total` (Int): 总记录数
- `context` (String): 本次请求的上下文信息

#### Scenario: Query bill detail by month

- **WHEN** 用户配置 `month = "2024-01"` 并调用数据源
- **THEN** 数据源调用 `DescribeBillDetail` 以该月份查询，返回该月所有账单明细记录到 `detail_set`，并返回 `total` 总记录数

#### Scenario: Query bill detail by time range

- **WHEN** 用户配置 `begin_time` 和 `end_time` 为同月的起止时间
- **THEN** 数据源以 begin_time/end_time 作为查询条件调用 `DescribeBillDetail`，返回该时间区间（整月）的账单明细

#### Scenario: Query bill detail with filters

- **WHEN** 用户同时配置 `resource_id`、`pay_mode`、`business_code` 等过滤条件
- **THEN** 数据源将所有非空过滤条件透传给 `DescribeBillDetail`，返回满足全部条件的账单明细子集

#### Scenario: Automatic pagination

- **WHEN** 账单明细总记录数超过单次 API 最大返回量（300）
- **THEN** 数据源内部自动以 Offset 递增分页循环调用 `DescribeBillDetail`，直到累计获取所有记录，用户无需关心分页参数

#### Scenario: Empty result

- **WHEN** 查询条件下没有账单明细数据（API 返回空 DetailSet）
- **THEN** 数据源返回 `NonRetryableError`，不清空 id，上层任务以「重试耗尽」形式失败，便于人工介入排障

#### Scenario: Output to file

- **WHEN** 用户配置 `result_output_file` 参数
- **THEN** 数据源在读取完成后将结果写入指定文件路径

### Requirement: Provider Registration

Provider SHALL 注册 `tencentcloud_billing_bill_detail` 数据源，使其 MUST 可通过 Terraform 配置引用。

#### Scenario: Data source registered in provider

- **WHEN** Provider 初始化
- **THEN** `tencentcloud/provider.go` 中包含 `"tencentcloud_billing_bill_detail": billing.DataSourceTencentCloudBillingBillDetail()` 注册项
- **AND** 运行 `make doc` 后 `tencentcloud/provider.md` 中包含该数据源声明

### Requirement: API Retry and Error Handling

数据源 Read 方法 SHALL 对云 API 调用进行重试包装，并 MUST 正确处理错误与空响应。

#### Scenario: API transient failure retried

- **WHEN** `DescribeBillDetail` 调用因瞬时错误失败
- **THEN** 使用 `tccommon.RetryError` 包装错误，外层 `resource.Retry(tccommon.ReadRetryTimeout)` 自动重试

#### Scenario: API success with empty response

- **WHEN** `DescribeBillDetail` 返回的 `Response` 为 nil 或 `DetailSet` 为空
- **THEN** Read 方法返回 `NonRetryableError`，保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 日志，不调用 `d.SetId("")`

