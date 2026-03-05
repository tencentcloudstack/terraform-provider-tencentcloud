# Proposal: 为 tencentcloud_cfs_file_system 添加 Timeout 块支持

## Metadata

- **Change ID**: add-cfs-file-system-timeout-block
- **Status**: Draft
- **Author**: AI Assistant
- **Created**: 2026-03-05
- **Category**: Enhancement
- **Affected Resource**: `tencentcloud_cfs_file_system`

## Why

### 问题陈述

当前 `tencentcloud_cfs_file_system` 资源使用固定的超时时间进行状态等待,无法满足不同使用场景的需求:

1. **Turbo 系列文件系统创建慢**:
   - Standard/High-Performance 系列: 1-2 分钟完成
   - Turbo 系列: **10-20 分钟**才能完成
   - 当前固定超时: `3 * tccommon.ReadRetryTimeout` = 9 分钟
   - **问题**: Turbo 系列经常超时失败

2. **无法自定义超时配置**:
   - 用户无法根据实际情况调整超时时间
   - CI/CD 场景中可能需要更长或更短的超时

### 用户影响

- ❌ **当前**: Turbo 文件系统创建经常超时,需要手动重试
- ✅ **期望**: 用户可以根据文件系统类型自定义创建超时

### 用户体验示例

```hcl
resource "tencentcloud_cfs_file_system" "turbo" {
  name              = "turbo-filesystem"
  availability_zone = "ap-guangzhou-3"
  access_group_id   = "pgroup-xxxxxx"
  protocol          = "TURBO"
  storage_type      = "TB"
  net_interface     = "CCN"
  capacity          = 40960  # 40 TiB
  
  # 自定义创建超时
  timeouts {
    create = "30m"  # Turbo 系列创建需要 10-20 分钟
  }
}
```

## Technical Design

### 1. 添加 Timeouts 字段

在 `ResourceTencentCloudCfsFileSystem()` 的 schema 中添加:

```go
Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(20 * time.Minute),
},
```

**注意**: 只添加 Create timeout,因为:
- **Create**: 第 196-208 行有异步状态等待(`creating` -> `available`)
- **Update**: 所有操作都是同步 API 调用,无异步任务
- **Delete**: 删除是同步操作,无异步任务  
- **Read**: 只是查询操作,无异步任务

### 2. 修改 Create 函数使用 `d.Timeout()`

只修改 Create 函数的状态等待部分(第 196 行):

```go
// 修改前
err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
    // 等待状态变为 available
})

// 修改后  
err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
    // 等待状态变为 available
})
```

**不修改**第 172-189 行的 API 调用 retry,因为那是网络重试,不是异步等待。

## Implementation Plan

### 文件变更

| 文件 | 变更类型 | 说明 |
|------|---------|------|
| `tencentcloud/services/cfs/resource_tc_cfs_file_system.go` | 修改 | 添加 Timeouts 字段,修改 Create 状态等待 |
| `website/docs/r/cfs_file_system.html.markdown` | 修改 | 添加 Timeouts 文档 |

### 实施步骤

1. **代码修改** (0.5h)
   - 添加 Timeouts 字段到 schema
   - 修改 Create 函数状态等待使用 `d.Timeout()`

2. **测试验证** (0.5h)
   - 代码格式化和静态检查
   - 运行单元测试
   - 可选: 手动验证创建 Turbo 文件系统

3. **文档更新** (0.3h)
   - 更新用户文档添加 Timeouts 说明
   - 创建变更日志

4. **代码审查和合并** (0.2h)
   - 创建 PR
   - 代码审查
   - 合并到主分支

**总计**: 约 1.5 小时

## Risks & Mitigations

### 风险分析

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|---------|
| 破坏现有行为 | 低 | 高 | 默认值保持向后兼容,只添加新功能 |
| 文档不完整 | 中 | 中 | 添加详细使用示例 |
| 测试覆盖不足 | 低 | 中 | 运行现有测试确保无回归 |

### 向后兼容性

- ✅ **完全向后兼容**: 不修改默认行为,只添加可选配置
- ✅ **默认值合理**: 20 分钟足够覆盖 Turbo 系列创建时间
- ✅ **无破坏性变更**: 现有配置无需修改

## Testing Plan

### 单元测试

```bash
cd tencentcloud/services/cfs
go test -v -run TestAccTencentCloudCfsFileSystem
```

### 手动验证 (可选)

```hcl
# 1. 测试默认超时
resource "tencentcloud_cfs_file_system" "test1" {
  name              = "test-default-timeout"
  protocol          = "TURBO"
  storage_type      = "TB"
  # ... 其他必需配置
}

# 2. 测试自定义超时
resource "tencentcloud_cfs_file_system" "test2" {
  name              = "test-custom-timeout"
  protocol          = "TURBO"
  storage_type      = "TB"
  # ... 其他必需配置
  
  timeouts {
    create = "30m"
  }
}
```

## Success Criteria

- [ ] 代码编译通过
- [ ] 代码格式符合规范
- [ ] 单元测试通过
- [ ] 文档完整
- [ ] 代码审查通过
- [ ] Turbo 文件系统可以使用自定义超时成功创建

## References

- [Terraform ResourceTimeout 文档](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#ResourceTimeout)
- [CFS API 文档](https://cloud.tencent.com/document/product/582)
- 项目中其他资源的 timeouts 实现 (如 `tencentcloud_instance`)
