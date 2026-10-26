# Technical Design: TKE Addon Values Data Source

## Context

TKE（容器服务）addon 通过 `tencentcloud_kubernetes_addon_config` 资源安装/更新，其 `raw_values` 需要用户自行编写 base64 编码的 JSON。用户在编写前通常需要知道两件事：

1. 该 addon **当前生效的参数**（云 API 在已安装时会用「已设置参数 + chart 默认参数」渲染后返回）
2. 该 addon chart **支持的全部参数及其默认值**（即 `DefaultValues`，据此可以知道可覆盖哪些 key）

云 API `DescribeAddonValues`（`tke/v20180525`）正好提供这两个值（均为 JSON 字符串），Provider 服务层已封装 `TkeService.DescribeAddonValuesList`，但尚未有 Terraform data source 暴露，用户无法用 `data.tencentcloud_tke_addon_values.x.values` 的方式在配置中引用。

现状与约束：

- SDK 版本已满足：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525` 中已存在 `DescribeAddonValuesRequest/Response`，**无需升级 SDK / 修改 go.mod**。
- 服务层方法已存在，签名固定为 `DescribeAddonValuesList(ctx context.Context, clusterId, addonName string) (ret *tke.DescribeAddonValuesResponseParams, errRet error)`，本次不改动。
- 该接口为**同步查询接口**（无 FlowId/TaskId），不涉及异步轮询、不需要 Timeouts 块。
- data source 无 Create/Update/Delete，只有 Read。

## Goals / Non-Goals

**Goals:**

- 新增 data source `tencentcloud_tke_addon_values`，Read 时调用 `DescribeAddonValues`，把 `Values` / `DefaultValues` 平铺暴露为顶层 computed 字段 `values` / `default_values`。
- 遵循 Provider data source 代码风格（参照 `tencentcloud_igtm_instance_list`）：`resource.Retry(tccommon.ReadRetryTimeout, ...)` + `tccommon.RetryError()`、`defer tccommon.LogElapsed` + `defer tccommon.InconsistentCheck`、`result_output_file` 支持。
- 提供 gomonkey mock 云 API 的单元测试，覆盖正常返回、nil 字段、空响应、schema 断言四类用例。

**Non-Goals:**

- 不改动 `tencentcloud_kubernetes_addon_config`、`tencentcloud_kubernetes_addons` 等已有资源/数据源。
- 不对 `values` / `default_values` 的 JSON 内容做结构化解码（保持云 API 原样字符串，交由用户用 `jsondecode()` 解析）。
- 不新增分页参数（接口无分页）。
- 不实现异步等待逻辑（接口同步返回）。

## Decisions

### 决策 1：资源命名与文件命名使用 `tke` 前缀，而非 `kubernetes`

- **选择**：资源名 `tencentcloud_tke_addon_values`；文件 `data_source_tc_tke_addon_values.go` / `.md` / `_test.go`；导出函数 `DataSourceTencentCloudTkeAddonValues()` / `dataSourceTencentCloudTkeAddonValuesRead()`。
- **备选**：`tencentcloud_kubernetes_addon_values` + `data_source_tc_kubernetes_addon_values.go`（TKE 历史文件多为 `kubernetes_` 前缀）。
- **理由**：本次需求明确指定资源名为 `tencentcloud_tke_addon_values`，且文件命名规范为 `data_source_tc_<Product>_<资源名>.go`（Product = `tke`）。tke 服务目录下两类前缀并存（如 `data_source_tc_eks_clusters.go`），不冲突。
- 仍需放在 `tencentcloud/services/tke/` 包内，包名为 `tke`。

### 决策 1.1：输入参数加非空校验

- **选择**：`cluster_id` 与 `addon_name` 都加 `ValidateFunc: validation.StringIsNotEmpty`。
- **理由（实测）**：`Required: true` **不拦空串**。`addon_name = ""` 能通过
  `terraform validate`，一路打到云端才返回
  `ResourceUnavailable: resource name may not be empty`，并被重试满 3 分钟。
  加校验后在 plan 阶段就失败，零 API 调用。

### 决策 2：直接复用服务层 `DescribeAddonValuesList`，不在 data source 内重复 retry

- **选择**：data source Read 内用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹对 `service.DescribeAddonValuesList(ctx, clusterId, addonName)` 的调用。
- **备选**：在 data source 里直接构造 `tke.NewDescribeAddonValuesRequest()` 调用 SDK。
- **理由**：
  - 服务层方法已存在且已含内层 retry + `Response == nil` 保护，符合"不要重复 retry"的约束：外层 retry 只负责把可重试错误（限流、超时）继续重试，语义与 `igtm_instance_list` 一致（那里也是 `service.DescribeIgtmInstanceListByFilter` 内含 retry，data source 外层再包一层）。
  - 避免修改服务层，缩小改动面。

### 决策 3：`values` / `default_values` 均为 `schema.TypeString`，不做 JSON 展开

- **选择**：两个出参都声明为 `TypeString, Computed`，值原样来自云 API。
- **备选**：解析 JSON 后展开成 `TypeMap` / `TypeList` 子结构。
- **理由**：
  - 云 API 注释明确"是一个 json 格式的字符串"，且不同 addon 的参数结构差异极大、层级不固定，Schema 无法静态描述。
  - 展开会引入"未知字段无法表达"与"state 漂移"风险；保持字符串可让用户用 `jsondecode(data.tencentcloud_tke_addon_values.this.values)` 自行消费。
  - 与已有 `tencentcloud_kubernetes_addons` 中 `raw_values` 的处理方式一致（字符串透传）。

### 决策 3.1：`values` 在「未安装」时**不是空**（实测修正）

原 design 假设「未安装 addon → 云端 `Values` 为 nil → `values` 为空字符串」。
真机实测（2026-10-26，`cls-dhhfi9n8`）证明该假设**只在 chart 都不存在时才成立**：

| 情形 | 实测 |
|---|---|
| 已安装（`rdma-agent`） | `Values` 与 `DefaultValues` 均为 1247B，`global.cluster.id = cls-dhhfi9n8` |
| **chart 存在但未安装**（`tke-log-agent` / `cfs` / `qgpu`） | `Values` **非空**（8865 / 2868 / 1629B），且与 `DefaultValues` **逐字节相同**；`global.cluster.id` 为**空串**，`kubeversion` 是 chart 里的硬编码值 |
| chart 不存在（名字写错） | `ResourceUnavailable: not found chart` |

结论与做法：

- provider 侧**保持原样透传**（不改代码逻辑），因为云端返回的就是这份渲染结果；
- `nil` 分支保留为防御性代码，但**不再作为文档承诺**；
- schema `Description` 与 website 文档改为实测语义，并明确指向
  `tencentcloud_kubernetes_addons` 判断安装状态；
- 单测新增 `ReadNotInstalled`，锁住「未安装也要原样透传、不清空」这一行为。

### 决策 4：`d.SetId()` 使用联合 id `cluster_id#addon_name`

