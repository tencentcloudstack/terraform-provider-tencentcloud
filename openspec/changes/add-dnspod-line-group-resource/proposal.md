# Proposal: 新增 DNSPod 线路分组资源 (tencentcloud_dnspod_line_group)

## 概述

**提案 ID**: `add-dnspod-line-group-resource`  
**变更类型**: Feature Addition (功能新增)  
**影响范围**: DNSPod 服务  
**实施优先级**: Medium

---

## 背景与动机

### 问题描述

DNSPod 提供了自定义线路分组功能，允许用户将多个解析线路（如"电信"、"移动"、"联通"等）组合成一个分组，以便在配置 DNS 解析记录时更灵活地进行线路管理。目前 Terraform Provider 中缺少对线路分组资源的管理能力，用户无法通过 IaC 方式创建、修改和删除线路分组。

### 用户需求

用户需要通过 Terraform 管理 DNSPod 的线路分组，实现：
1. 创建自定义线路分组
2. 查询现有线路分组信息
3. 修改线路分组的名称和线路列表
4. 删除不再需要的线路分组
5. 支持导入现有线路分组到 Terraform 状态

### 业务价值

- **自动化**: 实现 DNS 线路分组的自动化管理，减少手动操作
- **标准化**: 通过代码定义线路分组，保证环境一致性
- **可追溯**: 通过版本控制追踪线路分组的变更历史
- **集成性**: 与现有 DNSPod 资源（如解析记录）无缝集成

---

## 设计方案

### API 映射

腾讯云 DNSPod 提供了完整的线路分组管理 API：

| Terraform 操作 | 云 API 接口 | API 文档 |
|----------------|-------------|----------|
| Create | `CreateLineGroup` | https://cloud.tencent.com/document/api/1427/112219 |
| Read | `DescribeLineGroupList` | https://cloud.tencent.com/document/api/1427/112216 |
| Update | `ModifyLineGroup` | https://cloud.tencent.com/document/api/1427/112215 |
| Delete | `DeleteLineGroup` | https://cloud.tencent.com/document/api/1427/112217 |

**API 版本**: `2021-03-23`  
**请求域名**: `dnspod.tencentcloudapi.com`  
**频率限制**: 20次/秒

### 资源定义

#### 资源名称
```hcl
resource "tencentcloud_dnspod_line_group" "example" {
  ...
}
```

#### Schema 设计

**输入参数**:

| 字段名 | 类型 | 必填 | ForceNew | 说明 |
|--------|------|------|----------|------|
| `domain` | String | 是 | 是 | 域名（如 `dnspod.cn`）|
| `name` | String | 是 | 否 | 线路分组名称（1-17个字符）|
| `lines` | List(String) | 是 | 否 | 线路列表（如 `["电信", "移动"]`），最多120个 |
| `domain_id` | Integer | 否 | 是 | 域名 ID（优先级高于 domain）|

**输出参数（Computed）**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `line_group_id` | Integer | 线路分组 ID（云端生成）|
| `created_on` | String | 创建时间 |
| `updated_on` | String | 更新时间 |

#### 资源 ID 格式

采用复合 ID 格式：`{domain}#{line_group_id}`

示例：`dnspod.cn#123`

**设计理由**:
- 域名是资源的必要上下文（API 要求）
- 线路分组 ID 是云端唯一标识
- 使用 `#` 分隔符是项目标准

### 生命周期管理

#### Create 流程
1. 从 schema 读取 `domain`、`name`、`lines`
2. 构建 `CreateLineGroup` 请求
3. 调用 API 创建线路分组
4. 从响应获取 `line_group_id`
5. 构造资源 ID：`domain#line_group_id`
6. 调用 Read 函数刷新状态

#### Read 流程
1. 解析资源 ID 获取 `domain` 和 `line_group_id`
2. 调用 `DescribeLineGroupList` 查询指定域名的所有线路分组
3. 遍历结果，找到匹配 `line_group_id` 的分组
4. 将字段设置到 Terraform State
5. 如果分组不存在，清空 ID（触发重建）

**注意**: API 不支持按 ID 查询单个分组，需要过滤列表。

#### Update 流程
1. 检查 `domain` 是否变更（不可变字段）
2. 如果 `name` 或 `lines` 变更：
   - 构建 `ModifyLineGroup` 请求
   - 调用 API 更新线路分组
