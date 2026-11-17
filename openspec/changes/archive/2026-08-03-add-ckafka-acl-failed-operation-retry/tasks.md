## 1. 常量定义

- [x] 1.1 在 `tencentcloud/services/ckafka/extension_ckafka.go` 中新增重试常量：`CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_TIMES = 3` 与 `CKAFKA_CREATE_ACL_FAILED_OPERATION_RETRY_INTERVAL = 5 * time.Second`，并补充 `time` 包的 import

## 2. 核心逻辑修改

- [x] 2.1 修改 `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go` 的 `CkafkaService.CreateAcl` 方法：将原有 `resource.Retry` 调用包裹进「首次调用 + 最多 3 次重试」的 `for` 循环，当错误为 `Code == CkafkaFailedOperation` 时，每次重试前 `time.Sleep(5 * time.Second)`，并输出重试日志
- [x] 2.2 保持循环内 `resource.Retry` + `tccommon.RetryError` 对非 FailedOperation 错误的原有处理不变；非 FailedOperation 错误直接跳出循环返回
- [x] 2.3 确认 `CreateAcl` 成功路径（`OperateStatusCheck` 等）与返回逻辑保持与修改前一致

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/ckafka/resource_tc_ckafka_acl_test.go` 中新增 gomonkey 单元测试：mock `CreateAcl` 云 API 连续返回 3 次 `FailedOperation` 后第 4 次成功，断言 `CreateAcl` 返回 nil 且 API 被调用 4 次
- [x] 3.2 新增用例：mock `CreateAcl` 始终返回 `FailedOperation`，断言最终返回 `FailedOperation` 错误且 API 被调用 4 次（1 次首次 + 3 次重试）
- [x] 3.3 新增用例：mock `CreateAcl` 返回非 `FailedOperation` 错误，断言立即返回该错误且 API 仅被调用 1 次（不触发固定次数重试）

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/ckafka/resource_tc_ckafka_acl.md`，在文档中说明创建 ACL 时对 `FailedOperation` 错误最多重试 3 次、每次间隔 5 秒的行为（保持示例与 Import 段落不变）

## 5. 验证

- [x] 5.1 确认修改后的代码可编译（不执行 `go build`，由后续流程统一校验）
- [x] 5.2 检查单元测试文件 `resource_tc_ckafka_acl_test.go` 中新增用例与既有 `TestAccTencentCloudCkafkaAclResource` 等用例可共存、import 完整
