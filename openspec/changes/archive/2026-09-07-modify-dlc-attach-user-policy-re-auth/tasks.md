## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policy_attachment.go` 的 `ResourceTencentCloudDlcAttachUserPolicyAttachment` schema 中，将 `policy_set.re_auth` 字段从 `Computed: true` 修改为 `Optional: true, Computed: true`，保留 `Type: schema.TypeBool` 和现有 `Description`，使该字段变为用户可设置的入参，同时仍从云 API 响应中回填。

## 2. Create 逻辑修改

- [x] 2.1 在 `resourceTencentCloudDlcAttachUserPolicyAttachmentCreate` 函数中，于 `policy_set` 遍历构建 `dlc.Policy` 的逻辑内（与其他子字段并列），增加 `re_auth` 的映射：`if v, ok := dMap["re_auth"]; ok { policy.ReAuth = helper.Bool(v.(bool)) }`，确保用户配置的 `re_auth` 被转发到 `AttachUserPolicy` 请求的 `PolicySet[0].ReAuth`。
  - 确认：vendored SDK `dlc.Policy` 结构体已包含 `ReAuth *bool` 字段，`AttachUserPolicyRequest.PolicySet` 接受 `[]*Policy`，云 API 支持该入参。

## 3. 文档修改

- [x] 3.1 更新 `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policy_attachment.md`，在示例用法中补充 `re_auth` 字段的使用示例（如 `re_auth = true`），并确保文档描述反映该字段现为可选入参。
  - 注意：`website/docs/` 下的文档由 `make doc` 自动生成，本任务仅修改 `tencentcloud/services/dlc/` 下的 `.md` 源文件。

## 4. 单元测试补充

- [x] 4.1 在 `tencentcloud/services/dlc/resource_tc_dlc_attach_user_policy_attachment_test.go` 中补充单元测试，覆盖 Create 路径设置 `re_auth` 的场景：在 mock 的 `AttachUserPolicyWithContext` 中捕获请求，断言 `capturedRequest.PolicySet[0].ReAuth` 等于配置值（`true`）；并在 mock 的 `DescribeUserInfoWithContext` 响应中带回 `ReAuth`，断言 state 中 `policy_set.0.re_auth` 被正确回填。
  - 使用 gomonkey mock 云 API，不使用 terraform 验收测试套件。

## 5. 代码正确性校验

- [x] 5.1 校验 `re_auth` 字段在 schema 中为 `Optional: true, Computed: true`，且 Create 中映射的 `dlc.Policy.ReAuth` 字段名与 vendored SDK 一致；确认 Read 路径（`flattenDlcAttachUserPolicyAttachmentPolicySet`）已处理 `policy.ReAuth` 回填，无需额外修改。
- [x] 5.2 校验资源仍为 `RESOURCE_KIND_ATTACHMENT`（仅 CRD，无 Update），`policy_set` 保持 `ForceNew`，`re_auth` 作为 `policy_set` 子字段随父级 `ForceNew` 行为变更触发重建，符合现有不可变契约。
