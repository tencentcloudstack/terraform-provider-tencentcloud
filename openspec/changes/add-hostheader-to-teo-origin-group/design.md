## Context

当前 `tencentcloud_teo_origin_group` 资源已经在 schema 中定义了 `host_header` 字段（Optional 类型），并且在 Update 和 Read 函数中都有处理该参数：
- Update 函数在调用 ModifyOriginGroup API 时会传递 `request.HostHeader`
- Read 函数会从 API 响应中读取并设置 `respData.HostHeader`

然而，Create 函数在调用 CreateOriginGroup API 时没有传递 `HostHeader` 参数，导致用户无法在创建源站组时直接设置该参数，只能创建后再通过 update 操作来设置。

从 SDK 的 CreateOriginGroupRequest 结构体可以看到，该 API 支持 `HostHeader *string` 参数，说明底层 API 是支持在创建时设置该参数的。

## Goals / Non-Goals

**Goals:**
- 在 `resourceTencentCloudTeoOriginGroupCreate` 函数中添加对 `host_header` 参数的处理
- 将用户配置的 `host_header` 值传递给 CreateOriginGroup API 的 `request.HostHeader` 参数
- 确保实现方式与现有的 Update 函数保持一致，使用相同的 helper 函数和错误处理模式

**Non-Goals:**
- 不修改 schema 定义（字段已存在）
- 不修改 Update 和 Read 函数（已正确处理）
- 不新增验证逻辑或错误处理（遵循现有模式）
- 不修改文档或测试（除非发现遗漏）

## Decisions

**1. 在 Create 函数中添加 HostHeader 参数处理**

选择在 Create 函数中添加 HostHeader 参数处理，而不是在 schema 中新增一个强制字段，原因如下：
- HostHeader 参数在 API 中是可选的（omitnil 标签）
- 保持与 Update 函数的处理方式一致
- 不影响现有配置，向后兼容

实现方式：在处理完 `records` 参数后，添加对 `host_header` 参数的处理代码：
```go
if v, ok := d.GetOk("host_header"); ok {
    request.HostHeader = helper.String(v.(string))
}
```

**2. 使用与 Update 函数相同的参数获取方式**

使用 `d.GetOk("host_header")` 而不是 `d.Get("host_header")`，原因：
- GetOk 会检查字段是否设置，避免传递空字符串
- 与 Update 函数的实现保持一致
- 符合 Terraform Plugin SDK v2 的最佳实践

## Risks / Trade-offs

**Risk 1: API 参数验证**
- 风险：如果用户传入无效的 HostHeader 值，API 可能返回错误
- 缓解：依赖底层 API 的参数验证，不添加额外的验证逻辑，保持与现有代码一致

**Risk 2: 与 Update 函数的同步性**
- 风险：Create 和 Update 函数的 HostHeader 处理逻辑可能不同步
- 缓解：复制 Update 函数的实现代码，确保完全一致

**Trade-off: 参数处理位置**
- 选择：在 Create 函数中直接处理，而不是提取为公共函数
- 理由：该参数仅在 Create 和 Update 中使用，提取公共函数会增加代码复杂度，收益不大
