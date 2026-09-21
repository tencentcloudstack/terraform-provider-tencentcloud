## Context

`DescribeTCLakeMetaInstance`（dlc v20210125）是腾讯云数据湖计算中心（DLC）提供的查询接口，用于查询 TCLake 元数据实例的开通状态。该接口**无入参**，返回 `response.Response.Status`（`*string`），枚举值如 `Running` 表示开通成功。

接口语义关键点：

- 入参：无（`DescribeTCLakeMetaInstanceRequest` 仅继承 `*tchttp.BaseRequest`，无业务字段）。
- 出参：`response.Response.Status`（`*string`，开通状态枚举值）。
- 该接口为同步查询接口，非异步接口，调用后立即返回当前开通状态，无需轮询 Read 接口。

参考实现：本 change 的数据源代码风格严格对齐 `tencentcloud_dlc_describe_user_type` —— 一个同样"无入参/单出参查询"型数据源，在 Provider 内已有标准实现样板（service 层封装 + Read 方法 + result_output_file）。

## Goals / Non-Goals

**Goals:**

- 新增 Data Source `tencentcloud_dlc_tc_lake_meta_instance`，文件命名 `data_source_tc_dlc_tc_lake_meta_instance.go`。
- 在服务层 `service_tencentcloud_dlc.go` 中封装 `DescribeTCLakeMetaInstance` 方法，调用 SDK 的 `DescribeTCLakeMetaInstanceWithContext`。
- Read 方法使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 API 调用，错误用 `tccommon.RetryError(e)` 包装。
- 将 `response.Response.Status` 映射到 schema 的 `status` computed 字段（`TypeString`），设置前判断非 nil。
- 由于接口无入参也无唯一标识字段，Data Source 的 `id` 使用固定值（如 `status` 的值或固定字符串）作为 `d.SetId()` 的值。
- 支持 `result_output_file` 可选参数，使用 `tccommon.WriteToFile` 保存查询结果。
- 遵循 datasource Read 的空返回保护：retry 块内检查 `response == nil || response.Response == nil`，若为空返回 `NonRetryableError`，避免清空 state。
- 单元测试使用 mock（gomonkey）方法对云 API 进行 mock，不使用 Terraform 测试套件。

**Non-Goals:**

- 不实现 Create/Update/Delete 操作（本资源为 RESOURCE_KIND_DATASOURCE，仅有 Read）。
- 不暴露分页参数（接口无分页字段）。
- 不修改任何既有资源/数据源/service 方法。
- 不在 Read 中再次 retry（外层 retry 已覆盖）。

## Decisions

### D1 — Schema 字段映射

| HCL 字段 | SDK 字段 | 类型 | 必填 | Computed | 说明 |
|---|---|---|---|---|---|
| `status` | `response.Response.Status` | `TypeString` | No | **Yes** | 开通状态枚举值（如 `Running`） |
| `result_output_file` | - | `TypeString` | Optional | No | 用于保存结果 |

**理由**：
- 该接口无入参，因此无需 Optional 输入字段。
- `status` 是接口唯一的业务出参，映射到 computed 字段供用户读取。
- `result_output_file` 遵循 Provider 内数据源通用约定，用于将查询结果输出到文件。

### D2 — 服务层方法封装

在 `service_tencentcloud_dlc.go` 中新增 `DescribeDlcTCLakeMetaInstance` 方法，参考 `DescribeDlcDescribeUserTypeByFilter` 的写法：

```go
func (me *DlcService) DescribeDlcTCLakeMetaInstance(ctx context.Context) (status *string, errRet error) {
    var (
        logId   = tccommon.GetLogId(ctx)
        request = dlc.NewDescribeTCLakeMetaInstanceRequest()
    )
    defer func() {
        if errRet != nil {
            log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
        }
    }()
    ratelimit.Check(request.GetAction())
    response, err := me.client.UseDlcClient().DescribeTCLakeMetaInstance(request)
    if err != nil {
        errRet = err
        return
    }
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
    if response == nil || response.Response == nil || response.Response.Status == nil {
        return
    }
    status = response.Response.Status
    return
}
```

