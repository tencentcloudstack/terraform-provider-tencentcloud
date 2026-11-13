## Context

当前 `tencentcloud_lite_hbase_instance` 资源（`tencentcloud/services/emr/resource_tc_lite_hbase_instance.go`）的 Update 函数将 `disk_size` 列入 `immutableArgs`（第 380 行），用户修改该字段会直接返回 `argument cannot be changed` 错误，无法通过 Terraform 扩容磁盘。

现有 Update 函数已支持 `zone_settings.node_num` 变更：当检测到 node_num 变更时，构造 `ModifySLInstanceRequest`（设置 `InstanceId`、`Zone`、`NodeNum`），调用 `ModifySLInstanceWithContext`，随后通过 `SLInstanceStateRefreshFunc` 等待实例状态变为 `2`（运行中）。

`ModifySLInstanceRequest` 结构体（vendor SDK，models.go:11216）除 `InstanceId`、`Zone`、`NodeNum`、`ClientToken` 外，还包含 `DiskSize` 字段（`*int64`），其注释明确：存储空间大小，需是 100 的倍数，且不允许比当前存储空间小（只允许扩容），单位 GB。`ModifySLInstanceWithContext` 方法已存在于 client.go:6351，无需升级 SDK。

异步状态轮询机制已就绪：`SLInstanceStateRefreshFunc`（service_tencentcloud_emr.go:925）通过 `GetSLInstanceStatus` 调用 `DescribeSLInstanceList` 获取实例状态码。`DescribeSLInstance` 响应（models.go:6503）的 `Status` 字段注释明确状态码含义：-2 TERMINATED、2 RUNNING、14 TERMINATING、19 ISOLATING、22 ADJUSTING、201 ISOLATED。本次扩容等待目标为 `2`，变配中状态为 `22`。

**关键文件：**
- `tencentcloud/services/emr/resource_tc_lite_hbase_instance.go:380` - immutableArgs 列表
- `tencentcloud/services/emr/resource_tc_lite_hbase_instance.go:372-447` - Update 函数
- `tencentcloud/services/emr/service_tencentcloud_emr.go:925-946` - SLInstanceStateRefreshFunc
- `vendor/.../emr/v20190103/models.go:11216-11233` - ModifySLInstanceRequest 结构体

## Goals / Non-Goals

**Goals:**
- 支持通过修改 `disk_size` 参数对 lite_hbase_instance 进行磁盘扩容
- 复用现有 `ModifySLInstanceWithContext` 接口与 `SLInstanceStateRefreshFunc` 异步等待机制
- 保持与现有 `node_num` 变更逻辑一致的处理风格（resource.Retry + BuildStateChangeConf 等待状态 2）
- 保持向后兼容，不影响未修改 disk_size 的现有行为

**Non-Goals:**
- 不支持缩容操作（云 API 仅允许扩容，由云 API 校验）
- 不在前端做 100 倍数的步长校验（由云 API 返回错误处理）
- 不修改 Schema 中 `disk_size` 字段的类型与 Required 属性（保持 `TypeInt, Required`）
- 不新增 Timeouts 块（现有代码使用 `10*tccommon.ReadRetryTimeout` 固定等待时长，本次保持一致，不引入 schema.Timeouts）

## Decisions

### Decision 1: 从 immutableArgs 移除 disk_size 而非修改 Schema
**决策：** 仅从 Update 函数的 `immutableArgs` 数组中删除 `"disk_size"` 字符串，不修改 Schema 定义。

**理由：**
- `disk_size` 在 Schema 中已是 `TypeInt, Required`，无需调整
- immutableArgs 是现有代码控制可变性的机制，移除一项即可启用更新
- 不引入新的 ForceNew/Computed 变更，保持向后兼容

**替代方案：**
- 在 Schema 中显式设置 `ForceNew: false` → 当前 Schema 未设置 ForceNew（默认即 false），无需改动

### Decision 2: disk_size 变更处理块的放置位置
**决策：** 在 Update 函数中，保留现有 `zone_settings.node_num` 变更处理逻辑不变，在其后（`needChange` 块之后）新增独立的 `d.HasChange("disk_size")` 处理块。

