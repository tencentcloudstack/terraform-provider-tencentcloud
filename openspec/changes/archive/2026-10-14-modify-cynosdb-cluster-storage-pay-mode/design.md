## Context

当前 `tencentcloud_cynosdb_cluster` 资源将 `storage_pay_mode` 参数列为不可修改（在 Update 函数的 `immutableArgs` 列表中，约第 653 行），导致用户无法通过 Terraform 将集群的预付费存储转为后付费存储。腾讯云 CynosDB 提供 `TransferStoragePrepayToPostpay` 接口（https://cloud.tencent.com/document/product/1003/135369）用于将预付费存储转为后付费存储，该接口为异步接口，调用成功后返回 `DealNames`，需要通过 `DescribeResourcesByDealName` 接口轮询等待转换完成。

**约束：**
- `TransferStoragePrepayToPostpay` 仅支持预付费存储转后付费（单向），云 API 无反向接口
- 接口为异步，调用成功后需通过 `DescribeResourcesByDealName` 轮询（计费回调发货前可能返回 `InvalidParameterValue.DealNameNotFound`，需重试）
- `storage_pay_mode` 取值：`0`-按量计费（后付费），`1`-包年包月（预付费）
- `storage_pay_mode` 在 Schema 中已是 `Optional, Computed`，无需修改 Schema 定义

**关键文件：**
- `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go:647-664` - immutableArgs 列表与校验逻辑
- `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go:249-277` - create 模块中 `DescribeResourcesByDealName` 的异步等待参考实现
- `tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go` - 服务层 API 封装

**云 API 接口说明：**
- `TransferStoragePrepayToPostpay`：请求参数 `ClusterId`（string），返回 `BigDealIds`、`DealNames`、`ResourceIds`、`ClusterIds`
- `DescribeResourcesByDealName`：请求参数 `DealName`/`DealNames`，返回 `BillingResourceInfos`（含 `ClusterId`、`InstanceIds`、`DealName`）

## Goals / Non-Goals

**Goals:**
- 支持通过修改 `storage_pay_mode` 参数将 CynosDB 集群的预付费存储（`1`）转为后付费存储（`0`）
- 复用 create 模块中 `DescribeResourcesByDealName` 的异步等待模式，等待存储付费模式转换完成
- 提供明确的单向限制错误提示（不支持后付费转预付费）
- 保持向后兼容性，不影响现有资源行为

**Non-Goals:**
- 不支持后付费存储转预付费存储（云 API 无此能力）
- 不修改 `storage_pay_mode` 的 Schema 定义（保持 `Optional, Computed`）
- 不修改 create 函数中 `storage_pay_mode` 的设置逻辑

## Decisions

### Decision 1: 移除 immutableArgs 限制并新增变更处理
**决策：** 从 Update 函数的 `immutableArgs` 列表中移除 `storage_pay_mode`，并在 Update 函数中新增 `d.HasChange("storage_pay_mode")` 的处理分支。

**理由：**
- 与项目中其他可修改字段（如 `storage_limit`、`cluster_name` 等）的处理模式一致
- immutableArgs 机制用于"只有 CRD 接口、无 update 接口"的字段，而本次有对应的转换接口
- 放宽限制符合用户预期，允许通过 Terraform 完成转换

**替代方案：**
- 保留 immutableArgs 并返回错误 → 无法满足需求，用户无法通过 Terraform 转换
- 使用 ForceNew → 不合理，转换存储付费模式不应重建整个集群

### Decision 2: 单向转换校验
**决策：** 在 Update 函数的 `storage_pay_mode` 处理分支中，使用 `d.GetChange("storage_pay_mode")` 获取新旧值，仅当旧值为 `1`（预付费）、新值为 `0`（后付费）时调用 `TransferStoragePrepayToPostpay`；当反向修改（`0`→`1`）时返回明确错误。

