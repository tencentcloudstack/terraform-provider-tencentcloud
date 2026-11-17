# fwm-cfw-managed-config Specification

## Purpose
TBD - created by archiving change add-fwm-cfw-managed-config. Update Purpose after archive.
## Requirements
### Requirement: 资源定义与 Schema
资源 `tencentcloud_fwm_cfw_managed` SHALL 为配置型资源，管理单个集团成员的云防火墙纳管开关。Schema SHALL 包含三个字段：
- `member_id`：目标成员的 MemberId，Type=String，Required，用于标识资源；import 时提供。
- `app_id`：目标成员的 AppId，Type=String，Required + ForceNew，用于 Read 时按 AppId 过滤定位成员。
- `cfw_managed`：纳管开关，Type=Int，Required，取值 1（开启）/ 0（关闭）。

字段 description SHALL 使用英文。资源 description SHALL 在一句话描述中带上所属云产品名称（fwm）。

#### Scenario: Schema 字段定义正确
- **WHEN** 定义 `ResourceTencentCloudFwmCfwManaged()` 的 Schema
- **THEN** `member_id` 为 TypeString/Required，`app_id` 为 TypeString/Required/ForceNew，`cfw_managed` 为 TypeInt/Required

#### Scenario: 资源一句话描述包含云产品名
- **WHEN** 生成资源 .md 文档
- **THEN** 描述格式为 "Provides a resource to ..." 且包含 fwm 云产品名称

### Requirement: Read 读取纳管状态
Read SHALL 通过 `DescribeOrganMembers` 接口，使用 `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]`（1=EQUAL）按 AppId 过滤成员列表，再从返回的 Members 中按 MemberId 精确匹配目标成员，读取其 `cfw_managed` 字段并 set 到 state。
Read SHALL NOT 传 `Filters[{Name:"MemberId"}]`（MemberId 不在过滤映射表内，会导致「filter name does not meet the requirements」错误）。
当返回列表为空或匹配不到 MemberId 时，Read SHALL 视为成员已离开集团/资源不存在，打印日志后 `d.SetId("")` 移除 state，不报错。
Read SHALL 在调用 set 字段前判断 Response 中字段是否为 nil，nil 则不调用 set 方法。

#### Scenario: 按 AppId 过滤并按 MemberId 匹配读取
- **WHEN** Read 从 `d.Id()` 解析出 app_id 与 member_id，调用 DescribeOrganMembers 传入 AppId 过滤
- **THEN** 从返回 Members 中按 MemberId 匹配目标成员，set cfw_managed 到 state

#### Scenario: 成员不存在时移除 state
- **WHEN** DescribeOrganMembers 返回空列表或匹配不到 MemberId
- **THEN** 打印日志后 `d.SetId("")`，不报错

#### Scenario: 不传 MemberId 过滤
- **WHEN** 构建 DescribeOrganMembers 请求
- **THEN** Filters 中只含 AppId，不出现 Name="MemberId"

### Requirement: Update 设置纳管开关
Update SHALL 通过 `ModifyCfwManaged` 接口设置 `CfwManaged`（0/1）与 `MemberId`。
Update SHALL 在 cfw_managed 与当前值相同时跳过写操作（避免无谓写入与操作日志噪音），不同才调 ModifyCfwManaged。
ModifyCfwManaged 为同步接口（无 FlowId/TaskId），Update SHALL NOT 生成 waitReady/轮询逻辑。
调用云 API 失败时，Update SHALL 使用 `tccommon.RetryError()` 包装错误，并在 retry 块外处理错误。设置 id 等成功操作 SHALL 放到 retry 块外。
当关闭纳管（cfw_managed 1→0）被服务端以 InvalidParameter 拒绝时，Update SHALL 原样透传错误，不静默改 state。

#### Scenario: cfw_managed 变更时调用 ModifyCfwManaged
- **WHEN** Update 检测到 cfw_managed 与当前值不同
- **THEN** 调用 ModifyCfwManaged 传入 member_id 与新 cfw_managed 值

#### Scenario: cfw_managed 未变更时跳过写操作
- **WHEN** Update 检测到 cfw_managed 与当前值相同
- **THEN** 不调用 ModifyCfwManaged，直接返回

