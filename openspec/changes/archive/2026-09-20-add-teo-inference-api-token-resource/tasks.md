## 1. 资源 Schema 与 CRUD 实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_inference_api_token.go`，定义 `ResourceTencentCloudTeoInferenceAPIToken()` 返回 `schema.Resource`，包含 Create/Read/Update/Delete 与 Importer
- [x] 1.2 定义 schema 字段：`zone_id`（Required+ForceNew）、`name`（Required+ForceNew）、`token_id`（Computed）、`content`（Computed+Sensitive）、`create_time`（Computed）
- [x] 1.3 实现 `resourceTencentCloudTeoInferenceAPITokenCreate`：构造 `CreateInferenceAPITokenRequest` 填充 `ZoneId`、`Name`，用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 API，retry 块内检查 response/TokenId 为空则返回 `NonRetryableError`，retry 块外打印 logId 与 `d.Id()` 后设置复合 ID `zone_id#token_id`，最后调用 Read
- [x] 1.4 实现 `resourceTencentCloudTeoInferenceAPITokenRead`：从 `d.Id()` 解析 `zone_id` 与 `token_id`，调用 `DescribeInferenceAPITokens`（Limit=100），retry 块内若 Tokens 为空或未匹配到 token_id 则返回 `NonRetryableError`，retry 块外失败路径打印 `[CRUD] teo inference_api_token id=%s` 后 `d.SetId("")`，设置各字段前检查 nil
- [x] 1.5 实现 `resourceTencentCloudTeoInferenceAPITokenUpdate`：声明 `immutableArgs := []string{"name"}`，检测变更返回 error，无变更则调用 Read
- [x] 1.6 实现 `resourceTencentCloudTeoInferenceAPITokenDelete`：从 `d.Id()` 解析 `zone_id` 与 `token_id`，构造 `DeleteInferenceAPITokenRequest` 填充 `ZoneId`、`TokenId`，用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 API

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoInferenceAPIToken` 方法，接收 `ctx`、`zoneId`、`tokenId`，调用 `DescribeInferenceAPITokens`（Limit=100）并在客户端遍历匹配 `tokenId`，返回匹配的 `*teo.InferenceAPIToken`

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `"tencentcloud_teo_inference_api_token": teo.ResourceTencentCloudTeoInferenceAPIToken()`
- [x] 3.2 在 `tencentcloud/provider.md` 中追加 `tencentcloud_teo_inference_api_token` 资源名称

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_inference_api_token.md`，包含一句话描述（带上 TEO 产品名）、Example Usage、Import 部分（说明使用联合 id `zone_id#token_id`）

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_inference_api_token_test.go`，使用 gomonkey mock 云 API，复用 teo 测试包已有的 `newMockMeta`、`ptrString` 辅助函数
- [x] 5.2 新增 `TestTeoInferenceAPIToken_Create_Success` 测试用例：mock `CreateInferenceAPIToken` 返回 TokenId，mock `DescribeInferenceAPITokens` 返回匹配的 Token，验证 Create 成功且 ID 为 `zone_id#token_id`
- [x] 5.3 新增 `TestTeoInferenceAPIToken_Create_EmptyTokenId` 测试用例：mock `CreateInferenceAPIToken` 返回空 TokenId，验证 Create 返回 error
- [x] 5.4 新增 `TestTeoInferenceAPIToken_Create_APIError` 测试用例：mock `CreateInferenceAPIToken` 返回 error，验证 Create 返回 error
- [x] 5.5 新增 `TestTeoInferenceAPIToken_Read_Success` 测试用例：mock `DescribeInferenceAPITokens` 返回匹配 Token，验证 Read 正确设置各字段
- [x] 5.6 新增 `TestTeoInferenceAPIToken_Read_NotFound` 测试用例：mock `DescribeInferenceAPITokens` 返回空列表，验证 Read 调用 `d.SetId("")`
- [x] 5.7 新增 `TestTeoInferenceAPIToken_Update_ImmutableError` 测试用例：模拟 `name` 变更，验证 Update 返回 error
- [x] 5.8 新增 `TestTeoInferenceAPIToken_Delete_Success` 测试用例：mock `DeleteInferenceAPIToken` 成功，验证 Delete 无 error
- [x] 5.9 新增 `TestTeoInferenceAPIToken_Delete_APIError` 测试用例：mock `DeleteInferenceAPIToken` 返回 error，验证 Delete 返回 error
- [x] 5.10 新增 `TestTeoInferenceAPIToken_Schema` 测试用例：验证 schema 中各字段的 Required/ForceNew/Computed/Sensitive 属性正确

## 6. 验证

- [x] 6.1 执行 `gofmt` 格式化生成的 go 文件（由收尾阶段 tfpacer-finalize skill 执行）
- [x] 6.2 执行 `make doc` 生成 website/docs/ 下的资源文档（由收尾阶段 tfpacer-finalize skill 执行）
- [x] 6.3 创建 `.changelog/` 下的 changelog 文件（由收尾阶段 tfpacer-finalize skill 执行）