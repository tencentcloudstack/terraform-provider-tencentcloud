## 1. Schema 调整

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_security_client_attester.go` 的 `client_attesters` 嵌套 schema 中，于 `attester_duration` 之后新增 `usage_limit`（Optional, TypeString）字段，Description 说明取值 `on`/`off`、默认 `off`
- [x] 1.2 在 `client_attesters` 嵌套 schema 中，于 `usage_limit` 之后新增 `max_usage_count`（Optional, TypeInt）字段，Description 说明 `usage_limit=on` 时必填、取值范围 [1, 100000]
- [x] 1.3 在 `tc_eo_captcha_option` 子块 schema 中，于 `captcha_mode` 之后新增 `mobile_captcha_mode`（Optional, TypeString）字段，Description 说明取值 `Invisible`/`Visible`、默认 `Invisible`

## 2. CRUD 逻辑扩展

- [x] 2.1 在 `buildClientAttesterFromMap` 中新增对 `usage_limit`（string 非空时 `attester.UsageLimit = helper.String(val)`）和 `max_usage_count`（int 非零时 `attester.MaxUsageCount = helper.IntInt64(val)`）的处理，保证未配置时不传给云 API
- [x] 2.2 在 `buildClientAttesterFromMap` 的 `tc_eo_captcha_option` 分支中新增对 `mobile_captcha_mode`（string 非空时 `eoOpt.MobileCaptchaMode = helper.String(v)`）的处理
- [x] 2.3 在 `resourceTencentCloudTeoSecurityClientAttesterRead` 的 `attesterMap` 构造中，新增对 `respData.UsageLimit`（非 nil 时 `attesterMap["usage_limit"]`）和 `respData.MaxUsageCount`（非 nil 时 `attesterMap["max_usage_count"]`）的回填
- [x] 2.4 在 Read 的 `TCEOCaptchaOption` 分支中，新增对 `respData.TCEOCaptchaOption.MobileCaptchaMode`（非 nil 时 `opt["mobile_captcha_mode"]`）的回填
- [x] 2.5 确认 Create/Update 复用 `buildClientAttesterFromMap`，无需额外改动（新增字段自动覆盖两个接口）

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_security_client_attester_test.go` 中新增测试用例，使用 gomonkey mock `CreateSecurityClientAttesterWithContext`、`DescribeSecurityClientAttester`，验证包含 `usage_limit`、`max_usage_count` 的 Create 流程
- [x] 3.2 新增测试用例验证 `tc_eo_captcha_option` 中 `mobile_captcha_mode` 字段的 Create/Read 流程
- [x] 3.3 新增测试用例验证未配置新字段时 `buildClientAttesterFromMap` 不设置对应指针（向后兼容）
- [x] 3.4 保证已有测试用例（无新字段场景）继续通过（不修改其行为）

## 4. 文档同步

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_security_client_attester.md` 的 Example Usage 中补充 `usage_limit`、`max_usage_count` 字段示例（在 TC-EO-CAPTCHA 场景补充 `mobile_captcha_mode`）
- [x] 4.2 由收尾阶段执行 `make doc` 重新生成 `website/docs/r/teo_security_client_attester.html.markdown`（禁止手改 website/ 目录）

## 5. 验证

- [x] 5.1 检查所有函数返回的 error 已正确处理，无未使用变量
- [x] 5.2 核对新增字段在 Create/Update 接口的 `ClientAttester` 结构体中存在（已在 design.md 摘录确认），CRUD 参数一致性通过
- [x] 5.3 核对 Read 回填逻辑符合"Response 字段为 nil 时不调用 set"的要求