- **选择**：`d.SetId(strings.Join([]string{clusterId, addonName}, tccommon.FILED_SP))`。
- **备选**：`helper.BuildToken()`（随机 token，`igtm_instance_list` 的做法）；仅用 `clusterId`（`tencentcloud_kubernetes_addons` 的做法）。
- **理由**：本次查询的唯一键是 `(cluster_id, addon_name)` 二元组，联合 id 可在日志/state 中直接定位是哪次查询，且 `tccommon.FILED_SP`（值为 `#`）是 Provider 约定的分隔符。注意 `FILED_SP` 是既有常量（拼写如此），不要写成 `FIELD_SP`。

### 决策 5：空响应不 `SetId("")`，改为 `NonRetryableError`

`values` 在 addon 未安装时云 API 返回空值（`Values` 为 nil），这是**合法的业务语义**，因此：

- `Response == nil` / `Response.Response == nil` → 由服务层返回 `NonRetryableError`，data source 直接向上返回错误。
- `Response.Values == nil` 或 `DefaultValues == nil` → 不报错，跳过对应 `d.Set()`（遵循"setXX 前先判 nil"的约束），仍设置 `d.SetId()`。
- 绝不调用 `d.SetId("")`（data source 的 id 是查询结果标识，清空会导致 state 中丢失定位信息）。

