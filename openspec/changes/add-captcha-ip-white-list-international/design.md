## Context

腾讯云验证码（Captcha）国际站 IP 白名单 4 个接口（CreateIpWhiteListInternational / DescribeIpWhiteListInternational / ModifyIpWhiteListInternational / DeleteIpWhiteListInternational）位于国际站 SDK `tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722`，适合做成 Terraform CRUD 资源。资源命名遵守 provider 惯例 `tencentcloud_<service>_<resource>`，归属到 `tencentcloud/services/captcha/`（该包已存在，含上一个 `tencentcloud_captcha_info_international` 资源的 service 层与客户端）。

`CreateIpWhiteListInternational` 入参共 4 个标量字段（无嵌套结构体）。本文档解决 schema 字段映射、复合 ID 约定（`CaptchaAppid + IdList[0]`）、Update 行为、Read 查询路径等关键决策。

参考资源：`tencentcloud_igtm_monitor`（代码风格）、`tencentcloud_config_compliance_pack`（文档与测试命名）、上一个 `tencentcloud_captcha_info_international`（captcha 服务的 service 层与 `UseCaptchaClient()` 客户端约定）。

## Goals / Non-Goals

**Goals:**

- Schema 字段名与 `CreateIpWhiteListInternational` 接口入参 1:1 映射（snake_case 化），字段集不多不少。
- 资源 ID = `CaptchaAppid + IdList[0]` 复合 ID（`#` 分隔），CRUD 全程解析/拼接。
- Read 通过 `DescribeIpWhiteListInternational` + `CaptchaAppid` 过滤，再按 `Id` 精确匹配；分页 `PageSize` 取接口文档最大值。
- Create / Read / Update / Delete 的全部 SDK 调用均用 `resource.Retry(...)` 包装。
- 全部接口返回值取值都做空指针保护。
- 代码风格严格对齐 `tencentcloud_igtm_monitor`（变量声明顺序、`defer LogElapsed + InconsistentCheck` 模式）。
- 资源 doc 命名 `resource_tc_captcha_ip_white_list_international.md`，测试命名 `resource_tc_captcha_ip_white_list_international_test.go`。

**Non-Goals:**

- 不实现对应的 data source —— 本 change 只交付 resource。
- 不在 schema 暴露 `status`（Create 入参无此字段，见 D1），IP 白名单的启用/禁用状态不在本资源管理范围内。
- 不修改任何已存在资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射（与 CreateIpWhiteListInternational 入参 1:1）

| HCL 字段 | SDK 字段 | 类型 | Required | ForceNew | 说明 |
|---|---|---|---|---|---|
| `name` | `Name` | `TypeString` | Required | No | IP 白名单名称 |
| `captcha_appid` | `CaptchaAppid` | `TypeInt` | Required | **Yes** | 验证码 appid；资源 ID 组成部分 |
| `ip` | `Ip` | `TypeString` | Required | **Yes** | IP 数据；Modify 接口无此字段 |
| `comment` | `Comment` | `TypeString` | Optional | No | 备注信息 |

**Required/Optional 约定**（已按接口文档确认）：`name`、`captcha_appid`、`ip` 为必填（`Required`），`comment` 为选填（`Optional`）。类型映射：`*string → TypeString`，`*int64 → TypeInt`。

**关于 `status`**：`CreateIpWhiteListInternationalRequest` 无 `Status` 字段，`ModifyIpWhiteListInternationalRequest` 有 `Status`（选填，0 开启/1 关闭）。按「schema 与 Create 入参 1:1」约束，本资源**不暴露 `status`**；Update 不传 `Status`，维持后端默认。IP 白名单的启用/禁用状态不在本资源管理范围内（后续若需要可单独提供 data source / 操作资源）。

### D2 — 资源 ID = `CaptchaAppid + IdList[0]`（复合 ID）

Create 响应返回 `IdList []*int64`（IP 白名单资源 ID 列表，示例 `[2100000000]`）。资源 ID 由 `CaptchaAppid` 与 `IdList[0]` 复合而成，使用项目分隔符 `#`：

```go
captchaAppid := int64(v.(int)) // 来自 d.GetOkExists("captcha_appid")
...
d.SetId(fmt.Sprintf("%d#%d", captchaAppid, *response.Response.IdList[0]))
```