3. 调用 Read 函数刷新状态

#### Delete 流程
1. 解析资源 ID 获取 `domain` 和 `line_group_id`
2. 构建 `DeleteLineGroup` 请求
3. 调用 API 删除线路分组
4. 错误处理：如果分组已不存在，视为成功

#### Import 流程
支持导入现有线路分组，格式：`{domain}#{line_group_id}`

示例：
```bash
terraform import tencentcloud_dnspod_line_group.example dnspod.cn#123
```

### 字段映射细节

#### Lines 字段处理

**API 格式**: 字符串，逗号分隔（`"电信,移动"`）  
**Terraform 格式**: 字符串列表（`["电信", "移动"]`）

**转换逻辑**:
- **Create/Update**: `strings.Join(lines, ",")`
- **Read**: `strings.Split(apiLines, ",")`

**验证规则**:
- 不能为空列表
- 最多 120 个线路
- 每个线路名称为非空字符串

#### Name 字段验证

- 长度：1-17 个字符
- 不能与现有分组重名
- 不能使用系统内置线路名称

#### Domain 与 DomainId

- 两者任选其一，优先使用 `domain_id`
- `domain` 标记为 `ForceNew`（变更触发重建）
- `domain_id` 标记为 `ForceNew`

### 错误处理

#### 常见错误及处理策略

| 错误码 | 描述 | 处理策略 |
|--------|------|----------|
| `InvalidParameter.GroupNameOccupied` | 分组名已存在 | 返回错误，提示用户修改 |
| `InvalidParameter.LineGroupOverCounted` | 分组数量超限 | 返回错误，提示升级套餐 |
| `InvalidParameter.LineInAnotherGroup` | 线路已存在于其他分组 | 返回错误，提示先移除 |
| `InvalidParameter.LineNotExist` | 分组不存在（删除时）| 视为成功（幂等性）|
| `InvalidParameter.LineInUse` | 分组正在使用 | 返回错误，提示先移除依赖 |
| `RequestLimitExceeded` | 频率超限 | 自动重试（RetryError）|

#### 重试策略

- **写操作** (Create/Update/Delete): 使用 `tccommon.WriteRetryTimeout` (5分钟)
- **读操作** (Read): 使用 `tccommon.ReadRetryTimeout` (1分钟)
- **速率限制**: 通过 `ratelimit.Check()` 控制

### 数据一致性

#### Read 操作的一致性保证

由于 `DescribeLineGroupList` 返回列表，需要确保：
1. 列表为空时，正确设置 `d.SetId("")`
2. 分组不存在时，记录警告日志
3. 字段为 nil 时，跳过 `d.Set()`

#### 状态更新策略

所有 CRUD 操作完成后调用 Read 函数，确保状态与云端一致。

---

## 实施计划

### 阶段划分

**Phase 1: 核心功能** (预计 2 天)
- 创建资源文件和 Schema 定义
- 实现 CRUD 操作
- 实现 Import 功能

**Phase 2: Service 层** (预计 0.5 天)
- 添加 Service 层辅助方法（如果需要）
- 实现错误处理和重试逻辑

**Phase 3: 测试** (预计 1 天)
- 编写验收测试
- 测试各种场景（创建、更新、删除、导入）
- 边缘情况测试

**Phase 4: 文档与集成** (预计 0.5 天)
- 编写资源文档
- 在 provider.go 中注册
- 代码质量检查

### 文件清单

| 文件 | 说明 |
|------|------|
| `resource_tc_dnspod_line_group.go` | 资源实现（Create/Read/Update/Delete）|
| `resource_tc_dnspod_line_group_test.go` | 验收测试 |
| `resource_tc_dnspod_line_group.md` | 资源文档 |
| `service_tencentcloud_dnspod.go` | Service 层方法（如需添加）|
| `provider.go` | 注册资源到 Provider |

---

## 示例代码

### 基础示例

```hcl
# 创建线路分组
resource "tencentcloud_dnspod_line_group" "example" {
  domain = "example.com"
  name   = "telecom-group"
  lines  = ["电信", "移动"]
}

# 输出线路分组 ID
output "line_group_id" {
  value = tencentcloud_dnspod_line_group.example.line_group_id
}
```

