# clb-target-group-query Specification

## Purpose
TBD - created by archiving change replace-clb-target-group-api. Update Purpose after archive.
## Requirements
### Requirement: REQ-CLB-TG-QUERY-001 - 使用 DescribeTargetGroupList API 查询目标组

The `tencentcloud_clb_target_group` resource and related services MUST use the `DescribeTargetGroupList` API instead of `DescribeTargetGroups` API when querying target groups.

`tencentcloud_clb_target_group` 资源和相关服务在查询目标组时必须使用 `DescribeTargetGroupList` API 替代 `DescribeTargetGroups` API。

**理由**: 
- `DescribeTargetGroupList` 是腾讯云推荐的查询接口
- 提供更好的性能和一致性
- 两个接口返回相同的数据结构,可平滑迁移

**影响范围**:
- `ClbService.DescribeTargetGroups()` 方法
- `ClbService.DescribeClbTargetGroupAttachmentsById()` 方法
- 所有依赖这些方法的资源和数据源

#### Scenario: 通过 ID 查询单个目标组

**Given**: 用户需要查询特定 ID 的目标组详情  
**When**: 调用 `ClbService.DescribeTargetGroups(ctx, targetGroupId, filters)`  
**Then**: 
- 底层使用 `DescribeTargetGroupList` API
- 请求参数: `TargetGroupIds = [targetGroupId]`
- 返回目标组详情: `[]*clb.TargetGroupInfo`
- 如果不存在,返回空数组

**验收标准**:
- ✅ API 请求日志显示 `DescribeTargetGroupList`
- ✅ 返回数据结构与旧接口一致
- ✅ 单元测试通过

---

#### Scenario: 通过过滤器查询目标组列表

**Given**: 用户需要查询特定 VPC 或名称的目标组  
**When**: 调用 `ClbService.DescribeTargetGroups(ctx, "", filters)`  
**Then**:
- 底层使用 `DescribeTargetGroupList` API
- 请求参数: `Filters = [{Name: "TargetGroupVpcId", Values: [vpcId]}]`
- 返回匹配的目标组列表: `[]*clb.TargetGroupInfo`
- 支持分页查询 (Offset/Limit)

**验收标准**:
- ✅ 支持 `TargetGroupVpcId` 过滤
- ✅ 支持 `TargetGroupName` 过滤
- ✅ 分页逻辑正确(自动遍历所有页)

---

#### Scenario: 查询目标组绑定信息

**Given**: 需要查询目标组的 CLB 绑定关系  
**When**: 调用 `ClbService.DescribeClbTargetGroupAttachmentsById(ctx, targetGroups, associationsSet)`  
**Then**:
- 底层使用 `DescribeTargetGroupList` API
- 请求参数: `TargetGroupIds = targetGroups`
- 返回目标组的 `AssociatedRule` 信息
- 过滤出匹配 `associationsSet` 的绑定关系

**验收标准**:
- ✅ 批量查询支持
- ✅ 绑定关系解析正确
- ✅ 过滤逻辑正确

---

#### Scenario: 资源 Read 操作查询目标组

**Given**: Terraform 执行 `terraform refresh` 或 `terraform plan`  
**When**: 触发 `tencentcloud_clb_target_group` 资源的 Read 操作  
**Then**:
- 调用 `ClbService.DescribeTargetGroups(ctx, id, filters)`
- 底层使用 `DescribeTargetGroupList` API
- 更新 Terraform state 中的目标组属性

**验收标准**:
- ✅ `terraform refresh` 成功
- ✅ State 数据完整准确
- ✅ 不影响现有配置

---

#### Scenario: 数据源查询目标组列表

**Given**: Terraform 配置使用 `data.tencentcloud_clb_target_groups`  
**When**: 执行 `terraform plan` 或 `terraform apply`  
**Then**:
- 调用 `ClbService.DescribeTargetGroups(ctx, targetGroupId, filters)`
- 底层使用 `DescribeTargetGroupList` API
- 返回符合条件的目标组列表

**验收标准**:
- ✅ 数据源查询成功
- ✅ 过滤条件生效
- ✅ 输出数据完整

---

### Requirement: REQ-CLB-TG-QUERY-002 - API 调用兼容性

The new and old APIs MUST maintain full compatibility in request and response structures.

新旧 API 必须保持请求和响应结构的完全兼容,确保平滑迁移,不引入破坏性变更。

#### Scenario: 请求参数兼容

**Given**: 现有代码使用的请求参数  
**When**: 替换为新 API  
**Then**:
- `TargetGroupIds` 参数保持不变
- `Filters` 参数保持不变
- `Offset` 参数保持不变
- `Limit` 参数保持不变

**验收标准**:
- ✅ 无需修改参数设置代码
- ✅ 编译通过

---

#### Scenario: 响应结构兼容

**Given**: 现有代码解析的响应结构  
**When**: 替换为新 API  
**Then**:
- `TotalCount` 字段保持不变
- `TargetGroupSet` 字段保持不变
- `TargetGroupInfo` 结构保持不变

**验收标准**:
- ✅ 无需修改响应解析代码
- ✅ 数据完整性保持

---

### Requirement: REQ-CLB-TG-QUERY-003 - 错误处理一致性

The error handling logic of the new API MUST be consistent with the old API.

新 API 的错误处理逻辑必须与旧 API 保持一致。

#### Scenario: 目标组不存在

**Given**: 查询一个不存在的目标组 ID  
**When**: 调用 API  
**Then**:
- 返回空结果数组
- 不抛出错误
- 日志记录查询操作

**验收标准**:
- ✅ 不影响资源删除逻辑
- ✅ `terraform destroy` 幂等

---

#### Scenario: API 调用失败

**Given**: API 调用因网络或权限问题失败  
**When**: 重试机制触发  
**Then**:
- 使用现有的重试逻辑
- 错误日志正确记录
- 最终返回错误给调用方

**验收标准**:
- ✅ 重试机制正常
- ✅ 错误信息清晰

---

### Requirement: REQ-CLB-TG-QUERY-004 - 性能和日志

The API replacement MUST NOT cause performance degradation, and logs MUST correctly reflect the new API calls.

API 替换不应导致性能下降,日志应正确反映新 API 调用。

#### Scenario: 分页查询性能

**Given**: 需要查询大量目标组(超过 100 个)  
**When**: 触发分页查询  
**Then**:
- 自动分页遍历所有结果
- 每页最多查询 20 条(CLB_PAGE_LIMIT)
- 总查询时间不超过旧 API

**验收标准**:
- ✅ 分页逻辑正确
- ✅ 无性能退化

---

#### Scenario: API 调用日志

**Given**: 启用详细日志(TF_LOG=DEBUG)  
**When**: 执行资源操作  
**Then**:
- 日志显示 `DescribeTargetGroupList` Action
- 请求和响应 JSON 正确记录
- 错误日志包含完整上下文

**验收标准**:
- ✅ 日志格式一致
- ✅ 可追踪 API 调用
- ✅ 便于问题排查

---