### 决策 6：`result_output_file` 写入**真实数据**

- **选择**：`tccommon.WriteToFile(output.(string), map[string]interface{}{"values": d.Get("values"), "default_values": d.Get("default_values")})`。
- **理由（实测修正）**：原方案传 `d`，而 `WriteToFile` 内部是
  `json.MarshalIndent(data, ...)`；`*schema.ResourceData` 没有导出字段、也没有实现
  `MarshalJSON`，落盘结果只能是 `{}`（实测文件 2 字节）。同包
  `data_source_tc_kubernetes_addons.go` 传的是真实 list，能正常落盘。
- 返回的 error 需检查。

## 云 API 参考（vendor 摘录，实施阶段可直接使用，无需回读大文件）

文件：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525/models.go`

### Request（行 7294-7310）

```go
// Predefined struct for user
type DescribeAddonValuesRequestParams struct {
	// 集群ID，请从容器服务控制台集群列表中获取（https://console.cloud.tencent.com/tke2/cluster）。
	ClusterId *string `json:"ClusterId,omitnil,omitempty" name:"ClusterId"`

	// addon名称，请从容器服务控制台进入集群详情页后，在组件管理页面中获取。
	AddonName *string `json:"AddonName,omitnil,omitempty" name:"AddonName"`
}