#### Scenario: 关闭纳管被拒时透传错误
- **WHEN** ModifyCfwManaged 设置 CfwManaged=0 被服务端拒绝（存在已下发规则等）
- **THEN** 原样返回 InvalidParameter 错误，不改 state

### Requirement: Create 为 empty 实现
Create SHALL 为 empty 空实现：从 schema 读取 member_id 并 `d.SetId(member_id)`，随后复用 Update 调用 ModifyCfwManaged 设开关。Create MUST 在 `d.SetId()` 后再执行后续逻辑。

#### Scenario: Create 设置 ID 后复用 Update
- **WHEN** Create 被调用
- **THEN** 先 `d.SetId(member_id)`，再调用 Update 设置 cfw_managed 开关

### Requirement: Delete 为 empty 实现
Delete SHALL 为 empty 空实现：仅 `return nil`，从 state 移除即可，SHALL NOT 调用任何云 API。理由：不能「删除纳管」，只能关开关。

#### Scenario: Delete 不调云 API
- **WHEN** Delete 被调用
- **THEN** 直接返回 nil，不调用 ModifyCfwManaged 或任何其他云 API

### Requirement: Import 两段式 ID
资源 SHALL 支持 import，import 使用 `app_id` + `member_id` 两段式复合 ID，分隔符为 `tccommon.FILED_SP`（即 `#`）。Read/Update/各方法 SHALL 从 `d.Id()` 用 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析出 app_id 与 member_id，解析失败（len < 2）时返回错误。

#### Scenario: import 两段式 ID 解析
- **WHEN** 用户执行 `terraform import tencentcloud_fwm_cfw_managed.xxx "appId#memberId"`
- **THEN** Read 从 d.Id() 解析出 app_id=appId、member_id=memberId 并完成读取

#### Scenario: ID 格式错误时报错
- **WHEN** d.Id() 不含分隔符（len < 2）
- **THEN** 返回 "id is broken" 错误

### Requirement: connectivity 层接入 fwm 客户端
SHALL 在 `tencentcloud/connectivity/client.go` 中新增：
- import `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- `TencentCloudClient` 结构体新增 `fwmConn *fwm.Client` 字段
- `UseFwmClient() *fwm.Client` 方法（参考 UseCfwClient 模式）

#### Scenario: UseFwmClient 可用
- **WHEN** 资源代码调用 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmClient()`
- **THEN** 返回有效的 fwm.Client 实例

### Requirement: Provider 注册与文档
SHALL 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `"tencentcloud_fwm_cfw_managed": fwm.ResourceTencentCloudFwmCfwManaged()`。
SHALL 在 `tencentcloud/provider.md` 中新增 fwm 产品章节条目。
SHALL 生成 `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config.md` 文档样例（含一句话描述、Example Usage、Import 部分；不含 Argument Reference / Attribute Reference）。

#### Scenario: provider.go 注册资源
- **WHEN** 查看 provider.go
- **THEN** ResourcesMap 中存在 `tencentcloud_fwm_cfw_managed` 条目

#### Scenario: 文档样例含 Import 部分
- **WHEN** 生成 .md 文档
- **THEN** 含一句话描述（带 fwm）、Example Usage、Import 部分（说明使用 app_id#member_id 两段式 ID）

### Requirement: 单元测试（mock/gomonkey）
SHALL 在 `tencentcloud/services/fwm/resource_tc_fwm_cfw_managed_config_test.go` 中使用 mock（gomonkey）方法对云 API 进行 mock，只进行业务代码逻辑的单元测试，SHALL NOT 使用 terraform 验收测试套件。
测试 SHALL 覆盖 Create、Read、Update、Delete 的业务逻辑（包括 cfw_managed 变更/未变更分支、成员不存在分支、ID 解析）。

#### Scenario: mock 单测覆盖核心逻辑
- **WHEN** 运行单元测试
- **THEN** 使用 gomonkey mock DescribeOrganMembers / ModifyCfwManaged，验证 Create/Read/Update/Delete 业务逻辑正确