**实现：**
```go
if d.HasChange("disk_size") {
    request := emr.NewModifySLInstanceRequest()
    request.InstanceId = helper.String(instanceId)
    if v, ok := d.GetOkExists("disk_size"); ok {
        request.DiskSize = helper.IntInt64(v.(int))
    }
    err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifySLInstanceWithContext(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success, ...", logId, request.GetAction(), ...)
        return nil
    })
    if err != nil {
        log.Printf("[CRITAL]%s update lite_hbase_instance failed, reason:%+v", logId, err)
        return err
    }
    emrService := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
    conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.SLInstanceStateRefreshFunc(instanceId, []string{}))
    if _, e := conf.WaitForState(); e != nil {
        return e
    }
}
```

**理由：**
- 与现有 node_num 变更逻辑完全一致的处理风格（resource.Retry + BuildStateChangeConf 等待状态 2）
- disk_size 是实例级参数，不依赖 zone_settings 的遍历，独立成块更清晰
- 放在 node_num 处理之后，确保两类变更有序执行、各自独立等待状态恢复

**替代方案：**
- 合并到现有 needChange 循环中 → disk_size 不属于 zone_settings 维度，强行合并会破坏现有结构清晰度

### Decision 3: 异步等待复用 SLInstanceStateRefreshFunc
**决策：** 调用 ModifySLInstance 后，复用现有的 `emrService.SLInstanceStateRefreshFunc(instanceId, []string{})` 配合 `tccommon.BuildStateChangeConf([]string{}, []string{"2"}, ...)` 等待状态变为 2。

**理由：**
- 该函数已通过 `DescribeSLInstanceList` 获取实例 Status 字段，Status=22（ADJUSTING）时不会匹配目标 2，会继续轮询；Status=2（RUNNING）时匹配成功返回
- 与 Create、node_num 变更、Delete 的等待机制完全一致
- failStates 传空数组，与现有 node_num 变更保持一致

**替代方案：**
- 直接调用 DescribeSLInstance 接口轮询 Status → 现有刷新函数已封装好，重复实现无意义

### Decision 4: 不引入 schema Timeouts 块
**决策：** 本次变更不新增 `schema.Timeouts` 块，等待时长继续使用 `10*tccommon.ReadRetryTimeout`（与 Create、node_num 变更、Delete 保持一致）。

**理由：**
- 现有资源的所有异步等待（Create、Update node_num、Delete）均使用固定 `10*tccommon.ReadRetryTimeout`，未声明 Timeouts 块
- 引入 Timeouts 块属于额外增强，超出本次需求范围，且会改变现有 Create/Delete 的等待行为
- 保持与现有代码风格一致，降低变更风险

## Risks / Trade-offs

### Risk 1: 扩容失败导致状态不一致
**风险：** ModifySLInstance 调用成功，但等待 Status 变为 2 超时。

**缓解措施：**
- 使用 `10*tccommon.ReadRetryTimeout` 作为等待时长，与现有 node_num 变更一致
- Read 函数会从云端同步实际 disk_size 与状态，用户可 `terraform refresh` 更新状态
- 下次 apply 会基于云端实际状态继续操作

### Risk 2: 用户误操作缩容
**风险：** 用户将 disk_size 改小，期望缩容。

**缓解措施：**
- 云 API `ModifySLInstance` 的 `DiskSize` 参数明确"不允许比当前存储空间小"，API 会返回错误
- 错误经 `tccommon.RetryError` 包装后返回给用户，提示扩容限制
- 文档中说明仅支持扩容

### Risk 3: disk_size 与 node_num 同时变更时的串行等待
**风险：** 两类变配串行执行且各自等待状态 2，总耗时较长。

**缓解措施：**
- 这是必要的顺序约束：两次变配都需等待实例恢复运行中才能进行下一次
- 符合云 API 的变配串行约束
- 实际场景中同时变更两类参数较少见

## Migration Plan

**对现有资源的影响：**
1. 已创建的资源：
   - 不修改 disk_size → 无影响
   - 修改 disk_size → 从"返回错误"变为"原地扩容更新"（正向增强）
2. 用户迁移步骤：
   - 无需特殊迁移步骤
   - 用户首次修改 disk_size 时会触发扩容而非报错
3. 回滚计划：
   - 如发现严重问题，可将 `disk_size` 重新加入 immutableArgs
   - 已扩容的资源不受影响（只是后续修改会再次报错）
