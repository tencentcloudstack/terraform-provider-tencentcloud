## ADDED Requirements

### Requirement: 资源注册

Provider SHALL 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册资源 `tencentcloud_vod_aigc_api_token`，并在 `tencentcloud/provider.md` 的 VOD 分组下追加对应条目。

#### Scenario: 资源出现在 Provider Schema 中

- **WHEN** 用户执行 `terraform providers schema -json` 查询 tencentcloud provider 的 ResourcesMap
- **THEN** 输出中包含 `tencentcloud_vod_aigc_api_token` 资源定义

#### Scenario: 文档列出该资源

- **WHEN** 用户查看 `tencentcloud/provider.md`
- **THEN** 在 VOD 服务分组中能找到 `tencentcloud_vod_aigc_api_token` 条目

### Requirement: Schema 字段定义

资源 SHALL 声明以下 schema 字段：

- `sub_app_id` (Int, Required, ForceNew)：点播子应用 ID，创建后不可变
- `api_token` (String, Optional, Computed, ForceNew, Sensitive)：API Token 值；Create 时由云端返回并写入 state；Import 场景允许用户传入现有 token

资源 MUST NOT 声明其他输入字段（API 不支持）。

#### Scenario: 用户仅指定必填字段即可创建

- **WHEN** 用户在 `.tf` 中仅指定 `sub_app_id = 251006666`
- **THEN** `terraform plan` 成功，`api_token` 显示为 `(known after apply) (sensitive)`

#### Scenario: api_token 变更触发 ForceNew

- **WHEN** 用户修改 state 中的 `api_token` 值或 `sub_app_id` 值
- **THEN** `terraform plan` 显示资源将被销毁并重建（`-/+`）

#### Scenario: api_token 在输出中被标记为 Sensitive

- **WHEN** 用户执行 `terraform plan` 或 `terraform apply`
- **THEN** 任何涉及 `api_token` 的输出被替换为 `(sensitive value)`，不明文显示 token 字符串

### Requirement: 创建（Create）行为

Create 函数 SHALL：

1. 从 ResourceData 读取 `sub_app_id`
2. 调用 `VodService.CreateVodAigcApiToken(ctx, subAppId)` 触发云端创建，获取返回的 `api_token`
3. 使用 `tccommon.FILED_SP`（`#`）拼装复合 ID：`{sub_app_id}#{api_token}` 并调用 `d.SetId(...)`
4. 轮询 `DescribeAigcApiTokens` 直至列表中出现新 token（使用 `helper.Retry` 包装，超时上限为 `tccommon.ReadRetryTimeout`）
5. 调用 Read 函数以填充最终 state

Create 函数 MUST 将云端返回的 `ApiToken` 通过 `_ = d.Set("api_token", apiToken)` 写回 ResourceData。

#### Scenario: 创建成功

- **WHEN** `CreateAigcApiToken` 返回 `ApiToken = "hGjH1dsBbjUbT9Cu"`，后续轮询 30s 内发现该 token 在列表中
- **THEN** 资源 ID 被设置为 `{sub_app_id}#hGjH1dsBbjUbT9Cu`，`api_token` state 值为 `"hGjH1dsBbjUbT9Cu"`

#### Scenario: 创建后同步延迟导致列表暂未包含 token

- **WHEN** Create 成功但前几次 `DescribeAigcApiTokens` 返回的列表中不包含新 token
- **THEN** 在 `ReadRetryTimeout` 超时前持续重试；期间出现可识别的瞬时错误也可重试；最终成功发现 token 或超时报错

#### Scenario: 创建 API 返回不可重试错误

- **WHEN** `CreateAigcApiToken` 返回如参数非法类错误
- **THEN** 资源创建失败，错误原样返回；state 不写入

### Requirement: 读取（Read）行为

Read 函数 SHALL：

1. 解析 `d.Id()` 获取 `sub_app_id` 与 `api_token`（按 `#` 拆分）
2. 将解析结果通过 `_ = d.Set(...)` 写回（以支持 import 场景下首次 Read 时的字段填充）
3. 调用 `VodService.DescribeVodAigcApiTokenById(ctx, subAppId, apiToken)`：
   - 若返回 `exists=true`：保持 state
   - 若返回 `exists=false`：调用 `d.SetId("")` 并打印 WARN 日志，触发 Terraform 视为资源已消失
4. 在函数开头使用 `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(...)`

#### Scenario: Token 仍存在

- **WHEN** `DescribeAigcApiTokens` 返回的列表中包含当前资源的 token
- **THEN** state 保持不变，`sub_app_id` 与 `api_token` 字段被刷新

#### Scenario: Token 已被外部删除

- **WHEN** `DescribeAigcApiTokens` 返回的列表中不包含当前资源的 token
- **THEN** `d.SetId("")` 被调用，下一轮 `terraform apply` 会重新创建资源

#### Scenario: Resource ID 格式非法

- **WHEN** `d.Id()` 无法按 `#` 拆成 `[sub_app_id, api_token]` 两段，或 `sub_app_id` 部分无法解析为 uint64
- **THEN** Read 返回错误并提示期望的 ID 格式

