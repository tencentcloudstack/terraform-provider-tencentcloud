# Change: add-dnspod-domain-computed-fields

## Quick Summary

为 `tencentcloud_dnspod_domain_instance` 资源添加四个 computed 参数，以暴露 DNSPod API 返回的完整域名信息。

**Change Type**: Enhancement + Breaking Change  
**Status**: Proposal Phase ✅ Validated  
**Estimated Effort**: 4 hours  

---

## What's Being Changed

### New Computed Fields (3)
- ✨ `record_count` (int) - 域名下的 DNS 解析记录总数
- ✨ `grade` (string) - 域名套餐等级（DP_Free, DP_Plus, etc.）
- ✨ `updated_on` (string) - 域名最后更新时间

### Modified Field (1)
- ⚠️ `status` - **从 Optional 改为 Computed-only** (BREAKING CHANGE)

---

## Why This Change?

### Current Problem
用户无法通过 Terraform 获取域名的完整状态信息：
- ❌ 不知道域名下有多少条解析记录
- ❌ 不知道域名当前的套餐等级
- ❌ 不知道域名最后的修改时间
- ❌ `status` 字段语义混乱（既是配置参数又是状态输出）

### Solution
从 DNSPod API 的 `DomainInfo` 响应中读取并暴露这些字段，所有字段均为只读（Computed）。

---

## Breaking Changes

### ⚠️ `status` Field Behavior Change

**Before**: Optional + Computed (用户可以设置)
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
  status = "enable"  # ❌ 这在新版本中会报错
}
```

**After**: Computed-only (只读)
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
  # status 现在是只读字段，不能配置
}

# 通过输出访问状态
output "domain_status" {
  value = tencentcloud_dnspod_domain_instance.example.status
}
```

### Migration Required

如果你的配置中使用了 `status` 参数：
1. 移除资源块中的 `status = "..."` 行
2. 使用 `resource.status` 读取实际状态
3. 如需控制域名状态，等待未来的专用资源或使用其他方式

---

## User Benefits

✅ **完整的域名信息** - 一次查询获取所有域名元数据  
✅ **监控能力** - 追踪解析记录数量和域名更新时间  
✅ **计费可见性** - 了解域名当前的套餐等级  
✅ **语义清晰** - `status` 字段不再混淆配置和状态  
✅ **零额外成本** - 所有信息都在现有 API 响应中  

---

## Implementation Overview

### Files Modified
- `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go` - 核心实现
- `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance_test.go` - 测试更新
- `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.md` - 文档
- `website/docs/r/dnspod_domain_instance.html.markdown` - 网站文档

### Code Changes
1. **Schema**: 添加 3 个新字段，修改 status 为 Computed-only
2. **Read**: 映射 4 个字段从 API 响应到资源状态
3. **Create**: 移除 status 设置逻辑
4. **Update**: 移除 status 更新逻辑
5. **Tests**: 更新测试验证新字段

### Lines of Code
- Added: ~30 lines
- Removed: ~20 lines
- Net: +10 lines

---

## Technical Details

### API Mapping

| Terraform Field | DNSPod API Field | Type       | Null-safe |
|-----------------|------------------|------------|-----------|
| `record_count`  | `RecordCount`    | uint64→int | ✅         |
| `grade`         | `Grade`          | string     | ✅         |
| `status`        | `Status`         | string     | ✅         |
| `updated_on`    | `UpdatedOn`      | string     | ✅         |

### Example API Response
```json
{
  "DomainInfo": {
    "DomainId": 12345,
    "Domain": "example.com",
    "Status": "enable",
    "Grade": "DP_Free",
    "RecordCount": 5,
    "UpdatedOn": "2024-01-15 10:30:00",
    ...
  }
}
```

### Example Resource State After Read
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
  
  # Computed fields (read from API)
  domain_id    = 12345
  status       = "enable"
  grade        = "DP_Free"
  record_count = 5
  updated_on   = "2024-01-15 10:30:00"
  create_time  = "2024-01-01 08:00:00"
  slave_dns    = "no"
}
```

---

## Testing Strategy

### Unit Tests
- ✅ 验证新字段在资源创建后被设置
- ✅ 验证配置 `status` 字段会导致错误
- ✅ 验证 nil API 字段不会导致 panic

### Acceptance Tests
- ✅ 完整的 CRUD 生命周期测试
- ✅ 验证所有 computed 字段有实际值
- ✅ 验证 `record_count` 反映真实记录数

### Manual Tests
- ✅ 创建域名并检查新字段
- ✅ 添加解析记录后验证 `record_count` 变化
- ✅ 尝试设置 `status` 参数并确认报错

---

## Risks & Mitigation

### Risk: Breaking Change Impact
**Severity**: Medium  
**Affected Users**: 仅使用 `status` 参数的用户（预计较少）  
**Mitigation**:
- 清晰的 CHANGELOG 条目标记 BREAKING CHANGES
- 详细的迁移指南
- 示例代码展示迁移路径

### Risk: API Field Null Values
**Severity**: Low  
**Impact**: 如果 API 返回 nil 字段可能导致空值  
**Mitigation**:
- 所有字段映射都有 nil 检查
- 遵循 Terraform SDK 的 nil 值处理规范

### Risk: Type Conversion Errors
**Severity**: Low  
**Impact**: `uint64` → `int` 转换可能溢出  
**Mitigation**:
- DNSPod 域名记录数不会超过 int 最大值
- 实际使用中不会遇到溢出问题

---

## Timeline

| Phase | Tasks | Duration |
|-------|-------|----------|
| **Proposal** ✅ | 编写 proposal, tasks, spec | 完成 |
| **Implementation** | 修改代码、更新测试 | 2 小时 |
| **Testing** | 运行测试、手动验证 | 1 小时 |
| **Documentation** | 更新文档、写 CHANGELOG | 0.5 小时 |
| **Review** | 代码审查、最终验证 | 0.5 小时 |

**Total**: ~4 hours

---

## Next Steps

### To Implement (use `openspec apply`)
```bash
openspec apply add-dnspod-domain-computed-fields
```

### To Review Proposal
```bash
openspec show add-dnspod-domain-computed-fields
```

### To Validate
```bash
openspec validate add-dnspod-domain-computed-fields --strict
```

---

## Related Resources

**None** - 这是一个独立的资源增强，不依赖或影响其他资源。

未来可能的相关变更：
- 添加 `tencentcloud_dnspod_domain_status` 资源用于显式控制域名状态
- 为其他 DNSPod 资源添加类似的 computed 字段

---

## References

- **Proposal**: [proposal.md](./proposal.md)
- **Tasks**: [tasks.md](./tasks.md)
- **Spec**: [specs/dnspod-domain-instance-fields/spec.md](./specs/dnspod-domain-instance-fields/spec.md)
- **DNSPod API**: DescribeDomain / CreateDomain interfaces
- **SDK Source**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323/models.go:5171-5270`
- **Resource Code**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`

---

## Questions?

如有疑问，请参考：
1. [Proposal 详细设计](./proposal.md)
2. [Task 分解和验收标准](./tasks.md)
3. [Spec 详细需求](./specs/dnspod-domain-instance-fields/spec.md)
4. [OpenSpec AGENTS 指南](../../AGENTS.md)
