## 1. Framework Action 代码实现

- [x] 1.1 创建 `tencentcloud/services/kms/action_tc_kms_instance_operation.go`，实现 `action.Action` + `action.ActionWithConfigure`：
  - Metadata 设置 TypeName 为 `tencentcloud_kms_instance_operation`
  - Schema 定义：instance_charge_type（StringAttribute Required）、product_version（StringAttribute Optional）、client_token（StringAttribute Optional）为顶层 Attributes；instance_charge_prepaid（SingleNestedBlock，含 period Int64Attribute Required、renew_flag Int64Attribute Required）和 tags（ListNestedBlock，含 tag_key StringAttribute Required、tag_value StringAttribute Required）为 Blocks
  - Invoke 方法：从 config 读取 model，校验 instance_charge_type 非空，校验 client 非空，构造 `kms.CreateKmsInstanceRequest`，调用 service 方法，日志输出 instance_id/deal_name/bill_id
- [x] 1.2 创建 `tencentcloud/services/kms/service_tc_kms_instance.go`，实现 `KmsService.CreateKmsInstance`：
  - 用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `UseKmsClient().CreateKmsInstanceWithContext(ctx, request)` 调用，retry 块内只调用接口，失败用 `tccommon.RetryError(e)` 包装
  - 检查返回的 Response 是否为空，为空则返回 error
  - period/renew_flag 用 `helper.IntUint64()` 转换为 `*uint64`

## 2. 单元测试

- [x] 2.1 创建 `tencentcloud/services/kms/action_tc_kms_instance_operation_test.go`，使用 gomonkey mock `UseKmsClient` 与 `CreateKmsInstanceWithContext`，通过 tftypes 原始值构建 InvokeRequest：
  - 测试 Invoke 成功路径（验证 request 字段填充正确）
  - 测试 API 错误路径
  - 测试空响应路径
  - 测试缺少必填输入路径
  - 测试 client 未配置路径
  - 测试 Metadata 与 Schema 定义

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/kms/action_tc_kms_instance_operation.md`：
  - 一句话描述带云产品名 KMS
  - ~> NOTE 提示需 Terraform 1.14+
  - Example Usage（action 语法，含 config block、instance_charge_prepaid block、tags block）
  - 不手写 Argument Reference / Attribute Reference（由 make doc 生成）
- [x] 3.2 在 `tencentcloud/provider.md` 的 Key Management Service(KMS) 区块新增 `Action` 子区段，追加 `tencentcloud_kms_instance_operation`
- [x] 3.3 更新 `website/tencentcloud.erb`，在 KMS 产品下新增 Actions 侧边栏区段
- [x] 3.4 创建 `website/docs/actions/kms_instance_operation.html.markdown`

## 4. Registry 注册

- [x] 4.1 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 切片追加 `kms.NewKmsInstanceOperation`，并添加 kms 包 import
- [x] 4.2 确认不在 `tencentcloud/provider.go` 的 ResourcesMap 中注册（action 由 framework provider 管理）

## 5. Changelog

- [x] 5.1 更新 `.changelog/4616.txt`，保留 new-resource 条目并补充 enhancement 条目说明改为 framework action 实现

## 6. 代码正确性检查

- [x] 6.1 核对 CreateKmsInstance 入参（InstanceChargeType、ProductVersion、InstanceChargePrepaid、ClientToken、Tags）与 action schema 字段一一对应，确认云 API 支持所有映射参数
- [x] 6.2 核对出参 InstanceId、DealName、BillId 通过日志输出
- [x] 6.3 确认所有函数返回的 error 已检查；无未使用变量；import 包名与现有 action 一致