Read/Update/Delete 中解析：

```go
parts := strings.Split(d.Id(), tccommon.FILED_SP) // FILED_SP = "#"
captchaAppid := helper.StrToInt64(parts[0])
id := helper.StrToInt64(parts[1])
```

**理由**：`Id` 仅在某 `CaptchaAppid` 下唯一，必须与 `CaptchaAppid` 组合才能全局唯一定位一条白名单记录。`Delete` / `Modify` 需要 `CaptchaAppid + Id`，`Describe` 需要 `CaptchaAppid` 过滤，与复合 ID 天然对应。

### D3 — `captcha_appid` 与 `ip` 标记为 ForceNew

- `captcha_appid`：是资源 ID 的组成部分，且 `ModifyIpWhiteListInternational` 仅用它定位记录、不修改它 → `ForceNew: true`。
- `ip`：`ModifyIpWhiteListInternational` **没有 `Ip` 字段**，改 IP 无法原地更新 → `ForceNew: true`。

### D4 — Update 字段集合（mutableArgs）

`ModifyIpWhiteListInternational` 支持的可变字段为 `Name`、`Comment`（`Status` 选填但不在 schema 内）。因此：

```go
mutableArgs := []string{"name", "comment"}
```

任意字段变化触发一次 `ModifyIpWhiteListInternational` 调用，全量提交 `Name`、`Comment`（全量覆盖），并始终带上 `CaptchaAppid` + `Id`（从复合 ID 解析）。`ip`、`captcha_appid` 是 ForceNew，不在 mutableArgs 内。

### D5 — Read 路径

`DescribeIpWhiteListInternational` 支持 `CaptchaAppid` 过滤，返回 `Data.DataList []*DescribeCaptchaWhiteListItem`。service 层实现：

```go
func (me *CaptchaService) DescribeCaptchaIpWhiteListInternationalById(ctx context.Context, captchaAppid, id int64) (ret *captchaintl.DescribeCaptchaWhiteListItem, errRet error)
```

内部：
- `request.CaptchaAppid = helper.Int64(captchaAppid)`
- `request.PageIndex = helper.IntInt64(1)`
- `request.PageSize = helper.IntInt64(100)`（接口文档最大值，实现时二次确认）
- 用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装调用
- 遍历 `Response.Data.DataList`，返回 `Id == id` 的项；无匹配返回 `(nil, nil)`

Read 回调拿到 `nil` → `d.SetId("")` 标记资源已删除；否则把 `name` / `ip` / `comment` / `captcha_appid` 映射回 schema（nil-safe）。

### D6 — 同步操作，无需 waitForTaskFinish

Create / Modify / Delete 是同步接口（SDK 模型无 TaskId），无需 waitForTaskFinish。

### D7 — 文档与测试命名

- `resource_tc_captcha_ip_white_list_international.md`：包含 1 个 HCL example + Import 段。
- `resource_tc_captcha_ip_white_list_international_test.go`：基础 Create/Update/ImportStateVerify 三步骤，包名 `captcha_test`，测试函数 `TestAccTencentCloudCaptchaIpWhiteListInternationalResource_basic`。

均与 `resource_tc_config_compliance_pack.md` / `_test.go` 同款结构。

### D8 — Provider 注册位置

在 `tencentcloud/provider.go` 中 `tencentcloud_captcha_info_international` 附近追加一行：

```go
"tencentcloud_captcha_ip_white_list_international": captcha.ResourceTencentCloudCaptchaIpWhiteListInternational(),
```

## Risks / Trade-offs

- **Risk**: 复合 ID 解析失败（用户 import 传错格式）→ Mitigation: 解析时校验 `len(parts) == 2` 且两部分可转 int64，否则返回明确错误。
- **Risk**: `IdList[0]` 为 nil 或空 → Mitigation: Create 校验 `len(IdList) > 0 && IdList[0] != nil`，否则 NonRetryableError。
- **Risk**: Describe 返回列表里找不到 `Id`（白名单被外部删除）→ Mitigation: service 返回 `(nil, nil)`，Read 调 `d.SetId("")`。
- **Trade-off**: 不暴露 `status`，用户无法原地启停白名单。理由见 D1；若需要可后续单独交付。
