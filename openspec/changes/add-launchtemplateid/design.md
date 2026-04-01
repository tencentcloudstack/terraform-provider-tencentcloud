## Context

当前 `tencentcloud_cvm_launch_template` 资源已实现基本的创建和读取功能，但是缺少对 `LaunchTemplateId` 字段的支持。该字段在 CreateLaunchTemplate 和 DescribeLaunchTemplates API 响应中返回，表示实例启动模板的唯一标识符。

该资源文件位于 `tencentcloud/services/cvm/resource_tencentcloud_cvm_launch_template.go`，使用 Terraform Plugin SDK v2 和 tencentcloud-sdk-go 与腾讯云 API 交互。

## Goals / Non-Goals

**Goals:**
- 在 Schema 中新增 `launch_template_id` 字段，设置为 Computed 属性
- 在 Create 函数中，从 CreateLaunchTemplate API 响应的 `LaunchTemplateId` 字段读取并设置到 state
- 在 Read 函数中，从 DescribeLaunchTemplates API 响应的 `LaunchTemplateId` 字段读取并更新 state
- 更新单元测试和验收测试，确保字段正确读取和展示
- 保持向后兼容，不破坏现有用户配置和 state

**Non-Goals:**
- 不修改现有的 ID 字段或主键逻辑
- 不修改 Delete 或 Update 函数的行为
- 不添加新的 API 调用或外部依赖

## Decisions

**Schema 设计**
- 字段名称：`launch_template_id`（遵循 Terraform 命名约定，使用 snake_case）
- 字段类型：`schema.TypeString`
- 属性：`Computed: true`（因为该字段由 API 返回，用户不需要设置）
- 不设置 `Optional: true`，因为用户不应该手动设置此值

**数据流设计**
- Create 阶段：在调用 `CreateLaunchTemplate` API 后，从响应中提取 `LaunchTemplateId` 并使用 `d.Set("launch_template_id", value)` 设置到 state
- Read 阶段：在调用 `DescribeLaunchTemplates` API 后，从响应中提取 `LaunchTemplateId` 并使用 `d.Set("launch_template_id", value)` 更新 state
- 使用 `tccommon.GetSdkValue()` 辅助函数从 API 响应中读取字段值

**错误处理**
- 如果 API 响应中不包含 `LaunchTemplateId` 字段，忽略该字段（因为该字段为可选）
- 不因为缺少该字段而报错，保持与 API 定义的一致性

**测试策略**
- 单元测试：mock API 响应，验证 `launch_template_id` 字段正确设置
- 验收测试：实际调用 API，验证字段在创建和读取时正确返回
- 测试用例覆盖：创建后立即读取，验证 ID 字段值

## Risks / Trade-offs

**风险 1：API 响应中字段名称变化**
- 风险：如果腾讯云 API 修改了响应字段名称，会导致字段读取失败
- 缓解：使用 `tccommon.GetSdkValue()` 辅助函数，该函数支持灵活的字段路径，如果字段名称变化，只需修改字段路径即可

**风险 2：向后兼容性**
- 风险：新增字段可能影响现有用户配置
- 缓解：字段设置为 Computed（而非 Optional），用户不需要在配置中设置该字段，现有配置不受影响

**权衡：是否在文档中强调该字段**
- 选择：在 resource 样例文件中添加注释说明该字段是 Computed 属性
- 理由：帮助用户理解该字段的性质，避免用户尝试手动设置

## Migration Plan

**部署步骤：**
1. 修改 `resource_tencentcloud_cvm_launch_template.go` 文件
2. 更新 `resource_tencentcloud_cvm_launch_template_test.go` 测试文件
3. 运行单元测试验证修改
4. 运行验收测试验证实际 API 调用
5. 更新资源样例文档（可选，添加字段说明）

**回滚策略：**
- 如果测试失败或有兼容性问题，可以直接删除新增的 `launch_template_id` 字段定义和相关的 Set 调用
- 由于该字段是 Computed 且不参与资源 ID 计算，删除该字段不会影响现有 state

## Open Questions

无（变更范围清晰，无需进一步决策）
