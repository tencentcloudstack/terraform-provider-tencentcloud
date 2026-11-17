## Why

用户需要通过 Terraform 查询腾讯云账单明细数据（L3-明细账单）。当前 Provider 的 billing 产品缺少账单明细数据源，用户无法在 Terraform 配置中查询指定月份/时间段的账单明细，也无法根据资源 ID、产品、交易类型等条件过滤账单明细记录，导致自动化成本分析、账单核对场景缺乏数据支撑。

## What Changes

- 新增 Data Source: `tencentcloud_billing_bill_detail`
- 实现对 Billing API `DescribeBillDetail`（billing/v20180709）接口的调用，获取账单明细列表
- 支持以下查询入参（均为可选）:
  - `month`: 月份，格式 yyyy-mm
  - `begin_time` / `end_time`: 周期起止时间（与 month 二选一）
  - `need_record_num`: 是否返回总记录数
  - `pay_mode`: 付费模式（prePay/postPay）
  - `resource_id`: 资源 ID
  - `action_type`: 交易类型名称
  - `project_id`: 项目 ID
  - `business_code`: 产品名称代码
  - `context`: 上下文信息（用于翻页加速）
  - `payer_uin`: 支付者账号 ID
- 返回账单明细列表 `detail_set`，每条记录包含产品、资源、交易、组件明细、标签、关联订单、折扣优惠等完整字段
- 内部实现自动分页获取所有数据，不暴露 limit/offset 给用户

## Capabilities

### New Capabilities
- `billing-bill-detail-datasource`: 通过 `DescribeBillDetail` 接口查询腾讯云账单明细列表数据源

### Modified Capabilities
<!-- 无需修改已有能力 -->

## Impact

- **新增能力**: 账单明细列表查询（只读）
- **受影响的服务**: Billing (tencentcloud/services/billing)
- **新增文件**:
  - `tencentcloud/services/billing/data_source_tc_billing_bill_detail.go`
  - `tencentcloud/services/billing/data_source_tc_billing_bill_detail.md`
  - `tencentcloud/services/billing/data_source_tc_billing_bill_detail_test.go`
  - Provider 注册代码需新增此 data source
  - provider.md 需新增此 data source 声明（由 make doc 生成）
- **API 依赖**:
  - Billing API v20180709: `DescribeBillDetail`
  - 文档: https://cloud.tencent.com/document/product/555/35762
- **兼容性**: 无破坏性变更，纯新增只读数据源功能
