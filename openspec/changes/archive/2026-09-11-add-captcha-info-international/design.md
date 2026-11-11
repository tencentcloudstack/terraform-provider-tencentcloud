## Context

腾讯云验证码（Captcha）国际站 4 个接口（CreateCaptchaInfoInternational / DescribeCaptchaInfoListInternational / ModifyCaptchaInfoInternational / RemoveCaptchaInfoInternational）位于国际站 SDK `tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722`，适合做成 Terraform CRUD 资源。资源命名遵守 provider 惯例 `tencentcloud_<service>_<resource>`，归属到 `tencentcloud/services/captcha/`。

`CreateCaptchaInfoInternational` 入参共 13 个字段，其中 12 个为标量（`*string` / `*int64`），1 个为字符串列表（`Tags []*string`），无嵌套结构体。本文档解决 schema 字段映射、资源 ID 约定（`Data` 是 int64、`CaptchaAppId` 是 string 的类型转换）、Update 行为、Read 查询路径等关键决策。

参考资源：`tencentcloud_igtm_monitor`（代码风格）、`tencentcloud_config_compliance_pack`（文档与测试命名）、provider 中现有资源的 `UseXxxClient()` 客户端调用约定与 service 层方法签名。

## Goals / Non-Goals

**Goals:**

- Schema 字段名与 `CreateCaptchaInfoInternational` 接口入参 1:1 映射（snake_case 化），字段集不多不少。
- 资源 ID = Create 响应 `Data`（`*int64`），转 string 后 `d.SetId`；同步映射 `CaptchaAppId`（`*string`）入参。
- Read 通过 `DescribeCaptchaInfoListInternational` + `CaptchaAppId` 精确查询单个；分页 `PageSize` 取接口文档最大值。
- Create / Read / Update / Delete 的全部 SDK 调用均用 `resource.Retry(...)` 包装。
- 全部接口返回值取值都做空指针保护。
- 代码风格严格对齐 `tencentcloud_igtm_monitor`（变量声明顺序、`defer LogElapsed + InconsistentCheck` 模式）。
- 资源 doc 命名 `resource_tc_captcha_info_international.md`，测试命名 `resource_tc_captcha_info_international_test.go`，与参考样板一致。

**Non-Goals:**

- 不实现对应的 data source —— 本 change 只交付 resource。
- 不在 provider 层做字段取值合法性强校验（如 `verify_rank` 枚举值），让 API 层报错引导用户；仅在文档说明。
- 不修改任何已存在资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射（与 CreateCaptchaInfoInternational 入参 1:1）

| HCL 字段 | SDK 字段 | 类型 | ForceNew | 说明 |
|---|---|---|---|---|
| `app_name` | `AppName` | `TypeString` | No | 验证码应用名称 |
| `channel_info` | `ChannelInfo` | `TypeString` | **Yes** | 渠道信息；Modify 接口不含该字段，见 D3 |
| `verify_rank` | `VerifyRank` | `TypeString` | No | 验证码安全等级 |
| `user_set_cap_type` | `UserSetCapType` | `TypeString` | No | 用户设置容量类型 |
| `defend_mode` | `DefendMode` | `TypeString` | No | 防御模式 |
| `tags` | `Tags` | `TypeList`(`TypeString`) | No | 标签，字符串列表 |
| `disable_invisible_switch` | `DisableInvisibleSwitch` | `TypeString` | No | 是否关闭无感验证开关 |
| `verify_domain` | `VerifyDomain` | `TypeString` | No | 验证域名 |
| `verify_bundle_id` | `VerifyBundleId` | `TypeString` | No | 验证包 ID |
| `verify_package` | `VerifyPackage` | `TypeString` | No | 验证套餐 |
| `check_appid_switch` | `CheckAppidSwitch` | `TypeInt` | No | 检查 AppId 开关（0/1） |
| `check_iv_switch` | `CheckIvSwitch` | `TypeInt` | No | 检查 IV 开关（0/1） |
| `check_box_style` | `CheckBoxStyle` | `TypeString` | No | 验证码弹窗样式 |

**Required/Optional 约定**（已按接口文档确认）：`app_name`、`channel_info` 为必填（`Required`），其余 11 个字段为选填（`Optional`）。类型映射规则：`*string → TypeString`，`*int64 → TypeInt`，`[]*string → TypeList{TypeString}`。

**理由**：`tags` 是 `[]*string`（字符串列表），不是嵌套结构体，因此用 `TypeList` + `Elem: TypeString`，而非 `TypeMap` 或嵌套 Resource。

### D2 — 资源 ID = `Data`（int64 → string）

**问题**：Create 响应返回的 `Data` 是 `*int64`（验证码 AppId），而 `RemoveCaptchaInfoInternational` / `ModifyCaptchaInfoInternational` / `DescribeCaptchaInfoListInternational` 的查询入参是 `CaptchaAppId`（`*string`）。二者类型不一致，需要转换。

