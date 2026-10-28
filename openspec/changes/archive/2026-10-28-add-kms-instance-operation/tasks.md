## 1. 资源代码实现

- [x] 1.1 创建 `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation.go`，定义 `ResourceTencentCloudKmsInstanceOperation()` 的 schema：
  - `instance_charge_type`（TypeString, Required, ForceNew）
  - `product_version`（TypeString, Optional, ForceNew）
  - `instance_charge_prepaid`（TypeList, Optional, ForceNew, MaxItems=1, elem=Resource，含 `period` TypeInt Required、`renew_flag` TypeInt Required）
  - `client_token`（TypeString, Optional, ForceNew）
  - `tags`（TypeList, Optional, ForceNew, elem=Resource，含 `tag_key` TypeString Required、`tag_value` TypeString Required）
  - `instance_id`（TypeString, Computed）
  - `deal_name`（TypeString, Computed）
  - `bill_id`（TypeString, Computed）
  - 仅注册 Create/Read/Delete，无 Update
- [x] 1.2 实现 `resourceTencentCloudKmsInstanceOperationCreate`：
  - 构造 `kms.NewCreateKmsInstanceRequest()`，从 schema 读取各字段填充 request（instance_charge_type、product_version、client_token、instance_charge_prepaid block、tags 列表）
  - period/renew_flag 用 `helper.IntUint64()` 转换为 `*uint64`
  - 用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `UseKmsClient().CreateKmsInstanceWithContext(ctx, request)` 调用，retry 块内只调用接口，失败用 `tccommon.RetryError(e)` 包装
  - retry 块外检查返回的 Response/InstanceId 是否为空（Response nil / InstanceId nil / 空串），为空则打印 logId 并返回 `resource.NonRetryableError`
  - 成功后 `d.SetId(instanceId)`，并 set computed 字段 instance_id、deal_name、bill_id（set 前判空）
- [x] 1.3 实现 `resourceTencentCloudKmsInstanceOperationRead` 与 `resourceTencentCloudKmsInstanceOperationDelete` 为空实现（仅 defer LogElapsed/InconsistentCheck，return nil）

## 2. 单元测试

- [x] 2.1 创建 `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation_test.go`，使用 gomonkey mock `UseKmsClient` 与 `CreateKmsInstanceWithContext`，测试 create 成功路径（验证 id 设置为返回的 InstanceId）及空响应返回 NonRetryableError 路径，不调用真实云 API，不使用 terraform 测试套件

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/kms/resource_tc_kms_kms_instance_operation.md`：
  - 一句话描述带云产品名 KMS（Provides a resource to create a kms kms_instance operation）
  - Example Usage（hcl 示例，含 instance_charge_type、product_version、instance_charge_prepaid block、tags）
  - 无 Import 部分（OPERATION 类型）
  - 不手写 Argument Reference / Attribute Reference

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 的 KMS 资源区块（`tencentcloud_kms_overwrite_white_box_device_fingerprints` 之后）追加 `"tencentcloud_kms_instance_operation": kms.ResourceTencentCloudKmsInstanceOperation(),`
- [x] 4.2 在 `tencentcloud/provider.md` 的 Key Management Service(KMS) Resource 区块追加 `tencentcloud_kms_instance_operation`（最终由 `make doc` 生成，本步仅追加 provider.md 列表条目）

## 5. 代码正确性检查

- [x] 5.1 核对 CreateKmsInstance 入参（InstanceChargeType、ProductVersion、InstanceChargePrepaid、ClientToken、Tags）与 schema 字段一一对应，确认云 API 支持所有映射参数
- [x] 5.2 核对出参 InstanceId、DealName、BillId 与 schema computed 字段一一对应
- [x] 5.3 确认所有函数返回的 error 已检查；无未使用变量；import 包名与现有 kms 资源一致