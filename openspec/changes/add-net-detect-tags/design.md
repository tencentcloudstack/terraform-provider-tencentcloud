## Context

当前 `tencentcloud_vpc_net_detect` 资源实现了腾讯云网络探测功能，包括创建、读取、更新和删除操作。然而，该资源未支持腾讯云的标签（Tags）功能。腾讯云 API 的 `CreateNetDetect` 接口已经支持 Tags 参数，允许用户在创建网络探测资源时指定标签。同时，腾讯云提供了通用的标签管理服务（Tag Service），可以通过 `ModifyTags` API 来更新资源标签。

用户需要在 Terraform 中管理网络探测资源的标签，以便进行资源分组、权限控制和成本管理等。这是很多企业用户的重要需求。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_vpc_net_detect` 资源的 Schema 中添加 `tags` 参数（类型为 `map[string]string`，Computed 和 Optional）
- 在 Create 操作中支持通过 `CreateNetDetect` API 传入 Tags
- 在 Read 操作中读取并设置资源标签
- 在 Update 操作中支持标签更新，使用腾讯云标签服务的 `ModifyTags` API
- 确保向后兼容性，不影响现有用户配置
- 添加相应的文档和测试用例

**Non-Goals:**
- 不修改现有的资源行为，只新增 Tags 功能
- 不支持批量标签管理或标签策略功能（这些属于更高级的功能）
- 不修改资源的其他参数或行为

## Decisions

### 1. Tags 参数定义
**决策：** 在 Schema 中添加 `tags` 参数，类型为 `TypeMap`，设置为 `Optional` 和 `Computed`。

**理由：**
- 遵循腾讯云 Terraform Provider 中其他资源的一致模式（参考 `tencentcloud_cfs_file_system` 等资源）
- `Optional` 允许用户选择性添加标签
- `Computed` 确保从 API 读取的标签能够正确设置到 Terraform state

**替代方案考虑：**
- 将标签定义为 `TypeList`（List of Objects）：这是其他云厂商的常见做法，但在腾讯云 Provider 中，TypeMap 是标准模式，为了保持一致性，我们选择 TypeMap。

### 2. 创建操作中的 Tags 处理
**决策：** 在 `resourceTencentCloudVpcNetDetectCreate` 函数中，使用 `helper.GetTags()` 获取标签，并将其转换为 `CreateNetDetect` API 所需的格式。

**理由：**
- `helper.GetTags()` 是 Provider 中标准化的标签获取函数，已广泛使用
- 腾讯云 SDK 的 `CreateNetDetect` API 支持 Tags 参数，格式为 `[]*vpc.Tag`
- 这种方式与其他资源的创建逻辑一致，便于维护

### 3. 更新操作中的 Tags 处理
**决策：** 在 `resourceTencentCloudVpcNetDetectUpdate` 函数中，使用腾讯云标签服务的 `ModifyTags` API 来更新标签，而不是通过 `ModifyNetDetect` API。

**理由：**
- `ModifyNetDetect` API 可能不支持标签更新，或者标签更新有单独的 API
- 腾讯云提供了统一的标签管理服务（Tag Service），`ModifyTags` 是标准的标签更新方式
- 参考其他资源（如 CFS）的实现模式，使用 `svctag.NewTagService(tcClient)` 和 `tagService.ModifyTags()` 方法
- 需要使用 `tccommon.BuildTagResourceName()` 构建资源名称，格式为 `tencentcloud/vpc/vpc/netDetectId`（根据实际资源类型调整）

**替代方案考虑：**
- 在 `ModifyNetDetect` API 中添加标签更新：需要检查 API 文档确认是否支持。如果不支持，则必须使用标签服务 API。

### 4. 读取操作中的 Tags 处理
**决策：** 在 `resourceTencentCloudVpcNetDetectRead` 函数中，从响应中获取标签并设置到 d 中。

**理由：**
- 需要确保 Terraform state 与云资源保持同步
- 如果 API 返回的标签与 state 中的标签不一致，会触发 update 操作

### 5. 标签更新触发条件
**决策：** 在 `resourceTencentCloudVpcNetDetectUpdate` 函数中，检查 `d.HasChange("tags")` 来判断是否需要更新标签。

**理由：**
- 只在标签发生变化时才调用标签更新 API，避免不必要的 API 调用
- 减少网络开销和 API 配额消耗

## Risks / Trade-offs

### 风险 1：标签更新 API 不支持或格式不正确
**风险：** 腾讯云标签服务的 `ModifyTags` API 可能不支持网络探测资源，或者资源名称格式不正确。

**缓解措施：**
- 在实施前查看腾讯云 API 文档，确认网络探测资源是否支持标签管理
- 如果不支持，需要寻找替代方案，例如通过 `DescribeNetDetect` API 读取标签，然后通过 `CreateNetDetect` API 重新创建资源（但这会导致资源替换，不是理想方案）

### 风险 2：向后兼容性问题
**风险：** 新增 Tags 参数可能影响现有用户的配置或 state。

**缓解措施：**
- Tags 参数设置为 Optional，不设置该参数的用户不受影响
- 确保不修改现有字段的 Required 或 ForceNew 属性
- 在 Read 操作中正确处理标签缺失的情况（如 API 不返回标签时，state 中的标签应保持不变）

### 风险 3：标签读取不一致
**风险：** API 返回的标签可能与 Create 时传入的标签格式不一致，或者某些标签被过滤。

**缓解措施：**
- 在测试用例中验证标签的创建、读取和更新全流程
- 使用 `helper.GetTags()` 和 `helper.SetTags()` 等辅助函数，确保标签处理的一致性

### 权衡：性能与一致性
- 选择在 Update 操作中单独更新标签（使用标签服务 API）而不是通过 ModifyNetDetect API，可能导致更多的 API 调用
- 但这种方式更灵活，能够独立管理标签，与其他资源保持一致
- 这是一种合理的权衡，因为标签更新通常是低频操作

## Migration Plan

### 部署步骤
1. 修改 `tencentcloud/services/vpc/resource_tc_vpc_net_detect.go` 文件
   - 在 Schema 中添加 `tags` 参数
   - 在 `resourceTencentCloudVpcNetDetectCreate` 函数中添加标签处理逻辑
   - 在 `resourceTencentCloudVpcNetDetectRead` 函数中添加标签读取逻辑
   - 在 `resourceTencentCloudVpcNetDetectUpdate` 函数中添加标签更新逻辑
2. 更新 `tencentcloud/services/vpc/resource_tc_vpc_net_detect_test.go` 文件
   - 添加标签相关的测试用例
3. 更新 `website/docs/r/vpc_net_detect.md` 文档
   - 添加 Tags 参数的说明和示例

### 回滚策略
- 如果出现问题，可以回滚代码更改，移除 Tags 参数
- 由于 Tags 参数是 Optional 的，回滚不会影响现有用户配置
- 已创建的资源标签不会被删除（因为标签存储在云端，与 Provider 版本无关）

### 验证步骤
1. 运行单元测试，确保所有测试通过
2. 运行验收测试（`TF_ACC=1 go test ...`），验证标签功能
3. 在测试环境中创建带有标签的网络探测资源，验证创建、读取和更新操作

## Open Questions

1. **标签更新 API：** 腾讯云的 `ModifyTags` API 是否支持网络探测资源？需要确认资源类型名称（如 `vpc/netDetect` 或 `vpc/netdetect`）。

2. **标签创建 API：** `CreateNetDetect` API 的 Tags 参数是否直接支持？还是需要使用其他方式（如 `CreateTags` API）？

3. **资源类型名称：** 使用 `tccommon.BuildTagResourceName()` 构建资源名称时，服务名称是 "vpc"，但资源类型名称是什么？需要查阅 API 文档确认。

4. **并发更新：** 如果用户同时更新多个网络探测资源的标签，是否会触发 API 限流？需要考虑批量操作场景。
