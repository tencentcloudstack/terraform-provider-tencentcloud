## Context

tencentcloud_teo_origin_group 资源已经支持 HostHeader 参数的读取和更新（通过 Update API），但在创建资源时未传递该参数到 CreateOriginGroup API。SDK 中的 CreateOriginGroupRequest 结构体已包含 HostHeader 字段，且 Schema 中已定义该参数，只需要在 Create 函数中补充参数传递逻辑即可。

## Goals / Non-Goals

**Goals:**
- 在 resourceTencentCloudTeoOriginGroupCreate 函数中添加 HostHeader 参数的读取和设置
- 确保与现有 Update 和 Read 逻辑保持一致
- 保持向后兼容性（HostHeader 为 Optional 参数）

**Non-Goals:**
- 不修改 Schema 定义（已存在且正确）
- 不修改 Update 和 Read 逻辑（已正确实现）
- 不添加额外的验证逻辑（遵循 SDK 的参数验证）

## Decisions

**实现方式：**
在 CreateOriginGroupRequest 构建完成后、API 调用前，添加对 host_header 参数的处理逻辑：
```go
if v, ok := d.GetOk("host_header"); ok {
    request.HostHeader = helper.String(v.(string))
}
```
此位置与其他参数处理逻辑一致，确保代码风格统一。

**为什么选择这种方式：**
1. 与 Update 函数中的 HostHeader 处理逻辑完全一致（第 449-451 行）
2. 符合其他参数（如 name, type）的处理模式
3. 最小化代码变更，降低引入错误的风险

## Risks / Trade-offs

**风险：**
- 无显著风险（仅补全缺失的参数支持）

**权衡：**
- 无显著权衡（纯粹的参数补全，不影响其他逻辑）

## Migration Plan

无需迁移计划。这是一个非破坏性的参数补全变更：
1. 代码修改后，用户可以在资源创建时指定 HostHeader
2. 现有不指定 HostHeader 的配置继续正常工作
3. 现有通过 Update 设置 HostHeader 的配置不受影响

## Open Questions

无。