**实现：**
```go
if d.HasChange("storage_pay_mode") {
    oldPayMode, newPayMode := d.GetChange("storage_pay_mode")
    if oldPayMode.(int) != 1 || newPayMode.(int) != 0 {
        return fmt.Errorf("storage_pay_mode only supports transfer from prepaid(1) to postpaid(0)")
    }
    // 调用 TransferStoragePrepayToPostpay
}
```

**理由：**
- 云 API `TransferStoragePrepayToPostpay` 仅支持预付费转后付费，无反向接口
- 明确的错误信息帮助用户理解限制，避免向云 API 发送无效请求
- 相等值（无变化）不会触发 `d.HasChange`，无需特殊处理

**替代方案：**
- 仅检查新值是否为 `0` → 不够严谨，可能误处理 `0`→`0` 之外的其他场景（实际 `d.HasChange` 已过滤）
- 不做校验直接调用 API 让云侧报错 → 错误信息不友好，增加无效 API 调用

### Decision 3: 异步等待复用 create 模块的 DescribeResourcesByDealName 模式
**决策：** 调用 `TransferStoragePrepayToPostpay` 成功后，从返回值 `DealNames` 中取第一个 DealName，调用 `DescribeResourcesByDealName` 进行重试轮询，直到查询成功（计费回调发货完成），参考 create 模块（第 249-277 行）的实现方式。

**实现：**
```go
// 调用 TransferStoragePrepayToPostpay（在服务层封装，使用 WriteRetryTimeout 重试）
dealNames, err := cynosdbService.TransferStoragePrepayToPostpay(ctx, clusterId)
if err != nil {
    return err
}
if len(dealNames) < 1 {
    return fmt.Errorf("TransferStoragePrepayToPostpay response DealNames is empty, cluster id=%s", clusterId)
}

// 异步等待：DescribeResourcesByDealName 轮询
dealName := dealNames[0]
dealReq := cynosdb.NewDescribeResourcesByDealNameRequest()
dealReq.DealName = dealName
err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    dealRes, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeResourcesByDealName(dealReq)
    if e != nil {
        if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
            if sdkErr.Code == "InvalidParameterValue.DealNameNotFound" {
                return resource.RetryableError(fmt.Errorf("DealName[%s] Not Found, retry... reason: %s", *dealName, e.Error()))
            }
        }
        return tccommon.RetryError(e)
    }
    if dealRes == nil || dealRes.Response == nil || len(dealRes.Response.BillingResourceInfos) < 1 {
        return resource.RetryableError(fmt.Errorf("waiting for storage pay mode transfer, DealName[%s]", *dealName))
    }
    return nil
})
if err != nil {
    return err
}
```

**理由：**
- 与 create 模块的异步等待实现完全一致，降低维护成本
- `DescribeResourcesByDealName` 在计费未回调时会返回 `InvalidParameterValue.DealNameNotFound`，需作为可重试错误处理
- 使用 `tccommon.ReadRetryTimeout` 作为超时时间，符合项目规范

**替代方案：**
- 使用 `DescribeFlow` 轮询 → `TransferStoragePrepayToPostpay` 不返回 FlowId，无法使用此方式
- 使用 `DescribeClusterDetail` 轮询 `StoragePayMode` 字段 → 转换过程中该字段状态不确定，且需求明确要求参考 create 模块用 `DescribeResourcesByDealName`

### Decision 4: 服务层方法封装
**决策：** 在 `service_tencentcloud_cynosdb.go` 中新增 `TransferStoragePrepayToPostpay` 服务方法，封装 API 调用与 `WriteRetryTimeout` 重试，返回 `DealNames` 列表。

