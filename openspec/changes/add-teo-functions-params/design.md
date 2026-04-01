## Context

当前 `tencentcloud_teo_function` 资源缺少从 DescribeFunctions API 读取完整函数信息的能力。TEO (Tencent Edge One) 服务提供了 DescribeFunctions API，可以返回包含函数 ID、站点 ID、函数名、描述、内容、域名、创建时间和修改时间等完整信息的 Functions 列表。为了支持在 Terraform 中完整管理和查询 TEO 函数配置，需要将这些字段映射到资源 Schema 中。

当前的资源实现在 `tencentcloud/services/teo/resource_tc_teo_function.go`，使用 Terraform Plugin SDK v2 和 tencentcloud-sdk-go 调用云 API。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_function` 资源 Schema 中新增 `functions` 列表字段及其所有子字段
- 更新 Read 函数以从 DescribeFunctions API 响应中读取 Functions 列表并正确映射到资源字段
- 确保新字段的属性（Computed/Optional）与 CAPI 接口定义一致
- 更新单元测试和验收测试以验证新字段的功能
- 保持向后兼容性，不破坏现有 Terraform 配置和 state

**Non-Goals:**
- 不修改现有资源的核心 CRUD 逻辑（除了新增字段的读取）
- 不改变 Create/Update/Delete API 的调用方式
- 不引入新的外部依赖或架构变更

## Decisions

### 1. 字段属性设计
**决策：** 将 `functions` 列表及其所有子字段设置为 Computed 字段（只读），因为这些字段由云 API 返回，不需要用户在 Terraform 配置中设置。

**理由：**
- DescribeFunctions API 返回的函数信息是只读的元数据
- 用户不需要手动设置这些字段，它们由云服务自动生成
- 符合 Terraform Provider 的最佳实践，将只读信息标记为 Computed

### 2. 数据映射方式
**决策：** 在 Read 函数中，将 DescribeFunctions API 响应中的 Functions 列表直接映射到 `functions` Schema 字段，使用 `d.Set()` 方法设置每个子字段的值。

**理由：**
- 简单直接的数据映射，易于维护
- 利用 Terraform Plugin SDK v2 的列表类型支持
- 保持代码清晰，便于测试和调试

### 3. 测试策略
**决策：** 更新现有单元测试 `resource_tc_teo_function_test.go`，添加对新字段的断言；更新验收测试以验证实际 API 调用正确返回新字段。

**理由：**
- 确保新字段在单元测试层面被正确映射
- 验证与真实云 API 的集成
- 遵循项目的测试规范

## Risks / Trade-offs

**风险 1：** API 响应字段变更导致映射失败
- **缓解：** 使用 `helper.Retry()` 处理 API 调用的最终一致性和错误情况；在测试中验证字段映射的正确性

**风险 2：** 新字段名称与现有字段冲突
- **缓解：** 使用 `functions` 作为列表字段名，子字段使用 snake_case 命名（如 `function_id`, `zone_id`），遵循 Terraform Provider 的命名规范

**权衡：** 将所有子字段设置为 Computed 可能限制某些高级用例，但这是基于当前 API 设计的最合理选择，未来可以根据需求调整。

## Migration Plan

由于本次变更仅为新增 Computed 字段，不需要特殊的迁移步骤：

1. 更新 `resource_tc_teo_function.go` 中的 Schema 定义
2. 更新 Read 函数以映射新字段
3. 更新单元测试和验收测试
4. 提交代码并发布新版本

**回滚策略：** 如果出现问题，可以通过删除新增字段来快速回滚，因为这些字段是 Computed 的，不会影响现有用户的配置。

## Open Questions

无。当前需求和实现方案已明确。