type DescribeAddonValuesRequest struct {
	*tchttp.BaseRequest

	// 集群ID
	ClusterId *string `json:"ClusterId,omitnil,omitempty" name:"ClusterId"`

	// addon名称
	AddonName *string `json:"AddonName,omitnil,omitempty" name:"AddonName"`
}
```

### Response（行 7332-7347）

```go
// Predefined struct for user
type DescribeAddonValuesResponseParams struct {
	// 参数列表，如果addon已安装，会使用已设置的参数和chart里的默认参数做渲染，是一个json格式的字符串，未安装addon时返回为空值。
	Values *string `json:"Values,omitnil,omitempty" name:"Values"`

	// addon支持的参数列表，值为chart的默认值，是一个json格式的字符串。
	DefaultValues *string `json:"DefaultValues,omitnil,omitempty" name:"DefaultValues"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeAddonValuesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAddonValuesResponseParams `json:"Response"`
}
```

### Client 方法（client.go 行 5744-5792）

```go
func NewDescribeAddonValuesRequest() (request *DescribeAddonValuesRequest)
func NewDescribeAddonValuesResponse() (response *DescribeAddonValuesResponse)
func (c *Client) DescribeAddonValues(request *DescribeAddonValuesRequest) (response *DescribeAddonValuesResponse, err error)
func (c *Client) DescribeAddonValuesWithContext(ctx context.Context, request *DescribeAddonValuesRequest) (response *DescribeAddonValuesResponse, err error)
```

可能返回的错误码：`RESOURCEUNAVAILABLE = "ResourceUnavailable"`。

### 已存在的服务层封装（`tencentcloud/services/tke/service_tencentcloud_tke_addons.go` 行 151-178）

```go
func (me *TkeService) DescribeAddonValuesList(ctx context.Context, clusterId, addonName string) (ret *tke.DescribeAddonValuesResponseParams, errRet error) {
	request := tke.NewDescribeAddonValuesRequest()
	response := tke.NewDescribeAddonValuesResponse()
	request.ClusterId = &clusterId
	request.AddonName = &addonName

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseTkeClient().DescribeAddonValues(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	ret = response.Response
	return
}
```

> 结论：**不需要新增/修改服务层代码**。实施时 data source 只需调用此方法。

## 代码风格参照

### 参照 A：`tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`（data source 骨架标准）

关键结构（按文件顺序）：

1. `package <service>` — 无文件头注释
2. import 分两组：SDK/terraform 第三方包一组，`tccommon` / `helper` 项目包一组
3. `func DataSourceTencentCloudXxx() *schema.Resource` 返回 `&schema.Resource{Read: ..., Schema: map[string]*schema.Schema{...}}`
   - 入参在前（`Required`），计算字段其次（`Computed`），`result_output_file` 固定在 schema 末尾
   - 每个字段之间空一行；每个字段都写 `Description`（英文，取自云 API 注释）
4. `func dataSourceTencentCloudXxxRead(d *schema.ResourceData, meta interface{}) error`：
   ```go
   defer tccommon.LogElapsed("data_source.tencentcloud_xxx.read")()
   defer tccommon.InconsistentCheck(d, meta)()

   var (
       logId   = tccommon.GetLogId(nil)
       ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
       service = XxxService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
   )
   // 取入参
   // resource.Retry(tccommon.ReadRetryTimeout, ...) 包裹 service 调用
   // if reqErr != nil { return reqErr }
   // 判 nil 后 d.Set
   // d.SetId(...)
   // result_output_file 处理
   return nil
   ```

### 参照 B：`tencentcloud/services/tke/data_source_tc_kubernetes_addons.go`（同服务内写法）

- 服务实例化：`service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`
- 取参：`if v, ok := d.GetOk("cluster_id"); ok { clusterId = v.(string) }`
- `d.GetOk("result_output_file")` + `tccommon.WriteToFile(output.(string), d)`

> 注意：该文件首行为 `// Code generated by iacg; DO NOT EDIT.`，是历史生成文件；**新文件不要加这行注释**（遵循"不要在资源 go 文件开头添加注释"的约束）。

## 目标文件设计

### `tencentcloud/services/tke/data_source_tc_tke_addon_values.go`

```go
package tke

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudTkeAddonValues() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTkeAddonValuesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {Type: schema.TypeString, Required: true, Description: "..."},
			"addon_name": {Type: schema.TypeString, Required: true, Description: "..."},
			"values":         {Type: schema.TypeString, Computed: true, Description: "..."},
			"default_values": {Type: schema.TypeString, Computed: true, Description: "..."},
			"result_output_file": {Type: schema.TypeString, Optional: true, Description: "Used to save results."},
		},
	}
}

func dataSourceTencentCloudTkeAddonValuesRead(d *schema.ResourceData, meta interface{}) error {
	// defer / ctx / service
	// d.Get("cluster_id"), d.Get("addon_name")  (均为 Required，可直接 d.Get)
	// resource.Retry -> service.DescribeAddonValuesList(ctx, clusterId, addonName)
	// if resp != nil { if resp.Values != nil { _ = d.Set("values", *resp.Values) } ... }
	// d.SetId(strings.Join([]string{clusterId, addonName}, tccommon.FILED_SP))
	// result_output_file -> tccommon.WriteToFile(output.(string), d)
	return nil
}
```

### Schema 字段表

| Terraform 字段 | 类型 | 必填/计算 | 对应云 API |
|---|---|---|---|
| `cluster_id` | TypeString | Required | `request.ClusterId` |
| `addon_name` | TypeString | Required | `request.AddonName` |
| `values` | TypeString | Computed | `response.Response.Values` |
| `default_values` | TypeString | Computed | `response.Response.DefaultValues` |
| `result_output_file` | TypeString | Optional | Provider 通用字段 |

### 文档 `data_source_tc_tke_addon_values.md`

格式要求（`make doc` 的输入）：

1. 首行一句话描述，带云产品名：`Use this data source to query detailed information of TKE addon values`
2. `Example Usage` 段落（HCL），`values` / `default_values` 是 JSON 字符串场景，示例必须用 `jsondecode()`
3. **无** `Import` 段（data source 没有 import）
4. **不写** `Argument Reference` / `Attribute Reference`（由工具生成）

### 测试 `data_source_tc_tke_addon_values_test.go`

- 包名：`package tke_test`（外部测试包）
- import 别名约定（两者包名都叫 `tke`，**必须**起别名避免冲突）：
  ```go
  tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
  tkesvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"
  ```
- 用 `gomonkey` mock `DescribeAddonValuesList`：
  ```go
  patches.ApplyMethodFunc(&tkesvc.TkeService{}, "DescribeAddonValuesList",
      func(_ context.Context, _ string, _ string) (*tkev20180525.DescribeAddonValuesResponseParams, error) { ... })
  ```
- 需要实现一个 `tccommon.ProviderMeta` 的 mock（该接口只有 `GetAPIV3Conn() *connectivity.TencentCloudClient` 一个方法，参照 `tencentcloud/services/cam/data_source_tc_cam_accounts_test.go` 的 `mockMetaForCamAccountsDS`：内嵌 `*connectivity.TencentCloudClient` 并返回它）。
- 被测入口：`res := tkesvc.DataSourceTencentCloudTkeAddonValues()` + `d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"cluster_id": "...", "addon_name": "..."})`，然后 `res.Read(d, meta)`。
- 用例：
  1. `ReadBasic`：返回 `Values` + `DefaultValues` 都有值 → 断言 `d.Get("values")`、`d.Get("default_values")`、`d.Id() == "cls-xxx#addon-name"`
  2. `ReadWithNilValues`：`Values` 为 nil（addon 未安装场景）→ 不报错，`default_values` 正常
  3. `ReadError`：mock 返回 error → Read 返回 error
  4. `Schema`：断言 5 个字段存在、`cluster_id`/`addon_name` 为 Required、`values`/`default_values` 为 Computed
- 不在测试中执行 `go test`；保证可编译（`-gcflags="all=-l"` 才能让 gomonkey 生效，写在测试文件注释里）。

## 注册点

### `tencentcloud/provider.go`

在 `DataSourcesMap` 中 TKE 相关块（约 601-620 行附近）追加一行，按字母序排在该块内合适位置：

```go
"tencentcloud_tke_addon_values": tke.DataSourceTencentCloudTkeAddonValues(),
```

> `provider.go` 中该 map 使用 gofmt 对齐，新增行会导致整块重新对齐，由收尾阶段 `gofmt` 处理。

### `tencentcloud/provider.md`

该文件按云产品分节，TKE 的节标题为 `Tencent Kubernetes Engine(TKE)`（约 712 行），其下有 `Data Source` 与 `Resource` 两个列表。需在 TKE 的 **`Data Source` 列表**（约 714-729 行，当前包含 `tencentcloud_kubernetes_clusters` … `tencentcloud_kubernetes_cluster_available_extra_args`）末尾追加一行：

```
tencentcloud_tke_addon_values
```

注意：不要写到 `Resource` 列表里；不要改变已有的 `Data Source` / `Resource` 行标题。

## Risks / Trade-offs

- **[`values` 在 addon 未安装时不为空（实测）]** → 云端对「chart 存在但未安装」的 addon
  返回纯 chart 默认值的渲染结果（与 `DefaultValues` 逐字节相同、不含集群信息）。
  Read 中判 nil 跳过 set 只作为防御；文档明确写「非空」并提示用
  `tencentcloud_kubernetes_addons` 判断安装状态，避免用户误判。
- **[`ResourceUnavailable` 是通用可重试码]** → `not found chart` /
  `resource name may not be empty` 这类确定性失败会被重试满 3 分钟（实测 22 次调用）。
  Read 中对这两种 message 直接返回 `resource.NonRetryableError`，其余仍走
  `tccommon.RetryError`。
- **[JSON 字符串在 state 中的格式漂移]**（云 API 每次渲染 key 顺序可能不同）→ data source 每次 read 都重新拉取、不做 diff 抑制；由于是纯 computed，不会触发 update 计划，风险仅体现在 `values` 值变化导致下游 `raw_values` 变化，属预期语义。
- **[复用服务层方法，若后续服务层签名变更]** → 影响面限于本文件一处调用；且本次不改服务层，风险低。
- **[gomonkey 依赖编译期打桩]** → 测试需 `-gcflags="all=-l"`；已在测试文件头注释说明，不影响正常运行。
- **[`tccommon.FILED_SP` 拼写]** → 常量确实是 `FILED_SP`（非 `FIELD_SP`），误用会导致编译失败；已在决策 4 明确。

## Migration Plan

新增能力，无迁移：

1. 新增 data source 文件 + 文档 + 测试
2. `provider.go` / `provider.md` 注册
3. 收尾阶段执行 `gofmt` 与 `make doc`（生成 `website/docs/d/tke_addon_values.html.markdown`）与 changelog
4. 回滚：删除新增文件并移除两处注册行即可，无 state 迁移成本（用户未引用该 data source 时零影响）

## Open Questions

无。所需信息（接口字段、服务层方法、命名、风格参照）均已在上文确定。