**理由**：
- 接口无入参，服务层方法无需接收 paramMap，直接调用即可。
- 空指针保护覆盖 `response == nil`、`response.Response == nil`、`response.Response.Status == nil` 三层。

### D3 — Read 方法实现

```go
func dataSourceTencentCloudDlcTCLakeMetaInstanceRead(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("data_source.tencentcloud_dlc_tc_lake_meta_instance.read")()
    defer tccommon.InconsistentCheck(d, meta)()

    logId := tccommon.GetLogId(tccommon.ContextNil)
    ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

    service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
    var status *string
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        result, e := service.DescribeDlcTCLakeMetaInstance(ctx)
        if e != nil {
            return tccommon.RetryError(e)
        }
        if result == nil {
            return resource.NonRetryableError(fmt.Errorf("dlc tc_lake_meta_instance read empty, skip SetId"))
        }
        status = result
        return nil
    })
    if err != nil {
        return err
    }

    if status != nil {
        _ = d.Set("status", status)
    }

    d.SetId("tc_lake_meta_instance")
    output, ok := d.GetOk("result_output_file")
    if ok && output.(string) != "" {
        if e := tccommon.WriteToFile(output.(string), status); e != nil {
            return e
        }
    }
    return nil
}
```

**理由**：
- retry 块内检查 `result == nil`（即服务层返回空 status），返回 `NonRetryableError` 并打印 `[DATASOURCE] read empty, skip SetId` 提示，避免短暂波动导致 state 被清空。
- 由于接口无入参也无唯一标识字段，`d.SetId("tc_lake_meta_instance")` 使用固定字符串作为 id（参考其他无入参数据源的做法）。
- 设置 `status` 前判断非 nil。

### D4 — Provider 注册位置

在 `tencentcloud/provider.go` 既有 dlc 数据源注册段追加：

```go
"tencentcloud_dlc_tc_lake_meta_instance": dlc.DataSourceTencentCloudDlcTCLakeMetaInstance(),
```

同时在 `tencentcloud/provider.md` 的 DLC Data Source 列表追加资源名一行，确保 gendoc 能扫描并生成 website doc。

### D5 — 文档与测试命名

- `data_source_tc_dlc_tc_lake_meta_instance.md`：含一句话描述（带上 DLC 云产品名称）+ Example Usage，不包含 Argument Reference/Attribute Reference（由工具自动生成），非 datasource 类型才有的 Import 部分不适用。
- `data_source_tc_dlc_tc_lake_meta_instance_test.go`：使用 mock（gomonkey）方法对云 API 进行 mock 处理，只进行业务代码逻辑的单元测试，不使用 Terraform 测试套件。

### D6 — 测试策略

由于本资源为新增的 Terraform 数据源，单元测试文件 `data_source_tc_dlc_tc_lake_meta_instance_test.go` 使用 mock（gomonkey）方法对云 API 进行 mock 处理，只进行业务代码逻辑的单元测试，不调用真实的云 API，不使用 Terraform 的测试套件。

## Risks / Trade-offs

- **Risk**: 接口无入参，Data Source 无法基于输入参数区分多个实例 → Mitigation: TCLake 开通状态为账户级单例，固定 id `tc_lake_meta_instance` 符合业务语义；同一账户仅有一个状态查询结果。
- **Risk**: 云 API 短暂波动返回空 status 时，若直接 `d.SetId("")` 会清空 state → Mitigation: retry 块内返回 `NonRetryableError`，让外层 retry 继续尝试，避免数据丢失。
- **Trade-off**: 使用固定字符串作为 id 而非真实资源 id → 因接口本身无资源 id 概念（仅返回开通状态），这是合理的简化；不影响 Terraform 的 state 管理与刷新逻辑。
