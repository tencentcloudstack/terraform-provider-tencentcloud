## Context

当前 `tencentcloud_teo_origin_group` 资源已经定义了 `host_header` schema 参数（Optional），并且在 Read 和 Update 操作中实现了该参数的处理逻辑：
- Read 函数从 API 响应中读取 `HostHeader` 并设置到 state
- Update 函数将 `host_header` 传递给 `ModifyOriginGroupRequest.HostHeader`

但是在 Create 函数中，虽然 Schema 中定义了该参数，且云 API 的 `CreateOriginGroupRequest` 结构体支持 `HostHeader` 字段，但代码中缺少将 `host_header` schema 参数映射到 API 请求的逻辑。这导致用户无法在创建资源时设置 `host_header`，必须在创建后再通过 update 操作来设置。

当前代码位置：`tencentcloud/services/teo/resource_tc_teo_origin_group.go`

## Goals / Non-Goals

**Goals:**
- 在 `resourceTencentCloudTeoOriginGroupCreate` 函数中添加 `HostHeader` 参数的设置逻辑
- 确保与现有的 Read 和 Update 逻辑保持一致
- 验证云 API 确实支持该参数（已在 vendor 目录确认）

**Non-Goals:**
- 不修改 Schema 定义（Schema 已正确）
- 不修改 Read 和 Update 逻辑（已正确实现）
- 不添加新的 API 调用或错误处理
- 不更新文档（文档已存在）

## Decisions

**1. 在 Create 函数中添加 HostHeader 参数映射**

在 `resourceTencentCloudTeoOriginGroupCreate` 函数中，在构建 `CreateOriginGroupRequest` 时添加以下代码（在设置 `records` 之后）：

```go
if v, ok := d.GetOk("host_header"); ok {
    request.HostHeader = helper.String(v.(string))
}
```

**理由：**
- 与 Update 函数中的实现模式保持一致（第 449-451 行）
- 使用 `d.GetOk()` 检查参数是否设置，避免传递 nil 值
- 使用 `helper.String()` 将 Go string 转换为 *string，匹配 API 请求类型
- 代码位置选择在设置 `records` 之后（第 220 行之后），保持与参数在 Schema 中的顺序一致

**2. 不添加额外的错误处理或验证**

**理由：**
- 云 API 自身的参数验证会处理无效值
- 现有的错误处理机制（`defer tccommon.InconsistentCheck()`）已足够
- Read 函数会读取 API 响应并更新 state，确保一致性
- 其他参数（如 `name`、`type`）也没有额外的验证逻辑

## Risks / Trade-offs

**Risk 1: 用户可能设置了无效的 HostHeader 值**
- **Mitigation:** 云 API 会验证参数，返回错误会被现有的错误处理机制捕获并返回给用户

**Risk 2: 可能有其他未发现的依赖或副作用**
- **Mitigation:** 通过单元测试和集成测试验证完整 CRUD 流程；代码变更范围小，风险可控

**Trade-off:**
- **简单性 vs 验证：** 选择最小化代码变更，不添加额外的本地验证，依赖云 API 的验证。这减少了维护负担，但意味着错误信息可能不够友好。

## Migration Plan

此变更不需要数据迁移：
- 纯代码变更，不涉及数据库或持久化结构
- 对已有资源无影响（已创建的资源可通过 update 操作设置 host_header）
- 向后兼容：不修改已有 API 或 schema

**部署步骤：**
1. 提交代码变更
2. 运行测试验证
3. 发布新版本的 provider

**Rollback 策略：**
- 如果发现问题，直接回滚代码提交即可
- 不会导致已创建资源的状态损坏

## Open Questions

无
