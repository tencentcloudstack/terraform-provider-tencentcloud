## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/ssm/resource_tc_ssm_secret.go` 的 `ResourceTencentCloudSsmSecret()` Schema 中新增 `encrypt_type` 字段
  - Type: `schema.TypeInt`
  - Optional: true
  - Computed: true
  - ForceNew: true
  - Description: "Encryption type of secret. `0`: KMS key encryption. `1`: soft key encryption. Default is `0`. Once created, it cannot be modified."

## 2. Service 层

- [x] 2.1 在 `tencentcloud/services/ssm/service_tencentcloud_ssm.go` 的 `SecretInfo` 结构体中新增 `encryptType uint64` 字段
- [x] 2.2 在 `DescribeSecretByName` 方法中，对 `response.Response.EncryptType` 做 nil 判断后填充 `secretInfo.encryptType`（nil 时使用 `0`）

## 3. Create 操作

- [x] 3.1 在 `resourceTencentCloudSsmSecretCreate` 中，使用 `d.GetOk("encrypt_type")` 读取用户配置值，存在时设置 `request.EncryptType = helper.IntUint64(v.(int))`，放置在其他 request 字段赋值之后、`resource.Retry` 发送请求之前

## 4. Read 操作

- [x] 4.1 在 `resourceTencentCloudSsmSecretRead` 中，`secretInfo` 非空且字段设置前判断 `secretInfo.encryptType`，调用 `d.Set("encrypt_type", secretInfo.encryptType)` 设置状态值

## 5. Update 操作

- [x] 5.1 在 `resourceTencentCloudSsmSecretUpdate` 的 `immutableArgs` 数组中新增 `"encrypt_type"`（防御性措施，因 schema 已 ForceNew）

## 6. 单元测试

- [x] 6.1 在 `tencentcloud/services/ssm/resource_tc_ssm_secret_test.go` 中补充 `encrypt_type` 参数的单元测试用例（使用 gomonkey mock 云API，覆盖 create 传入 encrypt_type=0/1、read 回写、force-new 行为）

## 7. 文档

- [x] 7.1 在 `tencentcloud/services/ssm/resource_tc_ssm_secret.md` 中补充 `encrypt_type` 参数说明与示例（最终由收尾阶段 `make doc` 生成 website/docs 文档，本任务仅更新资源 .md 示例文件）

## 8. 验证

- [x] 8.1 检查所有 CRUD 接口参数与云API一致：CreateSecret 含 EncryptType、DescribeSecret 返回 EncryptType、UpdateSecret 不含 EncryptType（确认不可变）
- [x] 8.2 执行收尾阶段：gofmt 格式化、`make doc` 生成文档、生成 changelog 文件（由 tfpacer-finalize skill 在收尾阶段执行）