### Requirement: 删除（Delete）行为

Delete 函数 SHALL：

1. 解析 `d.Id()` 获取 `sub_app_id` 与 `api_token`
2. 调用 `VodService.DeleteVodAigcApiToken(ctx, subAppId, apiToken)`
3. 遇到 `ResourceNotFound.UserNotExist` 或等价"资源已不存在"的错误时，视为删除成功（幂等）
4. 成功后轮询 `DescribeAigcApiTokens`，直至列表中不再包含该 token（使用 `helper.Retry`，超时 `tccommon.ReadRetryTimeout`）

#### Scenario: 删除成功

- **WHEN** `DeleteAigcApiToken` 返回成功且后续轮询 30s 内列表不再包含 token
- **THEN** 资源被移除，state 清空

#### Scenario: 删除时 Token 已不存在

- **WHEN** `DeleteAigcApiToken` 返回 `ResourceNotFound.*` 类错误
- **THEN** Delete 函数不返回错误，state 被正常清理

#### Scenario: 删除后同步延迟

- **WHEN** `DeleteAigcApiToken` 调用成功但立即查询列表仍包含该 token
- **THEN** 在 `ReadRetryTimeout` 超时前持续重试直到列表中消失或超时

### Requirement: 资源导入（Import）

资源 SHALL 支持通过 `terraform import` 导入现有 Token，导入 ID 格式为 `{sub_app_id}#{api_token}`。

#### Scenario: 导入成功

- **WHEN** 用户执行 `terraform import tencentcloud_vod_aigc_api_token.foo 251006666#hGjH1dsBbjUbT9Cu` 且该 token 在 SubAppId 251006666 下真实存在
- **THEN** state 被成功创建，后续 `terraform plan` 显示 no changes

#### Scenario: 导入未知 Token

- **WHEN** 用户导入一个不在列表中的 token
- **THEN** 导入后首次 Read 将 `d.SetId("")`，Terraform 提示资源不存在

### Requirement: Service 层封装

`tencentcloud/services/vod/service_tencentcloud_vod.go` SHALL 新增以下方法，全部使用 `me.client.UseVodClient()` 调用 SDK，并遵循项目统一的 Request 构造、日志、错误处理模式：

- `CreateVodAigcApiToken(ctx context.Context, subAppId uint64) (string, error)`
- `DescribeVodAigcApiTokens(ctx context.Context, subAppId uint64) ([]string, error)`
- `DescribeVodAigcApiTokenById(ctx context.Context, subAppId uint64, apiToken string) (bool, error)`
- `DeleteVodAigcApiToken(ctx context.Context, subAppId uint64, apiToken string) error`

#### Scenario: Service 方法被资源 CRUD 调用

- **WHEN** 资源 Create/Read/Delete 执行
- **THEN** 均通过 VodService 封装方法间接调用 SDK，不直接在资源文件内操作 SDK client

#### Scenario: Service 层日志不打印 Token 明文

- **WHEN** Service 方法打印 request/response 日志
- **THEN** 日志中 `ApiToken` 字段被掩码为 `***` 或不出现

### Requirement: 文档模板

`tencentcloud/services/vod/resource_tc_vod_aigc_api_token.md` SHALL 包含：

- 资源功能简介
- 至少 1 个 HCL 示例（含 `sub_app_id`）
- `terraform import` 示例，ID 格式说明为 `sub_app_id#api_token`
- 关于 token 为敏感值、云端约有 30 秒同步延迟的注意事项

Provider SHALL 通过 `make doc` 生成 `website/docs/r/vod_aigc_api_token.html.markdown`，禁止手写该文件。

#### Scenario: make doc 生成 website 文档

- **WHEN** 开发者在仓库根目录执行 `make doc`
- **THEN** 在 `website/docs/r/` 下生成 `vod_aigc_api_token.html.markdown`，包含 schema、示例、import 章节

### Requirement: 验收测试

`resource_tc_vod_aigc_api_token_test.go` SHALL 提供至少 1 个验收测试用例 `TestAccTencentCloudVodAigcApiTokenResource_basic`，覆盖 Create + Destroy 主流程；用例使用 `TF_ACC=1` 开关并依赖 `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY` 环境变量。

#### Scenario: 基础验收用例通过

- **WHEN** 开发者在配置了真实云凭证的环境中执行 `TF_ACC=1 go test ./tencentcloud/services/vod/ -run TestAccTencentCloudVodAigcApiTokenResource_basic -v`
- **THEN** 测试完成 Create → Read → Destroy 全流程，最终资源被销毁，测试通过

### Requirement: 向后兼容与依赖

变更 MUST NOT:

- 修改任何现有资源或数据源的 schema
- 引入新的 SDK 依赖（复用已 vendored 的 `tencentcloud-sdk-go/tencentcloud/vod/v20180717`）
- 修改 `go.mod` / `go.sum`（AIGC Token 相关 SDK 已在 vendor 中）

#### Scenario: 依赖未被修改

- **WHEN** 开发者在 change 分支执行 `git diff master -- go.mod go.sum`
- **THEN** 与本 change 相关的 AIGC Token 资源代码不引起新的模块/版本变动
