## Context

腾讯云 BDRC（多活容灾）产品的"容灾站点对 VPC 映射"用于在双活/容灾架构下建立源端 VPC/子网与目标端 VPC/子网之间的网络映射关系。本次新增的资源对应三个云 API（SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330`，已 vendor）：

- `CreateDisasterRecoveryVpcMapping`：入参 `SourceVpcId`、`SourceSubnetId`、`TargetVpcId`、`TargetSubnetId`、`SitePairId`（全部 `*string`）；**出参仅 `RequestId`，不返回映射 ID**。
- `DescribeVpcMappings`：入参 `SitePairId`（`*string`）、`Filters`（`[]*FilterModel`，支持 `source-vpc-id` / `target-vpc-id` / `source-subnet-id` / `target-subnet-id`）、`Offset`、`Limit`（注释标注最大值 100）；出参 `TotalCount` 与 `VpcMappingSet []*VpcMapping`，其中 `VpcMapping` 含 `Id *uint64`、`SitePairId`、`SourceVpc`、`SourceSubnet`、`TargetVpc`、`TargetSubnet`、`Status`、`LifeState`（可能返回 null）。
- `DeleteDisasterRecoveryVpcMapping`：入参 `VpcMappingIds []*uint64`；出参仅 `RequestId`。

当前 Provider 内 **不存在** `tencentcloud/services/bdrc/` 目录，且 `tencentcloud/connectivity/client.go` 中没有 bdrc 客户端，需要一并新增。

代码风格严格参考 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（RESOURCE_KIND_GENERAL 样板）。

## Goals / Non-Goals

**Goals:**

- 新增通用型资源 `tencentcloud_bdrc_disaster_recovery_vpc_mapping`，实现 Create / Read / Delete 三个 CRUD 回调；Update 回调不做任何云 API 调用（无可用接口），直接转 Read。
- Schema 字段与云 API 参数映射：
  - 入参映射：`source_vpc_id`→`SourceVpcId`、`source_subnet_id`→`SourceSubnetId`、`target_vpc_id`→`TargetVpcId`、`target_subnet_id`→`TargetSubnetId`、`site_pair_id`→`SitePairId`；
  - 出参平铺映射（不引入 `vpc_mapping_set` 嵌套层）：`id`→`VpcMappingSet[].Id`、`source_vpc`→`VpcMappingSet[].SourceVpc`、`source_subnet`→`VpcMappingSet[].SourceSubnet`、`target_vpc`→`VpcMappingSet[].TargetVpc`、`target_subnet`→`VpcMappingSet[].TargetSubnet`、`status`→`VpcMappingSet[].Status`、`life_state`→`VpcMappingSet[].LifeState`；`site_pair_id` 同为入参与出参字段。
- 资源 ID 采用联合 ID `sitePairId + tccommon.FILED_SP + vpcMappingId`，Read/Delete 从 `d.Id()` 拆分填充请求参数；支持 import（RESOURCE_KIND_GENERAL），import 示例说明使用联合 ID。
- Create 后回查：因 Create 响应不返回 ID，Create 成功后用 `DescribeVpcMappings` 以 `source-vpc-id` + `source-subnet-id` 过滤查询新建映射，取 `VpcMappingSet` 中匹配项的 `Id` 作为 `vpcMappingId` 写入联合 ID。
- 所有云 API 调用（Create/Describe/Delete）均包裹 `resource.Retry(tccommon.WriteRetryTimeout / tccommon.ReadRetryTimeout, ...)`，错误经 `tccommon.RetryError()` 包装。
- 新增 `UseBdrcV20260330Client()` 与 `bdrcV20260330Conn` 字段到 `tencentcloud/connectivity/client.go`。
- 在 `tencentcloud/provider.go` / `tencentcloud/provider.md` 注册资源。
- 新增 `.md` 资源文档（一句话描述 + Example Usage + Import 说明，不手写 Argument/Attribute Reference）。
- 新增 `*_test.go` 单元测试，使用 gomonkey mock 云 API（不使用 terraform 测试套件）。

**Non-Goals:**

- 不实现配套数据源（如 `tencentcloud_bdrc_vpc_mappings`）——本次需求仅为 RESOURCE_KIND_GENERAL 资源。
- 不实现 Update 逻辑——云 API 无 Update/Modify 类接口；全部字段 ForceNew，变更即销毁重建。
- 不修改任何既有资源/数据源/service 方法。
- 不在收尾阶段之外生成 `website/docs/` 文档或 `.changelog/` 文件。

## Decisions

### D1 — Schema 定义（CRD-only，全部字段 ForceNew）

| HCL 字段 | 云 API 来源 | Type | Required | Optional | Computed | ForceNew |
|---|---|---|---|---|---|---|
| `site_pair_id` | Create/Describe 入参 `SitePairId`；Describe 出参 `VpcMappingSet[].SitePairId` | TypeString | Yes | - | - | **Yes** |
| `source_vpc_id` | Create 入参 `SourceVpcId` | TypeString | Yes | - | - | **Yes** |
| `source_subnet_id` | Create 入参 `SourceSubnetId` | TypeString | Yes | - | - | **Yes** |
| `target_vpc_id` | Create 入参 `TargetVpcId` | TypeString | Yes | - | - | **Yes** |
| `target_subnet_id` | Create 入参 `TargetSubnetId` | TypeString | Yes | - | - | **Yes** |
| `id`（出参平铺） | Describe 出参 `VpcMappingSet[].Id`（uint64） | TypeInt | - | - | Yes | - |
| `source_vpc`（出参平铺） | Describe 出参 `VpcMappingSet[].SourceVpc` | TypeString | - | - | Yes | - |
| `source_subnet`（出参平铺） | Describe 出参 `VpcMappingSet[].SourceSubnet` | TypeString | - | - | Yes | - |
| `target_vpc`（出参平铺） | Describe 出参 `VpcMappingSet[].TargetVpc` | TypeString | - | - | Yes | - |
| `target_subnet`（出参平铺） | Describe 出参 `VpcMappingSet[].TargetSubnet` | TypeString | - | - | Yes | - |
| `status`（出参平铺） | Describe 出参 `VpcMappingSet[].Status` | TypeString | - | - | Yes | - |
| `life_state`（出参平铺） | Describe 出参 `VpcMappingSet[].LifeState` | TypeString | - | - | Yes | - |

**理由**：

- 该资源只有 CRD 接口，按项目规则第 7 条："只将 Id() 字段设置成 ForceNew，并在资源 update 方法中将其余顶层字段加入 immutableArgs 数组"。但由于本项目该规则同时要求所有入参字段在无 Update 接口时不可变，且 terraform SDK 的惯用做法是：CRD-only 资源将全部业务字段设为 `ForceNew: true`，Update 回调直接 `return Read`，从而任何字段变更自动触发销毁重建，无需额外 immutableArgs 校验（immutableArgs 模式用于"有顶层可变字段但无 Update 接口"的场景；本资源所有字段均 ForceNew，`d.HasChange` 永远不会走到 Update，销毁重建由 SDK 自动完成）。为满足规则第 7 条的防御性要求，Update 回调中仍声明 `immutableArgs` 数组并校验，命中即返回 error，双保险。
- 出参字段按项目规则第 13 条平铺展开：云 API 响应本身是列表（`VpcMappingSet` 数组），按"列表中第一项"的字段结构定义 schema，禁止出现 `vpc_mapping_set` 嵌套层。
- `id`（出参平铺）命名与 terraform 保留的 Resource ID 不冲突：schema 中声明 `id` 作为 Computed 字段是本项目平铺模式的常规做法（如 `resource.Response.XxxSet` 列表资源），Resource ID 仍由 SDK 内部管理（联合 ID）。注意 Read 中 set 该字段时使用 `d.Set("id", ...)`，与 `d.Id()`（Resource ID）互不干扰。
- Create 入参 5 个字段在描述中使用小写蛇形资源名 `bdrc_disaster_recovery_vpc_mapping` 风格的错误信息与日志。

### D2 — Create 流程：创建 + 回查取 ID

`CreateDisasterRecoveryVpcMapping` 响应不含业务 ID，采用"创建后回查"模式：

1. 组装请求（5 个入参从 schema 读取，`helper.String(...)` 包装）。
2. `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `CreateDisasterRecoveryVpcMappingWithContext`；失败用 `tccommon.RetryError(e)` 包装；成功后校验 `result == nil || result.Response == nil` → `NonRetryableError`（规则第 9 条：必须检查返回值是否为空）。
3. Create 成功后（retry 块外），用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `DescribeVpcMappingsWithContext` 回查：
   - `SitePairId` 填入；
   - `Filters` 填两个过滤器：`source-vpc-id` → `[sourceVpcId]`、`source-subnet-id` → `[sourceSubnetId]`；
   - `Limit` 填 100（注释标注的最大值，规则第 5 条）；
   - 校验 `result == nil || result.Response == nil || len(result.Response.VpcMappingSet) == 0` → `NonRetryableError`（回查失败宁可让上层重试耗尽失败，也不写入空 ID）；
   - 在 `VpcMappingSet` 中找到 `SourceVpc`/`SourceSubnet`/`TargetVpc`/`TargetSubnet` 与请求一致的项，取其 `Id`；若 `Id == nil` → `NonRetryableError`。
