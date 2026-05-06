## 1. Service 层扩展

- [x] 1.1 在 `tencentcloud/services/vod/service_tencentcloud_vod.go` 新增 `CreateVodAigcApiToken(ctx, subAppId) (string, error)`，调用 SDK `CreateAigcApiToken`，使用 `resource.Retry(WriteRetryTimeout, ...)` 包装，日志打印不包含 token 明文（仅 SubAppId 与 RequestId）
- [x] 1.2 新增 `DescribeVodAigcApiTokens(ctx, subAppId) ([]string, error)`，调用 SDK `DescribeAigcApiTokens`，返回 `ApiTokens` 字符串切片；使用 `resource.Retry(ReadRetryTimeout, ...)` 包装
- [x] 1.3 新增 `DescribeVodAigcApiTokenById(ctx, subAppId, apiToken) (bool, error)`，内部调用 `DescribeVodAigcApiTokens` 并线性查找 apiToken
- [x] 1.4 新增 `DeleteVodAigcApiToken(ctx, subAppId, apiToken) error`，调用 SDK `DeleteAigcApiToken`，对 `ResourceNotFound.*` 错误做幂等处理（返回 nil）；其他错误走 `resource.Retry(WriteRetryTimeout, ...)`

## 2. Resource 实现

- [x] 2.1 创建 `tencentcloud/services/vod/resource_tc_vod_aigc_api_token.go`：声明 `ResourceTencentCloudVodAigcApiToken()`，包含 Create/Read/Delete 三个函数与 `Importer: schema.ImportStatePassthrough`
- [x] 2.2 定义 schema：`sub_app_id` (Int, Required, ForceNew)、`api_token` (String, Optional, Computed, ForceNew, Sensitive)
- [x] 2.3 实现 `resourceTencentCloudVodAigcApiTokenCreate`：读取 sub_app_id → 调 `CreateVodAigcApiToken` 拿回 token → `d.SetId(fmt.Sprintf("%d%s%s", subAppId, tccommon.FILED_SP, apiToken))` → `_ = d.Set("api_token", apiToken)` → 使用 `helper.Retry` 轮询 `DescribeVodAigcApiTokenById` 直到 exists=true（超时 `ReadRetryTimeout`）→ 调 Read
- [x] 2.4 实现 `resourceTencentCloudVodAigcApiTokenRead`：`defer tccommon.LogElapsed`、`defer tccommon.InconsistentCheck` → 解析 `d.Id()`（按 `tccommon.FILED_SP` 拆分；非法格式返回 fmt.Errorf）→ `_ = d.Set("sub_app_id", subAppId)`、`_ = d.Set("api_token", apiToken)` → 调 `DescribeVodAigcApiTokenById`；未找到则 `d.SetId("")` 并 log.WARN 后 return nil
- [x] 2.5 实现 `resourceTencentCloudVodAigcApiTokenDelete`：`defer tccommon.LogElapsed` → 解析 `d.Id()` → 调 `DeleteVodAigcApiToken` → `helper.Retry` 轮询直到 `DescribeVodAigcApiTokenById` 返回 exists=false（超时 `ReadRetryTimeout`）
- [x] 2.6 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `"tencentcloud_vod_aigc_api_token": vod.ResourceTencentCloudVodAigcApiToken()`，按字母序插入 VOD 资源段
- [x] 2.7 在 `tencentcloud/provider.md` VOD 资源分组中追加 `tencentcloud_vod_aigc_api_token` 条目

## 3. 文档模板

- [x] 3.1 创建 `tencentcloud/services/vod/resource_tc_vod_aigc_api_token.md`，包含：资源说明段、至少一个 HCL 示例（仅 `sub_app_id`）、import 示例（`terraform import tencentcloud_vod_aigc_api_token.foo {sub_app_id}#{api_token}`）、Sensitive + 30 秒同步延迟提示
- [x] 3.2 在仓库根执行 `make doc`，生成 `website/docs/r/vod_aigc_api_token.html.markdown`，确认生成成功且包含 schema/Example Usage/Import 段

## 4. 测试

- [x] 4.1 创建 `tencentcloud/services/vod/resource_tc_vod_aigc_api_token_test.go`，提供 `TestAccTencentCloudVodAigcApiTokenResource_basic`：使用基础 `testAccVodAigcApiTokenBasic` HCL 配置（`sub_app_id` 取测试账号下的有效 SubAppId），CheckDestroy 校验列表中不再包含对应 token
- [x] 4.2 在 `TestAccTencentCloudVodAigcApiTokenResource_basic` 中添加 ImportStep（`ImportStateIdFunc` 使用复合 ID、`ImportStateVerifyIgnore` 为空），验证导入后 no diff

## 5. 质量保障

- [x] 5.1 在仓库根执行 `go build ./...`，确认无编译错误
- [x] 5.2 执行 `go vet ./tencentcloud/services/vod/...`，确认无警告
- [x] 5.3 执行 `gofmt -l tencentcloud/services/vod/ | grep -v _test.go || true`，再 `gofmt -w tencentcloud/services/vod/resource_tc_vod_aigc_api_token*.go tencentcloud/services/vod/service_tencentcloud_vod.go`，确认格式化无输出
- [x] 5.4 执行 `go test ./tencentcloud/services/vod/ -run TestAccTencentCloudVodAigcApiTokenResource -count=1` 不带 `TF_ACC` 以确认编译与基础断言无误（验收测试会 Skip）

## 6. 验收测试（需要真实云凭证，由开发者人工运行）

- [ ] 6.1 配置环境变量：`export TF_ACC=1 TENCENTCLOUD_SECRET_ID=... TENCENTCLOUD_SECRET_KEY=...`，并准备一个有效的测试 SubAppId（在测试文件中以常量或 env var 提供）
- [ ] 6.2 执行 `go test ./tencentcloud/services/vod/ -run TestAccTencentCloudVodAigcApiTokenResource_basic -v -timeout 30m`，确认 Create → Import → Destroy 全流程通过
- [ ] 6.3 通过云控制台或 CLI 确认测试结束后没有残留 AIGC API Token

## 7. 交付收尾

- [x] 7.1 检查 `go.mod` / `go.sum` 未因本 change 修改（AIGC Token 相关 SDK 已 vendor）
- [x] 7.2 在 `.changelog/` 下新增一个新特性条目（文件名如 `<PR号>.txt`），类型 `release-note:new-resource`，内容 `tencentcloud_vod_aigc_api_token`
- [x] 7.3 自检 `website/docs/r/vod_aigc_api_token.html.markdown` 与 `resource_tc_vod_aigc_api_token.md` 一致（应由 make doc 自动保证）
