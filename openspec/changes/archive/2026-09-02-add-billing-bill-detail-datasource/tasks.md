## 1. 准备工作

- [x] 1.1 确认 vendor 目录下 billing SDK（v20180709）已包含 `DescribeBillDetail` 接口及 `BillDetail`、`BillDetailComponent`、`BillDetailComponentConfig`、`BillDetailAssociatedOrder`、`BillTagInfo` 等结构定义
- [x] 1.2 确认 `tencentcloud/services/billing/` 目录存在，并查看是否有既有 service 层文件 `service_tencentcloud_billing.go`

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/billing/service_tencentcloud_billing.go` 中新增（或补充）`DescribeBillingBillDetailByFilter` 方法
- [x] 2.2 方法内部封装 `DescribeBillDetail` 调用，将入参 map 转换为 `DescribeBillDetailRequest`（period_type、month、begin_time、end_time、need_record_num、product_code、pay_mode、resource_id、action_type、project_id、business_code、context、payer_uin）
- [x] 2.3 实现自动分页：使用 Offset 递增、Limit=300 循环调用，当某次返回条数 < Limit 时停止，合并所有 DetailSet（不依赖 Total 作为终止条件，避免 Total 缺失时提前中断丢数据）
- [x] 2.4 同时返回 Total、Context 顶层字段

## 3. Data Source Schema 与 Read 实现

- [x] 3.1 创建 `tencentcloud/services/billing/data_source_tc_billing_bill_detail.go`，定义 `DataSourceTencentCloudBillingBillDetail()` 返回 `*schema.Resource`
- [x] 3.2 定义可选输入参数 schema：period_type、month、begin_time、end_time、need_record_num(int)、product_code、pay_mode、resource_id、action_type、project_id(int)、business_code、context、payer_uin、result_output_file
- [x] 3.3 定义输出 `detail_set`（TypeList）schema，展开 BillDetail 全部字段；其中 component_set 为 TypeList（含 component_config 子 TypeList）、tags 为 TypeList（含 tag_key/tag_value）、associated_order 为 TypeList（含 6 个订单字段）、price_info 为 TypeList of String
- [x] 3.4 定义顶层输出 `total`(Int)、`context`(String)
- [x] 3.5 实现 `dataSourceTencentCloudBillingBillDetailRead`：defer LogElapsed/InconsistentCheck，组装 paramMap 调用 service 层
- [x] 3.6 在 `resource.Retry(tccommon.ReadRetryTimeout)` 内调用 service 方法，失败用 `tccommon.RetryError` 包装；`DetailSet` 为空（查询无结果）时正常返回空列表、不报错，`total=0`
- [x] 3.7 retry 块外遍历 DetailSet 填充 `detail_set` 列表（逐字段 nil 检查后再 set），设置 `total`、`context`，`d.SetId(helper.BuildToken())`，处理 `result_output_file` 输出

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 中添加 `"tencentcloud_billing_bill_detail": billing.DataSourceTencentCloudBillingBillDetail()` 注册项
- [x] 4.2 确认注册位置参考 `tencentcloud_igtm_strategy` 等既有资源的注册风格

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/billing/data_source_tc_billing_bill_detail.md`
- [x] 5.2 一句话描述：`Use this data source to query ...`（带上云产品 billing）
- [x] 5.3 Example Usage：按 month 查询、按 begin_time/end_time 查询两个示例
- [x] 5.4 不添加 Argument Reference / Attribute Reference（由工具自动生成）
- [x] 5.5 不添加 Import 部分（DATASOURCE 类型无 Import）

## 6. 测试

- [x] 6.1 创建 `tencentcloud/services/billing/data_source_tc_billing_bill_detail_test.go`
- [x] 6.2 使用 gomonkey mock `DescribeBillDetail` 接口，测试 Read 业务逻辑（正常返回、空返回、分页合并）
- [x] 6.3 不使用 terraform 测试套件，只测试业务代码逻辑

## 7. 验证

- [x] 7.1 运行 `openspec validate add-billing-bill-detail-datasource --strict` 验证提案完整性
- [x] 7.2 检查所有新增文件存在且格式正确
- [x] 7.3 收尾阶段通过 tfpacer-finalize skill 执行 gofmt、make doc、生成 changelog