4. 回查成功后（retry 块外）`d.SetId(strings.Join([]string{sitePairId, vpcMappingIdStr}, tccommon.FILED_SP))`，最后 `return Read(d, meta)`。

**为什么用 source-vpc-id + source-subnet-id 过滤**：Describe 的 Filters 支持这四个键，同一源端 VPC+子网在同一站点对下通常唯一；再叠加代码内比对四个字段做兜底匹配，避免命中同源不同目标的映射。

**备选方案（否决）**：Create 后直接 `d.SetId(sitePairId)`、Read 时遍历全量映射靠源端字段匹配 —— 会导致 state 与真实资源对应关系弱，destroy 时可能删错/漏删，不可取。

### D3 — Read 流程：联合 ID 拆分 + 按 ID 过滤查询

1. `idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`，`len(idSplit) != 2` 时返回 `fmt.Errorf("id is broken,%s", d.Id())`。
2. `sitePairId := idSplit[0]`；`vpcMappingId := idSplit[1]`（字符串形式的 uint64）。
3. `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `DescribeVpcMappingsWithContext`：
   - `SitePairId` 填入；
   - `Limit` 填 100；
   - 不使用 Filters（直接全量拉取后按 `Id` 匹配，避免过滤键遗漏）；
   - 校验 `result == nil || result.Response == nil` → `RetryError`；
   - 在 `VpcMappingSet` 中按 `strconv.FormatUint(*item.Id, 10) == vpcMappingId` 匹配；找不到 → 视为已删除。
4. retry 块外：若未找到，**先** `log.Printf("[CRUD] bdrc_disaster_recovery_vpc_mapping id=%s", d.Id())` 保留现场，**再** `d.SetId("")`（规则第 8 条）；找到则逐字段 set（set 前判 nil，规则第 8 条）：
   - `item.Id != nil` → `d.Set("id", *item.Id)`；
   - `item.SitePairId != nil` → `d.Set("site_pair_id", ...)`；
   - `item.SourceVpc != nil` → `d.Set("source_vpc", ...)`；
   - `item.SourceSubnet != nil` → `d.Set("source_subnet", ...)`；
   - `item.TargetVpc != nil` → `d.Set("target_vpc", ...)`；
   - `item.TargetSubnet != nil` → `d.Set("target_subnet", ...)`；
   - `item.Status != nil` → `d.Set("status", ...)`；
   - `item.LifeState != nil` → `d.Set("life_state", ...)`。

**说明**：Read 是资源（非数据源）回调，映射不存在时按规则第 8 条处理（打印现场后 `d.SetId("")`）；规则第 14 条"返回空时 NonRetryableError"仅适用于 RESOURCE_KIND_DATASOURCE，此处不适用——但 `result == nil || result.Response == nil` 属于接口调用异常，仍在 retry 内返回 `RetryError` 继续重试。

### D4 — Update 流程：immutableArgs 校验 + 转 Read

云 API 无 Update 接口。schema 中全部业务字段已设 `ForceNew`，terraform 不会把字段变更送进 Update。按项目规则第 7 条做双保险：

```go
needChange := false
immutableArgs := []string{"site_pair_id", "source_vpc_id", "source_subnet_id", "target_vpc_id", "target_subnet_id"}
for _, v := range immutableArgs {
    if d.HasChange(v) {
        needChange = true
        break
    }
}
if needChange {
    return fmt.Errorf("Update bdrc_disaster_recovery_vpc_mapping is not supported, all business arguments are immutable (CRD-only API), please recreate the resource.")
}
return resourceTencentCloudBdrcDisasterRecoveryVpcMappingRead(d, meta)
```

### D5 — Delete 流程

1. 拆分联合 ID 得到 `sitePairId`、`vpcMappingId`；`helper.StrToUint64Point(vpcMappingId)` 转换。
2. `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `DeleteDisasterRecoveryVpcMappingWithContext`：
   - `request.VpcMappingIds = []*uint64{vpcMappingIdPoint}`；
   - 失败 `tccommon.RetryError(e)`；成功校验 `result == nil || result.Response == nil` → `NonRetryableError`。
