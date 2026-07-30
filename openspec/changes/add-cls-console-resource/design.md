## Context

CLS DataSight 控制台 4 个新接口（CreateConsole/DescribeConsoles/ModifyConsole/DeleteConsole）发布于 cls v20201016 SDK，适合做成 Terraform CRUD 资源。资源命名遵守 provider 现有惯例 `tencentcloud_<service>_<resource>`，归属到 `tencentcloud/services/cls/`。

CreateConsole 入参字段较多（15 个），其中 4 个为嵌套结构体 list（`Accounts`、`AuthRoles`、`AccessControlRules`、`Tags`），1 个为嵌套结构体（`AnonymousLogin`）。3 个用于"账号密码鉴权 / 匿名登录 / 第三方认证登录"三种登录模式的条件字段。本文档解决 schema 字段映射、嵌套结构体的 schema 形态、Update 行为、Read 中的查询路径以及登录模式互斥校验等关键决策。

参考资源：`tencentcloud_igtm_monitor`（代码风格）、`tencentcloud_config_compliance_pack`（文档与测试命名）、provider 中现有 cls 资源（`UseClsClient()` 客户端调用约定、service 层方法签名）。

## Goals / Non-Goals

**Goals:**

- Schema 字段名与 CreateConsole 接口入参 1:1 映射（snake_case 化）。
- 资源 ID = `ConsoleId`，由 Create 响应给出。
- Read 通过 `DescribeConsoles` + `Filters: [{Key: "ConsoleId", Values: [d.Id()]}]` 精确查询单个实例；分页 limit 取接口最大值 100。
- Create / Read / Update / Delete 的全部 SDK 调用均用 `resource.Retry(WriteRetryTimeout, ...)` 或 `resource.Retry(ReadRetryTimeout, ...)` 包装。
- 全部接口返回值取值都做空指针保护。
- 代码风格严格对齐 `tencentcloud_igtm_monitor`（变量声明顺序、注释、defer LogElapsed + InconsistentCheck 模式）。
- 资源 doc 命名 `resource_tc_cls_console.md`，测试命名 `resource_tc_cls_console_test.go`，皆与参考样板一致。

**Non-Goals:**

- 不实现对应的 data source（`tencentcloud_cls_consoles`）——本 change 只交付 resource。
- 不在 provider 层做登录模式与字段的强校验（如"匿名登录必传 AnonymousLogin"），让 API 层报错引导用户；仅在文档说明。
- 不修改任何已存在资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射（与 CreateConsole 入参 1:1）

| HCL 字段 | SDK 字段 | 类型 | 必填 | ForceNew | 说明 |
|---|---|---|---|---|---|
| `access_mode` | `AccessMode` | `TypeList`(`TypeString`) | Required | No | `public` / `internal` |
| `login_mode` | `LoginMode` | `TypeInt` | Required | No | 0/1/2 |
| `domain_prefix` | `DomainPrefix` | `TypeString` | Required | No | 自定义域名前缀 |
| `accounts` | `Accounts` | `TypeList`(`TypeMap`) | Optional | No | 登录方式 0 必传；嵌套 `user_name`、`password`、`secret_id`、`secret_key`、`email` |
| `anonymous_login` | `AnonymousLogin` | `TypeList`(`TypeMap`, MaxItems:1) | Optional | No | 登录方式 1 必传；嵌套 `secret_id`、`secret_key` |
| `intranet_type` | `IntranetType` | `TypeInt` | Optional | No | 默认 0 |
| `intranet_region` | `IntranetRegion` | `TypeString` | Optional | No |  |
| `vpc_id` | `VpcId` | `TypeString` | Optional | No |  |
| `subnet_id` | `SubnetId` | `TypeString` | Optional | No |  |
| `auth_roles` | `AuthRoles` | `TypeList`(`TypeMap`) | Optional | No | 登录方式 2 必传；嵌套 `role_name`、`secret_id`、`secret_key` |
| `tags` | `Tags` | `TypeList`(`TypeMap`) | Optional | **Yes** | 嵌套 `key`、`value`；ModifyConsole 不接受 |
| `hide_params` | `HideParams` | `TypeList`(`TypeString`) | Optional | No |  |
| `access_control_rules` | `AccessControlRules` | `TypeList`(`TypeMap`) | Optional | No | 嵌套 `access_mode`（这是 SDK `AccessControlRule` 唯一字段） |
| `remarks` | `Remarks` | `TypeString` | Optional | No |  |
| `menus` | `Menus` | `TypeList`(`TypeString`) | Optional | No |  |
| `console_id` | `ConsoleId` | `TypeString` | Computed | — | 资源 ID |
| `domain` | `Domain` | `TypeString` | Computed | — | 公网访问域名 |
| `intranet_domain` | `IntranetDomain` | `TypeString` | Computed | — | 内网访问域名 |