**决策**：

```go
if response.Response.Data == nil {
    return fmt.Errorf("Create captcha info international failed, Data is nil.")
}
d.SetId(helper.Int64ToStr(*response.Response.Data))
```

此后 CRUD 中一律用 `d.Id()`（string）作为 `CaptchaAppId` 的取值来源：

```go
request.CaptchaAppId = helper.String(d.Id())
```

**理由**：Terraform ID 必须是 string；`helper.Int64ToStr` 是 provider 现有的 int64→string 工具函数，保持与项目其它资源一致。

### D3 — `channel_info` 标记为 ForceNew

**问题**：`CreateCaptchaInfoInternationalRequest` 有 `ChannelInfo` 字段，但 `ModifyCaptchaInfoInternationalRequest` **没有** `ChannelInfo`。如果 `channel_info` 是 Optional+No-ForceNew，用户改 `channel_info` 会触发 Update，但 Modify 接口无法把新值推送到后端 → 漂移。

**决策**：`channel_info` 字段标 `ForceNew: true`。改 `channel_info` 触发资源销毁重建，避免静默漂移。

**替代方案**：在 Update 中忽略 `channel_info` 变化（`DiffSuppressFunc`），但这会让用户 plan 看到"零变化"却无法实际改，更糟糕。ForceNew 语义透明。

### D4 — Update 字段集合（mutableArgs）

参考 `tencentcloud_igtm_monitor` 的 mutableArgs 模式。Modify 接口支持的字段（去除 `CaptchaAppId` 与 `ChannelInfo`）：

```go
mutableArgs := []string{
    "app_name", "verify_rank", "user_set_cap_type", "defend_mode",
    "disable_invisible_switch", "verify_domain", "verify_bundle_id",
    "verify_package", "check_appid_switch", "check_iv_switch",
    "check_box_style", "tags",
}
// 注：channel_info 是 ForceNew（D3），不在 mutableArgs 内
```

任意字段变化触发一次 `ModifyCaptchaInfoInternational` 调用，把所有可变字段全量提交（全量覆盖，与 Modify API 行为一致）。

### D5 — Read 路径

`DescribeCaptchaInfoListInternational` 支持 `CaptchaAppId` 精确入参。Read 实现：

```go
request := captcha.NewDescribeCaptchaInfoListInternationalRequest()
request.CaptchaAppId = helper.String(captchaAppId)
request.PageIndex = helper.IntInt64(1)
request.PageSize = helper.IntInt64(100)  // 接口文档最大值，实现时二次确认
```

返回 `Response.Data.DataList []*DescribeCaptchaConsoleSubDataInternational`：
- `len == 0` → 资源已被外部删除 → `d.SetId("")` 并 return nil
- `len >= 1` → 取第一项

`DescribeCaptchaInfoInternationalById` 放在 `service_tencentcloud_captcha.go`，签名：

```go
func (me *CaptchaService) DescribeCaptchaInfoInternationalById(ctx context.Context, captchaAppId string) (ret *captcha.DescribeCaptchaConsoleSubDataInternational, errRet error)
```

### D6 — 同步操作，无需 waitForTaskFinish

Create / Modify / Remove 是同步接口（SDK 模型无 TaskId），无需 waitForTaskFinish。Delete 后立即调 Describe 拿到空列表即视为已删除。

### D7 — 文档与测试命名

- `resource_tc_captcha_info_international.md`：包含 1-2 个 HCL example + Import 段。
- `resource_tc_captcha_info_international_test.go`：基础 Create/Update/ImportStateVerify 三步骤，包名 `captcha_test`，测试函数 `TestAccTencentCloudCaptchaInfoInternationalResource_basic`。

均与 `resource_tc_config_compliance_pack.md` / `_test.go` 同款结构。

### D8 — Provider 注册位置

在 `tencentcloud/provider.go` 中追加一行：

```go
"tencentcloud_captcha_info_international": captcha.ResourceTencentCloudCaptchaInfoInternational(),
```

## Risks / Trade-offs

- **Risk**: SDK 未 vendor —— 当前 vendor 尚无 `tencentcloud/captcha` 包（intl-en 只 vendor 了 clb/cvm/dnspod/mdl/privatedns）。Mitigation: 实现前在代码 import `tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722` 并重新 `go mod vendor`。
- **Risk**: `Data`(int64) 与 `CaptchaAppId`(string) 的类型转换 —— Mitigation: D2 用 `helper.Int64ToStr` 统一转 string，CRUD 全程用 `d.Id()`。
- **Risk**: `channel_info` 在 Modify 缺失 —— Mitigation: D3 ForceNew，plan 阶段显式展示"销毁重建"。
- **Trade-off**: 不实现 data source，聚焦 resource；data source 后续单独 PR 跟进。