3. 不做删除后回查（接口为同步删除语义，出参仅 RequestId；若后续需要可依赖 terraform refresh 阶段的 Read）。

### D6 — connectivity 客户端新增

在 `tencentcloud/connectivity/client.go` 中（参考 `UseIgtmV20231024Client`）：

```go
import bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

// TencentCloudClient struct 增加字段:
bdrcV20260330Conn *bdrcv20260330.Client

// UseBdrcV20260330Client return BDRC client for service
func (me *TencentCloudClient) UseBdrcV20260330Client() *bdrcv20260330.Client {
    if me.bdrcV20260330Conn != nil {
        return me.bdrcV20260330Conn
    }
    cpf := me.NewClientProfile(300)
    me.bdrcV20260330Conn, _ = bdrcv20260330.NewClient(me.Credential, me.Region, cpf)
    me.bdrcV20260330Conn.WithHttpTransport(&LogRoundTripper{})
    return me.bdrcV20260330Conn
}
```

### D7 — Provider 注册

- `tencentcloud/provider.go`：新增 import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"`，在资源 map 中追加 `"tencentcloud_bdrc_disaster_recovery_vpc_mapping": bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()`。
- `tencentcloud/provider.md`：新增 BDRC 段落并登记资源名一行。

### D8 — 文档与测试命名