**理由**：嵌套结构体（`Accounts`、`AnonymousLogin`、`AuthRoles`、`Tags`、`AccessControlRules`）使用 `TypeList` 而非 `TypeSet`，因为：
- 与 igtm 参考资源一致（嵌套结构都用 List）。
- API 不要求顺序无关，user 顺序即 API 顺序。
- 简化 plan diff 渲染（List 按位置比较，更直观）。

`anonymous_login` 用 `TypeList + MaxItems:1` 模拟单个嵌套对象，是 SDK v2 表达 "single nested object" 的标准做法（v2 没有 `TypeObject`）。

### D2 — `tags` 标记为 ForceNew

**问题**：CreateConsole 接受 `Tags`，但 ModifyConsole 不接受。如果 `tags` 是 Optional+No-ForceNew，用户改 `tags` 会触发 Update，但 ModifyConsole 没法把新 tags 推送到后端 → 漂移。

**决策**：`tags` 字段标 `ForceNew: true`。改 tags 触发资源销毁重建，避免静默漂移。

**替代方案**：在 Update 中忽略 tags 变化（用 `DiffSuppressFunc` 或单独 helper），但这会让用户在 plan 中看到"零变化"却无法实际改 tags，更糟糕。ForceNew 的"销毁重建"语义对用户透明。

### D3 — Update 字段集合（mutableArgs）

参考 `tencentcloud_igtm_monitor` 的 mutableArgs 模式：

```go
mutableArgs := []string{
    "access_mode", "login_mode", "domain_prefix", "accounts",
    "anonymous_login", "intranet_type", "intranet_region", "vpc_id",
    "subnet_id", "auth_roles", "hide_params", "access_control_rules",
    "remarks", "menus",
}
// 注：tags 是 ForceNew（D2），不在 mutableArgs 内
// console_id 是 Computed，不可改
```

任意字段变化触发一次 ModifyConsole 调用，把所有可变字段全量提交（API 是 PUT 语义，全量覆盖；与 ModifyConsole 行为一致）。

### D4 — Read 路径

`DescribeConsoles` 不支持 GetById 形式，只能用 `Filters` 过滤。Read 实现：

```go
request.Filters = []*cls.Filter{{
    Key:    helper.String("ConsoleId"),
    Values: []*string{helper.String(consoleId)},
}}
request.Limit = helper.Int64(100)  // 接口最大值
```

返回 `Consoles []*Console`：
- `len == 0` → 资源已被外部删除 → `d.SetId("")` 并 return nil
- `len == 1` → 取第一项
- `len > 1` → 不应发生（Filter 精确匹配 ConsoleId），但兜底取第一项并 log Warning

`DescribeClsConsoleById` 放在 `service_tencentcloud_cls.go`，签名：

```go
func (me *ClsService) DescribeClsConsoleById(ctx context.Context, consoleId string) (ret *cls.Console, errRet error)
```

### D5 — 异步任务等待

CreateConsole / ModifyConsole / DeleteConsole 是同步接口（经查 SDK 模型，无 RequestId 之外的 task 字段），**无需 waitForTaskFinish**。Delete 后立即调 Describe 也能拿到 nil 即视为已删除（参考 igtm 资源同款处理）。

### D6 — 文档与测试命名

- `resource_tc_cls_console.md`：包含 1-2 个 HCL example + Import 段。
- `resource_tc_cls_console_test.go`：基础 Create/Update/ImportStateVerify 三步骤，包名 `cls_test`，测试函数 `TestAccTencentCloudClsConsoleResource_basic`。

均与 `resource_tc_config_compliance_pack.md` / `_test.go` 同款结构。

### D7 — Provider 注册位置

在 `tencentcloud/provider.go` 第 ~2049 行（既有 cls 资源注册段）中追加一行：

```go
"tencentcloud_cls_console": cls.ResourceTencentCloudClsConsole(),
```

字母序保持紧邻其他 `tencentcloud_cls_*` 资源。

## Risks / Trade-offs

- **Risk**: 用户用单一 HCL 同时设置三种登录模式互斥字段（如同时给 `accounts` 和 `anonymous_login`）→ Mitigation: API 层会报错；文档明确"按 `login_mode` 选择对应字段"。不在 schema 强加 `ConflictsWith`，避免限制后端未来可能的扩展。
- **Risk**: ModifyConsole 不接受 Tags，用户期待改 tags → Mitigation: D2 用 ForceNew 把行为显式化，plan 阶段就能让用户看到"销毁重建"。
- **Risk**: DescribeConsoles 的 Filters 列表没有官方"精确匹配 ConsoleId"的强保证 → Mitigation: D4 在 `len>1` 时取第一条 + log warning。实际从 API 行为看 ConsoleId 全局唯一，不会冲突。
- **Trade-off**: 不实现 data source（`tencentcloud_cls_consoles`）。理由是本 change 聚焦 resource；data source 可以后续单独 PR 跟进。