### 完整示例

```hcl
# 使用 domain_id
resource "tencentcloud_dnspod_line_group" "advanced" {
  domain_id = 1005
  name      = "custom-isp-group"
  lines     = [
    "电信",
    "联通",
    "移动",
    "铁通",
  ]
}

# 导入现有线路分组
# terraform import tencentcloud_dnspod_line_group.imported example.com#123
```

### 与解析记录集成（未来扩展）

```hcl
# 创建线路分组
resource "tencentcloud_dnspod_line_group" "isp" {
  domain = "example.com"
  name   = "main-isp"
  lines  = ["电信", "联通", "移动"]
}

# 使用线路分组创建解析记录（未来功能）
# resource "tencentcloud_dnspod_record" "example" {
#   domain      = "example.com"
#   sub_domain  = "www"
#   record_type = "A"
#   value       = "1.2.3.4"
#   record_line = tencentcloud_dnspod_line_group.isp.name
# }
```

---

## 风险评估

### 技术风险

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| API 限流导致操作失败 | 中 | 中 | 实现重试机制，添加速率限制检查 |
| Lines 字段格式解析错误 | 低 | 中 | 充分测试字符串分割和拼接逻辑 |
| 列表查询性能问题 | 低 | 低 | API 返回速度快，分组数量有限 |

### 兼容性风险

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| SDK 版本不兼容 | 低 | 高 | 已验证 SDK 包含所需 API |
| API 变更 | 低 | 中 | 遵循腾讯云 API 版本控制规范 |

### 用户体验风险

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 分组名重复导致创建失败 | 中 | 低 | 提供清晰的错误提示 |
| 线路在其他分组中导致冲突 | 中 | 中 | 返回详细错误信息，指导用户解决 |
| ForceNew 导致意外重建 | 低 | 中 | 文档中明确说明不可变字段 |

---

## 验收标准

### 功能验收

- [ ] 可以成功创建线路分组
- [ ] 可以查询线路分组信息
- [ ] 可以修改线路分组名称
- [ ] 可以修改线路分组的线路列表
- [ ] 可以删除线路分组
- [ ] 支持导入现有线路分组
- [ ] Domain 字段变更触发重建

### 质量验收

- [ ] 所有验收测试通过（`TF_ACC=1 go test`）
- [ ] 代码通过 `go fmt` 格式化
- [ ] 代码通过 `golangci-lint` 检查
- [ ] 文档包含完整示例和参数说明
- [ ] 错误处理覆盖所有 API 错误码

### 性能验收

- [ ] Create 操作在 10 秒内完成
- [ ] Read 操作在 5 秒内完成
- [ ] Update 操作在 10 秒内完成
- [ ] Delete 操作在 10 秒内完成

---

## 待决事项

1. **Service 层方法**: 是否需要在 `service_tencentcloud_dnspod.go` 中添加 `DescribeDnspodLineGroupById()` 方法？
   - **建议**: 需要，用于 Read 操作，简化逻辑
   
2. **Lines 字段排序**: 是否需要对 lines 列表排序以避免 drift？
   - **建议**: 不需要，保持用户定义的顺序，API 也不强制排序

3. **DomainId 字段**: 是否暴露给用户？
   - **建议**: 可选暴露，优先级高于 domain，方便高级用户使用

4. **与记录资源的集成**: 是否在此提案中实现？
   - **建议**: 不在此提案范围，作为未来扩展

---

## 参考资料

- [DNSPod CreateLineGroup API](https://cloud.tencent.com/document/api/1427/112219)
- [DNSPod DescribeLineGroupList API](https://cloud.tencent.com/document/api/1427/112216)
- [DNSPod ModifyLineGroup API](https://cloud.tencent.com/document/api/1427/112215)
- [DNSPod DeleteLineGroup API](https://cloud.tencent.com/document/api/1427/112217)
- [Terraform Plugin SDK v2 文档](https://developer.hashicorp.com/terraform/plugin/sdkv2)
- [项目现有资源示例](tencentcloud/services/dnspod/resource_tc_dnspod_record_group.go)

---

## 变更历史

| 日期 | 版本 | 作者 | 变更说明 |
|------|------|------|----------|
| 2026-01-29 | v1.0 | AI Agent | 初始提案创建 |