- `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.md`：
  - 一句话描述："Provides a resource to create a disaster recovery VPC mapping for site pair of BDRC."（带产品名 BDRC）；
  - Example Usage（HCL）；
  - Import 部分（RESOURCE_KIND_GENERAL 有 Import），说明使用联合 ID `sitePairId#vpcMappingId`；
  - 不手写 Argument Reference / Attribute Reference（由工具自动生成）。
- `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping_test.go`：使用 gomonkey mock `DescribeVpcMappingsWithContext` / `CreateDisasterRecoveryVpcMappingWithContext` / `DeleteDisasterRecoveryVpcMappingWithContext`，只测业务代码逻辑；禁止通过 `go test` 执行。

## Risks / Trade-offs

- **Risk**: Create 成功但回查窗口内映射尚未可见（最终一致性）→ Mitigation: 回查本身包裹 `resource.Retry(tccommon.ReadRetryTimeout, ...)`，空结果返回 `NonRetryableError` 前……不对——回查放在 retry 内，空列表时返回 `resource.NonRetryableError` 会立即终止；应返回 `RetryError` 风格的可重试错误以等待最终一致。**修正**：回查 retry 块内"列表为空"返回 `resource.RetryableError`（继续重试直到超时），"接口调用出错"返回 `tccommon.RetryError(e)`；仅"列表非空但无法匹配/Id 为 nil"这类确定性异常返回 `NonRetryableError`。
- **Risk**: 同一源端 VPC+子网在同一站点对下存在多条映射时回查误配 → Mitigation: 回查时逐条比对四个字段（SourceVpc/SourceSubnet/TargetVpc/TargetSubnet）与 Create 入参一致才算命中；命中多于一条时报错人工介入。
- **Risk**: Delete 接口入参是 ID 列表，误传多条会删除其他映射 → Mitigation: 严格只传 `[]*uint64{当前资源ID}` 单元素。
- **Trade-off**: 无 Update 接口 ⇒ 所有字段 ForceNew，变更成本高（销毁重建）；这是云 API 能力边界决定的，文档中明确说明。
- **Trade-off**: `id` 出参平铺字段与 Terraform Resource ID 语义易混淆 → 文档 Import 示例中明确说明 Resource ID 是 `sitePairId#vpcMappingId` 联合格式。