**实现：**
```go
func (me *CynosdbService) TransferStoragePrepayToPostpay(ctx context.Context, clusterId string) (dealNames []*string, errRet error) {
    logId := tccommon.GetLogId(ctx)
    request := cynosdb.NewTransferStoragePrepayToPostpayRequest()
    request.ClusterId = &clusterId

    errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        response, err := me.client.UseCynosdbClient().TransferStoragePrepayToPostpay(request)
        if err != nil {
            log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
            return tccommon.RetryError(err)
        }
        if response != nil && response.Response != nil {
            dealNames = response.Response.DealNames
        }
        return nil
    })
    if errRet != nil {
        return
    }
    return
}
```

**理由：**
- 与项目中其他服务层方法（如 `ModifyClusterStorage`、`SwitchClusterVpc` 等）的模式一致
- 将 API 调用与重试逻辑封装在服务层，Update 函数保持简洁
- 使用的 `DescribeResourcesByDealName` 轮询逻辑留在 Update 函数中（与 create 模块一致，create 中也是直接在 Create 函数内调用该接口）

### Decision 5: 变更处理位置
**决策：** 将 `storage_pay_mode` 的变更处理分支放在 Update 函数中现有字段处理逻辑之间，建议置于 `storage_limit` 处理之后（约第 1057 行之后），二者均与存储相关，逻辑聚集便于维护。

**理由：**
- 与存储相关的变更操作集中放置，提升代码可读性
- 不影响其他字段的更新顺序

## Risks / Trade-offs

### Risk 1: 异步转换超时导致状态不一致
**风险：** `TransferStoragePrepayToPostpay` 调用成功，但 `DescribeResourcesByDealName` 轮询超时。

**缓解措施：**
- 使用 `tccommon.ReadRetryTimeout` 作为轮询超时时间，与 create 模块一致
- 转换请求已提交到云侧，Read 函数会从云端同步实际 `storage_pay_mode` 值，用户可 `terraform refresh` 更新状态
- 下次 apply 会基于云端实际状态继续操作
- 在错误信息中包含 clusterId 与 DealName 便于排障

### Risk 2: 用户误操作后付费转预付费
**风险：** 用户将 `storage_pay_mode` 从 `0` 改为 `1`，期望执行反向转换。

**缓解措施：**
- 在 Update 函数中校验新旧值，仅允许 `1`→`0`，反向返回明确错误
- 错误信息明确说明"仅支持从预付费(1)转为后付费(0)"
- 文档中突出说明此单向限制

### Risk 3: TransferStoragePrepayToPostpay 返回空 DealNames
**风险：** API 调用成功但返回的 `DealNames` 为空。

**缓解措施：**
- 调用成功后检查 `len(dealNames) < 1`，若为空则返回错误并打印 clusterId
- 避免后续 `DescribeResourcesByDealName` 因空 DealName 产生异常

## Migration Plan

**对现有资源的影响：**
1. **已创建的资源**：
   - 不修改 `storage_pay_mode` → 无影响
   - 修改 `storage_pay_mode`（`1`→`0`）→ 执行转换（原行为是报错）
   - 修改 `storage_pay_mode`（`0`→`1`）→ 报错（原行为也是报错，错误信息更明确）

2. **用户迁移步骤**：
   - 无需特殊迁移步骤
   - 用户首次将 `storage_pay_mode` 从 `1` 改为 `0` 时触发转换

3. **回滚计划**：
   - 如发现严重问题，可将 `storage_pay_mode` 重新加入 `immutableArgs`
   - 已转换的集群不受影响（只是后续修改会触发报错）

## Open Questions

1. **Q**: 是否需要在转换完成后再次调用 Read 接口校验 `StoragePayMode` 已变为 `0`？
   **A**: `DescribeResourcesByDealName` 查询成功即代表计费回调发货完成，转换已生效。Update 函数末尾会调用 `resourceTencentCloudCynosdbClusterRead` 同步状态，无需额外校验。

2. **Q**: 是否需要处理 `BigDealIds`？
   **A**: 不需要。异步等待使用 `DealNames`（与 create 模块一致），`BigDealIds` 为预付费总订单号，不用于资源查询。
