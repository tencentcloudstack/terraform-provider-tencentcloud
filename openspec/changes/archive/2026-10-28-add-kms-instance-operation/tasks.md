## 1. Framework Action 代码实现

- [x] 1.1 创建 `tencentcloud/services/kms/action_tc_kms_instance.go`，实现 `action.Action` + `action.ActionWithConfigure`：
  - Metadata 设置 TypeName 为 `tencentcloud_kms_instance`
  - Schema 定义：instance_charge_type（StringAttribute Required）、product_version（StringAttribute Optional）为顶层 Attributes（不暴露 client_token）；instance_charge_prepaid（SingleNestedBlock，含 period Int64Attribute Required、renew_flag Int64Attribute Required）和 tags（ListNestedBlock，含 tag_key StringAttribute Required、tag_value StringAttribute Required）为 Blocks
  - Invoke 方法：从 config 读取 model，校验 instance_charge_type 非空，校验 client 非空，通过 `clientToken := helper.BuildToken()` 生成 client_token，构造 `kms.CreateKmsInstanceRequest`，调用 service 方法，校验响应非空且 instance_id 非空后日志输出 instance_id/deal_name/bill_id
- [x] 1.2 在既有 `tencentcloud/services/kms/service_tencentcloud_kms.go` 中新增 `KmsService.CreateKmsInstance`（不新增额外 service 文件）：
  - 用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `UseKmsClient().CreateKmsInstanceWithContext(ctx, request)` 调用，retry 块内只调用接口，失败用 `tccommon.RetryError(e)` 包装
  - 检查返回的 Response 是否为空，为空则返回 error
  - period/renew_flag 用 `helper.Int64Uint64()` 转换为 `*uint64`

## 2. 单元测试

- [x] 2.1 创建 `tencentcloud/services/kms/action_tc_kms_instance_test.go`，使用 gomonkey mock `UseKmsClient` 与 `CreateKmsInstanceWithContext`，通过 tftypes 原始值构建 InvokeRequest：
  - 测试 Invoke 成功路径（验证 request 字段填充正确）
  - 测试 API 错误路径
  - 测试空响应路径（含 InstanceId 为空）
  - 测试缺少必填输入路径
  - 测试 client 未配置路径
  - 测试 Metadata（TypeName 为 `tencentcloud_kms_instance`）与 Schema 定义

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/kms/action_tc_kms_instance.md`：
  - 一句话描述带云产品名 KMS
  - ~> NOTE 提示需 Terraform 1.14+
  - Example Usage（action 语法，含 config block、instance_charge_prepaid block、tags block），action 名称为 `tencentcloud_kms_instance`
  - 不手写 Argument Reference / Attribute Reference（由 make doc 生成）
- [x] 3.2 在 `tencentcloud/provider.md` 的 Key Management Service(KMS) 区块新增 `Action` 子区段，追加 `tencentcloud_kms_instance`
- [x] 3.3 `website/tencentcloud.erb` 与 `website/docs/actions/kms_instance.html.markdown` 由收尾阶段 `make doc` 生成（不手工编辑 website 目录）

## 4. Registry 注册

- [x] 4.1 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 切片追加 `kms.NewKmsInstance`，并添加 kms 包 import
- [x] 4.2 确认不在 `tencentcloud/provider.go` 的 ResourcesMap 中注册（action 由 framework provider 管理）

## 5. Changelog

- [x] 5.1 更新 `.changelog/4616.txt`，new-resource 条目为 `tencentcloud_kms_instance`

## 6. 代码正确性检查

- [x] 6.1 核对 CreateKmsInstance 入参（InstanceChargeType、ProductVersion、InstanceChargePrepaid、ClientToken、Tags）与 action schema 字段的映射关系，确认云 API 支持所有映射参数（ClientToken 由 provider 内部生成，不对应 schema 字段）
- [x] 6.2 核对出参 InstanceId、DealName、BillId 通过日志输出
- [x] 6.3 确认所有函数返回的 error 已检查；无未使用变量；import 包名与现有 action 一致

## 7. 命名对齐（framework action 约定）

- [x] 7.1 action 名称由 `tencentcloud_kms_instance_operation` 调整为 `tencentcloud_kms_instance`，与仓库内既有 framework action 命名（`tencentcloud_antiddos_unblock_resources`、`tencentcloud_mongodb_restore_db_instance`、`tencentcloud_dlc_initialize_tc_lake`）保持一致
- [x] 7.2 重命名代码文件：`action_tc_kms_instance_operation.go/.md/_test.go` → `action_tc_kms_instance.go/.md/_test.go`
- [x] 7.3 重命名 Go 标识符：`KmsInstanceOperation` → `KmsInstance`、`KmsInstanceOperationModel` → `KmsInstanceModel`、`NewKmsInstanceOperation` → `NewKmsInstance`，并同步 `registry.go`、`provider.md`、单测断言与 `.changelog`
- [x] 7.4 删除非约定的 `tencentcloud/services/kms/service_tc_kms_instance.go`，将 `CreateKmsInstance` 方法合并进 `service_tencentcloud_kms.go`
- [x] 7.5 同步更新 proposal.md、design.md、spec delta 中的 action 名称、文件名与 service 位置描述
- [x] 7.6 删除重命名前 `make doc` 遗留的 `website/docs/actions/kms_instance_operation.html.markdown`，避免旧 action 名的孤儿文档页残留（新页面 `kms_instance.html.markdown` 由收尾阶段 `make doc` 重新生成）

## 8. 入参赋值细化

- [x] 8.1 `product_version` 仅在非 null / 非 unknown / 非空字符串时写入 request，避免显式空字符串透传给云 API
- [x] 8.2 `tags` 中 `tag_key` 为空（null / unknown / 空串）的条目跳过并打印 `[WARN]` 日志，不作为请求标签下发
- [x] 8.3 新增单测 `TestKmsInstance_Invoke_MinimalInput`，验证仅配置必填项时 `ProductVersion`、`InstanceChargePrepaid`、`Tags` 均未下发，且 `ClientToken` 由 provider 内部生成非空
- [x] 8.4 修正 design.md 中 `helper.Int64Uint64` 函数名与 KMS client 调用示意代码，使其与实现一致

## 9. 不暴露 client_token 属性

- [x] 9.1 从 `action_tc_kms_instance.go` 的 Schema 中删除 `client_token` StringAttribute，并从 `KmsInstanceModel` 中删除 `ClientToken` 字段
- [x] 9.2 Invoke 中改为 `clientToken := helper.BuildToken()` 生成可重入客户端 ID，并始终 `request.ClientToken = helper.String(clientToken)` 下发，不再读取用户输入
- [x] 9.3 同步删除 `action_tc_kms_instance.md` 示例中的 `client_token` 配置项
- [x] 9.4 更新单测：请求构造不再包含 `client_token`；成功路径与最小入参路径断言 `ClientToken` 非空；Schema 断言 `client_token` 不存在
- [x] 9.5 同步更新 proposal.md、design.md、spec delta 中关于 client_token 的描述（不对外暴露、内部自动生成）
