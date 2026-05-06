## Why

腾讯云点播（VOD）新推出 AIGC API Token 能力，用于对 AIGC 相关 API 调用进行独立鉴权。云端已提供完整的 Create/Describe/Delete 三个接口，但当前 Terraform Provider 未覆盖。用户需要通过 IaC 方式为不同 VOD 应用（SubAppId）签发并管理 AIGC API Token，以实现 AIGC 业务的自动化配置。SDK（`tencentcloud-sdk-go/tencentcloud/vod/v20180717`）中对应的 `CreateAigcApiToken` / `DescribeAigcApiTokens` / `DeleteAigcApiToken` 已就绪。

## What Changes

- 新增 Terraform 资源 `tencentcloud_vod_aigc_api_token`，位于 `tencentcloud/services/vod/` 目录，用于管理 VOD 应用的 AIGC API Token 生命周期
- 资源映射三个云 API：
  - `CreateAigcApiToken`（C）：为指定 `SubAppId` 创建一个 AIGC API Token
  - `DescribeAigcApiTokens`（R）：拉取 `SubAppId` 下全部 Token 列表，Read 时据此判断本资源 Token 是否仍存在
  - `DeleteAigcApiToken`（D）：删除指定 Token
- 无 Update 接口：`sub_app_id` 为 `ForceNew`，Token 值为 `Computed + ForceNew`，任何变更触发资源重建
- 处理云端 ~30 秒的数据同步延迟：Create 后读取与 Delete 后读取均使用 `helper.Retry` 轮询 `DescribeAigcApiTokens` 直至 Token 可见/消失
- 在 `tencentcloud/provider.go` 注册资源，在 `tencentcloud/provider.md` 中加入资源列表
- 提供资源文档模板（`resource_tc_vod_aigc_api_token.md`）与 website 文档（通过 `make doc` 生成）
- 提供验收测试（`resource_tc_vod_aigc_api_token_test.go`），包含基础创建/销毁流程

## Capabilities

### New Capabilities

- `vod-aigc-api-token`: 管理 VOD 应用的 AIGC API Token 资源，封装 `CreateAigcApiToken` / `DescribeAigcApiTokens` / `DeleteAigcApiToken` 三个云 API

### Modified Capabilities

无

## Impact

- **新增代码**:
  - `tencentcloud/services/vod/resource_tc_vod_aigc_api_token.go`
  - `tencentcloud/services/vod/resource_tc_vod_aigc_api_token_test.go`
  - `tencentcloud/services/vod/resource_tc_vod_aigc_api_token.md`
- **扩展代码**（在现有 `service_tencentcloud_vod.go` 中新增方法）:
  - `VodService.CreateVodAigcApiToken(ctx, subAppId) (apiToken string, err error)`
  - `VodService.DescribeVodAigcApiTokens(ctx, subAppId) (tokens []string, err error)`
  - `VodService.DescribeVodAigcApiTokenById(ctx, subAppId, apiToken) (exists bool, err error)`
  - `VodService.DeleteVodAigcApiToken(ctx, subAppId, apiToken) error`
- **修改代码**:
  - `tencentcloud/provider.go`: 注册 `tencentcloud_vod_aigc_api_token`
  - `tencentcloud/provider.md`: 在 VOD 分组下添加资源条目
- **依赖**: 复用已 vendored 的 `tencentcloud-sdk-go/tencentcloud/vod/v20180717`，SDK 中 `CreateAigcApiToken` / `DescribeAigcApiTokens` / `DeleteAigcApiToken` 三个方法及其 Request/Response 结构均已就绪，无需升级 SDK
- **生成文件**: `website/docs/r/vod_aigc_api_token.html.markdown`（通过 `make doc` 生成）
- **API 频率**: `DescribeAigcApiTokens` 默认 20 次/秒；Create/Delete 未声明严格限频，由 Provider 统一重试机制兜底
- **向后兼容**: 仅新增资源与 service 方法，不修改任何现有 schema 或行为
- **安全注意**: Token 值属于敏感凭证，`api_token` 字段在 schema 中声明 `Sensitive: true`，避免在 Terraform 输出与日志中明文暴露
